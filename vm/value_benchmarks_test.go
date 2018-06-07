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

package vm

import (
	"testing"
)

func BenchmarkSimpleTypes(b *testing.B) {
	b.Run("undefined", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newUndefined()
		}
	})
	b.Run("null", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newNull()
		}
	})
	b.Run("bool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newBool(true)
		}
	})
	b.Run("number", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newNumber(0)
		}
	})
	b.Run("string_empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newString("")
		}
	})
	b.Run("string_1c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newString("a")
		}
	})
	b.Run("string_5c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newString("hello")
		}
	})
}

func BenchmarkObjectTypes(b *testing.B) {
	testfn := func(vm *vm, f value, args []value) value {
		return newUndefined()
	}
	b.Run("object", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newObject()
		}
	})
	b.Run("bool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newBooleanObject(true)
		}
	})
	b.Run("number", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newNumberObject(0)
		}
	})
	b.Run("string_empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newStringObject("")
		}
	})
	b.Run("string_1c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newStringObject("a")
		}
	})
	b.Run("string_5c", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newStringObject("hello")
		}
	})
	b.Run("function_call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newFunctionObject(testfn, nil)
		}
	})
	b.Run("function_ctor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newFunctionObject(nil, testfn)
		}
	})
	b.Run("function_both", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newFunctionObject(testfn, testfn)
		}
	})
}
