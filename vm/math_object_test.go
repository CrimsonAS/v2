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
