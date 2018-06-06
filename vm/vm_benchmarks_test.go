package vm

import (
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
