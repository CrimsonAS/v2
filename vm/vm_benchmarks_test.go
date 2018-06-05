package vm

import (
	"github.com/CrimsonAS/v2/parser"
	"testing"
)

func BenchmarkFibIterative(b *testing.B) {
	f := "function fibonacci(n) {\n"
	f += "	var a = 0, b = 1, f = 1;\n"
	f += "	for(var i = 2; i <= n; i++) {\n"
	f += "		f = a + b;\n"
	f += "		a = b;\n"
	f += "		b = f;\n"
	f += "	}\n"
	f += "	return f;\n"
	f += "};\n"
	f += "fibonacci(10)\n"

	ast := parser.Parse(f)

	for i := 0; i < b.N; i++ {
		vm := NewVM(ast)
		vm.Run()
	}
}

func BenchmarkFibRecursive(b *testing.B) {
	f := "function fibonacci(n) {\n"
	f += "    if (n < 1) {\n"
	f += "        return 0\n"
	f += "    } else if (n <= 2) {\n"
	f += "        return 1\n"
	f += "    } else {\n"
	f += "        return fibonacci(n - 1) + fibonacci(n - 2)\n"
	f += "    }\n"
	f += "}\n"
	f += "fibonacci(10)\n"

	ast := parser.Parse(f)

	for i := 0; i < b.N; i++ {
		vm := NewVM(ast)
		vm.Run()
	}
}
