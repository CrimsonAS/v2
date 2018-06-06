package parser

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIfStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&IfStatement{
			tok: token{tokenType: IF, value: "if"},
			ConditionExpr: &FalseLiteral{
				tok: token{tokenType: FALSE,
					value: "false",
					pos:   4,
					col:   4,
				},
			},
			ThenStmt: &ExpressionStatement{X: &TrueLiteral{
				tok: token{tokenType: TRUE,
					value: "true",
					pos:   11,
					col:   11,
				},
			}},
		},
	},
	}

	assert.Equal(t, Parse("if (false) true", false), ep1)

	ep2 := &Program{body: []Node{
		&IfStatement{
			tok: token{tokenType: IF, value: "if"},
			ConditionExpr: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER,
					value: "a",
					pos:   4,
					col:   4,
				},
			},
			ThenStmt: &ExpressionStatement{X: &TrueLiteral{
				tok: token{tokenType: TRUE,
					value: "true",
					pos:   7,
					col:   7,
				},
			}},
			ElseStmt: &ExpressionStatement{X: &FalseLiteral{
				tok: token{tokenType: FALSE,
					value: "false",
					pos:   17,
					col:   17,
				},
			}},
		},
	},
	}
	assert.Equal(t, Parse("if (a) true else false", false), ep2)
}

func TestReturnStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ReturnStatement{
			tok: token{tokenType: RETURN, value: "return"},
		},
	},
	}
	assert.Equal(t, Parse("return", false), ep1)

	ep2 := &Program{body: []Node{
		&ReturnStatement{
			tok: token{tokenType: RETURN, value: "return"},
			X: &FalseLiteral{
				tok: token{tokenType: FALSE,
					value: "false",
					pos:   7,
					col:   7,
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("return false", false), ep2)
}

func TestBlockStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&BlockStatement{
			tok: token{tokenType: LBRACE},
		},
	},
	}
	assert.Equal(t, Parse("{}", false), ep1)

	ep2 := &Program{body: []Node{
		&BlockStatement{
			tok: token{tokenType: LBRACE},
			Body: []Node{
				&ExpressionStatement{X: &FalseLiteral{
					tok: token{tokenType: FALSE,
						value: "false",
						pos:   2,
						col:   2,
					},
				}},
			},
		},
	},
	}
	assert.Equal(t, Parse("{ false }", false), ep2)

	ep3 := &Program{body: []Node{
		&BlockStatement{
			tok: token{tokenType: LBRACE},
			Body: []Node{
				&ExpressionStatement{X: &TrueLiteral{
					tok: token{tokenType: TRUE,
						value: "true",
						pos:   2,
						col:   2,
					},
				}},
				&ExpressionStatement{X: &FalseLiteral{
					tok: token{tokenType: FALSE,
						value: "false",
						pos:   7,
						col:   0,
						line:  1,
					},
				}},
			},
		},
	},
	}
	assert.Equal(t, Parse("{ true\nfalse }", false), ep3)
}

func TestEmptyStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&EmptyStatement{
			tok: token{tokenType: SEMICOLON, value: ""},
		},
	},
	}
	assert.Equal(t, Parse(";", false), ep1)

	ep2 := &Program{body: []Node{
		&EmptyStatement{
			tok: token{tokenType: SEMICOLON, value: ""},
		},
		&EmptyStatement{
			tok: token{tokenType: SEMICOLON, value: "",
				pos: 2,
				col: 2,
			},
		},
	},
	}
	assert.Equal(t, Parse("; ;", false), ep2)
}

func TestVariableStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&VariableStatement{
			tok: token{tokenType: VAR, value: "var"},
			Vars: []*IdentifierLiteral{
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 4, col: 4}},
			},
			Initializers: []Node{
				nil,
			},
		},
	},
	}
	assert.Equal(t, Parse("var x", false), ep1)

	ep2 := &Program{body: []Node{
		&VariableStatement{
			tok: token{tokenType: VAR, value: "var"},
			Vars: []*IdentifierLiteral{
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 4, col: 4}},
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "y", pos: 7, col: 7}},
			},
			Initializers: []Node{
				nil,
				nil,
			},
		},
	},
	}
	assert.Equal(t, Parse("var x, y", false), ep2)

	ep3 := &Program{body: []Node{
		&VariableStatement{
			tok: token{tokenType: VAR, value: "var"},
			Vars: []*IdentifierLiteral{
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 4, col: 4}},
			},
			Initializers: []Node{
				&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", pos: 8, col: 8}},
			},
		},
	},
	}
	assert.Equal(t, Parse("var x = a", false), ep3)
}

func TestDoWhileStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&DoWhileStatement{
			tok: token{tokenType: DO, value: "do"},
			X: &NumericLiteral{
				tok: token{tokenType: NUMERIC_LITERAL, value: "1", pos: 16, col: 16},
			},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 3, col: 3},
				Body: []Node{
					&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 5, col: 5}}},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("do { x } while (1)", false), ep1)
}

func TestWhileStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&WhileStatement{
			tok: token{tokenType: WHILE, value: "while"},
			X: &NumericLiteral{
				tok: token{tokenType: NUMERIC_LITERAL, value: "1", pos: 7, col: 7},
			},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 10, col: 10},
				Body: []Node{
					&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 12, col: 12}}},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("while (1) { x }", false), ep1)
}

func TestForStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ForStatement{
			tok: token{tokenType: FOR, value: "for"},
			Initializer: &NumericLiteral{
				tok: token{tokenType: NUMERIC_LITERAL, value: "1", pos: 5, col: 5},
			},
			Test: &NumericLiteral{
				tok: token{tokenType: NUMERIC_LITERAL, value: "2", pos: 7, col: 7},
			},
			Update: &NumericLiteral{
				tok: token{tokenType: NUMERIC_LITERAL, value: "3", pos: 9, col: 9},
			},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 12, col: 12},
				Body: []Node{
					&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 14, col: 14}}},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("for (1;2;3) { x }", false), ep1)

	ep2 := &Program{body: []Node{
		&ForStatement{
			tok: token{tokenType: FOR, value: "for"},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 9, col: 9},
				Body: []Node{
					&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 11, col: 11}}},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("for (;;) { x }", false), ep2)

	ep3 := &Program{body: []Node{
		&ForStatement{
			tok: token{tokenType: FOR, value: "for"},
			Initializer: &VariableStatement{
				Vars: []*IdentifierLiteral{
					&IdentifierLiteral{
						Node: nil,
						tok:  token{tokenType: 4, value: "a", pos: 9, line: 0, col: 9},
					},
				},
				Initializers: []Node{
					&NumericLiteral{
						Node: nil,
						tok:  token{tokenType: 3, value: "1", pos: 13, line: 0, col: 13},
					},
				},
				tok: token{tokenType: 52, value: "var", pos: 5, line: 0, col: 5},
			},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 18, col: 18},
				Body: []Node{
					&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 20, col: 20}}},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("for (var a = 1;;) { x }", false), ep3)
}
