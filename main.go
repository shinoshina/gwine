package main

import (
	"gwine/repl"
	"os"
)

func main() {
	repl.StartForVm(os.Stdin,os.Stdout)
}