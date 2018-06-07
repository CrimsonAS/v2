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

var booleanProto value

func newBooleanObject(b bool) value {
	v := newBool(b)
	v.vtype = OBJECT
	v.odata = &objectData{BOOLEAN_OBJECT, value{}, nil, nil, nil}
	v.odata.prototype = booleanProto
	return v
}

func defineBooleanCtor(vm *vm) value {
	booleanProto = newObject()
	booleanProto.defineDefaultProperty(vm, "toString", newFunctionObject(boolean_prototype_toString, nil), 0)
	booleanProto.defineDefaultProperty(vm, "valueOf", newFunctionObject(boolean_prototype_valueOf, nil), 0)

	boolO := newFunctionObject(boolean_call, boolean_ctor)
	boolO.odata.prototype = booleanProto

	booleanProto.defineDefaultProperty(vm, "constructor", boolO, 0)
	return boolO
}

func boolean_call(vm *vm, f value, args []value) value {
	return newBool(args[0].toBoolean())
}

func boolean_ctor(vm *vm, f value, args []value) value {
	return newBooleanObject(args[0].toBoolean())
}

func boolean_prototype_toString(vm *vm, f value, args []value) value {
	switch f.vtype {
	case BOOL:
		break
	case OBJECT:
		if f.odata.objectType == BOOLEAN_OBJECT {
			break
		}
		fallthrough
	default:
		panic(fmt.Sprintf("Not a boolean! %s", f)) // ### throw
	}

	if f.asBool() {
		return newString("true")
	} else {
		return newString("false")
	}
}

func boolean_prototype_valueOf(vm *vm, f value, args []value) value {
	switch f.vtype {
	case BOOL:
		break
	case OBJECT:
		if f.odata.objectType == BOOLEAN_OBJECT {
			break
		}
		fallthrough
	default:
		panic(fmt.Sprintf("Not a boolean! %s", f)) // ### throw
	}
	return newBool(f.asBool())
}
