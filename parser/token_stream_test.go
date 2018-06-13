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

type tokenStreamTest struct {
	panicReason string
	input       string
	output      []token
}

func runTokenStreamTests(t *testing.T, tests []tokenStreamTest) {
	for idx, test := range tests {
		_ = idx
		s := tokenStream{stream: &byteStream{code: test.input}}
		ret := []token{}
		for !s.eof() {
			ret = append(ret, s.next())
		}
		//t.Logf("Testing %d (%s) gives %+v want %+v", idx, escapeStringToPrint(test.input), ret, test.output)
		assert.Equal(t, ret, test.output)
		t.Logf("Pass %s", escapeStringToPrint(test.input))
	}
}

func TestTokenStreamSimple(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input:  "",
			output: []token{},
		},
		tokenStreamTest{
			input:  " ",
			output: []token{},
		},
		tokenStreamTest{
			input:  "\n\n",
			output: []token{},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestSingleLineComments(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "// Hello world",
			output: []token{
				token{
					tokenType: COMMENT,
					value:     " Hello world",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "// Hello world\n//How are you",
			output: []token{
				token{
					tokenType: COMMENT,
					value:     " Hello world",
					pos:       0,
					col:       0,
					line:      0,
				},
				token{
					tokenType: COMMENT,
					value:     "How are you",
					pos:       15,
					col:       0,
					line:      1,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestMultiLineComments(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "/* Hello world */",
			output: []token{
				token{
					tokenType: COMMENT,
					value:     " Hello world ",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "/* Hello\nworld */",
			output: []token{
				token{
					tokenType: COMMENT,
					value:     " Hello\nworld ",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestStringLiterals(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: `"how are you"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "how are you",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: `"how \"are\" you"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     `how "are" you`,
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: `'how are you'`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "how are you",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: `'how \'are\' you'`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "how 'are' you",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: `"\x41"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "A",
				},
			},
		},
		tokenStreamTest{
			input: `"\x414"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "A4",
				},
			},
		},
		tokenStreamTest{
			input: `"\x41\x42\x43\x44"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "ABCD",
				},
			},
		},
		tokenStreamTest{
			input: `"\u2665"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "♥",
				},
			},
		},
		tokenStreamTest{
			input: `"\u26651"`,
			output: []token{
				token{
					tokenType: STRING_LITERAL,
					value:     "♥1",
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestNumberLiterals(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "0",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "0",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "1",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "1",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "1234567890",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "1234567890",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "12345.1234",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "12345.1234",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "0x0",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "0x0",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "0x0123456789abcdef",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "0x0123456789abcdef",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestIdentifiers(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "a",
			output: []token{
				token{
					tokenType: IDENTIFIER,
					value:     "a",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "abcdefghijklmnopqrstuvwxyz",
			output: []token{
				token{
					tokenType: IDENTIFIER,
					value:     "abcdefghijklmnopqrstuvwxyz",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			output: []token{
				token{
					tokenType: IDENTIFIER,
					value:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "a_test",
			output: []token{
				token{
					tokenType: IDENTIFIER,
					value:     "a_test",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
		tokenStreamTest{
			input: "a1234567890",
			output: []token{
				token{
					tokenType: IDENTIFIER,
					value:     "a1234567890",
					pos:       0,
					col:       0,
					line:      0,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestOperators(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "+",
			output: []token{
				token{
					tokenType: PLUS,
				},
			},
		},
		tokenStreamTest{
			input: "^=",
			output: []token{
				token{
					tokenType: XOR_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "|=",
			output: []token{
				token{
					tokenType: OR_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "&=",
			output: []token{
				token{
					tokenType: AND_EQ,
				},
			},
		},
		tokenStreamTest{
			input: ">>=",
			output: []token{
				token{
					tokenType: RIGHT_SHIFT_EQ,
				},
			},
		},
		tokenStreamTest{
			input: ">>>=",
			output: []token{
				token{
					tokenType: UNSIGNED_RIGHT_SHIFT_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "<<=",
			output: []token{
				token{
					tokenType: LEFT_SHIFT_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "+=",
			output: []token{
				token{
					tokenType: PLUS_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "-",
			output: []token{
				token{
					tokenType: MINUS,
				},
			},
		},
		tokenStreamTest{
			input: "-=",
			output: []token{
				token{
					tokenType: MINUS_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "*",
			output: []token{
				token{
					tokenType: MULTIPLY,
				},
			},
		},
		tokenStreamTest{
			input: "*=",
			output: []token{
				token{
					tokenType: MULTIPLY_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "/",
			output: []token{
				token{
					tokenType: DIVIDE,
				},
			},
		},
		tokenStreamTest{
			input: "/=",
			output: []token{
				token{
					tokenType: DIVIDE_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "%",
			output: []token{
				token{
					tokenType: MODULUS,
				},
			},
		},
		tokenStreamTest{
			input: "%=",
			output: []token{
				token{
					tokenType: MODULUS_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "=",
			output: []token{
				token{
					tokenType: ASSIGNMENT,
				},
			},
		},
		tokenStreamTest{
			input: "==",
			output: []token{
				token{
					tokenType: EQUALS,
				},
			},
		},
		tokenStreamTest{
			input: "===",
			output: []token{
				token{
					tokenType: STRICT_EQUALS,
				},
			},
		},
		tokenStreamTest{
			input: "&",
			output: []token{
				token{
					tokenType: BITWISE_AND,
				},
			},
		},
		tokenStreamTest{
			input: "&&",
			output: []token{
				token{
					tokenType: LOGICAL_AND,
				},
			},
		},
		tokenStreamTest{
			input: "|",
			output: []token{
				token{
					tokenType: BITWISE_OR,
				},
			},
		},
		tokenStreamTest{
			input: "||",
			output: []token{
				token{
					tokenType: LOGICAL_OR,
				},
			},
		},
		tokenStreamTest{
			input: "<",
			output: []token{
				token{
					tokenType: LESS_THAN,
				},
			},
		},
		tokenStreamTest{
			input: "<=",
			output: []token{
				token{
					tokenType: LESS_EQ,
				},
			},
		},
		tokenStreamTest{
			input: "<<",
			output: []token{
				token{
					tokenType: LEFT_SHIFT,
				},
			},
		},
		tokenStreamTest{
			input: ">",
			output: []token{
				token{
					tokenType: GREATER_THAN,
				},
			},
		},
		tokenStreamTest{
			input: ">=",
			output: []token{
				token{
					tokenType: GREATER_EQ,
				},
			},
		},
		tokenStreamTest{
			input: ">>",
			output: []token{
				token{
					tokenType: RIGHT_SHIFT,
				},
			},
		},
		tokenStreamTest{
			input: ">>>",
			output: []token{
				token{
					tokenType: UNSIGNED_RIGHT_SHIFT,
				},
			},
		},
		tokenStreamTest{
			input: "^",
			output: []token{
				token{
					tokenType: BITWISE_XOR,
				},
			},
		},
		tokenStreamTest{
			input: "instanceof",
			output: []token{
				token{
					tokenType: INSTANCEOF,
				},
			},
		},
		tokenStreamTest{
			input: "in",
			output: []token{
				token{
					tokenType: IN,
				},
			},
		},
		tokenStreamTest{
			input: "new",
			output: []token{
				token{
					tokenType: NEW,
				},
			},
		},
		tokenStreamTest{
			input: "?",
			output: []token{
				token{
					tokenType: CONDITIONAL,
				},
			},
		},
		tokenStreamTest{
			input: "!",
			output: []token{
				token{
					tokenType: LOGICAL_NOT,
				},
			},
		},
		tokenStreamTest{
			input: "!=",
			output: []token{
				token{
					tokenType: NOT_EQUALS,
				},
			},
		},
		tokenStreamTest{
			input: "!==",
			output: []token{
				token{
					tokenType: STRICT_NOT_EQUALS,
				},
			},
		},
		tokenStreamTest{
			input: "~",
			output: []token{
				token{
					tokenType: BITWISE_NOT,
				},
			},
		},
		tokenStreamTest{
			input: "delete",
			output: []token{
				token{
					tokenType: DELETE,
				},
			},
		},
		tokenStreamTest{
			input: "typeof",
			output: []token{
				token{
					tokenType: TYPEOF,
				},
			},
		},
		tokenStreamTest{
			input: "void",
			output: []token{
				token{
					tokenType: VOID,
				},
			},
		},
		tokenStreamTest{
			input: "++",
			output: []token{
				token{
					tokenType: INCREMENT,
				},
			},
		},
		tokenStreamTest{
			input: "--",
			output: []token{
				token{
					tokenType: DECREMENT,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestPunctuation(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: ".",
			output: []token{
				token{
					tokenType: DOT,
				},
			},
		},
		tokenStreamTest{
			input: ",",
			output: []token{
				token{
					tokenType: COMMA,
				},
			},
		},
		tokenStreamTest{
			input: ";",
			output: []token{
				token{
					tokenType: SEMICOLON,
				},
			},
		},
		tokenStreamTest{
			input: "(",
			output: []token{
				token{
					tokenType: LPAREN,
				},
			},
		},
		tokenStreamTest{
			input: ")",
			output: []token{
				token{
					tokenType: RPAREN,
				},
			},
		},
		tokenStreamTest{
			input: "[",
			output: []token{
				token{
					tokenType: LBRACKET,
				},
			},
		},
		tokenStreamTest{
			input: "]",
			output: []token{
				token{
					tokenType: RBRACKET,
				},
			},
		},
		tokenStreamTest{
			input: "{",
			output: []token{
				token{
					tokenType: LBRACE,
				},
			},
		},
		tokenStreamTest{
			input: "}",
			output: []token{
				token{
					tokenType: RBRACE,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestSimpleExpression(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "1 + 2.25 = 3.25 // And that's all, folks",
			output: []token{
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "1",
				},
				token{
					tokenType: PLUS,
					pos:       2,
					col:       2,
				},
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "2.25",
					pos:       4,
					col:       4,
				},
				token{
					tokenType: ASSIGNMENT,
					pos:       9,
					col:       9,
				},
				token{
					tokenType: NUMERIC_LITERAL,
					value:     "3.25",
					pos:       11,
					col:       11,
				},
				token{
					tokenType: COMMENT,
					value:     " And that's all, folks",
					pos:       16,
					col:       16,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestTokens(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "if",
			output: []token{
				token{
					tokenType: IF,
					value:     "if",
				},
			},
		},
		tokenStreamTest{
			input: "else",
			output: []token{
				token{
					tokenType: ELSE,
					value:     "else",
				},
			},
		},
		tokenStreamTest{
			input: "this",
			output: []token{
				token{
					tokenType: THIS,
					value:     "this",
				},
			},
		},
		tokenStreamTest{
			input: "return",
			output: []token{
				token{
					tokenType: RETURN,
					value:     "return",
				},
			},
		},
		tokenStreamTest{
			input: "null",
			output: []token{
				token{
					tokenType: NULL,
					value:     "null",
				},
			},
		},
		tokenStreamTest{
			input: "true",
			output: []token{
				token{
					tokenType: TRUE,
					value:     "true",
				},
			},
		},
		tokenStreamTest{
			input: "false",
			output: []token{
				token{
					tokenType: FALSE,
					value:     "false",
				},
			},
		},
		tokenStreamTest{
			input: "function",
			output: []token{
				token{
					tokenType: FUNCTION,
					value:     "function",
				},
			},
		},
		tokenStreamTest{
			input: "var",
			output: []token{
				token{
					tokenType: VAR,
					value:     "var",
				},
			},
		},
		tokenStreamTest{
			input: "do",
			output: []token{
				token{
					tokenType: DO,
					value:     "do",
				},
			},
		},
		tokenStreamTest{
			input: "while",
			output: []token{
				token{
					tokenType: WHILE,
					value:     "while",
				},
			},
		},
		tokenStreamTest{
			input: "get",
			output: []token{
				token{
					tokenType: GET,
					value:     "get",
				},
			},
		},
		tokenStreamTest{
			input: "set",
			output: []token{
				token{
					tokenType: SET,
					value:     "set",
				},
			},
		},
		tokenStreamTest{
			input: "switch",
			output: []token{
				token{
					tokenType: SWITCH,
					value:     "switch",
				},
			},
		},
		tokenStreamTest{
			input: "case",
			output: []token{
				token{
					tokenType: CASE,
					value:     "case",
				},
			},
		},
		tokenStreamTest{
			input: "default",
			output: []token{
				token{
					tokenType: DEFAULT,
					value:     "default",
				},
			},
		},
		tokenStreamTest{
			input: "throw",
			output: []token{
				token{
					tokenType: THROW,
					value:     "throw",
				},
			},
		},
		tokenStreamTest{
			input: "try",
			output: []token{
				token{
					tokenType: TRY,
					value:     "try",
				},
			},
		},
		tokenStreamTest{
			input: "catch",
			output: []token{
				token{
					tokenType: CATCH,
					value:     "catch",
				},
			},
		},
		tokenStreamTest{
			input: "finally",
			output: []token{
				token{
					tokenType: FINALLY,
					value:     "finally",
				},
			},
		},
		tokenStreamTest{
			input: "for",
			output: []token{
				token{
					tokenType: FOR,
					value:     "for",
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}

func TestIfElse(t *testing.T) {
	tests := []tokenStreamTest{
		tokenStreamTest{
			input: "if (a) {\n} else if (b) {}",
			output: []token{
				token{
					tokenType: IF,
					value:     "if",
				},
				token{
					tokenType: LPAREN,
					pos:       3,
					col:       3,
				},
				token{
					tokenType: IDENTIFIER,
					value:     "a",
					pos:       4,
					col:       4,
				},
				token{
					tokenType: RPAREN,
					pos:       5,
					col:       5,
				},
				token{
					tokenType: LBRACE,
					pos:       7,
					col:       7,
				},
				token{
					tokenType: RBRACE,
					pos:       9,
					line:      1,
					col:       0,
				},
				token{
					tokenType: ELSE,
					value:     "else",
					pos:       11,
					line:      1,
					col:       2,
				},
				token{
					tokenType: IF,
					value:     "if",
					pos:       16,
					line:      1,
					col:       7,
				},

				token{
					tokenType: LPAREN,
					pos:       19,
					line:      1,
					col:       10,
				},
				token{
					tokenType: IDENTIFIER,
					value:     "b",
					pos:       20,
					line:      1,
					col:       11,
				},
				token{
					tokenType: RPAREN,
					pos:       21,
					line:      1,
					col:       12,
				},

				token{
					tokenType: LBRACE,
					pos:       23,
					line:      1,
					col:       14,
				},
				token{
					tokenType: RBRACE,
					pos:       24,
					line:      1,
					col:       15,
				},
			},
		},
	}
	runTokenStreamTests(t, tests)
}
