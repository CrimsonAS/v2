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
	"fmt"
	"testing"

	"github.com/stvp/assert"
)

func TestFunctionExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&BinaryExpression{
		tok:  token{tokenType: ASSIGNMENT, value: "", col: 2, pos: 2},
		Left: &ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}}},
		Right: &FunctionExpression{tok: token{tokenType: FUNCTION, value: "function", col: 4, pos: 4},
			Body: &BlockStatement{
				tok:  token{tokenType: LBRACE, col: 15, pos: 15},
				Body: []Node{&ExpressionStatement{X: &TrueLiteral{tok: token{tokenType: TRUE, value: "true", col: 17, pos: 17}}}},
			},
		},
	}}}
	// for some strange reason, these don't compare equal...?
	//assert.Equal(t, Parse("a = function() { true }", false), ep1)
	assert.Equal(t, fmt.Sprintf("%s", recursivelyPrint(Parse("a = function() { true }", false))), fmt.Sprintf("%s", recursivelyPrint(ep1)))

	ep2 := &Program{body: []Node{&BinaryExpression{
		tok:  token{tokenType: ASSIGNMENT, value: "", col: 2, pos: 2},
		Left: &ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}}},
		Right: &FunctionExpression{tok: token{tokenType: FUNCTION, value: "function", col: 4, pos: 4},
			Parameters: []*IdentifierLiteral{&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b"}}},
			Body: &BlockStatement{
				tok:  token{tokenType: LBRACE, col: 15, pos: 15},
				Body: []Node{&ExpressionStatement{X: &TrueLiteral{tok: token{tokenType: TRUE, value: "true", col: 17, pos: 17}}}},
			},
		},
	}}}
	// for some strange reason, these don't compare equal...?
	//assert.Equal(t, Parse("a = function() { true }", false), ep1)
	assert.Equal(t, fmt.Sprintf("%s", recursivelyPrint(Parse("a = function(b) { true }", false))), fmt.Sprintf("%s", recursivelyPrint(ep2)))

	ep3 := &Program{body: []Node{&BinaryExpression{
		tok:  token{tokenType: ASSIGNMENT, value: "", col: 2, pos: 2},
		Left: &ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}}},
		Right: &FunctionExpression{tok: token{tokenType: FUNCTION, value: "function", col: 4, pos: 4},
			Parameters: []*IdentifierLiteral{
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b"}},
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "c"}},
			},
			Body: &BlockStatement{
				tok:  token{tokenType: LBRACE, col: 15, pos: 15},
				Body: []Node{&ExpressionStatement{X: &TrueLiteral{tok: token{tokenType: TRUE, value: "true", col: 17, pos: 17}}}},
			},
		},
	}}}
	// for some strange reason, these don't compare equal...?
	//assert.Equal(t, Parse("a = function() { true }", false), ep1)
	assert.Equal(t, fmt.Sprintf("%s", recursivelyPrint(Parse("a = function(b, c) { true }", false))), fmt.Sprintf("%s", recursivelyPrint(ep3)))
}

func TestNewExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{
		X: &NewExpression{tok: token{tokenType: NEW, value: ""}, X: &TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 4, col: 4}}},
	}}}
	assert.Equal(t, Parse("new true", false), ep1)
}

func TestCallExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &CallExpression{
		tok:       token{tokenType: LPAREN, value: "", col: 1, pos: 1},
		X:         &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Arguments: []Node{},
	}}}}
	assert.Equal(t, Parse("a()", false), ep1)

	ep2 := &Program{body: []Node{&ExpressionStatement{X: &CallExpression{
		tok: token{tokenType: LPAREN, value: "", col: 1, pos: 1},
		X:   &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Arguments: []Node{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}},
		},
	}}}}
	assert.Equal(t, Parse("a(b)", false), ep2)

	ep3 := &Program{body: []Node{&ExpressionStatement{X: &CallExpression{
		tok: token{tokenType: LPAREN, value: "", col: 1, pos: 1},
		X:   &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Arguments: []Node{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}},
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "c", col: 5, pos: 5}},
		},
	}}}}
	assert.Equal(t, Parse("a(b, c)", false), ep3)
}

func TestDotMemberExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &DotMemberExpression{
		tok:  token{tokenType: DOT, value: "", col: 1, pos: 1},
		X:    &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Name: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}}}},
	}}
	assert.Equal(t, Parse("a.b", false), ep1)
}

func TestBracketMemberExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &BracketMemberExpression{
		tok:   token{tokenType: LBRACKET, value: "", col: 1, pos: 1},
		left:  &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		right: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}}}},
	}}
	assert.Equal(t, Parse("a[b]", false), ep1)
}

// ### consider merging with TestUnaryExpression
func TestPostfixExpression(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ExpressionStatement{X: &UnaryExpression{
			tok:     token{tokenType: INCREMENT, value: "", pos: 1, col: 1},
			postfix: true,
			X: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "i"},
			},
		},
		}}}
	assert.Equal(t, Parse("i++", false), ep1)

	ep2 := &Program{body: []Node{
		&ExpressionStatement{X: &UnaryExpression{
			tok:     token{tokenType: DECREMENT, value: "", pos: 1, col: 1},
			postfix: true,
			X: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "i"},
			},
		},
		}}}
	assert.Equal(t, Parse("i--", false), ep2)
}

