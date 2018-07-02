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
	"math"
	"strings"
)

type stringObjectData struct {
	*valueBasicObjectData
	primitiveData valueString
}

func (this *stringObjectData) Prototype() *valueBasicObject {
	return &stringProto
}

func newStringObject(s string) valueBasicObject {
	return valueBasicObject{&stringObjectData{&valueBasicObjectData{}, newString(s)}}
}

var stringProto valueBasicObject

func defineStringCtor(vm *vm) value {
	stringProto = valueBasicObject{&rootObjectData{&valueBasicObjectData{extensible: true}}}
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
	stringO.odata.(*functionObjectData).prototype = &stringProto

	stringProto.defineDefaultProperty(vm, "constructor", stringO, 0)

	return stringO
}

func string_call(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		return newString(args[0].ToString().String())
	} else {
		return newString("")
	}
}

func string_ctor(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		return newStringObject(args[0].ToString().String())
	} else {
		return newStringObject("")
	}
}

func string_prototype_toString(vm *vm, f value, args []value) value {
	switch o := f.(type) {
	case valueString:
		return newString(f.ToString().String())
	case valueBasicObject:
		if sd, ok := o.odata.(*stringObjectData); ok {
			return newString(sd.primitiveData.ToString().String())
		}
	default:
		panic(fmt.Sprintf("Not a string! %s", f)) // ### throw
	}
	panic("unreachable")
}

func string_prototype_valueOf(vm *vm, f value, args []value) value {
	switch o := f.(type) {
	case valueString:
		return newString(f.ToString().String())
	case valueBasicObject:
		if sd, ok := o.odata.(*stringObjectData); ok {
			return newString(sd.primitiveData.ToString().String())
		}
	default:
		panic(fmt.Sprintf("Not a string! %s", f)) // ### throw
	}
	panic("unreachable")
}

func string_prototype_charAt(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()
	pos := args[0].ToInteger()
	size := len(S)
	if pos < 0 || pos >= size {
		return newString("")
	}

	return newString(string(S[pos]))
}

func string_prototype_charCodeAt(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()
	pos := args[0].ToInteger()
	size := len(S)
	if pos < 0 || pos >= size {
		return newNumber(math.NaN())
	}

	return newNumber(float64(S[pos]))
}

func string_prototype_concat(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()

	for _, arg := range args {
		S += arg.ToString().String()
	}

	return newString(S)
}

func string_prototype_indexOf(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()
	searchStr := args[0].ToString().String()
	pos := 0
	if len(args) > 1 {
		pos = args[1].ToInteger()
	}

	if pos > len(S) {
		return newNumber(-1)
	}

	return newNumber(float64(strings.Index(S[pos:], searchStr) + pos))
}

func string_prototype_lastIndexOf(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()
	searchStr := args[0].ToString().String()
	pos := len(S)
	if len(args) > 1 {
		pos = args[1].ToInteger()
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
	checkObjectCoercible(vm, f)
	S := f.ToString().String()

	return newString(strings.ToLower(S))
}

// ### toLocaleLowerCase

func string_prototype_toUpperCase(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()

	return newString(strings.ToUpper(S))
}

// ### toLocaleUpperCase

func string_prototype_trim(vm *vm, f value, args []value) value {
	checkObjectCoercible(vm, f)
	S := f.ToString().String()

	return newString(strings.Trim(S, "\n "))
}
