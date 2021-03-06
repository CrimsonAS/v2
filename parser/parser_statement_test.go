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

	_ "github.com/kr/pretty"
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
			tok:  token{tokenType: LBRACE},
			Body: []Node{},
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
						tok:  token{tokenType: IDENTIFIER, value: "a", pos: 9, line: 0, col: 9},
					},
				},
				Initializers: []Node{
					&NumericLiteral{
						Node: nil,
						tok:  token{tokenType: NUMERIC_LITERAL, value: "1", pos: 13, line: 0, col: 13},
					},
				},
				tok: token{tokenType: VAR, value: "var", pos: 5, line: 0, col: 5},
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

func TestForInStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ForInStatement{
			tok: token{tokenType: FOR, value: "for"},
			X: &VariableStatement{
				Vars:         []*IdentifierLiteral{&IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "i", pos: 9, col: 9}}},
				Initializers: []Node{nil},
				tok:          token{tokenType: VAR, value: "var", pos: 5, col: 5},
			},
			Y: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "k", pos: 14, col: 14},
			},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 17, col: 17},
				Body: []Node{
					&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "x", pos: 19, col: 19}}},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("for (var i in k) { x }", false), ep1)
}

func TestSwitchStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&SwitchStatement{
			tok: token{tokenType: SWITCH, value: "switch"},
			X: &IdentifierLiteral{
				Node: nil,
				tok:  token{tokenType: IDENTIFIER, value: "a", pos: 8, col: 8},
			},
		},
	},
	}
	assert.Equal(t, Parse("switch (a) {}", false), ep1)

	ep2 := &Program{body: []Node{
		&SwitchStatement{
			tok: token{tokenType: SWITCH, value: "switch"},
			X: &IdentifierLiteral{
				Node: nil,
				tok:  token{tokenType: IDENTIFIER, value: "a", pos: 8, col: 8},
			},
			Cases: []*CaseStatement{
				&CaseStatement{
					X:    &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "1", pos: 18, col: 18}},
					Body: []Node{},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("switch (a) { case 1:}", false), ep2)

	ep3 := &Program{body: []Node{
		&SwitchStatement{
			tok: token{tokenType: SWITCH, value: "switch"},
			X: &IdentifierLiteral{
				Node: nil,
				tok:  token{tokenType: IDENTIFIER, value: "a", pos: 8, col: 8},
			},
			Cases: []*CaseStatement{
				&CaseStatement{
					X:    &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "1", pos: 18, col: 18}},
					Body: []Node{},
				},
				&CaseStatement{
					X: &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "2", pos: 25, col: 25}},
					Body: []Node{
						&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", pos: 28, col: 28}}},
					},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("switch (a) { case 1:case 2: a;}", false), ep3)

	ep4 := &Program{body: []Node{
		&SwitchStatement{
			tok: token{tokenType: SWITCH, value: "switch"},
			X: &IdentifierLiteral{
				Node: nil,
				tok:  token{tokenType: IDENTIFIER, value: "a", pos: 8, col: 8},
			},
			Cases: []*CaseStatement{
				&CaseStatement{
					X:    &NumericLiteral{tok: token{tokenType: NUMERIC_LITERAL, value: "1", pos: 18, col: 18}},
					Body: []Node{},
				},
				&CaseStatement{
					IsDefault: true,
					Body: []Node{
						&ExpressionStatement{X: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "a", pos: 29, col: 29}}},
					},
				},
			},
		},
	},
	}
	assert.Equal(t, Parse("switch (a) { case 1:default: a;}", false), ep4)
}

func TestFunctionWithCall(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ExpressionStatement{
			X: &CallExpression{
				tok:       token{tokenType: LPAREN, value: "", pos: 15, col: 15},
				Arguments: []Node{},
				X: &FunctionExpression{
					tok:        token{tokenType: FUNCTION, value: "function"},
					Parameters: []*IdentifierLiteral{},
					Identifier: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "f", pos: 9, col: 9}},
					Body:       &BlockStatement{tok: token{tokenType: LBRACE, pos: 13, col: 13}, Body: []Node{}},
				},
			},
		},
	},
	}

	assert.Equal(t, Parse("function f() {}()", false), ep1)
}

func TestThrowStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&ThrowStatement{
			tok: token{tokenType: THROW, value: "throw", pos: 0, col: 0},
			X: &IdentifierLiteral{
				tok: token{tokenType: IDENTIFIER, value: "a", pos: 6, col: 6},
			},
		},
	},
	}

	assert.Equal(t, Parse("throw a;", false), ep1)
}

func TestTryStatement(t *testing.T) {
	ep1 := &Program{body: []Node{
		&TryStatement{
			tok: token{tokenType: TRY, value: "try", pos: 0, col: 0},
			Body: &BlockStatement{
				tok: token{tokenType: LBRACE, pos: 4, col: 4},
				Body: []Node{
					&ExpressionStatement{
						X: &IdentifierLiteral{
							tok: token{tokenType: IDENTIFIER, value: "a", pos: 6, col: 6},
						},
					},
				},
			},
			Catch: &CatchStatement{
				tok:        token{tokenType: CATCH, value: "catch", pos: 10, col: 10},
				Identifier: &IdentifierLiteral{tok: token{tokenType: IDENTIFIER, value: "e", pos: 17, col: 17}},
				Body: &BlockStatement{
					tok: token{tokenType: LBRACE, pos: 20, col: 20},
					Body: []Node{
						&ExpressionStatement{
							X: &IdentifierLiteral{
								tok: token{tokenType: IDENTIFIER, value: "b", pos: 22, col: 22},
							},
						},
					},
				},
			},
			Finally: &FinallyStatement{
				tok: token{tokenType: FINALLY, value: "finally", pos: 26, col: 26},
				Body: &BlockStatement{
					tok: token{tokenType: LBRACE, pos: 34, col: 34},
					Body: []Node{
						&ExpressionStatement{
							X: &IdentifierLiteral{
								tok: token{tokenType: IDENTIFIER, value: "c", pos: 36, col: 36},
							},
						},
					},
				},
			},
		},
	},
	}

	assert.Equal(t, Parse("try { a } catch (e) { b } finally { c }", false), ep1)
}
