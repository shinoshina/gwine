package repl

import (
	"bufio"
	"fmt"
	"gwine/compiler"
	"gwine/evaluator"
	"gwine/lexer"
	"gwine/object"
	"gwine/parser"
	"gwine/vm"
	"io"
)

func StartForInterpreter(in io.Reader,out io.Writer){

	sc := bufio.NewScanner(in)
	env  := object.NewEnvironment()
	
	for {

		fmt.Fprintf(out,">> ")
		ok := sc.Scan()
		if !ok {
			return 
		}

		l := lexer.New(sc.Text())
		p := parser.New(l)
		program := p.ParseProgram()
		//fmt.Fprintln(out,program.String())
	
		evaluated := evaluator.Eval(program,env)
		if evaluated != nil {
			io.WriteString(out,evaluated.Inspect())
            io.WriteString(out,"\n")
		}
	}
}
func StartForVm(in io.Reader,out io.Writer){

	sc := bufio.NewScanner(in)
	
	
	for {

		fmt.Fprintf(out,">> ")
		ok := sc.Scan()
		if !ok {
			return 
		}

		l := lexer.New(sc.Text())
		p := parser.New(l)
		program := p.ParseProgram()
		//fmt.Fprintln(out,program.String())
	
		compiler := compiler.New()
		err := compiler.Compile(program)
		if err != nil{
			fmt.Fprintf(out,"compile fail\n %s \n",err)
			continue
		}

		vm := vm.New(compiler.ByteCode())
		err = vm.Run()
		if err != nil{
			fmt.Println(err)
		}


		// st := vm.Top()
		// io.WriteString(out,st.Inspect() + "\n")
		io.WriteString(out,vm.LastPoped().Inspect()+ "\n")

	}
}