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
	"log"
)

type parser struct {
	stream tokenStream
}

func (this *parser) parseArrayLiteral() *ArrayLiteral {
	tok := this.expect(LBRACKET)
	n := &ArrayLiteral{tok: tok}

	looping := true
	var exp Node
	for looping {
		tok = this.stream.peek()
		switch tok.tokenType {
		case RBRACKET:
			this.expect(RBRACKET)
			looping = false
		case COMMA:
			this.expect(COMMA)
			if exp == nil {
				n.Elements = append(n.Elements, nil)
			} else {
				exp = nil
			}
		default:
			exp = this.parseAssignmentExpression()
			n.Elements = append(n.Elements, exp)
		}

	}

	return n
}

func (this *parser) parseObjectProperty(currentObject *ObjectLiteral, propertyName Node, wantsGet bool, wantsSet bool) {
	if wantsGet {
		// get PropertyName() BlockStatement
		// ### or should this be a function literal?
		this.expect(LPAREN)
		this.expect(RPAREN)
		body := this.parseBlockStatement()
		wantsGet = false
		currentObject.Properties = append(currentObject.Properties, ObjectPropertyLiteral{Key: propertyName, Type: Get, X: body})
	} else if wantsSet {
		// set PropertyName(Identifier) BlockStatement
		// ### the identifier here gets dropped. we should make this a function literal, presumably.
		this.expect(LPAREN)
		this.expect(IDENTIFIER)
		this.expect(RPAREN)
		body := this.parseBlockStatement()
		wantsSet = false
		currentObject.Properties = append(currentObject.Properties, ObjectPropertyLiteral{Key: propertyName, Type: Set, X: body})
	} else {
		this.expect(COLON)
		x := this.parseAssignmentExpression()
		currentObject.Properties = append(currentObject.Properties, ObjectPropertyLiteral{Key: propertyName, Type: Normal, X: x})
	}

	if this.stream.peek().tokenType == COMMA {
		this.expect(COMMA)
	}
}

func (this *parser) parseObjectLiteral() *ObjectLiteral {
	tok := this.expect(LBRACE)
	n := &ObjectLiteral{tok: tok}

	wantsGet := false
	wantsSet := false
	parsingObject := true
	for parsingObject {
		switch this.stream.peek().tokenType {
		case GET:
			wantsGet = true
			this.expect(GET)
		case SET:
			wantsSet = true
			this.expect(SET)
		case IDENTIFIER:
			propertyName := &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
			this.parseObjectProperty(n, propertyName, wantsGet, wantsSet)
			wantsGet = false
			wantsSet = false
		case STRING_LITERAL:
			propertyName := &StringLiteral{tok: this.expect(STRING_LITERAL)}
			this.parseObjectProperty(n, propertyName, wantsGet, wantsSet)
			wantsGet = false
			wantsSet = false
		case NUMERIC_LITERAL:
			propertyName := &NumericLiteral{tok: this.expect(NUMERIC_LITERAL)}
			this.parseObjectProperty(n, propertyName, wantsGet, wantsSet)
			wantsGet = false
			wantsSet = false
		case RBRACE:
			this.expect(RBRACE)
			parsingObject = false
		}
	}

	return n
}

func (this *parser) parseMemberExpression() Node {
	if this.stream.peek().tokenType == FUNCTION {
		funcTok := this.expect(FUNCTION)
		var id *IdentifierLiteral
		switch this.stream.peek().tokenType {
		case IDENTIFIER:
			id = &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
		}

		this.expect(LPAREN)

		params := []*IdentifierLiteral{}
		for this.stream.peek().tokenType == IDENTIFIER {
			params = append(params, &IdentifierLiteral{tok: this.expect(IDENTIFIER)})
			if this.stream.peek().tokenType == COMMA {
				this.expect(COMMA)
			}
		}

		this.expect(RPAREN)

		body := this.parseBlockStatement()
		return &FunctionExpression{tok: funcTok, Identifier: id, Parameters: params, Body: body}
	}

	left := this.parsePrimaryExpression()
	return this.parseMemberOrCall(left)
}

