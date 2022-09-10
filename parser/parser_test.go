package parser

import (
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
		t.Fatalf("let statement nums wrong ,expected 3,got %v \n",len(pg.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"five"},
		{"x"},
		{"foobar"},
	}

	for i ,tt := range tests{
		
		stmt := pg.Statements[i]
		lstmt,_ := stmt.(*ast.LetStatement)
		
		if lstmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("test %v wrong , expected %v ,got %v ",i,tt.expectedIdentifier,stmt.TokenLiteral())
		}
	}

}