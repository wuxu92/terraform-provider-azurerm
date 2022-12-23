package schema

import (
	"go/ast"
	"go/token"
	"path"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/tools/go/packages"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

// FileForResource for typed sdk resource, the file is terraform-provider-azurerm/internal/sdk/wrapper_resource.go
func FileForResource(funcs ...interface{}) (file string) {
	for _, fn := range funcs {
		if file, _ = util.FuncFileLine(fn); file != "" {
			return file
		}
	}
	return
}

type Resource struct {
	FilePath     string
	ResourceType string // azurerm_xxx

	Package *Package // all information saved in Package

	TokenSet *token.FileSet
	ASTFile  *ast.File

	// one of Schema or SDKResource must use
	Schema      *schema.Resource `json:"-"`
	SDKResource sdk.Resource     `json:"-"`

	PossibleValues map[string][]string // possible values for key(property path)

	// all below fields deprecated

	Imports     map[string]*ast.ImportSpec
	cacheDep    map[string]*packages.Package
	cacheConsts map[string]map[string]string // packagename => [key]value

	cacheFuncValues map[string][]string // key pack.Funcname; value: possible values of funvtion return
}

func ResourceForSDKType(res sdk.Resource) *schema.Resource {
	r := sdk.NewResourceWrapper(res)
	ins, _ := r.Resource()
	return ins
}

// NewResourceByTyped NewResource ...
// r is Schema.Resource or Typed SDK Resource
func NewResourceByTyped(r sdk.Resource) *Resource {
	s := &Resource{}
	s.SDKResource = r
	s.Schema = ResourceForSDKType(r)
	s.ResourceType = r.ResourceType()
	s.Init()
	return s
}

func NewResourceByUntyped(r *schema.Resource, rType string) *Resource {
	s := &Resource{}
	s.Schema = r
	s.ResourceType = rType
	s.Init()
	return s
}

func (r *Resource) Init() {
	if r.SDKResource != nil {
		// SDKResource is a type of interface, have to get the real
		vd := reflect.ValueOf(r.SDKResource).Interface()
		vd = reflect.ValueOf(vd).MethodByName("Arguments")
		// this is not work if Read() defined in other file
		r.FilePath = FileForResource(r.SDKResource.Read().Func)
	} else {
		r.FilePath = FileForResource(r.Schema.Read, r.Schema.ReadContext)
	}
	dir := path.Dir(r.FilePath)
	svcIndx := strings.Index(dir, "terraform-provider-azurerm")
	packName := dir[svcIndx+len("terraform-provider-azurerm"):]

	importPath := "github.com/hashicorp/terraform-provider-azurerm" + packName
	r.Package = manager.GetPackForPath(importPath)
	r.PossibleValues = map[string][]string{}
	// double confirm source file path match
	if filePath, ok := r.Package.ResourceTypeFile[r.ResourceType]; ok {
		r.FilePath = filePath
	}
	//r.FindAllInSliceProp() // find all need check inslice validate values
	r.FindAllInSlicePropByMonkey()
}
