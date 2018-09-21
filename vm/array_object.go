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
	arrayProto.defineDefaultProperty(vm, "toString", newFunctionObject(array_prototype_toString, nil), 0)
	arrayProto.defineDefaultProperty(vm, "concat", newFunctionObject(array_prototype_concat, nil), 1)
	arrayProto.defineDefaultProperty(vm, "join", newFunctionObject(array_prototype_join, nil), 1)
	arrayProto.defineDefaultProperty(vm, "pop", newFunctionObject(array_prototype_pop, nil), 0)
	arrayProto.defineDefaultProperty(vm, "push", newFunctionObject(array_prototype_push, nil), 1)
	arrayProto.defineDefaultProperty(vm, "reverse", newFunctionObject(array_prototype_reverse, nil), 1)
	arrayProto.defineDefaultProperty(vm, "shift", newFunctionObject(array_prototype_shift, nil), 1)

	arrayO := newFunctionObject(array_call, array_ctor)
	arrayProto.defineDefaultProperty(vm, "constructor", arrayO, 0)
	arrayO.defineDefaultProperty(vm, "isArray", newFunctionObject(array_isArray, nil), 0)

	return arrayO
}

func array_call(vm *vm, f value, args []value) value {
	return array_ctor(vm, f, args)
}

func array_ctor(vm *vm, f value, args []value) value {
	return newArrayObject(args)
}

func array_isArray(vm *vm, f value, args []value) value {
	switch f.(type) {
	case arrayObject:
		return newBool(true)
	}

	return newBool(false)
}

func array_prototype_toString(vm *vm, f value, args []value) value {
	array := f.ToObject()
	funcJ := array.get(vm, newString("join"))
	switch typedJ := funcJ.(type) {
	case functionObject:
		return typedJ.call(vm, array, []value{newUndefined()})
	default:
		return object_prototype_toString(vm, array, []value{})
	}
}

// ### toLocaleString

func array_prototype_concat(vm *vm, f value, args []value) value {
	other := args[0].(arrayObject)
	switch typedJ := f.(type) {
	case arrayObject:
		ad := valueArrayData{values: make([]value, len(typedJ.primitiveData.values)+len(other.primitiveData.values))}
		for idx, v := range typedJ.primitiveData.values {
			ad.values[idx] = v
		}
		for idx, v := range other.primitiveData.values {
			ad.values[len(typedJ.primitiveData.values)+idx] = v
		}

		return arrayObject{valueBasicObject: newBasicObject(), primitiveData: &ad}
	default:
		panic("TypeError")
	}
}

func array_prototype_join(vm *vm, f value, args []value) value {
	var sep valueString
	if args[0] == newUndefined() {
		sep = ","
	} else {
		sep = args[0].ToString()
	}
	switch typedJ := f.(type) {
	case arrayObject:
		if len(typedJ.primitiveData.values) == 0 {
			return newString("")
		}

		element0 := typedJ.primitiveData.values[0]
		var R valueString
		if element0 == newUndefined() || element0 == newNull() {
			R = newString("")
		} else {
			R = element0.ToString()
		}

		k := 1
		for ; k < len(typedJ.primitiveData.values); k += 1 {
			S := R + sep
			element := typedJ.primitiveData.values[k]
			var next valueString
			if element == newUndefined() || element == newNull() {
				next = newString("")
			} else {
				next = element.ToString()
			}
			R = S + next
		}

		return R
	default:
		panic("TypeError")
	}
}

func array_prototype_pop(vm *vm, f value, args []value) value {
	switch typedJ := f.(type) {
	case arrayObject:
		if len(typedJ.primitiveData.values) == 0 {
			return newUndefined()
		}

		element := typedJ.primitiveData.values[len(typedJ.primitiveData.values)-1]
		typedJ.primitiveData.values = typedJ.primitiveData.values[:len(typedJ.primitiveData.values)-1]
		return element
	default:
		panic("TypeError")
	}
}

func array_prototype_push(vm *vm, f value, args []value) value {
	switch typedJ := f.(type) {
	case arrayObject:
		for _, v := range args {
			typedJ.primitiveData.values = append(typedJ.primitiveData.values, v)
		}
		return newNumber(float64(len(typedJ.primitiveData.values)))
	default:
		panic("TypeError")
	}
}

func array_prototype_reverse(vm *vm, f value, args []value) value {
	switch typedJ := f.(type) {
	case arrayObject:
		for i, j := 0, len(typedJ.primitiveData.values)-1; i < j; i, j = i+1, j-1 {
			typedJ.primitiveData.values[i], typedJ.primitiveData.values[j] = typedJ.primitiveData.values[j], typedJ.primitiveData.values[i]
		}
		return typedJ
	default:
		panic("TypeError")
	}
}

func array_prototype_shift(vm *vm, f value, args []value) value {
	switch typedJ := f.(type) {
	case arrayObject:
		if len(typedJ.primitiveData.values) == 0 {
			return newUndefined()
		}

		element := typedJ.primitiveData.values[0]
		typedJ.primitiveData.values = typedJ.primitiveData.values[1:]
		return element
	default:
		panic("TypeError")
	}
}

// ### slice
// sort
// splice
// unshift
// indexOf
// lastIndexOf
// every
// some
// forEach
// map
// filter
// reduce
// reduceRight
