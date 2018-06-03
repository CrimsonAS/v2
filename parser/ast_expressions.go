package parser

// ###
// these types are inspired by Go's AST types, but I think we have some things
// to fix.
// * use Node less often (instead: Expr, or whatever)
// * change left/right into X/Y (more abstract, and we can be consistent then)
// * Somehow expose token.Pos/token.Token like Go does?

type NewExpression struct {
	Node
	tok  token
	expr Node
}

func (this *NewExpression) token() token {
	return this.tok
}

type DotMemberExpression struct {
	Node
	tok   token
	left  Node
	right *IdentifierLiteral
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
	tok         token
	test        Node
	trueBranch  Node
	falseBranch Node
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