func TestUnaryExpression(t *testing.T) {
	type ut struct {
		tokenString string
		tokenType   TokenType
		ipos        int
		icol        int
	}

	tests := []ut{}
	tests = append(tests,
		ut{
			tokenString: "delete",
			tokenType:   DELETE,
			ipos:        7,
			icol:        7,
		},
		ut{
			tokenString: "typeof",
			tokenType:   TYPEOF,
			ipos:        7,
			icol:        7,
		},
		ut{
			tokenString: "void",
			tokenType:   VOID,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: "++",
			tokenType:   INCREMENT,
			ipos:        3,
			icol:        3,
		},
		ut{
			tokenString: "--",
			tokenType:   DECREMENT,
			ipos:        3,
			icol:        3,
		},
		ut{
			tokenString: "-",
			tokenType:   MINUS,
			ipos:        2,
			icol:        2,
		},
		ut{
			tokenString: "+",
			tokenType:   PLUS,
			ipos:        2,
			icol:        2,
		},
		ut{
			tokenString: "~",
			tokenType:   BITWISE_NOT,
			ipos:        2,
			icol:        2,
		},
		ut{
			tokenString: "!",
			tokenType:   LOGICAL_NOT,
			ipos:        2,
			icol:        2,
		},
	)

	for _, test := range tests {
		ep1 := &Program{body: []Node{
			&ExpressionStatement{X: &UnaryExpression{
				tok:     token{tokenType: test.tokenType, value: "", pos: 0, col: 0},
				postfix: false,
				X: &IdentifierLiteral{
					tok: token{tokenType: IDENTIFIER, value: "i", pos: test.ipos, col: test.icol},
				},
			},
			}}}
		assert.Equal(t, Parse(test.tokenString+" i", false), ep1)
		t.Logf("%s", fmt.Sprintf("Passed %s i", test.tokenString))
	}
}

func TestSimpleBinaryExpression(t *testing.T) {
	type ut struct {
		tokenString string
		tokenType   TokenType
		ipos        int
		icol        int
	}

	tests := []ut{}
	tests = append(tests,
		ut{
			tokenString: "*",
			tokenType:   MULTIPLY,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "/",
			tokenType:   DIVIDE,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "%",
			tokenType:   MODULUS,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "+",
			tokenType:   PLUS,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "-",
			tokenType:   MINUS,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "<<",
			tokenType:   LEFT_SHIFT,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: ">>",
			tokenType:   RIGHT_SHIFT,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: ">>>",
			tokenType:   UNSIGNED_RIGHT_SHIFT,
			ipos:        6,
			icol:        6,
		},
		ut{
			tokenString: "==",
			tokenType:   EQUALS,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: "!=",
			tokenType:   NOT_EQUALS,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: "===",
			tokenType:   STRICT_EQUALS,
			ipos:        6,
			icol:        6,
		},
		ut{
			tokenString: "!==",
			tokenType:   STRICT_NOT_EQUALS,
			ipos:        6,
			icol:        6,
		},
		ut{
			tokenString: "&",
			tokenType:   BITWISE_AND,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "^",
			tokenType:   BITWISE_XOR,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "|",
			tokenType:   BITWISE_OR,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "&&",
			tokenType:   LOGICAL_AND,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: "||",
			tokenType:   LOGICAL_OR,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: "=",
			tokenType:   ASSIGNMENT,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "<",
			tokenType:   LESS_THAN,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: ">",
			tokenType:   GREATER_THAN,
			ipos:        4,
			icol:        4,
		},
		ut{
			tokenString: "<=",
			tokenType:   LESS_EQ,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: ">=",
			tokenType:   GREATER_EQ,
			ipos:        5,
			icol:        5,
		},
		ut{
			tokenString: "instanceof",
			tokenType:   INSTANCEOF,
			ipos:        13,
			icol:        13,
		},
	)

	for _, test := range tests {
		ep1 := &Program{body: []Node{
			&ExpressionStatement{X: &BinaryExpression{
				tok: token{tokenType: test.tokenType, value: "", pos: 2, col: 2},
				Left: &IdentifierLiteral{
					tok: token{tokenType: IDENTIFIER, value: "i"},
				},
				Right: &IdentifierLiteral{
					tok: token{tokenType: IDENTIFIER, value: "i", pos: test.ipos, col: test.icol},
				},
			},
			}}}
		assert.Equal(t, Parse("i "+test.tokenString+" i", false), ep1)
		t.Logf("%s", fmt.Sprintf("Passed i %s i", test.tokenString))
	}
}

func TestConditionalExpression(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ExpressionStatement{X: &ConditionalExpression{
			tok: token{tokenType: CONDITIONAL, value: "", pos: 1, col: 1},
			X: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "a"},
			},
			Then: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "b", pos: 2, col: 2},
			},
			Else: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "c", pos: 4, col: 4},
			},
		},
		}}}
	assert.Equal(t, Parse("a?b:c", false), ep1)
}
