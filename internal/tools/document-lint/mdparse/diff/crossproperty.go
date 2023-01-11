package diff

import (
	"fmt"
	"strings"

	schema2 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
)

func crossCheckProperty(r *schema.Resource, md *model.ResourceDoc) (res []DiffItem) {
	// check property exists in r not md, or in md not in r

	// exist in tf schema but not in document
	docProps := md.AllProp()
	for key, val := range r.Schema.Schema {
		field := docProps[key]
		if field == nil && md.Attr != nil {
			field = md.Attr[key]
		}
		res = append(res, diffDocMiss(r.ResourceType, key, val, field)...)
	}

	for key, f := range docProps {
		sub := r.Schema.Schema[key]
		subDiff := diffCodeMiss(r.ResourceType, key, f, sub)
		res = append(res, subDiff...)
	}
	return
}

var diffDocSkip = map[string][]string{
	"azurerm_application_gateway": {
		"backend_http_settings.authentication_certificate.data",
	},
}

func _shouldSkip(m map[string][]string, rt, path string) bool {
	if m2, ok := m[rt]; ok {
		for _, v := range m2 {
			if v == "*" || v == path {
				return true
			}
		}
	}
	return false
}

func shouldSkipDocProp(rt, path string) bool {
	return _shouldSkip(diffDocSkip, rt, path)
}

func diffDocMiss(rt, path string, s *schema2.Schema, f *model.Field) (res []DiffItem) {
	// skip deprecated property
	if shouldSkipDocProp(rt, path) {
		return
	}

	if isSkipProp(rt, path) {
		return
	}

	if f == nil {
		if s.Deprecated == "" && !s.Computed && path != "id" {
			res = append(res, NewMissDiffItem(path, MissInDoc))
		}
		return res
	}
	if s == nil || s.Elem == nil {
		return nil
	}

	switch ele := s.Elem.(type) {
	case *schema2.Schema:
		return nil
	case *schema2.Resource:
		if f.Subs == nil {
			res = append(res, NewMissDiffItem(path+" not block", MissInDoc))
			return
		}
		for key, val := range ele.Schema {
			subField := f.Subs[key]
			res = append(res, diffDocMiss(rt, path+"."+key, val, subField)...)
		}
	default:
		return
	}
	return
}

var diffCodeSkip = map[string][]string{
	"azurerm_application_gateway": []string{
		"backend_http_settings.authentication_certificate.data",
	},
	"azurerm_vpn_server_configuration": {
		"*",
	},
}

func shouldSkipCodeProp(rt, path string) bool {
	return _shouldSkip(diffCodeSkip, rt, path)
}

func diffCodeMiss(rt, path string, f *model.Field, s *schema2.Schema) (res []DiffItem) {
	if shouldSkipCodeProp(rt, path) {
		return
	}
	if isSkipProp(rt, path) {
		return
	}
	if s == nil {
		if path != "id" { // id not defined in code
			if strings.TrimSpace(path) == "" {
				path = fmt.Sprintf("%s:L%d", f.Name, f.Line)
			}
			if strings.Contains(strings.ToLower(f.Content), "deprecated") {
				path += " deprecated"
			}
			res = append(res, NewMissDiffItem(path, MissInCode))
		}
		return res
	}
	if f == nil {
		return nil
	}

	// check optional. optional&computed property diff
	f.TFRequied = s.Required
	if (f.Required != model.Required) && s.Required {
		res = append(res, NewRequiredDiffItem(path, f, ShouldBeRequired))
	} else if s.Optional {
		if !s.Computed {
			if f.Required != model.Optional {
				res = append(res, NewRequiredDiffItem(path, f, ShouldBeOptional))
			}
		} else {
			// optional computed
			if f.SameNameAttr == nil {
				// todo log this kind of miss in attribute
				//res = append(res, NewMissDiffItem(path, MissInDocAttr))
			} else if f.SameNameAttr.Required > 0 { // attribute should not have requriedness spec
				res = append(res, NewRequiredDiffItem(path, f.SameNameAttr, ShouldBeComputed))
			}
		}
	}

	// check default values
	if s.Default != nil {
		defaultStr := fmt.Sprintf("%v", s.Default)
		if str, ok := s.Default.(string); ok && str == "" {
			defaultStr = `""` // empty string in document
		}
		// for many default value is `false`, just skip them for now
		if defaultStr != f.Default && defaultStr != "false" {
			// todo remove this if strings.Contains check, it's not accurate (but works)
			if !strings.Contains(f.Content, defaultStr) {
				res = append(res, NewDefaultDiff(path, f, defaultStr))
			}
		}
	} else if f.Default != "" && !s.Computed {
		// code no default and not computed/optional property, but document has
		// if default to is a long sentence, then skip it now. todo add more logic to analysis
		res = append(res, NewDefaultDiff(path, f, ""))
	}

	// check forceNew attribute
	if s.ForceNew != f.ForceNew && f.Name != "resource_group_name" {
		var forceNew = ForceNewDefault
		if s.ForceNew && !f.ForceNew {
			forceNew = ShouldBeForceNew
		} else if f.ForceNew && !s.ForceNew {
			forceNew = ShouldBeNotForceNew
		}
		res = append(res, NewFoceNewDiff(path, f, forceNew))
	}

	var subRes *schema2.Resource
	if res, ok := s.Elem.(*schema2.Resource); ok {
		subRes = res
	}
	// doc has sub-field but schema has no
	subTF := func(name string) *schema2.Schema {
		if subRes == nil || subRes.Schema == nil {
			return nil
		}
		return subRes.Schema[name]
	}

	for _, subField := range f.Subs {
		subPath := path + "." + subField.Name
		sub := subTF(subField.Name)
		res = append(res, diffCodeMiss(rt, subPath, subField, sub)...)
	}

	return
}