func (this *parser) parseMemberOrCall(left Node) Node {
	tok := this.stream.peek()
	for tok.tokenType == LBRACKET || tok.tokenType == DOT || tok.tokenType == LPAREN {
		if tok.tokenType == LBRACKET {
			this.expect(LBRACKET)
			right := this.parseExpression()
			this.expect(RBRACKET)
			left = &BracketMemberExpression{tok: tok, X: left, Y: right}
		} else if tok.tokenType == DOT {
			this.expect(DOT)
			member := &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
			left = &DotMemberExpression{tok: tok, X: left, Name: member}
		} else if tok.tokenType == LPAREN {
			this.expect(LPAREN)
			args := []Node{}
			for this.stream.peek().tokenType != RPAREN {
				arg := this.parseAssignmentExpression()
				args = append(args, arg)
				if this.stream.peek().tokenType == COMMA {
					this.expect(COMMA)
				}
			}
			this.expect(RPAREN)

			left = &CallExpression{tok: tok, X: left, Arguments: args}
		}
		tok = this.stream.peek()
	}
	return left
}

func (this *parser) parseNewExpression() Node {
	tok := this.expect(NEW)
	left := this.parseMemberExpression()
	return &NewExpression{tok: tok, X: left}
}

func (this *parser) parseLeftHandSideExpression() Node {
	tok := this.stream.peek()
	var left Node
	if tok.tokenType == NEW {
		left = this.parseNewExpression()
	} else {
		left = this.parseMemberExpression()
	}

	return this.parseMemberOrCall(left)
}

func (this *parser) parsePostfixExpression() Node {
	left := this.parseLeftHandSideExpression()
	tok := this.stream.peek()
	switch tok.tokenType {
	case INCREMENT:
		this.expect(INCREMENT)
		return &UnaryExpression{tok: tok, postfix: true, X: left}
	case DECREMENT:
		this.expect(DECREMENT)
		return &UnaryExpression{tok: tok, postfix: true, X: left}
	}
	return left
}

func (this *parser) parseUnaryExpression() Node {
	tok := this.stream.peek()

	tt := tok.tokenType

	switch tt {
	case PLUS:
		fallthrough
	case MINUS:
		fallthrough
	case BITWISE_NOT:
		fallthrough
	case LOGICAL_NOT:
		fallthrough
	case DELETE:
		fallthrough
	case TYPEOF:
		fallthrough
	case VOID:
		this.expect(tt)
		return &UnaryExpression{tok: tok, postfix: false, X: this.parseUnaryExpression()}
	case INCREMENT:
		fallthrough
	case DECREMENT:
		this.expect(tt)
		return &UnaryExpression{tok: tok, postfix: false, X: this.parseUnaryExpression()}
	}

	return this.parsePostfixExpression()
}

