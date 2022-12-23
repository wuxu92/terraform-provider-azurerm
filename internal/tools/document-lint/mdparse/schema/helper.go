package schema

import (
	"go/ast"
	"go/token"
	"strings"
)

func litStr(node *ast.BasicLit) string {
	return strings.Trim(node.Value, `"`)
}

func BasicLitValue(node ast.Expr) string {
	if n, ok := node.(*ast.BasicLit); ok {
		return litStr(n)
	}
	return ""
}

func IdentName(node ast.Expr) string {
	if n, ok := node.(*ast.Ident); ok {
		return n.Name
	}
	return ""
}

func SelectorXY(node *ast.SelectorExpr) (x, y string) {
	x = IdentName(node.X)
	return x, node.Sel.Name
}

func IsPartSchemaSel(node *ast.SelectorExpr) bool {
	x, y := SelectorXY(node)
	return (y == "Schema" || y == "Resource") && (x == "schema" || x == "pluginsdk")
}

// IsCompositeString check if compositeLit presents []string{}
func IsCompositeString(node *ast.CompositeLit) bool {
	if arr, ok := node.Type.(*ast.ArrayType); ok {
		if ident, ok := arr.Elt.(*ast.Ident); ok {
			return ident.Name == "string"
		}
	}
	return false
}

// IsCompositeMapOfResource check if a composite node if map[string]*pluginsdk.Resource
func IsCompositeMapOfResource(node *ast.CompositeLit) bool {
	if imp, ok := node.Type.(*ast.MapType); ok {
		if k, ok := imp.Key.(*ast.Ident); ok && k.Name == "string" {
			if v, ok := imp.Value.(*ast.UnaryExpr); ok && v.Op == token.MUL {
				if sel, ok := v.X.(*ast.SelectorExpr); ok {
					x, y := SelectorXY(sel)
					if (x == "pluginsdk" || x == "schema") && y == "Schema" {
						// top level schema define
						return true
					}
				}
			}
		}
	}
	return false

}

func IsBasicStringType(node ast.Expr) bool {
	if ident, ok := node.(*ast.Ident); ok {
		if ident.Name == "string" {
			return true
		}
		if obj := ident.Obj; obj != nil {
			if typeSpec, ok := obj.Decl.(*ast.TypeSpec); ok {
				if typeIdent, ok := typeSpec.Type.(*ast.Ident); ok {
					return typeIdent.Name == "string"
				}
			}
		}
	}
	return false
}
