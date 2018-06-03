package parser

import (
	"fmt"
	"testing"

	"github.com/stvp/assert"
)

func TestFunctionExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&BinaryExpression{
		tok:  token{tokenType: ASSIGNMENT, value: "", col: 2, pos: 2},
		Left: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Right: &FunctionExpression{tok: token{tokenType: FUNCTION, value: "function", col: 4, pos: 4},
			Body: &BlockStatement{
				tok:  token{tokenType: LBRACE, col: 15, pos: 15},
				Body: []Node{&TrueLiteral{tok: token{tokenType: TRUE, value: "true", col: 17, pos: 17}}},
			},
		},
	}}}
	// for some strange reason, these don't compare equal...?
	//assert.Equal(t, Parse("a = function() { true }"), ep1)
	assert.Equal(t, fmt.Sprintf("%s", recursivelyPrint(Parse("a = function() { true }"))), fmt.Sprintf("%s", recursivelyPrint(ep1)))

	ep2 := &Program{body: []Node{&BinaryExpression{
		tok:  token{tokenType: ASSIGNMENT, value: "", col: 2, pos: 2},
		Left: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Right: &FunctionExpression{tok: token{tokenType: FUNCTION, value: "function", col: 4, pos: 4},
			Parameters: []*IdentifierLiteral{&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b"}}},
			Body: &BlockStatement{
				tok:  token{tokenType: LBRACE, col: 15, pos: 15},
				Body: []Node{&TrueLiteral{tok: token{tokenType: TRUE, value: "true", col: 17, pos: 17}}},
			},
		},
	}}}
	// for some strange reason, these don't compare equal...?
	//assert.Equal(t, Parse("a = function() { true }"), ep1)
	assert.Equal(t, fmt.Sprintf("%s", recursivelyPrint(Parse("a = function(b) { true }"))), fmt.Sprintf("%s", recursivelyPrint(ep2)))

	ep3 := &Program{body: []Node{&BinaryExpression{
		tok:  token{tokenType: ASSIGNMENT, value: "", col: 2, pos: 2},
		Left: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Right: &FunctionExpression{tok: token{tokenType: FUNCTION, value: "function", col: 4, pos: 4},
			Parameters: []*IdentifierLiteral{
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b"}},
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "c"}},
			},
			Body: &BlockStatement{
				tok:  token{tokenType: LBRACE, col: 15, pos: 15},
				Body: []Node{&TrueLiteral{tok: token{tokenType: TRUE, value: "true", col: 17, pos: 17}}},
			},
		},
	}}}
	// for some strange reason, these don't compare equal...?
	//assert.Equal(t, Parse("a = function() { true }"), ep1)
	assert.Equal(t, fmt.Sprintf("%s", recursivelyPrint(Parse("a = function(b, c) { true }"))), fmt.Sprintf("%s", recursivelyPrint(ep3)))
}

func TestNewExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&NewExpression{tok: token{tokenType: NEW, value: ""}, expr: &TrueLiteral{tok: token{tokenType: TRUE, value: "true", pos: 4, col: 4}}}}}
	assert.Equal(t, Parse("new true"), ep1)
}

func TestCallExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&CallExpression{
		tok:       token{tokenType: LPAREN, value: "", col: 1, pos: 1},
		X:         &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Arguments: []Node{},
	}}}
	assert.Equal(t, Parse("a()"), ep1)

	ep2 := &Program{body: []Node{&CallExpression{
		tok: token{tokenType: LPAREN, value: "", col: 1, pos: 1},
		X:   &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Arguments: []Node{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}},
		},
	}}}
	assert.Equal(t, Parse("a(b)"), ep2)

	ep3 := &Program{body: []Node{&CallExpression{
		tok: token{tokenType: LPAREN, value: "", col: 1, pos: 1},
		X:   &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		Arguments: []Node{
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}},
			&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "c", col: 5, pos: 5}},
		},
	}}}
	assert.Equal(t, Parse("a(b, c)"), ep3)
}

func TestDotMemberExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&DotMemberExpression{
		tok:   token{tokenType: DOT, value: "", col: 1, pos: 1},
		left:  &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		right: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}}}},
	}
	assert.Equal(t, Parse("a.b"), ep1)
}

func TestBracketMemberExpression(t *testing.T) {
	ep1 := &Program{body: []Node{&BracketMemberExpression{
		tok:   token{tokenType: LBRACKET, value: "", col: 1, pos: 1},
		left:  &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a"}},
		right: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "b", col: 2, pos: 2}}}},
	}
	assert.Equal(t, Parse("a[b]"), ep1)
}

// ### consider merging with TestUnaryExpression
func TestPostfixExpression(t *testing.T) {
	ep1 := &Program{body: []Node{
		&UnaryExpression{
			tok:     token{tokenType: INCREMENT, value: "", pos: 1, col: 1},
			postfix: true,
			X: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "i"},
			},
		},
	}}
	assert.Equal(t, Parse("i++"), ep1)

	ep2 := &Program{body: []Node{
		&UnaryExpression{
			tok:     token{tokenType: DECREMENT, value: "", pos: 1, col: 1},
			postfix: true,
			X: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "i"},
			},
		},
	}}
	assert.Equal(t, Parse("i--"), ep2)
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
			&UnaryExpression{
				tok:     token{tokenType: test.tokenType, value: "", pos: 0, col: 0},
				postfix: false,
				X: &IdentifierLiteral{
					tok: token{tokenType: IDENTIFIER, value: "i", pos: test.ipos, col: test.icol},
				},
			},
		}}
		t.Logf("Testing %s", test.tokenString)
		assert.Equal(t, Parse(test.tokenString+" i"), ep1)
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
			&BinaryExpression{
				tok: token{tokenType: test.tokenType, value: "", pos: 2, col: 2},
				Left: &IdentifierLiteral{
					tok: token{tokenType: IDENTIFIER, value: "i"},
				},
				Right: &IdentifierLiteral{
					tok: token{tokenType: IDENTIFIER, value: "i", pos: test.ipos, col: test.icol},
				},
			},
		}}
		assert.Equal(t, Parse("i "+test.tokenString+" i"), ep1)
		t.Logf("%s", fmt.Sprintf("Passed i %s i", test.tokenString))
	}
}

func TestConditionalExpression(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ConditionalExpression{
			tok: token{tokenType: CONDITIONAL, value: "", pos: 1, col: 1},
			test: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "a"},
			},
			trueBranch: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "b", pos: 2, col: 2},
			},
			falseBranch: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "c", pos: 4, col: 4},
			},
		},
	}}
	assert.Equal(t, Parse("a?b:c"), ep1)
}
