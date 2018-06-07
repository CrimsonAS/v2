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
	"fmt"
	"log"
	"math"
	"strings"
)

func newStringObject(s string) value {
	v := newString(s)
	v.vtype = OBJECT
	v.odata = &objectData{STRING_OBJECT, value{}, nil, nil, nil, true}
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
	if len(args) > 0 {
		return newString(args[0].toString())
	} else {
		return newString("")
	}
}

func string_ctor(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		return newStringObject(args[0].toString())
	} else {
		return newStringObject("")
	}
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
