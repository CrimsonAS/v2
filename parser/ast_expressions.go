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

// ###
// these types are inspired by Go's AST types, but I think we have some things
// to fix.
// * use Node less often (instead: Expr, or whatever)
// * change left/right into X/Y (more abstract, and we can be consistent then)
// * Somehow expose token.Pos/token.Token like Go does?

type ExpressionStatement struct {
	Node
	tok token
	X   Node
}

func (this *ExpressionStatement) token() token {
	return this.tok
}

type NewExpression struct {
	Node
	tok token
	X   Node
}

func (this *NewExpression) token() token {
	return this.tok
}

type DotMemberExpression struct {
	Node
	tok  token
	X    Node
	Name *IdentifierLiteral
}

func (this *DotMemberExpression) token() token {
	return this.tok
}

type BracketMemberExpression struct {
	Node
	tok   token
	left  Node
	right Node
}

func (this *BracketMemberExpression) token() token {
	return this.tok
}

type UnaryExpression struct {
	Node
	tok     token
	postfix bool
	X       Node // ### Exp
}

func (this *UnaryExpression) token() token {
	return this.tok
}

func (this *UnaryExpression) Operator() TokenType {
	return TokenType(this.tok.tokenType)
}

func (this *UnaryExpression) IsPrefix() bool {
	return !this.postfix
}
func (this *UnaryExpression) IsPostfix() bool {
	return this.postfix
}

type BinaryExpression struct {
	Node
	tok   token
	Left  Node // ### Exp
	Right Node // ### Exp
}

func (this *BinaryExpression) Operator() TokenType {
	return TokenType(this.tok.tokenType)
}

func (this *BinaryExpression) token() token {
	return this.tok
}

type ConditionalExpression struct {
	Node
	tok  token
	X    Node
	Then Node
	Else Node
}

func (this *ConditionalExpression) token() token {
	return this.tok
}

type FunctionExpression struct {
	Node
	tok        token
	Identifier *IdentifierLiteral
	Parameters []*IdentifierLiteral
	Body       *BlockStatement
}

func (this *FunctionExpression) token() token {
	return this.tok
}

type CallExpression struct {
	Node
	tok       token
	X         Node
	Arguments []Node
}

func (this *CallExpression) token() token {
	return this.tok
}
