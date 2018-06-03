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

func TestPrefixOperators(t *testing.T) {
	type simpleTest struct {
		in  string
		out value
	}

	tests := []simpleTest{
		simpleTest{
			in:  "++1",
			out: newNumber(2),
		},
		simpleTest{
			in:  "--1",
			out: newNumber(0),
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

		// order of operations is broken.
		// something wrong with unary expressions?
		simpleTest{
			in:  "2+2*2+2",
			out: newNumber(8),
		},
	}

	for _, test := range tests {
		ast := parser.Parse(test.in)
		//t.Logf("Code %s\n%# v", test.in, pretty.Formatter(ast))
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
