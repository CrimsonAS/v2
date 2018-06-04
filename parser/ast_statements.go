package parser

import (
	"fmt"
)

type IfStatement struct {
	Node
	ConditionExpr Node
	ThenStmt      Node
	ElseStmt      Node
	tok           token
}

func (this *IfStatement) Condition() Node {
	return this.ConditionExpr
}

func (this *IfStatement) Then() Node {
	return this.ThenStmt
}

func (this *IfStatement) token() token {
	return this.tok
}

type ReturnStatement struct {
	Node
	X   Node
	tok token
}

func (this *ReturnStatement) token() token {
	return this.tok
}

type BlockStatement struct {
	Node
	Body []Node
	tok  token
}

func (this *BlockStatement) String() string {
	return fmt.Sprintf("{ %s }", this.Body)
}

func (this *BlockStatement) token() token {
	return this.tok
}

type EmptyStatement struct {
	Node
	tok token
}

func (this *EmptyStatement) token() token {
	return this.tok
}

type VariableStatement struct {
	Vars         []*IdentifierLiteral
	Initializers []Node
	tok          token
}

func (this *VariableStatement) token() token {
	return this.tok
}

type DoWhileStatement struct {
	Vars []*IdentifierLiteral
	X    Node
	Body Node // ### Statement
	tok  token
}

func (this *DoWhileStatement) token() token {
	return this.tok
}

type WhileStatement struct {
	Vars []*IdentifierLiteral
	X    Node
	Body Node // ### Statement
	tok  token
}

func (this *WhileStatement) token() token {
	return this.tok
}

type ForStatement struct {
	Vars        []*IdentifierLiteral
	Initializer Node
	Test        Node
	Body        Node // ### Statement
	Update      Node
	tok         token
}

func (this *ForStatement) token() token {
	return this.tok
}
