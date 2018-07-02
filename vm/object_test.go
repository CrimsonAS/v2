package vm

import (
	"testing"
)

func TestObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = {}; return a.toString()",
			out: newString("[object Object]"),
		},
		simpleVMTest{
			in:  "var a = {abc: 5}; return a.abc",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  `var a = {abc: 5, def: "test"}; return a.def`,
			out: newString("test"),
		},
	}

	runSimpleVMTestHelper(t, tests)
}
