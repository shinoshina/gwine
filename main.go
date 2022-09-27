package main

import (
	"gwine/repl"
	"os"
)

func main() {
	// file, err := os.Create("test.gwine")
	// if err != nil {
	// 	fmt.Errorf("kksk")
	// }
	repl.StartForVm(os.Stdin, os.Stdout)
	// repl.FromFile("test.gwine")
}
