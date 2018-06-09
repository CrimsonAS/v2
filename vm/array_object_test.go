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

func TestArrayObject(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var b = []; b",
			out: newArrayObject(nil),
		},
		simpleVMTest{
			in:  "var b = new Array(); b",
			out: newArrayObject(nil),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; b",
			out: newArrayObject([]value{newNumber(1), newNumber(2), newNumber(3), newNumber(4), newNumber(5)}),
		},
		simpleVMTest{
			in:  "var b = new Array(1, 2, 3, 4, 5); b",
			out: newArrayObject([]value{newNumber(1), newNumber(2), newNumber(3), newNumber(4), newNumber(5)}),
		},
		simpleVMTest{
			in:  "var b = Array(1, 2, 3, 4, 5); b",
			out: newArrayObject([]value{newNumber(1), newNumber(2), newNumber(3), newNumber(4), newNumber(5)}),
		},
	}
	runSimpleVMTestHelper(t, tests)
}
