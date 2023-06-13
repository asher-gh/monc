package evaluator

import (
	"interpreter_in_go/lexer"
	"interpreter_in_go/object"
	"interpreter_in_go/parser"
	"testing"
)

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

}

func testIntegerObject(t *testing.T, evaluated object.Object, i int64) bool {

	obj, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if obj.Value != i {

		t.Errorf("object has wrong value. got=%d, want=%d", obj.Value, i)
		return false
	}

	return true
}

func testEval(s string) object.Object {
	l := lexer.New(s)
	p := parser.New(l)
	prog := p.ParseProgram()

	return Eval(prog)
}
