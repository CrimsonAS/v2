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

// A tokenStream consumes a byteStream to genereate tokens.
type tokenStream struct {
	stream         *byteStream
	current        *token
	hasStarted     bool
	ignoreComments bool
}

type TokenType int

//go:generate stringer -type=TokenType
const (
	EOF TokenType = iota
	COMMENT
	STRING_LITERAL
	NUMERIC_LITERAL
	IDENTIFIER

	// Assignment operators
	ASSIGNMENT              // =
	PLUS_EQ                 // +=
	MINUS_EQ                // -=
	MULTIPLY_EQ             // *=
	DIVIDE_EQ               // /=
	MODULUS_EQ              // %=
	LEFT_SHIFT_EQ           // <<=
	RIGHT_SHIFT_EQ          // >>=
	UNSIGNED_RIGHT_SHIFT_EQ // >>>=
	AND_EQ                  // &=
	XOR_EQ                  // ^=
	OR_EQ                   // |=

	PLUS                 // +
	INCREMENT            // ++
	MINUS                // -
	DECREMENT            // --
	MULTIPLY             // *
	DIVIDE               // /
	MODULUS              // %
	EQUALS               // ==
	STRICT_EQUALS        // ===
	BITWISE_AND          // &
	LOGICAL_AND          // &&
	BITWISE_OR           // |
	LOGICAL_OR           // ||
	LESS_THAN            // <
	LESS_EQ              // <=
	LEFT_SHIFT           // <<
	GREATER_THAN         // >
	GREATER_EQ           // >=
	RIGHT_SHIFT          // >>
	UNSIGNED_RIGHT_SHIFT // .>>>
	BITWISE_XOR          // ^
	INSTANCEOF           // instanceof
	IN                   // in
	NEW                  // new
	CONDITIONAL          // ?
	LOGICAL_NOT          // !
	NOT_EQUALS           // !=          // !
	STRICT_NOT_EQUALS    // !==
	BITWISE_NOT          // ~
	DELETE               // delete
	TYPEOF               // typeof
	VOID                 // void

	// Punctuation.
	DOT       // .
	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACKET  // [
	RBRACKET  // ]
	LBRACE    // {
	RBRACE    // }

	// literals
	THIS
	NULL
	TRUE
	FALSE

	// keywords
	VAR
	RETURN
	FUNCTION
	DO
	WHILE
	FOR
	GET
	SET

	// Flow control
	IF
	ELSE
	SWITCH
	CASE
	DEFAULT
	THROW
)

type token struct {
	tokenType TokenType
	value     string
	pos       int
	line      int
	col       int
}

func (this *tokenStream) eof() bool {
	this.peek() // ensure this.current is set
	return this.current.tokenType == EOF
}

// Returns the current token without advancing the stream
func (this *tokenStream) peek() token {
	// Allow an initial read to get an EOF token
	if this.current == nil && !this.hasStarted {
		this.hasStarted = true
		this.readNext()
	}

	return *this.current
}

// Returns the current token and advances the stream
func (this *tokenStream) next() token {
	cur := this.peek()
	this.readNext()
	return cur
}

//////// private below this point ////////

func isWhitespace(c byte) bool {
	if c == ' ' || c == '\t' || c == '\n' {
		return true
	}
	return false
}

func (this *tokenStream) consumeWhitespace() {
	for !this.stream.eof() && isWhitespace(this.stream.peek()) {
		this.stream.next()
	}
}

func (this *tokenStream) consumeSingleLineComment() *token {
	c := this.createToken(COMMENT, "")
	// these are off-by-one, as we read the first / already
	c.pos -= 1
	c.col -= 1
	this.stream.next()
	for !this.stream.eof() && this.stream.peek() != '\n' {
		c.value += string(this.stream.next())
	}
	return c
}

func (this *tokenStream) consumeMultiLineComment() *token {
	c := this.createToken(COMMENT, "")
	// these are off-by-one, as we read the first / already
	c.pos -= 1
	c.col -= 1
	this.stream.next()
	for !this.stream.eof() {
		if this.stream.peek() == '*' {
			c.value += string(this.stream.next())
			if !this.stream.eof() && this.stream.peek() == '/' {
				this.stream.next()                    // eat /
				c.value = c.value[0 : len(c.value)-1] // strip * from text
				return c
			}
		}
		c.value += string(this.stream.next())
	}
	return c
}

