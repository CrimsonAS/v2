package parser

import (
	"fmt"
)

// This type serves as a simple iterator over source code bytes. It maintains a
// position inside the code (both in terms of bytes, but also line/column
// information).
type byteStream struct {
	code string // source
	pos  int    // where are we (as an index into code)

	// where are we (in the user sense)
	// note that these are 0-indexed
	line int
	col  int
}

func (this *byteStream) next() byte {
	if this.eof() {
		panic(fmt.Sprintf("stream is already eof at byte %d position %d:%d", this.pos, this.line, this.col))
	}
	ch := this.code[this.pos]
	if ch == '\n' {
		this.line++
		this.col = 0
	} else {
		this.col++
	}
	this.pos += 1
	return ch
}

func (this *byteStream) peek() byte {
	return this.code[this.pos]
}

func (this *byteStream) eof() bool {
	if this.pos >= len(this.code) {
		return true
	}
	return false
}
