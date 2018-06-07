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

package parser

import (
	"testing"

	"github.com/stvp/assert"
)

type byteStreamTest struct {
	panicReason  string
	input        string
	output       []byte
	expectedPos  int
	expectedLine int
	expectedCol  int
}

func escapeStringToPrint(str string) string {
	ret := ""
	for _, c := range str {
		if c == '\\' {
			ret += "\\\\"
		} else if c == '\n' {
			ret += "\\n"
		} else {
			ret += string(c)
		}
	}
	return ret
}

func runByteStreamTests(t *testing.T, tests []byteStreamTest) {
	for _, test := range tests {
		s := byteStream{code: test.input}
		ret := []byte{}
		for !s.eof() {
			ret = append(ret, s.next())
		}
		assert.Equal(t, ret, test.output)
		assert.Equal(t, s.pos, test.expectedPos)
		assert.Equal(t, s.line, test.expectedLine)
		assert.Equal(t, s.col, test.expectedCol)
		t.Logf("Pass %s", escapeStringToPrint(test.input))
	}
}

func TestByteStreamSimple(t *testing.T) {
	tests := []byteStreamTest{
		byteStreamTest{
			input:        "",
			output:       []byte{},
			expectedPos:  0,
			expectedCol:  0,
			expectedLine: 0,
		},
		byteStreamTest{
			input:        "a",
			output:       []byte{'a'},
			expectedPos:  1,
			expectedCol:  1,
			expectedLine: 0,
		},
		byteStreamTest{
			input:        "ab",
			output:       []byte{'a', 'b'},
			expectedPos:  2,
			expectedCol:  2,
			expectedLine: 0,
		},
		byteStreamTest{
			input:        "\n",
			output:       []byte{'\n'},
			expectedPos:  1,
			expectedCol:  0,
			expectedLine: 1,
		},
		byteStreamTest{
			input:        "\na",
			output:       []byte{'\n', 'a'},
			expectedPos:  2,
			expectedCol:  1,
			expectedLine: 1,
		},
		byteStreamTest{
			input:        "a\n",
			output:       []byte{'a', '\n'},
			expectedPos:  2,
			expectedCol:  0,
			expectedLine: 1,
		},
	}
	runByteStreamTests(t, tests)
}
