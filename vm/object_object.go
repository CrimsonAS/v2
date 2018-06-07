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
	"log"
)

var objectProto value

func defineObjectCtor(vm *vm) value {
	objectProto = newObject()
	objectProto.defineDefaultProperty(vm, "toString", newFunctionObject(object_prototype_toString, nil), 0)

	objectCtor := newFunctionObject(object_call, object_ctor)
	objectCtor.odata.prototype = objectProto
	objectProto.set(vm, "constructor", objectCtor)

	return objectCtor
}

func object_call(vm *vm, f value, args []value) value {
	return newBool(args[0].toBoolean())
}

func object_ctor(vm *vm, f value, args []value) value {
	o := newObject()
	o.odata.prototype = objectProto
	return o
}

func object_prototype_toString(vm *vm, f value, args []value) value {
	if f.vtype == UNDEFINED {
		return newString("[object Undefined]")
	} else if f.vtype == NULL {
		return newString("[object Null]")
	}

	o := f.toObject()
	switch o.odata.objectType {
	case OBJECT_PLAIN:
		return newString("[object Object]")
	case BOOLEAN_OBJECT:
		return newString("[object Boolean]")
	case NUMBER_OBJECT:
		return newString("[object Number]")
	case STRING_OBJECT:
		return newString("[object String]")
	case FUNCTION_OBJECT:
		return newString("[object Function]")
	}

	panic("unknown object type")
}

func object_get(vm *vm, f value, prop string, pd *propertyDescriptor) value {
	if objectDebug {
		log.Printf("object_get %s", prop)
	}
	return pd.value
}

func object_set(vm *vm, f value, prop string, pd *propertyDescriptor, v value) value {
	if objectDebug {
		log.Printf("object_set %s", prop)
	}
	pd.value = v
	return newUndefined()
}
