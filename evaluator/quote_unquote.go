package evaluator

import (
	"interpreter_in_go/ast"
	"interpreter_in_go/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
