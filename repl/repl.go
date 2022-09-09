package repl

import (
	"bufio"
	"fmt"
	"gwine/lexer"
	"gwine/token"
	"io"
)

func Start(in io.Reader,out io.Writer){

	sc := bufio.NewScanner(in)
	
	for {

		fmt.Fprintf(out,">> ")
		ok := sc.Scan()
		if !ok {
			return 
		}

		l := lexer.New(sc.Text())

		for t := l.NextToken() ; t.Type != token.EOF ; t = l.NextToken(){
			fmt.Fprintf(out,"%+v\n",t)
		}


	}
}