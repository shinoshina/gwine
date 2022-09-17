package repl

import (
	"bufio"
	"fmt"
	"gwine/evaluator"
	"gwine/lexer"
	"gwine/object"
	"gwine/parser"
	"io"
)

func Start(in io.Reader,out io.Writer){

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