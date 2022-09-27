package vm

import (
	"fmt"
	"gwine/compiler"
	"gwine/lexer"
	"gwine/object"
	"gwine/parser"
	"testing"
)

func TestArrayinHash(t *testing.T) {

	constants := []object.Object{}
	globals := make([]object.Object, GlobalsSize)
	symboltbl := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symboltbl.DefineBuiltin(i, v.Name)
	}

	l := lexer.New(`len([1,2,3])`)
	p := parser.New(l)
	program := p.ParseProgram()
	//fmt.Fprintln(out,program.String())

	comp := compiler.NewWithState(symboltbl, constants)
	err := comp.Compile(program)
	if err != nil {
		fmt.Println(err)
	}
	code := comp.ByteCode()
	vmm := NewWithGlobalStore(code, globals)
	err = vmm.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vmm.LastPoped().Inspect()+"\n")

}
