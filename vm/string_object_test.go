/*
 * Copyright 2018 Crimson AS <info@crimson.no>
 * Author: Robin Burchell <robin.burchell@crimson.no>
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package vm

import (
	"testing"
)

func TestStringObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var s = new String(\"hello\"); return s.toString()",
			out: newString("hello"),
		},
		simpleVMTest{
			in:  "var s = new String(\"hello\"); return s.valueOf()",
			out: newString("hello"),
		},
		simpleVMTest{
			in:  "var s = new String(); return s.valueOf()",
			out: newString(""),
		},
		simpleVMTest{
			in:  "var s = String(\"hello\"); return s.toString()",
			out: newString("hello"),
		},
		simpleVMTest{
			in:  "var s = new String(\"hello\"); return s.charAt(0)",
			out: newString("h"),
		},
		simpleVMTest{
			in:  "var s = new String(\"hi\"); return s.charAt(1)",
			out: newString("i"),
		},
		simpleVMTest{
			in:  "var s = new String(\"ABC\"); return s.charCodeAt(0)",
			out: newNumber(65),
		},
		simpleVMTest{
			in:  `var s = new String("I"); return s.concat(" am", " simply", " the", " best")`,
			out: newString("I am simply the best"),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); return s.indexOf("A")`,
			out: newNumber(0),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); return s.indexOf("B")`,
			out: newNumber(1),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); return s.indexOf("C")`,
			out: newNumber(2),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); return s.indexOf("A", 1)`,
			out: newNumber(2),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); return s.indexOf("N")`,
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); return s.lastIndexOf("A")`,
			out: newNumber(0),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); return s.lastIndexOf("B")`,
			out: newNumber(1),
		},
		simpleVMTest{
			in:  `var s = new String("ABC"); return s.lastIndexOf("C")`,
			out: newNumber(2),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); return s.lastIndexOf("A", 1)`,
			out: newNumber(0),
		},
		simpleVMTest{
			in:  `var s = new String("ABA"); return s.lastIndexOf("N")`,
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  `var s = new String("abba"); return s.toUpperCase()`,
			out: newString("ABBA"),
		},
		simpleVMTest{
			in:  `var s = new String("ABBA"); return s.toLowerCase()`,
			out: newString("abba"),
		},
		simpleVMTest{
			in:  `var s = new String("abcd"); return s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("  abcd"); return s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("abcd    "); return s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("    abcd    "); return s.trim()`,
			out: newString("abcd"),
		},
		simpleVMTest{
			in:  `var s = new String("    ab  cd    "); return s.trim()`,
			out: newString("ab  cd"),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[-1]`,
			out: newUndefined(),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[0]`,
			out: newString("h"),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[1]`,
			out: newString("e"),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[2]`,
			out: newString("l"),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[3]`,
			out: newString("l"),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[4]`,
			out: newString("o"),
		},
		simpleVMTest{
			in:  `var s = new String("hello"); return s[5]`,
			out: newUndefined(),
		},
	}

	runSimpleVMTestHelper(t, tests)
}
