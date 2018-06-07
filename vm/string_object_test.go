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
