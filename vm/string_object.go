package vm

import (
	"fmt"
	"log"
	"math"
	"strings"
)

func newStringObject(s string) value {
	v := newString(s)
	v.vtype = OBJECT
	v.odata = &objectData{STRING_OBJECT, value{}, nil, nil, nil}
	v.odata.prototype = stringProto
	return v
}

var stringProto value

func defineStringCtor(vm *vm) value {
	stringProto = newObject()
	stringProto.defineDefaultProperty(vm, "toString", newFunctionObject(string_prototype_toString, nil), 0)
	stringProto.defineDefaultProperty(vm, "valueOf", newFunctionObject(string_prototype_valueOf, nil), 0)
	stringProto.defineDefaultProperty(vm, "charAt", newFunctionObject(string_prototype_charAt, nil), 1)
	stringProto.defineDefaultProperty(vm, "charCodeAt", newFunctionObject(string_prototype_charCodeAt, nil), 1)
	stringProto.defineDefaultProperty(vm, "concat", newFunctionObject(string_prototype_concat, nil), 1)
	stringProto.defineDefaultProperty(vm, "indexOf", newFunctionObject(string_prototype_indexOf, nil), 1)
	stringProto.defineDefaultProperty(vm, "lastIndexOf", newFunctionObject(string_prototype_lastIndexOf, nil), 1)
	stringProto.defineDefaultProperty(vm, "toLowerCase", newFunctionObject(string_prototype_toLowerCase, nil), 0)
	stringProto.defineDefaultProperty(vm, "toUpperCase", newFunctionObject(string_prototype_toUpperCase, nil), 0)
	stringProto.defineDefaultProperty(vm, "trim", newFunctionObject(string_prototype_trim, nil), 0)

	stringO := newFunctionObject(string_call, string_ctor)
	stringO.odata.prototype = stringProto

	stringProto.defineDefaultProperty(vm, "constructor", stringO, 0)

	return stringO
}

func string_call(vm *vm, f value, args []value) value {
	return newString(args[0].toString())
}

func string_ctor(vm *vm, f value, args []value) value {
	return newStringObject(args[0].toString())
}

func string_prototype_toString(vm *vm, f value, args []value) value {
	switch f.vtype {
	case STRING:
		break
	case OBJECT:
		if f.odata.objectType == STRING_OBJECT {
			break
		}
		fallthrough
	default:
		panic(fmt.Sprintf("Not a string! %s", f)) // ### throw
	}
	return newString(f.asString())
}

func string_prototype_valueOf(vm *vm, f value, args []value) value {
	switch f.vtype {
	case STRING:
		break
	case OBJECT:
		if f.odata.objectType == STRING_OBJECT {
			break
		}
		fallthrough
	default:
		panic(fmt.Sprintf("Not a string! %s", f)) // ### throw
	}
	return newString(f.asString())
}

func string_prototype_charAt(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()
	pos := args[0].toInteger()
	size := len(S)
	if pos < 0 || pos >= size {
		return newString("")
	}

	return newString(string(S[pos]))
}

func string_prototype_charCodeAt(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()
	pos := args[0].toInteger()
	size := len(S)
	if pos < 0 || pos >= size {
		return newNumber(math.NaN())
	}

	return newNumber(float64(S[pos]))
}

func string_prototype_concat(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()

	for _, arg := range args {
		S += arg.toString()
	}

	return newString(S)
}

func string_prototype_indexOf(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()
	searchStr := args[0].toString()
	pos := 0
	if len(args) > 1 {
		pos = args[1].toInteger()
	}

	if pos > len(S) {
		return newNumber(-1)
	}

	return newNumber(float64(strings.Index(S[pos:], searchStr) + pos))
}

func string_prototype_lastIndexOf(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()
	searchStr := args[0].toString()
	pos := len(S)
	if len(args) > 1 {
		pos = args[1].toInteger()
	}

	if pos > len(S) {
		return newNumber(-1)
	}

	return newNumber(float64(strings.LastIndex(S[0:pos], searchStr)))
}

// ### localeCompare
// ### match
// ### replace
// ### search
// ### slice
// ### split
// ### substring

func string_prototype_toLowerCase(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()

	log.Printf(strings.ToLower(S))
	return newString(strings.ToLower(S))
}

// ### toLocaleLowerCase

func string_prototype_toUpperCase(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()

	return newString(strings.ToUpper(S))
}

// ### toLocaleUpperCase

func string_prototype_trim(vm *vm, f value, args []value) value {
	f.checkObjectCoercible(vm)
	S := f.toString()

	return newString(strings.Trim(S, "\n "))
}
