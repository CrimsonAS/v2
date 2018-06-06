package vm

import (
	"testing"
)

func TestMathObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "Math.round(4.7)",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "Math.round(4.4)",
			out: newNumber(4),
		},
		simpleVMTest{
			in:  "Math.sqrt(64)",
			out: newNumber(8),
		},
		simpleVMTest{
			in:  "Math.abs(1.2)",
			out: newNumber(1.2),
		},
		simpleVMTest{
			in:  "Math.abs(-1.2)",
			out: newNumber(1.2),
		},
		simpleVMTest{
			in:  "Math.ceil(4.4)",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "Math.ceil(3.2)",
			out: newNumber(4),
		},
		simpleVMTest{
			in:  "Math.floor(3.2)",
			out: newNumber(3),
		},
		simpleVMTest{
			in:  "Math.floor(4.9)",
			out: newNumber(4),
		},
		simpleVMTest{
			in:  "Math.sin(90*Math.PI/180)",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "Math.cos(0*Math.PI/180)",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "Math.min(-3, 4, 10, -9)",
			out: newNumber(-9),
		},
		simpleVMTest{
			in:  "Math.max(-3, 4, 10, -9)",
			out: newNumber(10),
		},
	}

	runSimpleVMTestHelper(t, tests)

	// untested: pow, atan2 (not implemented yet)
	// random() (no reliable return value, but we should ensure it returns 0..1
}
