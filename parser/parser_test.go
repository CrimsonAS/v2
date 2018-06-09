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

func TestEmptyParse(t *testing.T) {
	ep := &Program{}
	assert.Equal(t, Parse("", false), ep)
	assert.Equal(t, Parse(" ", false), ep)
	assert.Equal(t, Parse("\t", false), ep)
	assert.Equal(t, Parse("\n", false), ep)
}

func TestLiterals(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &StringLiteral{tok: token{tokenType: STRING_LITERAL, value: "use strict"}}}}}
	assert.Equal(t, Parse("\"use strict\"", false), ep1)

	ep2 := &Program{body: []Node{&ExpressionStatement{X: &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "123.45"}}}}}
	assert.Equal(t, Parse("123.45", false), ep2)
	assert.Equal(t, ep2.body[0].(*ExpressionStatement).X.(*NumericLiteral).Float64Value(), float64(123.45))

	ep3 := &Program{body: []Node{&ExpressionStatement{X: &TrueLiteral{tok: token{tokenType: TRUE, value: "true"}}}}}
	assert.Equal(t, Parse("true", false), ep3)

	ep4 := &Program{body: []Node{&ExpressionStatement{X: &FalseLiteral{tok: token{tokenType: FALSE, value: "false"}}}}}
	assert.Equal(t, Parse("false", false), ep4)

	ep5 := &Program{body: []Node{&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}}}}}
	assert.Equal(t, Parse("a", false), ep5)

	ep6 := &Program{body: []Node{&ExpressionStatement{X: &ThisLiteral{tok: token{tokenType: THIS, value: "this"}}}}}
	assert.Equal(t, Parse("this", false), ep6)

	ep7 := &Program{body: []Node{&ExpressionStatement{X: &NullLiteral{tok: token{tokenType: NULL, value: "null"}}}}}
	assert.Equal(t, Parse("null", false), ep7)

	ep8 := &Program{body: []Node{&ExpressionStatement{X: &ThisLiteral{tok: token{tokenType: THIS, value: "this", pos: 1, col: 1}}}}}
	assert.Equal(t, Parse("(this)", false), ep8)

	ep9 := &Program{body: []Node{&ExpressionStatement{X: &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "0xff"}}}}}
	assert.Equal(t, Parse("0xff", false), ep9)
	assert.Equal(t, ep9.body[0].(*ExpressionStatement).X.(*NumericLiteral).Float64Value(), float64(255))
}

func TestArrayLiterals(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}}}}}
	assert.Equal(t, Parse("[]", false), ep1)

	ep2 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, Elements: []Node{&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}}}}}}}
	assert.Equal(t, Parse("[true]", false), ep2)

	ep3 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, Elements: []Node{
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}},
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 7, col: 7}},
	}}}}}
	assert.Equal(t, Parse("[true, false]", false), ep3)

	ep4 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, Elements: []Node{
		nil,
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 3, col: 3}},
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 9, col: 9}},
	}}}}}
	assert.Equal(t, Parse("[, true, false]", false), ep4)

	ep5 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, Elements: []Node{
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}},
		nil,
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 9, col: 9}},
	}}}}}
	assert.Equal(t, Parse("[true, , false]", false), ep5)

	ep6 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, Elements: []Node{
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}},
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 9, col: 9}},
	}}}}}
	assert.Equal(t, Parse("[true,   false,]", false), ep6)
}

