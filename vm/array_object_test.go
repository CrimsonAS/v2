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
			in:  "var b = []; return b",
			out: newArrayObject(nil),
		},
		simpleVMTest{
			in:  "var b = new Array(); return b",
			out: newArrayObject(nil),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; return b",
			out: newArrayObject([]value{newNumber(1), newNumber(2), newNumber(3), newNumber(4), newNumber(5)}),
		},
		simpleVMTest{
			in:  "var b = new Array(1, 2, 3, 4, 5); return b",
			out: newArrayObject([]value{newNumber(1), newNumber(2), newNumber(3), newNumber(4), newNumber(5)}),
		},
		simpleVMTest{
			in:  "var b = Array(1, 2, 3, 4, 5); return b",
			out: newArrayObject([]value{newNumber(1), newNumber(2), newNumber(3), newNumber(4), newNumber(5)}),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayObjectReadWrite(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; return b[0]",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; return b[1]",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; return b[2]",
			out: newNumber(3),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; return b[3]",
			out: newNumber(4),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; return b[4]",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "var b = [1, 2, 3, 4, 5]; b[4] = 255; return b[4]",
			out: newNumber(255),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayToString(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a', 'b']; return a.toString()",
			out: newString("a,b"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayConcat(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a']; var b = ['b']; c = a.concat(b); return c.toString()",
			out: newString("a,b"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayJoin(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a', 'b']; return a.join('hello')",
			out: newString("ahellob"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayPop(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = []; var b = a.pop(); return b",
			out: newUndefined(),
		},
		simpleVMTest{
			in:  "var a = [1]; var b = a.pop(); return b",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = [1, 2]; var b = a.pop(); return b",
			out: newNumber(2),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayPush(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a']; return a.push('b')",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var a = ['a']; a.push('b'); return a.toString()",
			out: newString("a,b"),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b']; a.push('c'); return a.toString()",
			out: newString("a,b,c"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayReverse(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; a.reverse(); return a.toString()",
			out: newString("d,c,b,a"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayShift(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = []; return a.shift()",
			out: newUndefined(),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.shift()",
			out: newString("a"),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; a.shift(); return a.toString()",
			out: newString("b,c,d"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArraySlice(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; var b = a.slice(0, 2); return b.toString()",
			out: newString("a,b"),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; var b = a.slice(0, 3); return b.toString()",
			out: newString("a,b,c"),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayUnshift(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['c', 'd']; a.unshift('a', 'b'); return a.toString()",
			out: newString("a,b,c,d"),
		},
		simpleVMTest{
			in:  "var a = ['c', 'd']; return a.unshift('a', 'b')",
			out: newNumber(4),
		},
	}
	runSimpleVMTestHelper(t, tests)
}

func TestArrayIndexOf(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.indexOf('e')",
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.indexOf('a')",
			out: newNumber(0),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.indexOf('b')",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.indexOf('d')",
			out: newNumber(3),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'a']; return a.indexOf('a')",
			out: newNumber(0),
		},
	}
	runSimpleVMTestHelper(t, tests)
}
func TestArrayLastIndexOf(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.lastIndexOf('e')",
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.lastIndexOf('a')",
			out: newNumber(0),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.lastIndexOf('b')",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'd']; return a.lastIndexOf('d')",
			out: newNumber(3),
		},
		simpleVMTest{
			in:  "var a = ['a', 'b', 'c', 'a']; return a.lastIndexOf('a')",
			out: newNumber(3),
		},
	}
	runSimpleVMTestHelper(t, tests)
}
