package parser

import (
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if stlen := len(program.Statements); stlen != 1 {
		t.Fatalf("Program has not enough statements. got=%d", stlen)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if ln := len(program.Statements); ln != 1 {
		t.Fatalf("program has incorrect number of statements. got=%d", ln)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if tl := literal.TokenLiteral(); tl != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", tl)
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("program does not contain %d statements. got=%d", 1, len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prog.Statements[0] is not *ast.ExpressionStatement. got=%T", prog.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not of type *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.intValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {

	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.IntegerLiteral. got=%T", il)
	}

	if il.Value != value {
		t.Errorf("il.Value not %d. got=%d", value, il.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("il.TokenLiteral is not %d. got=%s", value, il.TokenLiteral())
		return false
	}

	return true
}