func (this *tokenStream) consumeComment() *token {
	switch this.stream.peek() {
	case '/':
		return this.consumeSingleLineComment()
	case '*':
		return this.consumeMultiLineComment()
	}

	panic("unreachable")
}

// ### string escaping, single quoted strings, etc (es5 7.8.4)
func (this *tokenStream) consumeString(char byte) *token {
	c := this.createToken(STRING_LITERAL, "")
	// these are off-by-one, as we read the " already
	c.pos -= 1
	c.col -= 1
	for !this.stream.eof() && this.stream.peek() != char {
		c.value += string(this.stream.next())
	}
	if !this.stream.eof() {
		this.stream.next() // consume ending "
	}
	return c
}

func isIdentifier(c byte, isFirstChar bool) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_') || (!isFirstChar && c >= '0' && c <= '9')
}

func classifyIdentifier(id string) (TokenType, bool) {
	switch id {
	case "if":
		return IF, false
	case "else":
		return ELSE, false
	case "switch":
		return SWITCH, false
	case "case":
		return CASE, false
	case "default":
		return DEFAULT, false
	case "throw":
		return THROW, false
	case "return":
		return RETURN, false
	case "this":
		return THIS, false
	case "null":
		return NULL, false
	case "true":
		return TRUE, false
	case "false":
		return FALSE, false
	case "instanceof":
		return INSTANCEOF, true
	case "in":
		return IN, true
	case "new":
		return NEW, true
	case "delete":
		return DELETE, true
	case "typeof":
		return TYPEOF, true
	case "void":
		return VOID, true
	case "function":
		return FUNCTION, false
	case "do":
		return DO, false
	case "while":
		return WHILE, false
	case "get":
		return GET, false
	case "set":
		return SET, false
	case "for":
		return FOR, false
	case "var":
		return VAR, false
	}

	return IDENTIFIER, false
}

