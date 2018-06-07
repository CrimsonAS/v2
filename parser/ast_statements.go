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
