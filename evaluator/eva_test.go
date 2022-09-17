package evaluator

import (
	"fmt"
	"gwine/lexer"
	"gwine/object"
	"gwine/parser"
	"testing"
)

func TestEva(t *testing.T) {

	// why cant -5 pass ,cause -5 seen as a prefix expression type,not a single integer literal type,
	// which is just a part of prefix expression : "-"
	input := `
	-5;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	env  := object.NewEnvironment()
	obj := Eval(program,env)


	obj2, _ := obj.(*object.Integer)

	fmt.Println(obj2.Value)

}
func TestBoolean(t *testing.T) {
	input := `
	-5;
	true;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	env  := object.NewEnvironment()
	obj := Eval(program,env)


	obj2, _ := obj.(*object.Boolean)

	fmt.Println(obj2.Value)

}
func TestPrefixExpressionBoolean(t *testing.T) {
	input := `
	!0;
	!!true`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	env  := object.NewEnvironment()
	obj := Eval(program,env)


	obj2, _ := obj.(*object.Boolean)

	fmt.Println(obj2.Value)

}
func TestPrefixExpressionMinus(t *testing.T) {
	input := `
	-5;
	-true`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
    env  := object.NewEnvironment()
	obj := Eval(program,env)

	obj2, _ := obj.(*object.Integer)

	fmt.Println(obj2.Value)

}
