package vm

import (
	"testing"
)

func TestBooleanObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var b = new Boolean(true); b.toString()",
			out: newString("true"),
		},
		simpleVMTest{
			in:  "var b = new Boolean(false); b.toString()",
			out: newString("false"),
		},
		simpleVMTest{
			in:  "var b = new Boolean(false); b.valueOf()",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "var b = new Boolean(true); b.valueOf()",
			out: newBool(true),
		},
	}
	runSimpleVMTestHelper(t, tests)
}