func (this *tokenStream) consumeIdentifier(firstCharacter byte) *token {
	c := this.createToken(IDENTIFIER, string(firstCharacter))
	// these are off-by-one, as we read the first character already
	c.pos -= 1
	c.col -= 1
	for !this.stream.eof() && isIdentifier(this.stream.peek(), false) {
		c.value += string(this.stream.next())
	}

	tt, emptyValue := classifyIdentifier(c.value)
	c.tokenType = tt
	if emptyValue {
		c.value = ""
	}
	return c
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// ### hex literals (es5 7.3.8) and probably more
func (this *tokenStream) consumeNumber(firstDigit byte) *token {
	c := this.createToken(NUMERIC_LITERAL, string(firstDigit))
	// these are off-by-one, as we read the first digit already
	c.pos -= 1
	c.col -= 1
	for !this.stream.eof() && isDigit(this.stream.peek()) {
		c.value += string(this.stream.next())
	}
	if c.value == "0" && this.stream.peek() == 'x' {
		c.value += string(this.stream.next()) // consume 'x'
		for !this.stream.eof() && isHexDigit(this.stream.peek()) {
			c.value += string(this.stream.next())
		}
	}

	if c.value == "0x" {
		panic("Unexpected hexadecimal value (no digits after 'x')")
	}

	if !this.stream.eof() && this.stream.peek() == '.' {
		c.value += "."
		this.stream.next() // consume dot
		if this.stream.eof() {
			panic("malformed: got a number with no decimal part")
		}
		for !this.stream.eof() && isDigit(this.stream.peek()) {
			c.value += string(this.stream.next())
		}
	}
	return c
}

func isOperator(c byte) bool {
	switch c {
	case '+':
		fallthrough
	case '-':
		fallthrough
	case '*':
		fallthrough
	case '/':
		fallthrough
	case '%':
		fallthrough
	case '=':
		fallthrough
	case '&':
		fallthrough
	case '|':
		fallthrough
	case '<':
		fallthrough
	case '>':
		fallthrough
	case '^':
		fallthrough
	case '?':
		fallthrough
	case '!':
		fallthrough
	case '~':
		return true
	}

	return false
}

func (this *tokenStream) consumeOperator(firstDigit byte) *token {
	c := this.createToken(EOF, "")
	// these are off-by-one, as we read the first digit already
	c.pos -= 1
	c.col -= 1
	switch firstDigit {
	case '+':
		c.tokenType = PLUS
	case '-':
		c.tokenType = MINUS
	case '*':
		c.tokenType = MULTIPLY
	case '/':
		c.tokenType = DIVIDE
	case '%':
		c.tokenType = MODULUS
	case '=':
		c.tokenType = ASSIGNMENT
	case '&':
		c.tokenType = BITWISE_AND
	case '|':
		c.tokenType = BITWISE_OR
	case '<':
		c.tokenType = LESS_THAN
	case '>':
		c.tokenType = GREATER_THAN
	case '^':
		c.tokenType = BITWISE_XOR
	case '?':
		c.tokenType = CONDITIONAL
	case '!':
		c.tokenType = LOGICAL_NOT
	case '~':
		c.tokenType = BITWISE_NOT
	default:
		panic("unknown operator")
	}

	if c.tokenType == MODULUS {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = MODULUS_EQ
		}
	}
	if c.tokenType == MULTIPLY {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = MULTIPLY_EQ
		}
	}
	if c.tokenType == DIVIDE {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = DIVIDE_EQ
		}
	}
	if c.tokenType == PLUS {
		if !this.stream.eof() && this.stream.peek() == '+' {
			this.stream.next()
			c.tokenType = INCREMENT
		} else if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = PLUS_EQ
		}
	}
	if c.tokenType == MINUS {
		if !this.stream.eof() && this.stream.peek() == '-' {
			this.stream.next()
			c.tokenType = DECREMENT
		} else if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = MINUS_EQ
		}
	}
	if c.tokenType == BITWISE_AND {
		if !this.stream.eof() && this.stream.peek() == '&' {
			this.stream.next()
			c.tokenType = LOGICAL_AND
		} else if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = AND_EQ
		}
	}
	if c.tokenType == BITWISE_OR {
		if !this.stream.eof() && this.stream.peek() == '|' {
			this.stream.next()
			c.tokenType = LOGICAL_OR
		} else if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = OR_EQ
		}
	}
	if c.tokenType == BITWISE_XOR {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = XOR_EQ
		}
	}
	if c.tokenType == ASSIGNMENT {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = EQUALS
			if !this.stream.eof() && this.stream.peek() == '=' {
				this.stream.next()
				c.tokenType = STRICT_EQUALS
			}
		}
	}
	if c.tokenType == LOGICAL_NOT {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = NOT_EQUALS
			if !this.stream.eof() && this.stream.peek() == '=' {
				this.stream.next()
				c.tokenType = STRICT_NOT_EQUALS
			}
		}
	}
	if c.tokenType == LESS_THAN {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = LESS_EQ
		} else if !this.stream.eof() && this.stream.peek() == '<' {
			this.stream.next()
			c.tokenType = LEFT_SHIFT
		}
	}
	if c.tokenType == GREATER_THAN {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = GREATER_EQ
		} else if !this.stream.eof() && this.stream.peek() == '>' {
			this.stream.next()
			c.tokenType = RIGHT_SHIFT

			if !this.stream.eof() && this.stream.peek() == '>' {
				this.stream.next()
				c.tokenType = UNSIGNED_RIGHT_SHIFT
			}
		}

	}

	if c.tokenType == LEFT_SHIFT {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = LEFT_SHIFT_EQ
		}
	}

	if c.tokenType == RIGHT_SHIFT {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = RIGHT_SHIFT_EQ
		}
	}

	if c.tokenType == UNSIGNED_RIGHT_SHIFT {
		if !this.stream.eof() && this.stream.peek() == '=' {
			this.stream.next()
			c.tokenType = UNSIGNED_RIGHT_SHIFT_EQ
		}
	}

	return c
}

// ### don't duplicate all these cases
func isPunctuation(c byte) bool {
	switch c {
	case '.':
		fallthrough
	case ',':
		fallthrough
	case ':':
		fallthrough
	case ';':
		fallthrough
	case '(':
		fallthrough
	case ')':
		fallthrough
	case '[':
		fallthrough
	case ']':
		fallthrough
	case '{':
		fallthrough
	case '}':
		fallthrough
	case '?':
		return true
	}

	return false
}

func (this *tokenStream) consumePunctuation(firstDigit byte) *token {
	c := this.createToken(EOF, "")
	// these are off-by-one, as we read the first digit already
	c.pos -= 1
	c.col -= 1
	switch firstDigit {
	case '.':
		c.tokenType = DOT
	case ',':
		c.tokenType = COMMA
	case ':':
		c.tokenType = COLON
	case ';':
		c.tokenType = SEMICOLON
	case '(':
		c.tokenType = LPAREN
	case ')':
		c.tokenType = RPAREN
	case '[':
		c.tokenType = LBRACKET
	case ']':
		c.tokenType = RBRACKET
	case '{':
		c.tokenType = LBRACE
	case '}':
		c.tokenType = RBRACE
	default:
		panic("unknown punctuation")
	}

	return c
}

