package vm

import (
	"testing"
)

func TestStringObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var s = new String(\"hello\"); s.toString()",
			out: newString("hello"),
		},
		simpleVMTest{
			in:  "var s = new String(\"hello\"); s.valueOf()",
			out: newString("hello"),
		},
		// ### es5 8.7.1 GetValue(), we need to promote primitives to object
		//simpleVMTest{
		//	in:  "var s = String(\"hello\"); s.toString()",
		//	out: newString("hello"),
		//},
		simpleVMTest{
			in:  "var s = new String(\"hello\"); s.charAt(0)",
			out: newString("h"),
		},
		simpleVMTest{
			in:  "var s = new String(\"hi\"); s.charAt(1)",
			out: newString("i"),
		},
		simpleVMTest{
			in:  "var s = new String(\"ABC\"); s.charCodeAt(0)",
			out: newNumber(65),
		},
		simpleVMTest{
			in:  `var s = new String("I"); s.concat(" am", " simply", " the", " best")`,
			out: newString("I am simply the best"),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); s.indexOf("A")`,
			out: newNumber(0),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); s.indexOf("B")`,
			out: newNumber(1),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); s.indexOf("C")`,
			out: newNumber(2),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); s.indexOf("A", 1)`,
			out: newNumber(2),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); s.indexOf("N")`,
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); s.lastIndexOf("A")`,
			out: newNumber(0),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); s.lastIndexOf("B")`,
			out: newNumber(1),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); s.lastIndexOf("C")`,
			out: newNumber(2),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); s.lastIndexOf("A", 1)`,
			out: newNumber(0),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); s.lastIndexOf("N")`,
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  `var s = new String("abba"); s.toUpperCase()`,
			out: newString("ABBA"),
		},
		simpleVMTest{
			in:  `var s = new String("ABBA"); s.toLowerCase()`,
			out: newString("abba"),
		},
		simpleVMTest{
			in:  `var s = new String("abcd"); s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("  abcd"); s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("abcd    "); s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("    abcd    "); s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("    ab  cd    "); s.trim()`,
			out: newString("ab  cd"),
		},
	}

	runSimpleVMTestHelper(t, tests)
}
