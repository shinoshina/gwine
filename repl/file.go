package repl

import (
	"fmt"
	"gwine/compiler"
	"gwine/lexer"
	"gwine/object"
	"gwine/parser"
	"gwine/vm"
	"io"
	"io/ioutil"
	"os"
)

func FromFile(file string) {

	f , err := ioutil.ReadFile(file)
	if err != nil{
		fmt.Println(err)
	}
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symboltbl := compiler.NewSymbolTable()

	l := lexer.New(string(f))
	p := parser.New(l)
	program := p.ParseProgram()
	//fmt.Fprintln(out,program.String())

	comp := compiler.NewWithState(symboltbl, constants)
	err = comp.Compile(program)
	if err != nil {
		fmt.Fprintf(os.Stdout, "compile fail")
	}
	code := comp.ByteCode()
	constants = code.Constants
	vmm := vm.NewWithGlobalStore(code, globals)
	err = vmm.Run()
	if err != nil {
		fmt.Println(err)
	}
	// st := vm.Top()
	// io.WriteString(out,st.Inspect() + "\n")
	io.WriteString(os.Stdout, vmm.LastPoped().Inspect()+"\n")
}
