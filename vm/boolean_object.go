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
)

var booleanProto valueObject

type booleanObjectData struct {
	*valueObjectData
	primitiveData bool
}

func (this *booleanObjectData) Prototype() *valueObject {
	return &booleanProto
}

func newBooleanObject(b bool) valueObject {
	return valueObject{&booleanObjectData{&valueObjectData{extensible: true}, b}}
}

func defineBooleanCtor(vm *vm) valueObject {
	booleanProto = valueObject{&rootObjectData{&valueObjectData{extensible: true}}}
	booleanProto.defineDefaultProperty(vm, "toString", newFunctionObject(boolean_prototype_toString, nil), 0)
	booleanProto.defineDefaultProperty(vm, "valueOf", newFunctionObject(boolean_prototype_valueOf, nil), 0)

	boolO := newFunctionObject(boolean_call, boolean_ctor)
	boolO.odata.(*functionObjectData).prototype = &booleanProto

	booleanProto.defineDefaultProperty(vm, "constructor", boolO, 0)
	return boolO
}

func boolean_call(vm *vm, f value, args []value) value {
	return newBool(args[0].ToBoolean())
}

func boolean_ctor(vm *vm, f value, args []value) value {
	return newBooleanObject(args[0].ToBoolean())
}

func boolean_prototype_toString(vm *vm, f value, args []value) value {
	b := false
	switch o := f.(type) {
	case valueBool:
		b = bool(o)
	case valueObject:
		b = o.odata.(*booleanObjectData).primitiveData
	default:
		panic(fmt.Sprintf("Not a boolean! %s", f)) // ### throw
	}

	if b {
		return newString("true")
	} else {
		return newString("false")
	}
}

func boolean_prototype_valueOf(vm *vm, f value, args []value) value {
	b := false
	switch o := f.(type) {
	case valueBool:
		b = bool(o)
	case valueObject:
		b = o.odata.(*booleanObjectData).primitiveData
	default:
		panic(fmt.Sprintf("Not a boolean! %s", f)) // ### throw
	}

	return newBool(b)
}
