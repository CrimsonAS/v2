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
