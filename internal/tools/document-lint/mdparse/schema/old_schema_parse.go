package schema

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/tools/go/packages"
)

// ATTENTION: this is the first version using AST Analysis to get the possible values.
// it's not used for now, we use gomokey to get possible values.

func (r *Resource) GetImportPath(name string) string {
	if dsc, ok := r.Imports[name]; ok {
		return litStr(dsc.Path)
	}
	return ""
}

// FindAllInSliceProp not used for current version , use gomonkey instead of ast parse
func (r *Resource) FindAllInSliceProp() {
	for name, item := range r.Schema.Schema {
		r.InSliceProp(name, item)
	}
}

func (r *Resource) InSliceProp(name string, item *schema.Schema) {
	if item.ValidateFunc != nil {
		// check if it is StringsInSlice
		pc := reflect.ValueOf(item.ValidateFunc).Pointer()
		fn := runtime.FuncForPC(pc)
		if strings.Contains(fn.Name(), "StringInSlice") {
			values := r.ValidateFor(name)
			r.PossibleValues[name] = values
		}
	}
	switch ele := item.Elem.(type) {
	case *schema.Resource:
		for subName, prop := range ele.Schema {
			r.InSliceProp(name+"."+subName, prop)
		}
	case *schema.Schema:
		r.InSliceProp(name, ele)
	}
}

// ValidateFor extract validate function from source file
// propName can be "name", "prop1.prop2"
func (r *Resource) ValidateFor(propName string) (values []string) {
	return r.Package.ValidateForPropPath(propName, r.FilePath)
}

func (r *Resource) ValidteForNode(parts []string, node ast.Node) string {
	var name string
	ast.Inspect(node, func(node ast.Node) bool {
		//ts := r.TokenSet
		switch val := node.(type) {
		case *ast.MapType:

		case *ast.KeyValueExpr:
			// if kv like `propName: abc` then try to extract it
			lit, ok := val.Key.(*ast.BasicLit)
			if !ok {
				return true
			}
			if litStr(lit) == parts[0] {
				if len(parts) == 1 {
					name = parts[0]
					possibleValues := r.InspectSchemaValidation(strings.Join(parts, "."), val.Value)
					log.Printf("possible value for %s: %v", name, possibleValues)
					return false
				} else {
					name = r.ValidteForNode(parts[1:], val.Value)
					return false
				}
			}
		case *ast.BasicLit:
			//if val.Value == propName {
			//	name = val.Value
			//}
		}
		return true
	})
	return name
}

// InspectSchemaValidation node is the Schema of specific property
// try to find the ValidateFunc and extract possible values from it
func (r *Resource) InspectSchemaValidation(name string, node ast.Expr) (values []string) {
	ast.Inspect(node, func(node ast.Node) bool {
		switch val := node.(type) {
		case *ast.KeyValueExpr:
			if lit, ok := val.Key.(*ast.Ident); ok {
				if lit.Name == "ValidateFunc" || lit.Name == "ValidateDiagFunc" {
					log.Printf("ast for %s: %s", name, lit.Name)
					// get first argument
					// can be function call, or slice of const or slice of liter string
					if call, ok := val.Value.(*ast.CallExpr); ok {
						if funSel, ok := call.Fun.(*ast.SelectorExpr); ok && funSel.Sel.Name == "StringInSlice" {
							arg0 := call.Args[0]
							caseSense := call.Args[1].(*ast.Ident).Name == "false"
							_ = caseSense
							//_ = ast.Print(r.TokenSet, arg0)
							values = r.InspectInSliceArgs(arg0)
							return false
						}
					}
					return false // stop walker
				}
			}
		}
		return true
	})
	return
}

// InspectInSliceArgs extract possible values from validateFunc args0
// node is the args[0] of function StringInSlice
// it can be a literal string slice, a slice of consts, of a function call
func (r *Resource) InspectInSliceArgs(node ast.Expr) (values []string) {
	ast.Inspect(node, func(node ast.Node) bool {
		switch val := node.(type) {
		case *ast.CompositeLit:
			for _, el := range val.Elts {
				// current package
				if val := r.EvaluateValue("", el); val != "" {
					values = append(values, val)
				}
			}
		case *ast.CallExpr:
			// call function to get all items
			_ = val
			values = r.EvaluateFunCall(val)
		}
		return false
	})
	return
}

