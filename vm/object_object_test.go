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

func TestObjectObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var b = true; var bo = Object(b); return bo.toString()",
			out: newString("true"),
		},
		simpleVMTest{
			in:  "var b = true; var bo = new Object(b); return bo.toString()",
			out: newString("true"),
		},
		// missing Number object
		//simpleVMTest{
		//	in:  "var b = 1; var bo = new Object(b); return bo.toString()",
		//	out: newString("1"),
		//},
		simpleVMTest{
			in:  "var b = null; var bo = new Object(b); return bo.toString()",
			out: newString("[object Object]"),
		},
		simpleVMTest{
			in:  "var b = undefined; var bo = new Object(b); return bo.toString()",
			out: newString("[object Object]"),
		},
		simpleVMTest{
			in:  "var bo = new Object(); return bo.toString()",
			out: newString("[object Object]"),
		},
		simpleVMTest{
			in:  `var bo = new Object(); return bo.hasOwnProperty("foo")`,
			out: newBool(false),
		},
		simpleVMTest{
			in:  `var bo = new Object(); bo[55] = 66; return bo[55]`,
			out: newNumber(66),
		},
	}
	runSimpleVMTestHelper(t, tests)
}
