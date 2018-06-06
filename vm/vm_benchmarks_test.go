package vm

import (
	"github.com/CrimsonAS/v2/parser"
	"github.com/robertkrimen/otto"
	"testing"
)

func BenchmarkVM(b *testing.B) {
	type benchmark struct {
		name string
		code string
	}

	tests := []benchmark{
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

	for _, test := range tests {
		ast := parser.Parse(test.code)
		b.Run(test.name,
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vm := NewVM(ast)
					vm.Run()
				}
			})
		b.Run(test.name+"_o",
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					vm := otto.New()
					vm.Run(test.code)
				}
			})
	}
}
