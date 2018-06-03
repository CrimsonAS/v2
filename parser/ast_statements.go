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
	rval Node
	tok  token
}

func (this *ReturnStatement) ReturnValue() Node {
	return this.rval
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