// EvaluateValue evaluate one item: literal string or const define
// ele is one "value", string(abc.def), string(var), abc, or abc.var
func (r *Resource) EvaluateValue(pack string, ele ast.Expr) (res string) {
	switch val := ele.(type) {
	case *ast.BasicLit:
		// "abc"
		return litStr(val)
	case *ast.Ident:
		// const variable
		log.Printf("todo: wait ident for %s: %v", pack, ele)
	case *ast.CallExpr:
		// string(xxx)
		if fun, ok := val.Fun.(*ast.Ident); ok && fun.Name == "string" {
			switch arg0 := val.Args[0].(type) {
			case *ast.SelectorExpr:
				packName, sel := r.SelectorXYWithLoad(arg0)
				if p1, ok := r.cacheConsts[packName]; ok {
					return p1[sel]
				}
			case *ast.Ident:
				// const value
				if pc, ok := r.cacheConsts[pack]; ok {
					return pc[arg0.Name]
				} else {
					log.Printf("no const for pack: %v", pack)
				}
			}
		}
	}
	return
}

func (r *Resource) SelectorXYWithLoad(node *ast.SelectorExpr) (x, y string) {
	x, y = SelectorXY(node)
	r.GetDepPackage(x)
	return
}

// EvaluateFunCall ecaluate values from function call. can be like pack.PossibleValuesForXXX() or
// closure like func() []string {}()
func (r *Resource) EvaluateFunCall(node *ast.CallExpr) (values []string) {
	//_ = ast.Print(r.TokenSet, node)

	switch fn := node.Fun.(type) {
	case *ast.Ident:
		// direct call local package: `ValidateFunc: validation.StringInSlice(possibleForTest(), false),`
		// or local function call like stringInSlice(possibleValueForXXX(), false)
		return r.cacheFuncValues["."+fn.Name] // local function names
	case *ast.SelectorExpr:
		// call other functions `ValidateFunc: validation.StringInSlice(nginxdeployment.PossibleValuesForNginxPrivateIPAllocationMethod(), false),`
		pack, sel := r.SelectorXYWithLoad(fn)
		// get all values from possible function
		return r.cacheFuncValues[pack+"."+sel]
	case *ast.FuncLit:
		// closure function direct call
		// inspect body to find range possible values
		// 	ValidateFunc: validation.StringInSlice(func() (res []string) {
		//		for _, s := range nginxdeployment.PossibleValuesForNginxPrivateIPAllocationMethod() {
		//			res = append(res, s)
		//		}
		//		return res
		//	}(), false),
		values = append(values, r.ValuesFromFuncBody("", fn.Body)...)
	}
	return
}

// ValuesFromFuncBody extract possible values from function body
// 1. body contains range string
// 2. body is return []string{...}
func (r *Resource) ValuesFromFuncBody(pack string, body *ast.BlockStmt) (values []string) {
	_ = ast.Print(r.TokenSet, body)
	ast.Inspect(body, func(node ast.Node) bool {
		switch st := node.(type) {
		case *ast.UnaryExpr:
			// fetch all range values
			if st.Op == token.RANGE {
				switch x := st.X.(type) {
				case *ast.CallExpr:
					values = append(values, r.EvaluateFunCall(x)...)
				case *ast.CompositeLit:
					values = append(values, r.extractFromSlice(pack, x)...)
				}
			}
		case *ast.CompositeLit:
			values = append(values, r.extractFromSlice(pack, st)...)
		case *ast.ReturnStmt:
			if len(st.Results) == 0 {
				return true
			}
			switch ret := st.Results[0].(type) {
			case *ast.CompositeLit:
				values = append(values, r.extractFromSlice(pack, ret)...)
			}
		}
		return true
	})
	return
}

// extract values from []string(xxx,xxx), element can be function call, literal string, or both of them
func (r *Resource) extractFromSlice(pack string, node *ast.CompositeLit) (values []string) {
	for _, el := range node.Elts {
		values = append(values, r.EvaluateValue(pack, el))
	}
	return
}

func (r *Resource) InspectValidateValues(node ast.Expr) (values []string) {
	return
}

func (r *Resource) GetDepPackage(pack string) *packages.Package {
	if _, ok := r.cacheDep[pack]; !ok {
		packPath := r.GetImportPath(pack)
		if packs := r.loadPackage(packPath); len(packs) > 0 {
			// cache consts values
			r.cacheConsts[pack] = map[string]string{}
			r.cacheDep[pack] = packs[0]
			//cahce result of possible values funtion
			//r.CachePossibleFunction(pack, packs[0].Syntax)
		}
	}

	return r.cacheDep[pack]
}

func (r *Resource) loadPackage(patterns ...string) []*packages.Package {
	//conf := types.Config{
	//	Importer: importer.Default(),
	//}
	//conf.Check("", "")
	cfg := &packages.Config{
		Mode: packages.NeedName |
			//packages.NeedImports |
			packages.NeedSyntax |
			packages.NeedFiles |
			packages.NeedModule,
	}
	packs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Printf("import %v: %v", patterns, err)
		return nil
	}
	return packs
}
