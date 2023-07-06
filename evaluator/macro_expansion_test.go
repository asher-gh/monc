package evaluator

import (
	"interpreter_in_go/ast"
	"interpreter_in_go/lexer"
	"interpreter_in_go/object"
	"interpreter_in_go/parser"
	"testing"
)

func TestDefineMacros(t *testing.T) {
	input := `
   let number = 1;
   let function = fn(x, y) { x + y };
   let myMacro = macro(x, y) { x + y };
   `

	env := object.NewEnvironment()
	prog := testParseProgram(input)

	// takes the parsed program and add the macro defs from that to env
	DefineMacros(prog, env)

	if ln := len(prog.Statements); ln != 2 {
		t.Fatalf("wrong number of statements. got=%d", ln)
	}

	_, ok := env.Get("number")

	if ok {
		t.Fatalf("`number` should not be defined")
	}
	_, ok = env.Get("function")
	if ok {
		t.Fatalf("`function` should not be defined")
	}

	obj, ok := env.Get("myMacro")
	if !ok {
		t.Fatalf("macro not in environment.")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. got=%T (%+v)", obj, obj)
	}

	if ln := len(macro.Parameters); ln != 2 {
		t.Fatalf("wrong number of macro parameters. got=%d", ln)
	}

	if p1 := macro.Parameters[0]; p1.String() != "x" {
		t.Fatalf("parameters is not 'x'. got=%q", p1)
	}
	if p2 := macro.Parameters[1]; p2.String() != "y" {
		t.Fatalf("parameters is not 'x'. got=%q", p2)
	}

	expectedBody := "(x + y)"

	if mb := macro.Body.String(); mb != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, mb)
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
      let infixExpression = macro() { quote(1+2); };
      infixExpression();
      `,
			`(1 + 2)`,
		},
		{
			`
      let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
      reverse(2 + 2, 10 - 5);
      `,
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
		    let unless = macro(condition, consequence, alternative) {
		       quote(if (!(unquote(condition))) {
		          unquote(consequence);
		       } else {
		          unquote(alternative);
		       });
		    };

		    unless(10 > 5, puts("not greater"), puts("greater"));
		   `,
			`if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal. want=%q, got=%q", expected.String(), expanded.String())
		}
	}

}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
