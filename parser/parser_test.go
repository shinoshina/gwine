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
			t.Fatalf("ffffail")
		}
		fmt.Println(stmt.String())
	}

}

func TestTrueFalse(t *testing.T) {

	input := `
	true;
	true;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("failffffffail")
		}
		fmt.Println(stmt.String())
	}

}

func TestFunctionLiteral(t *testing.T) {
	input := `
	fn(x,y){x+y;};
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("failffffffail")
		}
		fmt.Println(stmt.String())
	}
}
func TestCall(t *testing.T) {
	input := `
	add(x,y);
	add(x,y+1,-1+add(x,y,z));
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("failffffffail")
		}
		fmt.Println(stmt.String())
	}
}
func TestLetAndReturn(t *testing.T) {
	input := `
	let a = add(x,y);
	let a = add(x,y+1,-1+add(x,y,z));
	let a = 1;
	let a = 1
	let b = fn(x,y){
		let a = b;
	};
	let c = if(x > y) {a;}else{b;};
	let c = if(x > y) {return a;} else{b;};
	`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.LetStatement)
		if !ok {
			t.Fatalf("failffffffail")
		}
		fmt.Println(stmt.String())
		fmt.Println(stmt)
	}
}

func TestEQNEQ(t *testing.T) {
	input := `
	a!=b
	a==b
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("failffffffail")
		}
		fmt.Println(stmt.String())
	}

}
func TestIFELSE(t *testing.T) {
	input := `
	if(x>y){return x;}else{y;}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for i, _ := range program.Statements {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("failffffffail")
		}
		fmt.Println(stmt.String())
	}

}

func TestStructDeclarion(t *testing.T) {
	input := `
	struct ff{
		fn a(){
			return 1;
		}
		fn b(){
			return 2;
		}
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	// for i, _ := range program.Statements {
	// 	stmt, ok := program.Statements[i].(*ast.StructDeclarionStatement)
	// 	if !ok {
	// 		t.Fatalf("failffffffail")
	// 	}
	// 	fmt.Println(stmt.String())
	// }
	for _, stmt := range program.Statements {
		fmt.Println(stmt.String())
	}
}

func TestFunctionDeclarion(t *testing.T) {
	input := `
	fn ff(a,b){
		let k = 1;
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for _, stmt := range program.Statements {
		fmt.Println(stmt.String())
	}

}
