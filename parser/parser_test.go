package parser

import (
	"fmt"
	"gwine/ast"
	"gwine/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
	let five = 5  ;
	let x = 6;
	let foobar = 123500;
	`

	l := lexer.New(input)
	p := New(l)

	pg := p.ParseProgram()
	if pg == nil {
		t.Fatalf("ParseProgram return nil")
	}
	if len(pg.Statements) != 3 {
		t.Fatalf("let statement nums wrong ,expected 3,got %v \n", len(pg.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"five"},
		{"x"},
		{"foobar"},
	}

	for i, tt := range tests {

		stmt := pg.Statements[i]
		lstmt, _ := stmt.(*ast.LetStatement)

		if lstmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("test %v wrong , expected %v ,got %v ", i, tt.expectedIdentifier, stmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {

	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wtf??")
	}
	fmt.Printf("expression statement token %+v\n", stmt.Token)

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("wtf2??")
	}
	fmt.Printf("identifier token %+v\nvalue %v\n", ident.Token, ident.Value)

}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wtf??")
	}
	fmt.Printf("expression statement token %+v\n", stmt.Token)
	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("wtf2??")
	}
	fmt.Printf("identifier token %+v\nvalue %v\n", ident.Token, ident.Value)
}
func TestPrefixExpression(t *testing.T) {
	input := "-115;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wtf??")
	}
	fmt.Printf("expression statement token %+v\n", stmt.Token)
	ident, ok := stmt.Expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("wtf2??")
	}
	fmt.Printf("identifier token %+v\noperator %v\n", ident.Token, ident.Operator)
	num, ok := ident.Right.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("wtf??3")
	}
	fmt.Printf("identifier token %+v\nvalue %v\n", num.Token, num.Value)
}

func TestNodeLiteral(t *testing.T) {

	input := `
	5+5;
	5+1+2;
	5*1+2;
	5+1*2;
	5/2*1;
	5/2+1-1;
	!2*2+1/2;
	-1*2+1-2;
	-1+2+3;
	-1*(5+5);
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("nmsl")
		}
		fmt.Println(stmt.String())
	}

}
