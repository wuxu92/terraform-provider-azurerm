package schema

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"path"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"golang.org/x/tools/go/packages"
)

/// types and functions to process dependent package consts

type Package struct {
	Path   string
	Name   string // import name of package
	pack   *packages.Package
	Tokens *token.FileSet

	// import spec in origin file
	OriginImportSpec *ast.ImportSpec

	ConstValues map[string]string // extract all const values of package
	// cache all values for possible value functions, of all functions like `func() []string`
	PossibleFuncValues map[string][]string // string slice of all PossibleValuesForXXX function

	// save all validate string in slice possible values. key is prop/propa.propb, values is the value of string slice
	ValidateInSliceValues map[string][]string

	// todo recursive import package consts
	Import map[string]*Package // full import path to Package

	// mapping resource type to filepath, for sdk typed resource create.Func May defined in other file
	ResourceTypeFile map[string]string

	// shared between methods, so Package's method is not safe for goroutines
	// todo it's really a bad design to use a shared pointer to current processing file. it's hard to maintain this
	// remove this, for shared package it cause mess
	//currentPath string    // cache current process file path for log
	//currentFile *ast.File // pointer to current process file

	PartSchemaFunc map[string]PartSchemaFunc // cache functions for create schemas, need file path too
}

type PartSchemaFunc struct {
	FilePath   string
	SchemaNode ast.Node
}

func NewPartSchemaFunc(path string, node ast.Node) PartSchemaFunc {
	return PartSchemaFunc{
		FilePath:   path,
		SchemaNode: node,
	}
}

func (p *Package) Printf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s", fmt.Sprintf(format, args...)))
}

func ImportPathByName(name string, fs *ast.File) (importPath string) {
	for _, spec := range fs.Imports {
		importPath = litStr(spec.Path)
		// import alias equals, return directly
		if spec.Name != nil && spec.Name.Name == name {
			return
		}
		if path.Base(importPath) == name {
			return importPath
		}
	}
	return ""
}

func NewPackage(path string) *Package {
	p := &Package{
		Path:               path,
		ConstValues:        map[string]string{},
		PossibleFuncValues: map[string][]string{},
	}
	cfg := &packages.Config{
		Mode: packages.NeedName |
			//packages.NeedImports |
			packages.NeedSyntax |
			packages.NeedFiles |
			packages.NeedModule,
	}
	packs, err := packages.Load(cfg, path)
	if err != nil {
		p.Printf("load package: %s: %v", path, err)
		return nil
	}
	p.pack = packs[0]
	p.Tokens = cfg.Fset
	p.ExtractConsts() // auto extact all consts
	p.Import = map[string]*Package{}
	p.PartSchemaFunc = map[string]PartSchemaFunc{}
	p.ResourceTypeFile = map[string]string{}
	p.CachePossibleFunction()
	return p
}

func (p *Package) ExtractConsts() {
	if p.pack == nil {
		return
	}
	for _, node := range p.pack.Syntax {
		p.extractConstFromAst(node)
	}
}

func (p *Package) extractConstFromAst(node *ast.File) {
	ast.Inspect(node, func(node ast.Node) bool {
		switch val := node.(type) {
		case *ast.GenDecl:
			// also it can be defined in a var decl block
			if val.Tok == token.CONST || val.Tok == token.VAR {
				p.extractConstDecl(val)
			}
		case *ast.AssignStmt:
		default:
		}
		return true
	})
}

func (p *Package) extractConstDecl(node *ast.GenDecl) {
	ast.Inspect(node, func(node ast.Node) bool {
		switch val := node.(type) {
		case *ast.ValueSpec:
			// get all values
			// only process like a = b, no iota consts define
			if len(val.Values) == 0 {
				return false
			}
			for idx, expr := range val.Names {
				key := expr.Name
				value := BasicLitValue(val.Values[idx]) // no matter string or other type const
				if value != "" {
					p.ConstValues[key] = value
				}
			}
		}
		return true
	})
}

