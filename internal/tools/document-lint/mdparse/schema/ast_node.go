package schema

import "go/ast"

type Node struct {
	node ast.Node
}

func NewNode(node ast.Node) *Node {
	return &Node{
		node: node,
	}
}

func (n *Node) UnaryX() *Node {
	if n == nil {
		return nil
	}
	if v, ok := n.node.(*ast.UnaryExpr); ok {
		return NewNode(v.X)
	}
	return nil
}

func (n *Node) FuncBody() *Node {
	if n != nil {
		if v, ok := n.node.(*ast.FuncLit); ok {
			return NewNode(v.Body)
		}
	}
	return nil
}

func (n *Node) CallExpreFun() *Node {
	if n != nil {
		if v, ok := n.node.(*ast.CallExpr); ok {
			return NewNode(v.Fun)
		}
	}
	return nil
}

func (n *Node) CallExpreArgs0() *Node {
	if n != nil {
		if v, ok := n.node.(*ast.CallExpr); ok && len(v.Args) > 0 {
			return NewNode(v.Args[0])
		}
	}
	return nil
}

func (n *Node) Ident() string {
	if n != nil {
		if v, ok := n.node.(*ast.Ident); ok {
			return v.Name
		}
	}
	return ""
}

func (n *Node) CompositeType() *Node {
	if n != nil {
		if v, ok := n.node.(*ast.CompositeLit); ok {
			return NewNode(v.Type)
		}
	}
	return nil
}

func (n *Node) CompositeN(idx int) *Node {
	if n != nil {
		if v, ok := n.node.(*ast.CompositeLit); ok {
			if idx < len(v.Elts) {
				return NewNode(v.Elts[idx])
			}
		}
	}
	return nil
}

func (n *Node) Selector() *Node {
	if n != nil {
		if v, ok := n.node.(*ast.SelectorExpr); ok {
			return NewNode(v)
		}
	}
	return nil
}

func (n *Node) X() *Node {
	x := n.Selector()
	if x != nil {
		x2 := x.node.(*ast.SelectorExpr).X
		if v, ok := x2.(*ast.SelectorExpr); ok {
			return NewNode(v.X)
		}
		return NewNode(x2)
	}
	return nil
}

func (n *Node) Sel() *ast.Ident {
	if v := n.Selector(); v != nil {
		return v.node.(*ast.SelectorExpr).Sel
	}
	return nil
}
