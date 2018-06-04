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
				n.vals = append(n.vals, nil)
			} else {
				exp = nil
			}
		default:
			exp = this.parseAssignmentExpression()
			n.vals = append(n.vals, exp)
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
	tok := this.stream.peek()
	if tok.tokenType == LBRACKET {
		this.expect(LBRACKET)
		right := this.parseExpression()
		this.expect(RBRACKET)
		return &BracketMemberExpression{tok: tok, left: left, right: right}
	} else if tok.tokenType == DOT {
		this.expect(DOT)
		member := &IdentifierLiteral{tok: this.expect(IDENTIFIER)}
		return &DotMemberExpression{tok: tok, left: left, right: member}
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

		return &CallExpression{tok: tok, X: left, Arguments: args}
	}

	return left
}

func (this *parser) parseNewExpression() Node {
	tok := this.expect(NEW)
	left := this.parseMemberExpression()
	return &NewExpression{tok: tok, expr: left}
}

func (this *parser) parseLeftHandSideExpression() Node {
	tok := this.stream.peek()
	var left Node
	if tok.tokenType == NEW {
		left = this.parseNewExpression()
	} else {
		left = this.parseMemberExpression()
	}
	return left
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
		return &ConditionalExpression{tok: tok, test: test, trueBranch: trueBranch, falseBranch: falseBranch}
	}

	return test
}

func (this *parser) parseAssignmentExpression() Node {
	left := this.parseConditionalExpression()
	tok := this.stream.peek()

	switch tok.tokenType {
	case ASSIGNMENT:
		this.expect(tok.tokenType)
		right := this.parseConditionalExpression()
		return &BinaryExpression{tok: tok, Left: left, Right: right}
	}
	return left
}

func (this *parser) parseExpression() Node {
	left := this.parseAssignmentExpression()

	tok := this.stream.peek()
	if tok.tokenType == COMMA {
		panic("sequence expression not implemented")
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
	case LPAREN:
		this.expect(LPAREN)
		expr := this.parseExpression()
		this.expect(RPAREN)
		return expr
	case EOF:
		return nil
	default:
		panic(fmt.Sprintf("unknown expression type %s", tok.tokenType))
	}
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
	return n
}

