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

import ()

func newArrayObject(s []value) valueObject {
	// Ensure we get consistent handling of this for testing's sake
	v := newObject()
	v.odata.objectType = ARRAY_OBJECT
	v.odata.primitiveData = newArrayData(s)
	v.odata.prototype = &arrayProto
	return v
}

var arrayProto valueObject

func defineArrayCtor(vm *vm) value {
	arrayProto = newObject()
	//arrayProto.defineDefaultProperty(vm, "toString", newFunctionObject(array_prototype_toString, nil), 0)

	arrayO := newFunctionObject(array_call, array_ctor)
	//arrayO.odata.prototype = &arrayProto

	arrayProto.defineDefaultProperty(vm, "constructor", arrayO, 0)

	return arrayO
}

func array_call(vm *vm, f value, args []value) value {
	return array_ctor(vm, f, args)
}

func array_ctor(vm *vm, f value, args []value) value {
	return newArrayObject(args)
}
