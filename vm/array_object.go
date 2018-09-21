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

type arrayObject struct {
	valueBasicObject
	primitiveData *valueArrayData
}

func (this *arrayObject) Prototype() *valueBasicObject {
	return &arrayProto
}

//////////////////////////////////////
// value methods
//////////////////////////////////////

func (this arrayObject) ToInteger() int {
	return int(this.ToNumber())
}

func (this arrayObject) ToNumber() float64 {
	panic("object conversion not implemented")
}

func (this arrayObject) ToBoolean() bool {
	return true
}

func (this arrayObject) ToString() valueString {
	return this.primitiveData.ToString()
}

func (this arrayObject) ToObject() valueObject {
	return this
}

func (this arrayObject) hasPrimitiveBase() bool {
	return true
}

func (this arrayObject) String() string {
	return string(this.primitiveData.ToString())
}

//////////////////////////////////////
// object methods
//////////////////////////////////////

func (this arrayObject) defineOwnProperty(vm *vm, prop value, desc *propertyDescriptor, throw bool) bool {
	return true
}

func (this arrayObject) getOwnProperty(vm *vm, prop value) *propertyDescriptor {
	return nil
}

func (this arrayObject) put(vm *vm, prop value, v value, throw bool) {
	if numIdx, ok := prop.(valueNumber); ok {
		idx := numIdx.ToInteger()
		if float64(idx) == float64(numIdx) {
			this.primitiveData.Set(idx, v)
		}
	}

	this.valueBasicObject.put(vm, prop, v, throw)
}

func (this arrayObject) get(vm *vm, prop value) value {
	// ### belongs in getOwnProperty perhaps?
	if numIdx, ok := prop.(valueNumber); ok {
		idx := numIdx.ToInteger()
		if float64(idx) == float64(numIdx) && idx >= 0 && idx < len(this.primitiveData.values) {
			return this.primitiveData.values[idx]
		}
	}

	if this.valueBasicObject.getOwnProperty(vm, prop) != nil {
		return this.valueBasicObject.get(vm, prop)
	} else {
		return arrayProto.get(vm, prop)
	}
}

//////////////////////////////////////
// array data
//////////////////////////////////////

// Copying is necessary, otherwise we'll end up with stack data, which is bad
func newArrayData(v []value) *valueArrayData {
	ad := valueArrayData{values: make([]value, len(v))}
	for idx, _ := range v {
		ad.values[idx] = v[idx]
	}
	return &ad
}

func (this valueArrayData) Get(idx int) value {
	return this.values[idx]
}

func (this valueArrayData) Set(idx int, v value) {
	this.values[idx] = v
}

type valueArrayData struct {
	values []value
}

func (this valueArrayData) ToInteger() int {
	panic("Should never happen")
}

func (this valueArrayData) ToNumber() float64 {
	panic("Should never happen")
}

func (this valueArrayData) ToBoolean() bool {
	panic("Should never happen")
}

func (this valueArrayData) ToString() valueString {
	return newString(fmt.Sprintf("ARRAY[%s]", this.values))
}

func (this valueArrayData) ToObject() valueObject {
	panic("Should never happen")
}

func (this valueArrayData) hasPrimitiveBase() bool {
	panic("Should never happen")
}

func (this valueArrayData) String() string {
	return this.ToString().String()
}

//////////////////////////////////////

func newArrayObject(s []value) valueObject {
	return arrayObject{valueBasicObject: newBasicObject(), primitiveData: newArrayData(s)}
}

var arrayProto valueBasicObject

func defineArrayCtor(vm *vm) value {
	arrayProto = valueBasicObject{&rootObjectData{&valueBasicObjectData{extensible: true}}}
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
