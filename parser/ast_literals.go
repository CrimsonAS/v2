package parser

import (
	"strconv"
)

type NumericLiteral struct {
	Node
	tok token
}

func (this *NumericLiteral) token() token {
	return this.tok
}

func (this *NumericLiteral) Float64Value() float64 {
	v, _ := strconv.ParseFloat(this.tok.value, 64)
	return v
}

type IdentifierLiteral struct {
	Node
	tok token
}

func (this *IdentifierLiteral) token() token {
	return this.tok
}

func (this *IdentifierLiteral) String() string {
	return this.tok.value
}

type StringLiteral struct {
	Node
	tok token
}

func (this *StringLiteral) token() token {
	return this.tok
}

func (this *StringLiteral) String() string {
	return this.tok.value
}

type FalseLiteral struct {
	Node
	tok token
}

func (this *FalseLiteral) token() token {
	return this.tok
}

type TrueLiteral struct {
	Node
	tok token
}

func (this *TrueLiteral) token() token {
	return this.tok
}

type ThisLiteral struct {
	Node
	tok token
}

func (this *ThisLiteral) token() token {
	return this.tok
}

type NullLiteral struct {
	Node
	tok token
}

func (this *NullLiteral) token() token {
	return this.tok
}

type ArrayLiteral struct {
	Node
	tok  token
	vals []Node
}

func (this *ArrayLiteral) token() token {
	return this.tok
}