// CachePossibleFunction extract values for PossibleValuesForXXX functioo
// should always be like return []string{"xxx", string(XXX), ...}
func (p *Package) CachePossibleFunction() {
	for idx, f := range p.pack.Syntax {
		currentPath := p.pack.GoFiles[idx]
		ast.Inspect(f, func(node ast.Node) bool {
			switch fn := node.(type) {
			// variable decl as schema
			case *ast.GenDecl:
				if fn.Tok == token.VAR {
					//if fn.Specs[0].(*ast.ValueSpec).Names[0].Name == "expressRoutePortSchema" {
					//	log.Printf("%v", fn)
					//}
					for _, spec := range fn.Specs {
						vs := spec.(*ast.ValueSpec)
						for _, v := range vs.Values {
							if sel := NewNode(v).UnaryX().CompositeType().Selector().Sel(); sel != nil {
								if sel.Name == "Schema" || sel.Name == "Resource" {
									p.PartSchemaFunc[vs.Names[0].Name] = NewPartSchemaFunc(currentPath, vs.Values[0])
								}
							}
						}
					}
				}
			case *ast.FuncDecl:
				// all function with no arg and returns []string
				// this function can have receiver, and ignore all parameters
				typ := fn.Type
				//if (typ.TypeParams == nil || len(typ.TypeParams.List) == 0) &&
				if typ.Results != nil && len(typ.Results.List) == 1 {
					// check reuslt type is []string
					// if function is resource type's ResourceType function definition
					if IdentName(fn.Name) == "ResourceType" {
						ast.Inspect(fn.Body, func(node ast.Node) bool {
							switch lit := node.(type) {
							case *ast.BasicLit:
								p.ResourceTypeFile[litStr(lit)] = currentPath
								return false
							}
							return true
						})
						return true
					}

					// other function call
					switch realTyp := typ.Results.List[0].Type.(type) {
					case *ast.ArrayType:
						// result type is string or basic type is string
						//if ident, ok := realTyp.Elt.(*ast.Ident); ok && ident.Name == "string" {
						if IsBasicStringType(realTyp.Elt) {
							// catch this function body
							values := p.ValuesForFuncBody(currentPath, fn.Body)
							p.PossibleFuncValues[fn.Name.Name] = values
						}
					case *ast.StarExpr:
						// functions return *schema.Schema/*pluginsdk.Schema
						if selector, ok := realTyp.X.(*ast.SelectorExpr); ok {
							if IsPartSchemaSel(selector) {
								p.PartSchemaFunc[fn.Name.Name] = NewPartSchemaFunc(currentPath, fn.Body)
							}
						}
					case *ast.MapType:
						// return part Resource define return map[string]*plugin.Schema:
						if key, ok := realTyp.Key.(*ast.Ident); ok && key.Name == "string" {
							if un, ok := realTyp.Value.(*ast.StarExpr); ok {
								if selector, ok := un.X.(*ast.SelectorExpr); ok {
									if IsPartSchemaSel(selector) {
										p.PartSchemaFunc[fn.Name.Name] = NewPartSchemaFunc(currentPath, fn.Body)
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return
}

// GetFileAST find schema definition node and find all StringInSlice validate function
// set currentFile as filePath
func (p *Package) GetFileAST(filePath string) *ast.File {
	var fileAst *ast.File
	for idx, fpath := range p.pack.GoFiles {
		if fpath == filePath {
			fileAst = p.pack.Syntax[idx]
			break
		}
	}
	//p.currentPath = filePath
	//p.currentFile = fileAst
	return fileAst
}

// ValidateForPropPath path is the property path format as: propa.propb.prop
// filepath is the full path of the source code file
// return the possible values for this path property
func (p *Package) ValidateForPropPath(path string, filePath string) (values []string) {
	fileAST := p.GetFileAST(filePath)

	pathSep := strings.Split(path, ".")
	file, node, pack := p.FindPropNode(filePath, pathSep, fileAST)
	if node == nil {
		p.Printf("no such node for path: %v in %s", path, strings.TrimPrefix(filePath, "/home/wuxu/terraform-provider-azurerm/internal/"))
		return
	}
	pathSep = strings.Split(path, ".")
	return pack.ValidateInSliceFunc(file, pathSep, node)
}

func (p *Package) GetPartialSchema(file string, call ast.Expr) (node PartSchemaFunc, pack *Package) {
	if f, ok := call.(*ast.CallExpr); ok {
		call = f.Fun
	}
	switch fn := call.(type) {
	case *ast.Ident:
		if partial, ok := p.PartSchemaFunc[fn.Name]; ok {
			node = partial
			p.GetFileAST(partial.FilePath)
			pack = p
		}
	case *ast.SelectorExpr:
		// or it can be x.y.z.a()
		x, y := SelectorXY(fn)
		if packSub := p.PackageByName(file, x); packSub != nil {
			if partial, ok := packSub.PartSchemaFunc[y]; ok {
				node = partial
				packSub.GetFileAST(partial.FilePath)
				pack = packSub
			}
		} else {
			// todo: may cause wrong logic
			// try with local package
			if partial, ok := p.PartSchemaFunc[y]; ok {
				node = partial
				p.GetFileAST(partial.FilePath)
				pack = p
			}
		}
	}
	return
}

// FindPropNode path is ["prop1", "subprop2", ...]
// it can reference to import package, so... return additional package pointer
// if it in current pack then pack is p
func (p *Package) FindPropNode(filepath string, path []string, root ast.Node) (file string, node ast.Node, pack *Package) {
	if len(path) == 0 {
		return
	}
	file = filepath // default use origin filepath

	//oldPath := p.currentPath
	//if filepath != oldPath {
	//	p.GetFileAST(filepath)
	//}
	// can be in other file as function call
	tmpKey := fmt.Sprintf(`"%s"`, path[0])
	ast.Inspect(root, func(sub ast.Node) bool {
		switch kv := sub.(type) {
		case *ast.KeyValueExpr:
			// "" prop
			switch key := kv.Key.(type) {
			case *ast.Ident:
				// if Schema and Elem Ndoe is function call
				if key.Name == "Schema" || key.Name == "Elem" {
					if call, ok := kv.Value.(*ast.CallExpr); ok {
						if subNode, subPack := p.GetPartialSchema(file, call); subPack != nil {
							file, node, pack = subPack.FindPropNode(subNode.FilePath, path, subNode.SchemaNode)
							return false
						}
					}
				}
			case *ast.BasicLit:
				// if schema it self is a function call , then dive into it
				if key.Value != tmpKey {
					// stop recursive if key not match
					return false
				}
				if len(path) == 1 {
					node = kv.Value
					return false // founded, stop search
				}
				// schema was defined in a separate function call
				switch call := kv.Value.(type) {
				case *ast.Ident: // variable directly
					// local variable define
					if subNode, subPack := p.GetPartialSchema(file, call); subPack != nil {
						file, node, pack = subPack.FindPropNode(subNode.FilePath, path[1:], subNode.SchemaNode)
						return false
					}
				case *ast.CallExpr:
					// current package or from import
					if subNode, subPack := p.GetPartialSchema(filepath, call); subPack != nil {
						file, node, pack = subPack.FindPropNode(subNode.FilePath, path[1:], subNode.SchemaNode)
						return false
					}
					//node, pack = p.GetPartialSchema(call)
					//if pack != nil {
					//	node, pack = pack.FindPropNode(path[1:], node) // recursive find in node
					//}
					if node != nil {
						return false
					}
				}
				// recursive, ele or ele.Schema can be a function call
				file, node, pack = p.FindPropNode(filepath, path[1:], kv.Value)
				if node != nil {
					return false
				}
			}
		}
		return true
	})
	if node == nil {
		ast.Inspect(root, func(sub ast.Node) bool {
			switch fn := sub.(type) {
			case *ast.FuncDecl:
				if IdentName(fn.Name) == "Arguments" && fn.Recv != nil {
					// if Arguments function call other functions
					ast.Inspect(fn.Body, func(sub2 ast.Node) bool {
						switch dcl := sub2.(type) {
						case *ast.CallExpr:
							//switch fn := dcl.Fun.(type) {
							//case *ast.Ident:
							//	if sch, ok := p.PartSchemaFunc[fn.Name]; ok {
							//		p.GetFileAST(sch.FilePath)
							//		node, pack = p.FindPropNode(path, sch.SchemaNode)
							//		return false
							//	}
							//case *ast.SelectorExpr:
							//	// it can r.base.XXX(). if so use local package function XXX
							//	x, y := SelectorXY(fn)
							//	if packSub := p.PackageByName(x); packSub != nil {
							//		if partial, ok := packSub.PartSchemaFunc[y]; ok {
							//			packSub.GetFileAST(partial.FilePath)
							//			node, pack = packSub.FindPropNode(path, partial.SchemaNode)
							//			return false
							//		}
							//	}
							//	// try with local function call for recv.base.arguments calls
							//	if sch, ok := p.PartSchemaFunc[y]; ok {
							//		p.GetFileAST(sch.FilePath)
							//		node, pack = p.FindPropNode(path, sch.SchemaNode)
							//	}
							//}
							if subNode, subPack := p.GetPartialSchema(filepath, dcl); subPack != nil {
								file, node, pack = subPack.FindPropNode(subNode.FilePath, path, subNode.SchemaNode)
								if node != nil {
									return false
								}
							}
							//subNode, subPack := p.GetPartialSchema(dcl)
							//if subPack != nil {
							//	if subNode == nil {
							//		log.Printf("got nil node for part: %v, subPack:  %v, path: %v", dcl, subPack, path)
							//	}
							//	node, pack = subPack.FindPropNode(path, subNode)
							//}
							if node != nil {
								return false
							}
						}
						return true
					})
					return false
				}
			}
			return true
		})
	}
	if pack == nil {
		pack = p
	}
	return
}

// ValidateInSliceFunc caller has ensured the property of path validate is StringInSlice, so skip check it here
func (p *Package) ValidateInSliceFunc(file string, path []string, node ast.Node) (values []string) {
	// find schema definition
	// find the Validate Func to get the possible value
	ast.Inspect(node, func(node ast.Node) bool {
		if len(values) > 0 {
			return false
		}
		switch val := node.(type) {
		case *ast.KeyValueExpr:
			if lit, ok := val.Key.(*ast.Ident); ok {
				if lit.Name == "ValidateFunc" || lit.Name == "ValidateDiagFunc" {
					//p.Printf("ast for %v: %s", path, lit.Name)
					// get first argument
					// can be function call, or slice of const or slice of liter string
					if call, ok := val.Value.(*ast.CallExpr); ok {
						// caller has checked this validate function is StringInSlice
						// argo can be literal []string, or function call
						// it can be a ValidateFunc: aaa.bbb // direct call a function
						if len(call.Args) == 0 {
							values = p.ValuesForCall(file, call)
							return false
						}

						// function like strings.StringInSlice(arg0, arg1)
						arg0 := call.Args[0]
						//caseSense := call.Args[1].(*ast.Ident).Name == "false"
						//_ = caseSense
						switch arg := arg0.(type) {
						case *ast.CompositeLit:
							values = p.ValuesForComposite(file, arg)
						case *ast.FuncLit:
							values = p.ValuesForFuncBody(file, arg.Body)
						case *ast.CallExpr:
							// both local function or import package
							values = p.ValuesForCall(file, arg)
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

// ValuesForFuncBody include some values like
// { return []string{"abc", string(constVar), ...)}
//
//	or for _, val := range xxx/xx.xxx {res -= append(res, stringaval)} return res
//	or for _, val := range []string{xxx} {res -= append(res, stringaval)} return res
func (p *Package) ValuesForFuncBody(file string, body *ast.BlockStmt) (values []string) {
	var inReturn bool
	ast.Inspect(body, func(node ast.Node) bool {
		switch st := node.(type) {
		case *ast.ReturnStmt:
			inReturn = true
		case *ast.RangeStmt:
			values = append(values, p.ValuesForRange(file, st)...)
		case *ast.CompositeLit:
			values = append(values, p.ValuesForComposite(file, st)...)
			if inReturn {
				return false
			}
		}
		return true
	})
	return
}

// ValuesForRange if range over []string/ or function call returns []string in function, work as literal
// if not range over []string, return nil
func (p *Package) ValuesForRange(file string, rng *ast.RangeStmt) (values []string) {
	switch x := rng.X.(type) {
	case *ast.CompositeLit:
		return p.ValuesForComposite(file, rng.X)
	case *ast.CallExpr:
		return p.ValuesForCall(file, x)
	}
	return
}

// ValuesForComposite return values of []string{xxx}, if not composite of array of string, return nil
func (p *Package) ValuesForComposite(file string, comp ast.Expr) (values []string) {
	if st, ok := comp.(*ast.CompositeLit); ok {
		for _, el := range st.Elts {
			values = append(values, p.EvaluateValue(file, el))
		}
	}
	return
}

// ValuesForCall return values returned from calling of a given function name
// call.Fun can be localFuncName or pack.FuncName
// or funclit
func (p *Package) ValuesForCall(file string, call *ast.CallExpr) (values []string) {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		if val, ok := p.PossibleFuncValues[fn.Name]; ok {
			return val
		} else {
			p.Printf("no possible function for: %s", fn.Name)
			return
		}
	case *ast.SelectorExpr:
		x, y := SelectorXY(fn) // import from package
		return p.ValuesFromPackageFunc(file, x, y)
	case *ast.FuncLit:
		return p.ValuesForFuncBody(file, fn.Body)
	}
	p.Printf("no support for %T: %v", call, call)
	return
}

func (p *Package) ValueFromPackageConst(file, packName, fn string) string {
	if pack := p.PackageByName(file, packName); pack != nil {
		return pack.ConstValues[fn]
	}
	return ""
}

func (p *Package) ValuesFromPackageFunc(file, packName, fn string) (values []string) {
	if pack := p.PackageByName(file, packName); pack != nil {
		return pack.PossibleFuncValues[fn]
	}
	return
}

func (p *Package) PackageByName(file, packName string) *Package {
	if packName == "" {
		return nil // do not import current pack
	}
	ast := p.GetFileAST(file)
	importPath := ImportPathByName(packName, ast)
	if importPath == "" {
		log.Printf("no import path for %v in %s", packName, file)
		return nil
	}
	if _, ok := p.Import[importPath]; !ok {
		pack := NewPackage(importPath)
		p.Import[importPath] = pack
	}
	return p.Import[importPath]
}

// EvaluateValue used in simple single value of []string
// ele can be ele + "-" + "ss"
func (p *Package) EvaluateValue(file string, ele ast.Expr) (res string) {
	switch val := ele.(type) {
	case *ast.BasicLit:
		// "abc"
		return litStr(val)
	case *ast.Ident:
		varName := val.Name
		return p.ConstValues[varName]
	case *ast.SelectorExpr:
		x, y := SelectorXY(val)
		return p.ValueFromPackageConst(file, x, y)
	case *ast.BinaryExpr:
		if val.Op == token.ADD {
			res = p.EvaluateValue(file, val.X) + p.EvaluateValue(file, val.Y)
		}
	case *ast.CallExpr:
		// string(xxx)
		// or xxx.xxx(string(xxx))
		switch fun := val.Fun.(type) {
		case *ast.SelectorExpr:
			x, y := SelectorXY(fun)
			if x == "azure" && y == "TitleCase" {
				res = azure.TitleCase(p.EvaluateValue(file, fun.X))
			}
		case *ast.Ident:
			if fun.Name == "string" {
				switch arg0 := val.Args[0].(type) {
				case *ast.SelectorExpr:
					packName, sel := SelectorXY(arg0)
					return p.ValueFromPackageConst(file, packName, sel)
				case *ast.Ident:
					// const value
					if pc, ok := p.ConstValues[arg0.Name]; ok {
						return pc
					} else {
						p.Printf("no value for const: %s", arg0.Name)
					}
				default:
					p.Printf("not support arg0 for function call: %T, arg0: %v", val.Fun, arg0)
				}
			} else {
				p.Printf("not support current function call: %T: %v", val.Fun, val.Args)
			}
		}
	default:
		p.Printf("not support current value: %T: %v of", ele, ele)
	}
	return
}