func (this *parser) parseMultiplicativeExpression() Node {
	left := this.parseUnaryExpression()
	tok := this.stream.peek()

	for tok.tokenType == MULTIPLY || tok.tokenType == DIVIDE || tok.tokenType == MODULUS {
		this.expect(tok.tokenType)
		right := this.parseUnaryExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseAdditiveExpression() Node {
	left := this.parseMultiplicativeExpression()
	tok := this.stream.peek()

	for tok.tokenType == PLUS || tok.tokenType == MINUS {
		this.expect(tok.tokenType)
		right := this.parseMultiplicativeExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseShiftExpression() Node {
	left := this.parseAdditiveExpression()
	tok := this.stream.peek()

	for tok.tokenType == RIGHT_SHIFT || tok.tokenType == LEFT_SHIFT || tok.tokenType == UNSIGNED_RIGHT_SHIFT {
		this.expect(tok.tokenType)
		right := this.parseAdditiveExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseRelationalExpression() Node {
	left := this.parseShiftExpression()
	tok := this.stream.peek()

	switch tok.tokenType {
	case LESS_THAN:
		fallthrough
	case GREATER_THAN:
		fallthrough
	case LESS_EQ:
		fallthrough
	case GREATER_EQ:
		fallthrough
	case IN:
		fallthrough
	case INSTANCEOF:
		this.expect(tok.tokenType)
		right := this.parseShiftExpression()
		return &BinaryExpression{tok: tok, Left: left, Right: right}
	}

	return left
}

func (this *parser) parseEqualityExpression() Node {
	left := this.parseRelationalExpression()
	tok := this.stream.peek()

	for tok.tokenType == EQUALS || tok.tokenType == STRICT_EQUALS || tok.tokenType == NOT_EQUALS || tok.tokenType == STRICT_NOT_EQUALS {
		this.expect(tok.tokenType)
		right := this.parseRelationalExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseBitwiseAndExpression() Node {
	left := this.parseEqualityExpression()
	tok := this.stream.peek()

	for tok.tokenType == BITWISE_AND {
		this.expect(tok.tokenType)
		right := this.parseEqualityExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseBitwiseXorExpression() Node {
	left := this.parseBitwiseAndExpression()
	tok := this.stream.peek()

	for tok.tokenType == BITWISE_XOR {
		this.expect(tok.tokenType)
		right := this.parseBitwiseAndExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseBitwiseOrExpression() Node {
	left := this.parseBitwiseXorExpression()
	tok := this.stream.peek()

	for tok.tokenType == BITWISE_OR {
		this.expect(tok.tokenType)
		right := this.parseBitwiseXorExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseLogicalAndExpression() Node {
	left := this.parseBitwiseOrExpression()
	tok := this.stream.peek()

	for tok.tokenType == LOGICAL_AND {
		this.expect(tok.tokenType)
		right := this.parseBitwiseOrExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseLogicalOrExpression() Node {
	left := this.parseLogicalAndExpression()
	tok := this.stream.peek()

	for tok.tokenType == LOGICAL_OR {
		this.expect(tok.tokenType)
		right := this.parseLogicalAndExpression()
		left = &BinaryExpression{tok: tok, Left: left, Right: right}
		tok = this.stream.peek()
	}

	return left
}

func (this *parser) parseConditionalExpression() Node {
	test := this.parseLogicalOrExpression()
	tok := this.stream.peek()

	switch tok.tokenType {
	case CONDITIONAL:
		this.expect(tok.tokenType)
		trueBranch := this.parseAssignmentExpression()
		this.expect(COLON)
		falseBranch := this.parseAssignmentExpression()
		return &ConditionalExpression{tok: tok, X: test, Then: trueBranch, Else: falseBranch}
	}

	return test
}

func (this *parser) parseAssignmentExpression() Node {
	left := this.parseConditionalExpression()
	tok := this.stream.peek()

	switch tok.tokenType {
	case PLUS_EQ:
		fallthrough
	case MINUS_EQ:
		fallthrough
	case MULTIPLY_EQ:
		fallthrough
	case DIVIDE_EQ:
		fallthrough
	case MODULUS_EQ:
		fallthrough
	case LEFT_SHIFT_EQ:
		fallthrough
	case RIGHT_SHIFT_EQ:
		fallthrough
	case UNSIGNED_RIGHT_SHIFT_EQ:
		fallthrough
	case AND_EQ:
		fallthrough
	case XOR_EQ:
		fallthrough
	case OR_EQ:
		fallthrough
	case ASSIGNMENT:
		this.expect(tok.tokenType)
		right := this.parseAssignmentExpression()
		return &AssignmentExpression{tok: tok, Left: left, Right: right}
	}
	return left
}

func (this *parser) parseExpression() Node {
	left := this.parseAssignmentExpression()

	tok := this.stream.peek()
	if tok.tokenType == COMMA {
		seq := []Node{left}
		for this.stream.peek().tokenType == COMMA {
			this.expect(COMMA)
			seq = append(seq, this.parseAssignmentExpression())
		}

		return &SequenceExpression{tok: tok, Seq: seq}
	}
	return left
}

func (this *parser) parsePrimaryExpression() Node {
	tok := this.stream.peek()
	switch tok.tokenType {
	case NUMERIC_LITERAL:
		return &NumericLiteral{tok: this.expect(NUMERIC_LITERAL)}
	case STRING_LITERAL:
		return &StringLiteral{tok: this.expect(STRING_LITERAL)}
	case THIS:
		return &ThisLiteral{tok: this.expect(THIS)}
	case IDENTIFIER:
		return &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
	case TRUE:
		return &TrueLiteral{tok: this.expect(TRUE)}
	case FALSE:
		return &FalseLiteral{tok: this.expect(FALSE)}
	case NULL:
		return &NullLiteral{tok: this.expect(NULL)}
	case LBRACKET:
		return this.parseArrayLiteral()
	case LBRACE:
		return this.parseObjectLiteral()
	case LPAREN:
		this.expect(LPAREN)
		expr := this.parseExpression()
		this.expect(RPAREN)
		return expr
	case DIVIDE:
		return this.parseRegExpLiteral(false)
	case DIVIDE_EQ:
		return this.parseRegExpLiteral(true)
	case EOF:
		return nil
	default:
		panic(fmt.Sprintf("unknown expression type %s %s", tok.tokenType, tok.value))
	}
}

func (this *parser) parseRegExpLiteral(eq bool) Node {
	tok := this.stream.peek()
	re, flags := this.stream.scanRegExp(eq)
	this.expect(tok.tokenType) // we ate it
	return &RegExpLiteral{tok: tok, RegExp: re, Flags: flags}
}

func (this *parser) parseIfStatement() *IfStatement {
	tok := this.expect(IF)
	n := &IfStatement{tok: tok}
	this.expect(LPAREN)
	n.ConditionExpr = this.parseExpression()
	this.expect(RPAREN)
	n.ThenStmt = this.parseStatement()
	tok = this.stream.peek()
	if tok.tokenType == ELSE {
		this.expect(ELSE)
		n.ElseStmt = this.parseStatement()
	}
	return n
}

func (this *parser) parseReturnStatement() *ReturnStatement {
	tok := this.expect(RETURN)
	n := &ReturnStatement{tok: tok}
	if this.stream.peek().tokenType == SEMICOLON {
		this.expect(SEMICOLON)
		return n
	}
	n.X = this.parseExpression()
	if this.stream.peek().tokenType == SEMICOLON {
		this.expect(SEMICOLON)
	}
	return n
}

func (this *parser) parseBlockStatementBody() []Node {
	ret := []Node{}

	// CASE/DEFAULT checks are because this is also used to read a switch case body.
	for this.stream.peek().tokenType != RBRACE && this.stream.peek().tokenType != CASE && this.stream.peek().tokenType != DEFAULT {
		ret = append(ret, this.parseStatement())
	}
	return ret
}

func (this *parser) parseBlockStatement() *BlockStatement {
	tok := this.expect(LBRACE)
	n := &BlockStatement{tok: tok, Body: this.parseBlockStatementBody()}
	this.expect(RBRACE)
	return n
}

func (this *parser) parseVariableStatement() Node {
	tok := this.expect(VAR)
	n := &VariableStatement{tok: tok}

	for this.stream.peek().tokenType == IDENTIFIER {
		id := &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
		var initializer Node = nil
		if this.stream.peek().tokenType == ASSIGNMENT {
			this.expect(ASSIGNMENT)
			initializer = this.parseAssignmentExpression()
		}

		n.Vars = append(n.Vars, id)
		n.Initializers = append(n.Initializers, initializer)

		if this.stream.peek().tokenType == COMMA {
			this.expect(COMMA)
		}
	}

	return n
}

func (this *parser) parseWhileStatement() Node {
	tok := this.expect(WHILE)
	this.expect(LPAREN)
	expr := this.parseExpression()
	this.expect(RPAREN)
	body := this.parseStatement()

	return &WhileStatement{tok: tok, X: expr, Body: body}
}

func (this *parser) parseDoWhileStatement() Node {
	tok := this.expect(DO)
	body := this.parseStatement()
	this.expect(WHILE)
	this.expect(LPAREN)
	expr := this.parseExpression()
	this.expect(RPAREN)

	return &DoWhileStatement{tok: tok, X: expr, Body: body}
}

func (this *parser) parseForStatement() Node {
	tok := this.expect(FOR)
	this.expect(LPAREN)

	var init Node
	if this.stream.peek().tokenType != SEMICOLON {
		if this.stream.peek().tokenType == VAR {
			init = this.parseVariableStatement()
		} else {
			init = this.parseExpression()
		}
	}

	if this.stream.peek().tokenType == IN {
		this.expect(IN)
		Y := this.parseExpression()
		this.expect(RPAREN)
		return &ForInStatement{tok: tok, X: init, Y: Y, Body: this.parseStatement()}
	} else {
		this.expect(SEMICOLON)
		var test Node
		if this.stream.peek().tokenType != SEMICOLON {
			test = this.parseExpression()
		}
		this.expect(SEMICOLON)
		var update Node
		if this.stream.peek().tokenType != RPAREN {
			update = this.parseExpression()
		}
		this.expect(RPAREN)
		return &ForStatement{tok: tok, Initializer: init, Test: test, Update: update, Body: this.parseStatement()}
	}
}

func (this *parser) parseIterationStatement() Node {
	tok := this.stream.peek()
	switch tok.tokenType {
	case DO:
		return this.parseDoWhileStatement()
	case WHILE:
		return this.parseWhileStatement()
	case FOR:
		return this.parseForStatement()
	}

	panic("unreachable")
}

func (this *parser) parseExpressionStatement() Node {
	r := &ExpressionStatement{X: this.parseExpression()}
	if this.stream.peek().tokenType == SEMICOLON {
		this.expect(SEMICOLON)
	}
	return r
}

func (this *parser) parseSwitchStatement() Node {
	r := &SwitchStatement{tok: this.expect(SWITCH)}
	this.expect(LPAREN)
	r.X = this.parseExpression()
	this.expect(RPAREN)
	this.expect(LBRACE)

	hasDefault := false
	for this.stream.peek().tokenType != RBRACE {
		isDefault := false
		var expr Node
		switch this.stream.peek().tokenType {
		case CASE:
			this.expect(CASE)
			expr = this.parseExpression()
		case DEFAULT:
			if hasDefault {
				panic("already got a default case in switch")
			}
			this.expect(DEFAULT)
			isDefault = true
		default:
			panic(fmt.Sprintf("unexpected token in switch %s %s", this.stream.peek().tokenType, this.stream.peek().value))
		}

		this.expect(COLON)

		// this will stop at CASE, DEFAULT or }
		body := this.parseBlockStatementBody()
		r.Cases = append(r.Cases, &CaseStatement{X: expr, Body: body, IsDefault: isDefault})
		if isDefault {
			hasDefault = true
		}
	}

	this.expect(RBRACE)
	return r
}

func (this *parser) parseThrowStatement() Node {
	tok := this.expect(THROW)
	x := this.parseExpression()
	if this.stream.peek().tokenType == SEMICOLON {
		this.expect(SEMICOLON)
	}
	return &ThrowStatement{tok: tok, X: x}
}

func (this *parser) parseTryStatement() Node {
	tb := &TryStatement{tok: this.expect(TRY), Body: this.parseBlockStatement()}

	switch this.stream.peek().tokenType {
	case CATCH:
		cb := &CatchStatement{tok: this.expect(CATCH)}
		this.expect(LPAREN)
		cb.Identifier = &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
		this.expect(RPAREN)
		cb.Body = this.parseBlockStatement()
		tb.Catch = cb
	case FINALLY:
		fb := &FinallyStatement{tok: this.expect(FINALLY), Body: this.parseBlockStatement()}
		tb.Finally = fb
	default:
		panic("expected catch or finally")
	}

	switch this.stream.peek().tokenType {
	case CATCH:
		panic("catch expected before finally")
	case FINALLY:
		if tb.Finally != nil {
			panic("only one finally block expected")
		}
		fb := &FinallyStatement{tok: this.expect(FINALLY), Body: this.parseBlockStatement()}
		tb.Finally = fb
	}

	return tb
}

func (this *parser) parseStatement() Node {
	tok := this.stream.peek()
	switch tok.tokenType {
	case VAR:
		return this.parseVariableStatement()
	case IF:
		return this.parseIfStatement()
	case RETURN:
		return this.parseReturnStatement()
	case LBRACE:
		return this.parseBlockStatement()
	case DO:
		return this.parseIterationStatement()
	case WHILE:
		return this.parseIterationStatement()
	case FOR:
		return this.parseIterationStatement()
	case TRY:
		return this.parseTryStatement()
	case THROW:
		return this.parseThrowStatement()
	case SWITCH:
		return this.parseSwitchStatement()
	case SEMICOLON:
		return &EmptyStatement{tok: this.expect(SEMICOLON)}
	}

	return this.parseExpressionStatement()
}

func (this *parser) expect(ttype TokenType) token {
	tok := this.stream.next()
	if tok.tokenType != ttype {
		panic(fmt.Sprintf("expected: %s, got %s", ttype, tok.tokenType))
	}
	return tok
}

func (this *parser) parseProgram() *Program {
	p := &Program{}

	for !this.stream.eof() {
		stmt := this.parseStatement()
		p.body = append(p.body, stmt)
	}

	return p
}

const parseDebug = false

func Parse(code string, ignoreComments bool) Node {
	np := parser{tokenStream{stream: &byteStream{code: code}, ignoreComments: ignoreComments}}
	ret := np.parseProgram()
	if parseDebug {
		log.Printf("%s", RecursivelyPrint(ret))
	}
	return ret
}

func RecursivelyPrint(node Node) string {
	if node == nil {
		return "(nil)"
	}
	switch n := node.(type) {
	case *Program:
		p := "program:\n"
		for _, c := range n.Body() {
			p += RecursivelyPrint(c)
		}
		return p
	case *IfStatement:
		if n.ElseStmt != nil {
			return fmt.Sprintf("If %s then %s else %s", RecursivelyPrint(n.ConditionExpr), RecursivelyPrint(n.ThenStmt), RecursivelyPrint(n.ElseStmt))
		} else {
			return fmt.Sprintf("If %s then %s", RecursivelyPrint(n.ConditionExpr), RecursivelyPrint(n.ThenStmt))
		}
	case *ReturnStatement:
		return fmt.Sprintf("Return(%s)", RecursivelyPrint(n.X))
	case *BlockStatement:
		p := "{:\n"
		for _, c := range n.Body {
			p += RecursivelyPrint(c) + "\n"
		}
		p += "}\n"
		return p
	case *ExpressionStatement:
		return fmt.Sprintf("(unused) %s", RecursivelyPrint(n.X))
	case *ArrayLiteral:
		p := "[:\n"
		for _, c := range n.Elements {
			p += RecursivelyPrint(c) + "\n"
		}
		p += "]\n"
		return p
	case *FunctionExpression:
		args := ""
		for _, arg := range n.Parameters {
			args += RecursivelyPrint(arg) + ", "
		}
		if len(args) > 0 {
			args = args[:len(args)-2]
		}
		if n.Identifier != nil {
			return fmt.Sprintf("function %s(%s) %s", RecursivelyPrint(n.Identifier), args, RecursivelyPrint(n.Body))
		} else {
			return fmt.Sprintf("function(%s) %s", args, RecursivelyPrint(n.Body))
		}
	case *NewExpression:
		return fmt.Sprintf("new %s", RecursivelyPrint(n.X))
	case *DotMemberExpression:
		return fmt.Sprintf("%s.%s", RecursivelyPrint(n.X), RecursivelyPrint(n.Name))
	case *BracketMemberExpression:
		return fmt.Sprintf("%s[%s]", RecursivelyPrint(n.X), RecursivelyPrint(n.Y))
	case *CallExpression:
		args := ""
		for _, arg := range n.Arguments {
			args += fmt.Sprintf("%s, ", RecursivelyPrint(arg))
		}
		if len(args) > 0 {
			args = args[:len(args)-2]
		}
		return fmt.Sprintf("%s(%s)", RecursivelyPrint(n.X), args)
	case *UnaryExpression:
		if n.postfix {
			if n.token().tokenType == INCREMENT {
				return fmt.Sprintf("%s++", RecursivelyPrint(n.X))
			} else if n.token().tokenType == DECREMENT {
				return fmt.Sprintf("%s--", RecursivelyPrint(n.X))
			} else {
				panic(fmt.Sprintf("unknown postfix op %s", n.token().tokenType))
			}
		} else {
			if n.token().tokenType == INCREMENT {
				return fmt.Sprintf("++%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == DECREMENT {
				return fmt.Sprintf("--%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == DELETE {
				return fmt.Sprintf("delete %s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == TYPEOF {
				return fmt.Sprintf("typeof %s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == VOID {
				return fmt.Sprintf("void %s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == MINUS {
				return fmt.Sprintf("-%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == PLUS {
				return fmt.Sprintf("+%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_NOT {
				return fmt.Sprintf("~%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_XOR {
				return fmt.Sprintf("^%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_AND {
				return fmt.Sprintf("&%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_OR {
				return fmt.Sprintf("|%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == LOGICAL_NOT {
				return fmt.Sprintf("!%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == LOGICAL_AND {
				return fmt.Sprintf("&&%s", RecursivelyPrint(n.X))
			} else if n.token().tokenType == LOGICAL_OR {
				return fmt.Sprintf("||%s", RecursivelyPrint(n.X))
			} else {
				panic(fmt.Sprintf("unknown prefix op %s", n.token().tokenType))
			}
		}
		return ";"
	case *ConditionalExpression:
		return fmt.Sprintf("%s ? %s : %s", RecursivelyPrint(n.X), RecursivelyPrint(n.Then), RecursivelyPrint(n.Else))
	case *SequenceExpression:
		buf := ""
		for _, a := range n.Seq {
			buf += fmt.Sprintf("%s, ", RecursivelyPrint(a))
		}
		if len(buf) > 0 {
			buf = buf[:len(buf)-2]
		}
		return buf
	case *AssignmentExpression:
		switch n.token().tokenType {
		case ASSIGNMENT:
			return fmt.Sprintf("%s = %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case PLUS_EQ:
			return fmt.Sprintf("%s += %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case MINUS_EQ:
			return fmt.Sprintf("%s -= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case MULTIPLY_EQ:
			return fmt.Sprintf("%s *= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case DIVIDE_EQ:
			return fmt.Sprintf("%s /= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case MODULUS_EQ:
			return fmt.Sprintf("%s %%= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case LEFT_SHIFT_EQ:
			return fmt.Sprintf("%s <<= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case RIGHT_SHIFT_EQ:
			return fmt.Sprintf("%s >>= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case UNSIGNED_RIGHT_SHIFT_EQ:
			return fmt.Sprintf("%s >>>= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case AND_EQ:
			return fmt.Sprintf("%s &= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case XOR_EQ:
			return fmt.Sprintf("%s ^= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case OR_EQ:
			return fmt.Sprintf("%s |= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		default:
			panic(fmt.Sprintf("unknown assignment expression %s", node.token().tokenType))
		}
	case *BinaryExpression:
		switch n.token().tokenType {
		case MULTIPLY:
			return fmt.Sprintf("%s * %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case DIVIDE:
			return fmt.Sprintf("%s / %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case MODULUS:
			return fmt.Sprintf("%s %% %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case PLUS:
			return fmt.Sprintf("%s + %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case MINUS:
			return fmt.Sprintf("%s - %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case LEFT_SHIFT:
			return fmt.Sprintf("%s << %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case RIGHT_SHIFT:
			return fmt.Sprintf("%s >> %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case UNSIGNED_RIGHT_SHIFT:
			return fmt.Sprintf("%s >>> %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case EQUALS:
			return fmt.Sprintf("%s == %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case NOT_EQUALS:
			return fmt.Sprintf("%s != %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case STRICT_EQUALS:
			return fmt.Sprintf("%s === %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case STRICT_NOT_EQUALS:
			return fmt.Sprintf("%s !== %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case BITWISE_AND:
			return fmt.Sprintf("%s & %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case BITWISE_XOR:
			return fmt.Sprintf("%s ^ %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case BITWISE_OR:
			return fmt.Sprintf("%s | %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case LOGICAL_AND:
			return fmt.Sprintf("%s && %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case LOGICAL_OR:
			return fmt.Sprintf("%s || %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case LESS_THAN:
			return fmt.Sprintf("%s < %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case GREATER_THAN:
			return fmt.Sprintf("%s > %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case LESS_EQ:
			return fmt.Sprintf("%s <= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case GREATER_EQ:
			return fmt.Sprintf("%s >= %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case INSTANCEOF:
			return fmt.Sprintf("%s instanceof %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		case IN:
			return fmt.Sprintf("%s in %s", RecursivelyPrint(n.Left), RecursivelyPrint(n.Right))
		default:
			panic(fmt.Sprintf("unknown binary expression %s", node.token().tokenType))
		}
	case *EmptyStatement:
		return ";"
	case *ObjectLiteral:
		buf := ""
		for _, prop := range n.Properties {
			if prop.Type == Get {
				buf += "get "
			} else if prop.Type == Set {
				buf += "set "
			}

			buf += RecursivelyPrint(prop.Key)
			buf += ": "
			buf += RecursivelyPrint(prop.X)
			buf += ", "
		}
		if len(buf) > 0 {
			buf = buf[:len(buf)-2]
		}
		return fmt.Sprintf("{ %s }", buf)
	case *TrueLiteral:
		return "true"
	case *FalseLiteral:
		return "false"
	case *NumericLiteral:
		return fmt.Sprintf("%s", n.tok.value)
	case *IdentifierLiteral:
		return fmt.Sprintf("%s", n.tok.value)
	case *StringLiteral:
		return fmt.Sprintf("\"%s\"", n.tok.value)
	case *ThisLiteral:
		return "this"
	case *NullLiteral:
		return "null"
	case *RegExpLiteral:
		return fmt.Sprintf("/%s/%s", n.RegExp, n.Flags)
	case *DoWhileStatement:
		return fmt.Sprintf("do %s while %s", RecursivelyPrint(n.Body), RecursivelyPrint(n.X))
	case *WhileStatement:
		return fmt.Sprintf("while %s %s", RecursivelyPrint(n.X), RecursivelyPrint(n.Body))
	case *ForStatement:
		return fmt.Sprintf("for (%s; %s; %s) %s", RecursivelyPrint(n.Initializer), RecursivelyPrint(n.Test), RecursivelyPrint(n.Update), RecursivelyPrint(n.Body))
	case *ForInStatement:
		return fmt.Sprintf("for %s in %s) %s", RecursivelyPrint(n.X), RecursivelyPrint(n.Y), RecursivelyPrint(n.Body))
	case *TryStatement:
		b := fmt.Sprintf("try %s", RecursivelyPrint(n.Body))
		if n.Catch != nil {
			b += "\n"
			b += fmt.Sprintf("catch (%s) %s", n.Catch.Identifier, n.Catch.Body)
		}
		if n.Finally != nil {
			b += "\n"
			b += fmt.Sprintf("finally %s", n.Finally.Body)
		}
		return b
	case *SwitchStatement:
		b := fmt.Sprintf("switch (%s) {\n", n.X)
		for _, cs := range n.Cases {
			if cs.IsDefault {
				b += "default:\n"
			} else {
				b += fmt.Sprintf("case %s:\n", cs.X)
			}

			for _, s := range cs.Body {
				b += fmt.Sprintf("%s\n", RecursivelyPrint(s))
			}
		}
		b += "}\n"
		return b
	case *ThrowStatement:
		return fmt.Sprintf("throw %s\n", n.X)
	case *VariableStatement:
		buf := "var "
		for idx, _ := range n.Vars {
			v := n.Vars[idx]
			i := n.Initializers[idx]
			buf += RecursivelyPrint(v)
			if i != nil {
				buf += " = "
				buf += RecursivelyPrint(i)
			}

			buf += ", "
		}
		if len(buf) > 0 {
			buf = buf[:len(buf)-2]
		}
		return buf
	default:
		panic(fmt.Sprintf("unknown node in print %T", node))
	}
}
