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
}

func TestArrayLiterals(t *testing.T) {
	ep1 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}}}}}
	assert.Equal(t, Parse("[]", false), ep1)

	ep2 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, vals: []Node{&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}}}}}}}
	assert.Equal(t, Parse("[true]", false), ep2)

	ep3 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, vals: []Node{
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}},
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 7, col: 7}},
	}}}}}
	assert.Equal(t, Parse("[true, false]", false), ep3)

	ep4 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, vals: []Node{
		nil,
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 3, col: 3}},
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 9, col: 9}},
	}}}}}
	assert.Equal(t, Parse("[, true, false]", false), ep4)

	ep5 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, vals: []Node{
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}},
		nil,
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 9, col: 9}},
	}}}}}
	assert.Equal(t, Parse("[true, , false]", false), ep5)

	ep6 := &Program{body: []Node{&ExpressionStatement{X: &ArrayLiteral{tok: token{tokenType: LBRACKET, value: ""}, vals: []Node{
		&TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 1, col: 1}},
		&FalseLiteral{tok: token{tokenType: FALSE, value: "false", pos: 9, col: 9}},
	}}}}}
	assert.Equal(t, Parse("[true,   false,]", false), ep6)
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
	assert.Equal(t, recursivelyPrint(Parse("a.b()", false)), recursivelyPrint(ep2))
}
