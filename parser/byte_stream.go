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