const tokenDebug = false

func (this *tokenStream) readNext() {
	if tokenDebug {
		defer func() {
			log.Printf("Read next token: %s %s %+v", this.current.tokenType, this.current.value, this.current)
		}()
	}
	this.consumeWhitespace()

	if this.stream.eof() {
		this.current = this.createToken(EOF, "")
		return
	}

	c := this.stream.next()
	var n byte
	if !this.stream.eof() {
		n = this.stream.peek()
	}
	//log.Printf("Looking at %+v %+v, %+v", c, n, false)

	if c == '/' && (n == '/' || n == '*') {
		this.current = this.consumeComment()
		if this.ignoreComments {
			this.readNext() // recurse until we hit EOF or something not a comment
		}
		return
	}

	if c == '"' || c == '\'' {
		this.current = this.consumeString(c)
		return
	}

	if isDigit(c) {
		this.current = this.consumeNumber(c)
		return
	}

	if isIdentifier(c, true) {
		this.current = this.consumeIdentifier(c)
		return
	}

	if isOperator(c) {
		this.current = this.consumeOperator(c)
		return
	}

	if isPunctuation(c) {
		this.current = this.consumePunctuation(c)
		return
	}

	if this.stream.eof() {
		this.current = this.createToken(EOF, "")
		return
	}

	panic("unknown token: " + string(c))
}

func (this *tokenStream) createToken(tokenType TokenType, value string) *token {
	return &token{
		pos:       this.stream.pos,
		line:      this.stream.line,
		col:       this.stream.col,
		tokenType: tokenType,
		value:     value,
	}
}

func regExpFlagFromChar(ch byte) RegExpFlag {
	switch ch {
	case 'g':
		return GlobalRegExp
	case 'i':
		return IgnoreCaseRegExp
	case 'm':
		return MultilineRegExp
	}
	return NoFlagsRegExp
}

// read a regular expression. this can't be done as a regular token, because it
// isn't following the token rules.. var a = /regex/ would be parsed as VAR EQ
// DIV IDENTIFIER DIV, but that's obviously unhelpful. so we rely on the parser
// to "pull" a regular expression out of us at the right time.
//
// pre-requisites: the '/' starting the regex *must* have already been consumed
// (i.e. by the parser)!
func (this *tokenStream) scanRegExp(eq bool) (tokenText string, patternFlags RegExpFlag) {
	if eq {
		tokenText = "="
	}

	for { // this will terminate when we find the following /
		currChar := this.stream.peek()
		switch currChar {
		case '\\':
			tokenText += string(currChar)
			currChar = this.stream.next()

			if this.stream.eof() || this.stream.peek() == '\n' {
				panic("Unterminated regular expression")
			}

			tokenText += string(currChar)
			currChar = this.stream.next()
			break

		case '[':
			tokenText += string(currChar)
			currChar = this.stream.next()

			for !this.stream.eof() && this.stream.peek() != '\n' {
				if currChar == ']' {
					break
				} else if currChar == '\\' {
					tokenText += string(currChar)
					currChar = this.stream.next()

					if this.stream.eof() || this.stream.peek() == '\n' {
						panic("Unterminated regular expression")
					}

					tokenText += string(currChar)
					currChar = this.stream.next()
				} else {
					tokenText += string(currChar)
					currChar = this.stream.next()
				}
			}

			if currChar != ']' {
				panic("Unterminated regular expression class")
			}

			tokenText += string(currChar)
			currChar = this.stream.next()
			break

		case '/': // terminating the regexp...
			patternFlags = NoFlagsRegExp
			currChar = this.stream.next() // skip /

			if !this.stream.eof() {
				currChar = this.stream.next()

				for isIdentifier(currChar, false) {
					flag := regExpFlagFromChar(currChar)
					if flag == NoFlagsRegExp || patternFlags&flag != 0 {
						panic(fmt.Sprintf("Bad regular expression flag %c", currChar))
					}
					patternFlags |= flag
					if this.stream.eof() {
						break
					} else {
						currChar = this.stream.next()
					}
				}
			}

			return tokenText, patternFlags

		default:
			if this.stream.eof() || this.stream.peek() == '\n' {
				panic("Unterminated regular expression")
			} else {
				tokenText += string(currChar)
				currChar = this.stream.next()
			}
		} // switch
	} // for

	panic("unreachable")
}
