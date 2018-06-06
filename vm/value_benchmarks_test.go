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