func (this *parser) parseBlockStatement() *BlockStatement {
	tok := this.expect(LBRACE)
	n := &BlockStatement{tok: tok}
	for this.stream.peek().tokenType != RBRACE {
		stmt := this.parseStatement()
		n.Body = append(n.Body, stmt)
	}
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

	if this.stream.peek().tokenType == VAR {
		panic("var declaration in for not yet supported")
	}

	var init Node
	if this.stream.peek().tokenType != SEMICOLON {
		init = this.parseExpression()
	}
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

	body := this.parseStatement()

	return &ForStatement{tok: tok, Initializer: init, Test: test, Update: update, Body: body}
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
	return &ExpressionStatement{X: this.parseExpression()}
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
		fallthrough
	case WHILE:
		fallthrough
	case FOR:
		return this.parseIterationStatement()
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

func Parse(code string) Node {
	np := parser{tokenStream{stream: &byteStream{code: code}}}
	ret := np.parseProgram()
	if parseDebug {
		log.Printf("%s", recursivelyPrint(ret))
	}
	return ret
}

// ### finish this and move it to parser.go
func recursivelyPrint(node Node) string {
	if node == nil {
		return "(nil)"
	}
	switch n := node.(type) {
	case *Program:
		p := "program:\n"
		for _, c := range n.Body() {
			p += recursivelyPrint(c)
		}
		return p
	case *IfStatement:
		if n.ElseStmt != nil {
			return fmt.Sprintf("If %s then %s else %s", recursivelyPrint(n.ConditionExpr), recursivelyPrint(n.ThenStmt), recursivelyPrint(n.ElseStmt))
		} else {
			return fmt.Sprintf("If %s then %s", recursivelyPrint(n.ConditionExpr), recursivelyPrint(n.ThenStmt))
		}
	case *ReturnStatement:
		return fmt.Sprintf("Return(%s)", recursivelyPrint(n.X))
	case *BlockStatement:
		p := "{:\n"
		for _, c := range n.Body {
			p += recursivelyPrint(c) + "\n"
		}
		p += "}\n"
		return p
	case *ExpressionStatement:
		return fmt.Sprintf("(unused) %s", recursivelyPrint(n.X))
	case *ArrayLiteral:
		p := "[:\n"
		for _, c := range n.vals {
			p += recursivelyPrint(c) + "\n"
		}
		p += "]\n"
		return p
	case *FunctionExpression:
		args := ""
		for _, arg := range n.Parameters {
			args += recursivelyPrint(arg) + ", "
		}
		if len(args) > 0 {
			args = args[:len(args)-2]
		}
		if n.Identifier != nil {
			return fmt.Sprintf("function %s(%s) %s", recursivelyPrint(n.Identifier), args, recursivelyPrint(n.Body))
		} else {
			return fmt.Sprintf("function(%s) %s", args, recursivelyPrint(n.Body))
		}
	case *NewExpression:
		return fmt.Sprintf("new %s", recursivelyPrint(n.expr))
	case *DotMemberExpression:
		return fmt.Sprintf("%s.%s", recursivelyPrint(n.left), recursivelyPrint(n.right))
	case *BracketMemberExpression:
		return fmt.Sprintf("%s[%s]", recursivelyPrint(n.left), recursivelyPrint(n.right))
	case *CallExpression:
		args := ""
		for _, arg := range n.Arguments {
			args += fmt.Sprintf("%s, ", recursivelyPrint(arg))
		}
		if len(args) > 0 {
			args = args[:len(args)-2]
		}
		return fmt.Sprintf("CALL %s(%s)", recursivelyPrint(n.X), args)
	case *UnaryExpression:
		if n.postfix {
			if n.token().tokenType == INCREMENT {
				return fmt.Sprintf("%s++", recursivelyPrint(n.X))
			} else if n.token().tokenType == DECREMENT {
				return fmt.Sprintf("%s--", recursivelyPrint(n.X))
			} else {
				panic(fmt.Sprintf("unknown postfix op %s", n.token().tokenType))
			}
		} else {
			if n.token().tokenType == INCREMENT {
				return fmt.Sprintf("++%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == DECREMENT {
				return fmt.Sprintf("--%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == DELETE {
				return fmt.Sprintf("delete %s", recursivelyPrint(n.X))
			} else if n.token().tokenType == TYPEOF {
				return fmt.Sprintf("typeof %s", recursivelyPrint(n.X))
			} else if n.token().tokenType == VOID {
				return fmt.Sprintf("void %s", recursivelyPrint(n.X))
			} else if n.token().tokenType == MINUS {
				return fmt.Sprintf("-%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == PLUS {
				return fmt.Sprintf("+%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_NOT {
				return fmt.Sprintf("~%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_XOR {
				return fmt.Sprintf("^%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_AND {
				return fmt.Sprintf("&%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == BITWISE_OR {
				return fmt.Sprintf("|%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == LOGICAL_NOT {
				return fmt.Sprintf("!%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == LOGICAL_AND {
				return fmt.Sprintf("&&%s", recursivelyPrint(n.X))
			} else if n.token().tokenType == LOGICAL_OR {
				return fmt.Sprintf("||%s", recursivelyPrint(n.X))
			} else {
				panic(fmt.Sprintf("unknown prefix op %s", n.token().tokenType))
			}
		}
		return ";"
	case *ConditionalExpression:
		return fmt.Sprintf("%s ? %s : %s", recursivelyPrint(n.test), recursivelyPrint(n.trueBranch), recursivelyPrint(n.falseBranch))
	case *BinaryExpression:
		switch n.token().tokenType {
		case MULTIPLY:
			return fmt.Sprintf("%s * %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case DIVIDE:
			return fmt.Sprintf("%s / %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case MODULUS:
			return fmt.Sprintf("%s %% %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case PLUS:
			return fmt.Sprintf("%s + %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case MINUS:
			return fmt.Sprintf("%s - %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case LEFT_SHIFT:
			return fmt.Sprintf("%s << %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case RIGHT_SHIFT:
			return fmt.Sprintf("%s >> %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case UNSIGNED_RIGHT_SHIFT:
			return fmt.Sprintf("%s >>> %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case EQUALS:
			return fmt.Sprintf("%s == %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case NOT_EQUALS:
			return fmt.Sprintf("%s != %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case STRICT_EQUALS:
			return fmt.Sprintf("%s === %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case STRICT_NOT_EQUALS:
			return fmt.Sprintf("%s !== %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case BITWISE_AND:
			return fmt.Sprintf("%s & %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case BITWISE_XOR:
			return fmt.Sprintf("%s ^ %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case BITWISE_OR:
			return fmt.Sprintf("%s | %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case LOGICAL_AND:
			return fmt.Sprintf("%s && %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case LOGICAL_OR:
			return fmt.Sprintf("%s || %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case ASSIGNMENT:
			return fmt.Sprintf("%s = %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case LESS_THAN:
			return fmt.Sprintf("%s < %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case GREATER_THAN:
			return fmt.Sprintf("%s > %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case LESS_EQ:
			return fmt.Sprintf("%s <= %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case GREATER_EQ:
			return fmt.Sprintf("%s >= %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		case INSTANCEOF:
			return fmt.Sprintf("%s instanceof %s", recursivelyPrint(n.Left), recursivelyPrint(n.Right))
		default:
			panic(fmt.Sprintf("unknown binary expression %s", node.token().tokenType))
		}

	case *EmptyStatement:
		return ";"
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
	case *DoWhileStatement:
		return fmt.Sprintf("do %s while %s", recursivelyPrint(n.Body), recursivelyPrint(n.X))
	case *WhileStatement:
		return fmt.Sprintf("while %s %s", recursivelyPrint(n.X), recursivelyPrint(n.Body))
	case *ForStatement:
		return fmt.Sprintf("for (%s; %s; %s) %s", recursivelyPrint(n.Initializer), recursivelyPrint(n.Test), recursivelyPrint(n.Update), recursivelyPrint(n.Body))
	case *VariableStatement:
		buf := "var "
		for idx, _ := range n.Vars {
			v := n.Vars[idx]
			i := n.Initializers[idx]
			buf += recursivelyPrint(v)
			if i != nil {
				buf += " = "
				buf += recursivelyPrint(i)
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
