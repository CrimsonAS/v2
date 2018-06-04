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
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		vm := NewVM(ast)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("Passed %s == %s", test.in, test.out)
	}
}

func TestBuiltin(t *testing.T) {
	testFunc := func(vm *vm, f value, args []value) value {
		t.Logf("I'm a message from a builtin function: %+v", args)
		return newString("Hello world")
	}

	ast := parser.Parse("testFunc()")
	vm := NewVM(ast)
	pf := newFunctionObject(testFunc)
	vm.defineVar("testFunc", pf)
	assert.Equal(t, vm.Run(), newString("Hello world"))
}
