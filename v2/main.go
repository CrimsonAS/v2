/*
 * Copyright 2018 Crimson AS <info@crimson.no>
 * Author: Robin Burchell <robin.burchell@crimson.no>
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/CrimsonAS/v2/vm"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	newCompiler := flag.Bool("new-compiler", false, "use new compiler")
	profile := flag.Bool("profile", false, "enable profiling")
	showBytecode := flag.Bool("show-bytecode", false, "show bytecode after code generation")
	flag.Parse()

	if *profile {
		log.Printf("Enabling profiling")
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	vm.NewCompiler = *newCompiler

	f := ""
	if flag.NArg() > 0 {
		f = flag.Args()[0]
		log.Printf("Running %s", f)
	} else {
		fmt.Printf("Need a file to run\n")
		os.Exit(0)
	}
	code, _ := ioutil.ReadFile(f)
	vm := vm.New(string(code))

	if *showBytecode {
		vm.DumpCode()
	}

	ret := vm.Run()
	log.Printf("Code returned %s", ret)
}
