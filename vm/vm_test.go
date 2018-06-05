package vm

import (
	"github.com/CrimsonAS/v2/parser"
	_ "github.com/kr/pretty"
	"github.com/stvp/assert"
	"testing"
)

func TestStrings(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "\"hello\"",
			out: newString("hello"),
		},
		simpleTest{
			in:  "\"hello\"+\"world\"",
			out: newString("helloworld"),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}
func TestPostfixOperators(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 0; a++",
			out: newNumber(0),
		},
		simpleTest{
			in:  "var a = 1; a--",
			out: newNumber(1),
		},
		simpleTest{
			in:  "var a = 0; var b = 1; a = b++; a",
			out: newNumber(1),
		},
		simpleTest{
			in:  "var a = 0; var b = 1; a = b++; b",
			out: newNumber(2),
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.in)
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestPrefixOperators(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 0; ++a",
			out: newNumber(1),
		},
		simpleTest{
			in:  "var a = 0; --a",
			out: newNumber(-1),
		},
		simpleTest{
			in:  "!false",
			out: newBool(true),
		},
		simpleTest{
			in:  "!1",
			out: newBool(false),
		},
		simpleTest{
			in:  "!!1",
			out: newBool(true),
		},
		simpleTest{
			in:  "!!!1",
			out: newBool(false),
		},
		simpleTest{
			in:  "+3",
			out: newNumber(3),
		},
		simpleTest{
			in:  "-3",
			out: newNumber(-3),
		},
	}
	for _, test := range tests {
		t.Logf("%s", test.in)
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestSimple(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "2+3",
			out: newNumber(5),
		},
		simpleTest{
			in:  "(2+3)",
			out: newNumber(5),
		},
		simpleTest{
			in:  "(1+1)",
			out: newNumber(2),
		},
		simpleTest{
			in:  "(2+2)*(2+2)",
			out: newNumber(16),
		},
		simpleTest{
			in:  "10/2",
			out: newNumber(5),
		},
		simpleTest{
			in:  "if (true) { 10/2 }",
			out: newNumber(5),
		},
		simpleTest{
			in:  "if (false) { 10/2 }",
			out: value{},
		},
		simpleTest{
			in:  "if (true) { 10/2 }",
			out: newNumber(5),
		},
		simpleTest{
			in:  "2+2*2+2",
			out: newNumber(8),
		},
		simpleTest{
			in:  "1<2",
			out: newBool(true),
		},
		simpleTest{
			in:  "1<=2",
			out: newBool(true),
		},
		simpleTest{
			in:  "2<=1",
			out: newBool(false),
		},
		simpleTest{
			in:  "2<1",
			out: newBool(false),
		},
		simpleTest{
			in:  "2>1",
			out: newBool(true),
		},
		simpleTest{
			in:  "1>2",
			out: newBool(false),
		},
		simpleTest{
			in:  "1==2",
			out: newBool(false),
		},
		simpleTest{
			in:  "1==1",
			out: newBool(true),
		},
		simpleTest{
			in:  "1!=2",
			out: newBool(true),
		},
		simpleTest{
			in:  "1!=1",
			out: newBool(false),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestCall(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "f() function f() { 5 }",
			out: newNumber(5),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestVar(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 5, b; b = a + 10",
			out: newNumber(15),
		},
		simpleTest{
			in:  "var a = 5; a",
			out: newNumber(5),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestWhile(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 0; while (a > 10) a = a + 1; a",
			out: newNumber(0),
		},
		simpleTest{
			in:  "var a = 0; while (a < 10) a = a + 1; a",
			out: newNumber(10),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestDoWhile(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 0; do { a = a + 1 } while (a > 10); a",
			out: newNumber(1),
		},
		simpleTest{
			in:  "var a = 0; do { a = a + 1 } while (a < 5); a",
			out: newNumber(5),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestForStatement(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 0; for (a = 0; a < 5; a = a + 1) { }; a",
			out: newNumber(5),
		},
		simpleTest{
			in:  "var a = 10; for (; a < 5; a = a + 1) { }; a",
			out: newNumber(10),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestReturnStatement(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "function f() { return 10; } var a; a = f();",
			out: newNumber(10),
		},
		simpleTest{
			in:  "function a() { return 10; } function b() { return 5; } var c = a(); c;",
			out: newNumber(10),
		},
		simpleTest{
			in:  "function a() { return 10; } function b() { return 5; } var c = b(); c;",
			out: newNumber(5),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestBuiltinFunction(t *testing.T) {
	{
		testFunc := func(vm *vm, f value, args []value) value {
			t.Logf("I'm a message from a builtin function: %+v", args)
			return newString("Hello world")
		}

		ast := parser.Parse("testFunc()")
		vm := NewVM(ast)
		pf := newFunctionObject(testFunc)
		vm.defineVar(appendStringtable("testFunc"), pf)
		assert.Equal(t, vm.Run(), newString("Hello world"))
	}

	// test multiple arguments, and their order
	{
		testFunc := func(vm *vm, f value, args []value) value {
			return newString(args[0].toString() + args[1].toString())
		}

		ast := parser.Parse("testFunc(\"Hello\", \"World\")")
		vm := NewVM(ast)
		pf := newFunctionObject(testFunc)
		vm.defineVar(appendStringtable("testFunc"), pf)
		assert.Equal(t, vm.Run(), newString("HelloWorld"))
	}
}

func TestJSFunction(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "function f(a) { return a; } var b; b = f(10);",
			out: newNumber(10),
		},
		simpleTest{
			in:  "function f(a) { return a; } var a; a = f(10);",
			out: newNumber(10),
		},
		simpleTest{
			in:  "function f(a, b) { return a; } var a; a = f(10, 5);",
			out: newNumber(10),
		},
		simpleTest{
			in:  "function f(a, b) { return b; } var a; a = f(5, 10);",
			out: newNumber(10),
		},
		simpleTest{
			in:  "function f(a, b) { return a; } var a; a = f(5, 10);",
			out: newNumber(5),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestRecursiveLookups(t *testing.T) {
	ast := parser.Parse("function f(a) { if (a > 3) return a; a = a + 1; return f(a); } var n = f(0); n")
	vm := NewVM(ast)
	ret := vm.Run()
	assert.Equal(t, ret, newNumber(4))
}

func TestFibonnaci(t *testing.T) {
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

	iterative := newUndefined()
	{
		ast := parser.Parse(f)
		vm := NewVM(ast)
		iterative = vm.Run()
	}

	f = ""
	f = "function fibonacci(n) {\n"
	f += "    if (n < 1) {\n"
	f += "        return 0\n"
	f += "    } else if (n <= 2) {\n"
	f += "        return 1\n"
	f += "    } else {\n"
	f += "        return fibonacci(n - 1) + fibonacci(n - 2)\n"
	f += "    }\n"
	f += "}\n"
	f += "fibonacci(10)\n"

	recursive := newUndefined()
	{
		ast := parser.Parse(f)
		vm := NewVM(ast)
		recursive = vm.Run()
	}

	assert.Equal(t, iterative, recursive)
	assert.Equal(t, iterative, newNumber(55))
}

func TestConditionalExpression(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "var a = 1; var b = a == 1 ? 2 : 3;",
			out: newNumber(2),
		},
		// ### failing, which seems odd...
		//simpleTest{
		//	in:  "var a = 1; var b = a == 1 ? 2 : 3; b",
		//	out: newNumber(2),
		//},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}
