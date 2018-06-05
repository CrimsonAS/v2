package main

import (
	"fmt"
	"github.com/CrimsonAS/v2/parser"
	"github.com/CrimsonAS/v2/vm"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
