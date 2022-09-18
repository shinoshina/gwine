package main

import (
	"gwine/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin,os.Stdout)
}