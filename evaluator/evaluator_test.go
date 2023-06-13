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

func testIntegerObject(t *testing.T, obj object.Object, i int64) bool {

	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if res.Value != i {

		t.Errorf("object has wrong value. got=%d, want=%d", res.Value, i)
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

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, b bool) bool {
	res, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if res.Value != b {
		t.Errorf("object has wrong value. got=%t, want=%t", res.Value, b)
		return false
	}

	return true
}
