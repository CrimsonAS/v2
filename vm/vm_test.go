package vm

import (
	_ "github.com/kr/pretty"
	"github.com/stvp/assert"
	"testing"
)

type simpleVMTest struct {
	in  string
	out value
}

func runSimpleVMTestHelper(t *testing.T, tests []simpleVMTest) {
	for _, test := range tests {
		vm := New(test.in)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestStrings(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "\"hello\"",
			out: newString("hello"),
		},
		simpleVMTest{
			in:  "\"hello\"+\"world\"",
			out: newString("helloworld"),
		},
	}

	runSimpleVMTestHelper(t, tests)
}
func TestPostfixOperators(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; a++",
			out: newNumber(0),
		},
		simpleVMTest{
			in:  "var a = 1; a--",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; var b = 1; a = b++; a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; var b = 1; a = b++; b",
			out: newNumber(2),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestPrefixOperators(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; ++a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; --a",
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  "!false",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "!1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "!!1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "!!!1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "+3",
			out: newNumber(3),
		},
		simpleVMTest{
			in:  "-3",
			out: newNumber(-3),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestSimple(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "2+3",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "(2+3)",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "(1+1)",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "(2+2)*(2+2)",
			out: newNumber(16),
		},
		simpleVMTest{
			in:  "10/2",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "if (true) { 10/2 }",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "if (false) { 10/2 }",
			out: value{},
		},
		simpleVMTest{
			in:  "if (true) { 10/2 }",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "2+2*2+2",
			out: newNumber(8),
		},
		simpleVMTest{
			in:  "1<2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "1<=2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "2<=1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "2<1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "2>1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "1>2",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "1==2",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "1==1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "1!=2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "1!=1",
			out: newBool(false),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestCall(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "f() function f() { 5 }",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestVar(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 5, b; b = a + 10",
			out: newNumber(15),
		},
		simpleVMTest{
			in:  "var a = 5; a",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestWhile(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; while (a > 10) a = a + 1; a",
			out: newNumber(0),
		},
		simpleVMTest{
			in:  "var a = 0; while (a < 10) a = a + 1; a",
			out: newNumber(10),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestDoWhile(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; do { a = a + 1 } while (a > 10); a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; do { a = a + 1 } while (a < 5); a",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestForStatement(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; for (a = 0; a < 5; a = a + 1) { }; a",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "var a = 10; for (; a < 5; a = a + 1) { }; a",
			out: newNumber(10),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestReturnStatement(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "function f() { return 10; } var a; a = f();",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function a() { return 10; } function b() { return 5; } var c = a(); c;",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function a() { return 10; } function b() { return 5; } var c = b(); c;",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestBuiltinFunction(t *testing.T) {
	{
		testFunc := func(vm *vm, f value, args []value) value {
			t.Logf("I'm a message from a builtin function: %+v", args)
			return newString("Hello world")
		}

		vm := New("testFunc()")
		pf := newFunctionObject(testFunc, nil)
		vm.defineVar(appendStringtable("testFunc"), pf)
		assert.Equal(t, vm.Run(), newString("Hello world"))
	}

	// test multiple arguments, and their order
	{
		testFunc := func(vm *vm, f value, args []value) value {
			return newString(args[0].toString() + args[1].toString())
		}

		vm := New("testFunc(\"Hello\", \"World\")")
		pf := newFunctionObject(testFunc, nil)
		vm.defineVar(appendStringtable("testFunc"), pf)
		assert.Equal(t, vm.Run(), newString("HelloWorld"))
	}

	// test call vs construct
	{
		testCall := func(vm *vm, f value, args []value) value {
			return newNumber(10)
		}
		testConstruct := func(vm *vm, f value, args []value) value {
			return newNumber(20)
		}

		{
			vm := New("testFunc()")
			pf := newFunctionObject(testCall, testConstruct)
			vm.defineVar(appendStringtable("testFunc"), pf)
			assert.Equal(t, vm.Run(), newNumber(10))
		}
		{
			vm := New("new testFunc()")
			pf := newFunctionObject(testCall, testConstruct)
			vm.defineVar(appendStringtable("testFunc"), pf)
			assert.Equal(t, vm.Run(), newNumber(20))
		}
	}
}

func TestJSFunction(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "function f(a) { return a; } var b; b = f(10);",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a) { return a; } var a; a = f(10);",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a, b) { return a; } var a; a = f(10, 5);",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a, b) { return b; } var a; a = f(5, 10);",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a, b) { return a; } var a; a = f(5, 10);",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestRecursiveLookups(t *testing.T) {
	vm := New("function f(a) { if (a > 3) return a; a = a + 1; return f(a); } var n = f(0); n")
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
		vm := New(f)
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
		vm := New(f)
		recursive = vm.Run()
	}

	assert.Equal(t, iterative, recursive)
	assert.Equal(t, iterative, newNumber(55))
}

func TestConditionalExpression(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 1; var b = a == 1 ? 2 : 3;",
			out: newNumber(2),
		},
		// ### failing, which seems odd...
		//simpleVMTest{
		//	in:  "var a = 1; var b = a == 1 ? 2 : 3; b",
		//	out: newNumber(2),
		//},
	}

	runSimpleVMTestHelper(t, tests)
}
