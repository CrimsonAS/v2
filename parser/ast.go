package parser

type Node interface {
	token() token
}

type Program struct {
	Node
	tok  token
	body []Node
}

func (this *Program) Body() []Node {
	return this.body
}

func (this *Program) token() token {
	return this.tok
}