func TestObjectLiterals(t *testing.T) {
	ep1 := &Program{body: []Node{&VariableStatement{
		tok: token{tokenType: VAR, value: "var"},
		Vars: []*IdentifierLiteral{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "v", pos: 4, col: 4}},
		},
		Initializers: []Node{
			&ObjectLiteral{tok: token{tokenType: LBRACE, value: "", pos: 8, col: 8}},
		},
	}}}
	assert.Equal(t, Parse("var v = {}", false), ep1)

	ep2 := &Program{body: []Node{&VariableStatement{
		tok: token{tokenType: VAR, value: "var"},
		Vars: []*IdentifierLiteral{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "v", pos: 4, col: 4}},
		},
		Initializers: []Node{
			&ObjectLiteral{
				tok: token{tokenType: LBRACE, value: "", pos: 8, col: 8},
				Properties: []ObjectPropertyLiteral{
					ObjectPropertyLiteral{
						Key: &StringLiteral{tok: token{tokenType: STRING_LITERAL, value: "a", col: 9, pos: 9}},
						X: &NumericLiteral{
							tok: token{tokenType: NUMERIC_LITERAL, value: "1", col: 14, pos: 14},
						},
					},
				},
			},
		},
	}}}
	assert.Equal(t, Parse(`var v = {"a": 1}`, false), ep2)

	ep3 := &Program{body: []Node{&VariableStatement{
		tok: token{tokenType: VAR, value: "var"},
		Vars: []*IdentifierLiteral{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "v", pos: 4, col: 4}},
		},
		Initializers: []Node{
			&ObjectLiteral{
				tok: token{tokenType: LBRACE, value: "", pos: 8, col: 8},
				Properties: []ObjectPropertyLiteral{
					ObjectPropertyLiteral{
						Key: &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "3", col: 9, pos: 9}},
						X: &NumericLiteral{
							tok: token{tokenType: NUMERIC_LITERAL, value: "1", col: 14, pos: 14},
						},
					},
				},
			},
		},
	}}}
	assert.Equal(t, Parse(`var v = {3  : 1}`, false), ep3)

	ep4 := &Program{body: []Node{&VariableStatement{
		tok: token{tokenType: VAR, value: "var"},
		Vars: []*IdentifierLiteral{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "v", pos: 4, col: 4}},
		},
		Initializers: []Node{
			&ObjectLiteral{
				tok: token{tokenType: LBRACE, value: "", pos: 8, col: 8},
				Properties: []ObjectPropertyLiteral{
					ObjectPropertyLiteral{
						Key: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", col: 9, pos: 9}},
						X: &NumericLiteral{
							tok: token{tokenType: NUMERIC_LITERAL, value: "1", col: 14, pos: 14},
						},
					},
				},
			},
		},
	}}}
	assert.Equal(t, Parse(`var v = {a  : 1}`, false), ep4)

	ep5 := &Program{body: []Node{&VariableStatement{
		tok: token{tokenType: VAR, value: "var"},
		Vars: []*IdentifierLiteral{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "v", pos: 4, col: 4}},
		},
		Initializers: []Node{
			&ObjectLiteral{
				tok: token{tokenType: LBRACE, value: "", pos: 8, col: 8},
				Properties: []ObjectPropertyLiteral{
					ObjectPropertyLiteral{
						Key: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", col: 9, pos: 9}},
						X: &NumericLiteral{
							tok: token{tokenType: NUMERIC_LITERAL, value: "1", col: 12, pos: 12},
						},
					},
					ObjectPropertyLiteral{
						Key: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 15, pos: 15}},
						X: &NumericLiteral{
							tok: token{tokenType: NUMERIC_LITERAL, value: "2", col: 18, pos: 18},
						},
					},
				},
			},
		},
	}}}
	assert.Equal(t, Parse(`var v = {a: 1, b: 2}`, false), ep5)

	ep6 := &Program{body: []Node{&VariableStatement{
		tok: token{tokenType: VAR, value: "var"},
		Vars: []*IdentifierLiteral{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "v", pos: 4, col: 4}},
		},
		Initializers: []Node{
			&ObjectLiteral{
				tok: token{tokenType: LBRACE, value: "", pos: 8, col: 8},
				Properties: []ObjectPropertyLiteral{
					ObjectPropertyLiteral{
						Type: Get,
						Key:  &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", col: 13, pos: 13}},
						X: &BlockStatement{
							tok:  token{tokenType: LBRACE, col: 17, pos: 17},
							Body: []Node{},
						},
					},
				},
			},
		},
	}}}
	assert.Equal(t, Parse(`var v = {get a() {}}`, false), ep6)
}

func TestDotExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &DotMemberExpression{tok: token{tokenType: DOT, value: "", pos: 1, col: 1},
		X:    &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", pos: 0, col: 0}},
		Name: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", pos: 2, col: 2}},
	}}}}
	assert.Equal(t, Parse("a.b", false), ep1)

	ep2 := &Program{body: []Node{&ExpressionStatement{X: &CallExpression{tok: token{tokenType: LPAREN, col: 3, pos: 3},
		X: &DotMemberExpression{tok: token{tokenType: DOT, value: "", pos: 1, col: 1},
			X:    &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", pos: 0, col: 0}},
			Name: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", pos: 2, col: 2}},
		},
		Arguments: nil,
	}}}}
	// not giving equal, for some reason
	assert.Equal(t, RecursivelyPrint(Parse("a.b()", false)), RecursivelyPrint(ep2))
}
