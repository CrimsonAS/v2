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
	"strconv"
)

type NumericLiteral struct {
	Node
	tok token
}

func (this *NumericLiteral) token() token {
	return this.tok
}

func (this *NumericLiteral) String() string {
	return this.tok.value
}

func (this *NumericLiteral) Float64Value() float64 {
	if len(this.tok.value) >= 3 && this.tok.value[0] == '0' && this.tok.value[1] == 'x' {
		v := rune(0)
		for i := 2; i < len(this.tok.value); i++ {
			val := hex2dec(this.tok.value[i])
			v = v<<4 | val
		}
		return float64(v)
	}
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
	tok      token
	Elements []Node
}

func (this *ArrayLiteral) token() token {
	return this.tok
}

type ObjectPropertyType int

const (
	Normal ObjectPropertyType = iota
	Get
	Set
)

type ObjectPropertyLiteral struct {
	Key  Node
	Type ObjectPropertyType
	X    Node
}

type ObjectLiteral struct {
	Node
	tok        token
	Properties []ObjectPropertyLiteral
}

func (this *ObjectLiteral) token() token {
	return this.tok
}

type RegExpFlag int

const (
	NoFlagsRegExp RegExpFlag = iota
	GlobalRegExp
	IgnoreCaseRegExp
	MultilineRegExp
)

type RegExpLiteral struct {
	Node
	tok    token
	RegExp string
	Flags  RegExpFlag
}

func (this *RegExpLiteral) token() token {
	return this.tok
}
