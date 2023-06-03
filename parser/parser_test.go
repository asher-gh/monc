package parser

import (
	"interpreter_in_go/ast"
	"interpreter_in_go/lexer"
	"testing"
)

// TODO: Add the remaining parser tests
func TestLetStatements(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 838383;
   `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if ln := len(program.Statements); ln != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", ln)
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {

	if tl := s.TokenLiteral(); tl != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", tl)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if identName := letStmt.Name.Value; identName != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, identName)
		return false
	}

	if tokenLiteral := letStmt.Name.TokenLiteral(); tokenLiteral != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, tokenLiteral)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
   return 5;
   return 10;
   return 993322;
   `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if rtl := returnStmt.TokenLiteral(); rtl != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got '%q'", rtl)
		}
	}

}

func TestIdentifierExpression(t *Testing) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
}
