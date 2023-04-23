package check

import (
	"fmt"
	"strings"

	schema2 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

func crossCheckProperty(r *schema.Resource, md *model.ResourceDoc) (res []Checker) {
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

func diffDocMiss(rt, path string, s *schema2.Schema, f *model.Field) (res []Checker) {
	// skip deprecated property
	if shouldSkipDocProp(rt, path) {
		return
	}

	if isSkipProp(rt, path) {
		return
	}

	if f == nil {
		if s.Deprecated == "" && !s.Computed && path != "id" {
			res = append(res, newMissInDoc(path, f))
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
			res = append(res, newMissInDoc(path+" not block", f))
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
	"azurerm_application_gateway": {
		"backend_http_settings.authentication_certificate.data",
	},
	"azurerm_vpn_server_configuration": {
		"*",
	},
}

func shouldSkipCodeProp(rt, path string) bool {
	return _shouldSkip(diffCodeSkip, rt, path)
}

func diffCodeMiss(rt, path string, f *model.Field, s *schema2.Schema) (res []Checker) {
	if shouldSkipCodeProp(rt, path) {
		return
	}
	if isSkipProp(rt, path) {
		return
	}

	if f != nil && f.FormatErr {
		res = append(res, newFormatErr(f.Content, newCheckBase(f.Line, path, f)))
		return
	}

	if s == nil {
		if path != "id" && f != nil { // id not defined in code
			if strings.TrimSpace(path) == "" {
				path = fmt.Sprintf("%s:L%d", f.Name, f.Line)
			}
			if strings.Contains(strings.ToLower(f.Content), "deprecated") {
				path += " deprecated"
			}
			// not available for some block
			if idx := strings.Index(strings.ToLower(f.Content), "not available for"); idx > 0 {
				if code := util.FirstCodeValue(f.Content[idx:]); code != "" && strings.Contains(path, code) {
					return res
				}
			}
			res = append(res, newMissInCode(path, f))
		}
		return res
	}

	if f == nil {
		return nil
	}
	base := newCheckBase(f.Line, path, f)

	// check optional. optional&computed property diff
	if (f.Required != model.Required) && s.Required {
		res = append(res, newRequireDiff(base, ShouldBeRequired))
	} else if s.Optional {
		if f.Required != model.Optional && f.Pos == model.PosArgs {
			res = append(res, newRequireDiff(base, ShouldBeOptional))
		}
		if s.Computed {
			// optional and computed, but not in attribute part
			if f.SameNameAttr == nil {
				// todo log this kind of miss in attribute
				// res = append(res, newMissItem(path, f, MissInDocAttr))
			} else if f.SameNameAttr.Required > 0 { // attribute should not have requriedness spec
				base2 := newCheckBase(f.SameNameAttr.Line, path, f.SameNameAttr)
				res = append(res, newRequireDiff(base2, ShouldBeComputed))
			}
		}
	}

	// check default values
	if s.Default != nil {
		defaultStr := fmt.Sprintf("%v", s.Default)
		if str, ok := s.Default.(string); ok && str == "" {
			defaultStr = `""` // empty string in code
		}
		// for many default value is `false`, just skip them for now
		if defaultStr != f.Default && defaultStr != "false" {
			// todo remove this if strings.Contains check, it's not accurate (but works)
			if !strings.Contains(f.Content, defaultStr) {
				res = append(res, newDefaultDiff(base, f.Default, defaultStr))
			}
		}
	} else if f.Default != "" && !s.Computed {
		// code no default and not computed/optional property, but document has
		// if default to is a long sentence, then skip it now. todo add more logic to analysis
		res = append(res, newDefaultDiff(base, f.Default, ""))
	}

	// check forceNew attribute
	if s.ForceNew != f.ForceNew && f.Name != "resource_group_name" {
		var forceNew = ForceNewDefault
		if s.ForceNew && !f.ForceNew {
			forceNew = ShouldBeForceNew
		} else if f.ForceNew && !s.ForceNew {
			forceNew = ShouldBeNotForceNew
		}
		res = append(res, newForceNewDiff(base, forceNew))
	}

	// if code schema is not list/set and md field is attr, then skip iterate sub-fields even exists
	// for we guess a md property as block if not found other block-type properties
	if s.Type != schema2.TypeList && f.Typ == model.FieldTypeAttr {
		return res
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

	return res
}
