package main

import (
	"fmt"
	"github.com/CrimsonAS/v2/parser"
	"github.com/CrimsonAS/v2/vm"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	f := ""
	if len(os.Args) > 1 {
		f = os.Args[1]
	} else {
		fmt.Printf("Need a file to run\n")
		os.Exit(0)
	}
	code, _ := ioutil.ReadFile(f)
	ast := parser.Parse(string(code))
	vm := vm.NewVM(ast)
	ret := vm.Run()
	log.Printf("Code returned %s", ret)
}
