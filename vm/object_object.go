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
)

type rootObjectData struct {
	*valueBasicObjectData
}

func (this *rootObjectData) Prototype() *valueBasicObject {
	return nil
}

type basicObjectData struct {
	*valueBasicObjectData
}

func (this *basicObjectData) Prototype() *valueBasicObject {
	return &objectProto
}

var objectProto valueBasicObject

// Keep in mind that this is not just used by this file.
func newBasicObject() valueBasicObject {
	v := valueBasicObject{&basicObjectData{&valueBasicObjectData{extensible: true}}}
	return v
}

func defineObjectCtor(vm *vm) value {
	objectProto = valueBasicObject{&rootObjectData{&valueBasicObjectData{extensible: true}}}
	objectProto.defineDefaultProperty(vm, "toString", newFunctionObject(object_prototype_toString, nil), 0)
	objectProto.defineDefaultProperty(vm, "valueOf", newFunctionObject(object_prototype_valueOf, nil), 0)
	objectProto.defineDefaultProperty(vm, "hasOwnProperty", newFunctionObject(object_prototype_hasOwnProperty, nil), 0)

	objectCtor := newFunctionObject(object_call, object_ctor)
	objectCtor.defineDefaultProperty(vm, "getPrototypeOf", newFunctionObject(object_ctor_getPrototypeOf, nil), 0)
	objectCtor.prototype = &objectProto

	return objectCtor
}

func object_call(vm *vm, f value, args []value) value {
	return args[0].ToObject()
}

func object_ctor(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		v := args[0]
		switch v.(type) {
		case valueBasicObject:
			return v
		case valueString:
			return v.ToObject()
		case valueBool:
			return v.ToObject()
		case valueNumber:
			return v.ToObject()
		}
	}

	o := newBasicObject()
	return o
}

func object_prototype_toString(vm *vm, f value, args []value) value {
	switch f.(type) {
	case valueUndefined:
		return newString("[object Undefined]")
	case valueNull:
		return newString("[object Null]")
	case valueString:
		return newString("[object String]")
	case valueNumber:
		return newString("[object Number]")
	}

	o := f.ToObject()
	switch o.(type) {
	case stringObject:
		return newString("[object String]")
	case functionObject:
		return newString("[object Function]")
	case numberObject:
		return newString("[object Number]")
	}
	switch o.objectData().(type) {
	case *basicObjectData:
		return newString("[object Object]")
	case *booleanObjectData:
		return newString("[object Boolean]")
	}
	panic(fmt.Sprintf("%T is an unknown object type", o.objectData()))
}

func object_prototype_valueOf(vm *vm, f value, args []value) value {
	o := f.ToObject()
	return o
}

func object_prototype_hasOwnProperty(vm *vm, f value, args []value) value {
	P := args[0].ToString()
	O := f.ToObject()

	pd := O.getOwnProperty(vm, newString(P.String()))
	if pd == nil {
		return newBool(false)
	} else {
		return newBool(true)
	}
}

func object_ctor_getPrototypeOf(vm *vm, f value, args []value) value {
	switch o := f.(type) {
	case valueBasicObject:
		return o.odata.Prototype()
	default:
		return vm.ThrowTypeError("")
	}
}

func object_get(vm *vm, f value, prop value, pd *propertyDescriptor) value {
	if objectDebug {
		log.Printf("object_get %s %+v", prop, pd)
	}
	return pd.value
}

func object_set(vm *vm, f value, prop value, pd *propertyDescriptor, v value) value {
	if objectDebug {
		log.Printf("object_set %s", prop)
	}
	pd.value = v
	return newUndefined()
}
