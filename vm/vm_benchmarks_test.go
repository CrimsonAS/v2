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

package vm

import (
	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
	"testing"
)

type benchmark struct {
	name string
	code string
}

func vmBenchmarks() []benchmark {
	return []benchmark{
		benchmark{
			name: "fib_it",
			code: `function fibonacci(n) {
					var a = 0, b = 1, f = 1;
					for(var i = 2; i <= n; i++) {
						f = a + b;
						a = b;
						b = f;
					}
					return f;
				};
				fibonacci(10)`,
		},
		benchmark{
			name: "fib_rec",
			code: `function fibonacci(n) {
					if (n < 1) {
						return 0
					} else if (n <= 2) {
						return 1
					} else {
						return fibonacci(n - 1) + fibonacci(n - 2)
					}
				}
				fibonacci(10)`,
		},
		benchmark{
			// The breakage (x never changing) in this benchmark is known and deliberate. Don't fix it.
			name: "sum10k",
			code: `function sum(n) {
					var x = 0;
					for (var i=0; i<n; ++i)
						x = x + x;
					return x;
				}
				sum(10000)`,
		},
		benchmark{
			// The above, but corrected.
			name: "sum10kf",
			code: `function sum(n) {
					var x = 0;
					for (var i=0; i<n; ++i)
						x = x + 1;
					return x;
				}
				sum(10000)`,
		},
	}
}

func BenchmarkVM(b *testing.B) {
	tests := vmBenchmarks()
	for _, test := range tests {
		b.Run(test.name,
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vm := New(test.code)
					vm.Run()
				}
			})
	}
}

func BenchmarkVM_Otto(b *testing.B) {
	tests := vmBenchmarks()
	for _, test := range tests {
		b.Run(test.name,
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vm := otto.New()
					vm.Run(test.code)
				}
			})
	}
}

func BenchmarkVM_Goja(b *testing.B) {
	tests := vmBenchmarks()
	for _, test := range tests {
		b.Run(test.name,
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vm := goja.New()
					vm.RunString(test.code)
				}
			})
	}
}
