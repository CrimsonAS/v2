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
)

var numberProto valueBasicObject

type numberObjectData struct {
	*valueBasicObjectData
	primitiveData float64
}

func (this *numberObjectData) Prototype() *valueBasicObject {
	return &numberProto
}

func newNumberObject(f float64) valueBasicObject {
	return valueBasicObject{&numberObjectData{&valueBasicObjectData{extensible: true}, f}}
}

func defineNumberCtor(vm *vm) functionObject {
	numberProto = valueBasicObject{&rootObjectData{&valueBasicObjectData{extensible: true}}}
	numberProto.defineDefaultProperty(vm, "toString", newFunctionObject(number_prototype_toString, nil), 0)

	numberO := newFunctionObject(number_call, number_ctor)
	numberO.prototype = &numberProto
	numberO.defineDefaultProperty(vm, "MAX_VALUE", newNumber(math.MaxFloat64), 0)
	numberO.defineDefaultProperty(vm, "MIN_VALUE", newNumber(math.SmallestNonzeroFloat64), 0)
	numberO.defineDefaultProperty(vm, "NaN", newNumber(math.NaN()), 0)
	numberO.defineDefaultProperty(vm, "NEGATIVE_INFINITY", newNumber(math.Inf(-1)), 0)
	numberO.defineDefaultProperty(vm, "POSITIVE_INFINITY", newNumber(math.Inf(+1)), 0)

	numberProto.defineDefaultProperty(vm, "constructor", numberO, 0)
	return numberO
}

func number_call(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		return newNumber(args[0].ToNumber())
	} else {
		return newNumber(+0)
	}
}

func number_ctor(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		return newNumberObject(args[0].ToNumber())
	} else {
		return newNumberObject(+0)
	}
}

func number_prototype_toString(vm *vm, f value, args []value) value {
	if len(args) > 0 {
		panic("Can't handle radix right now")
	}

	n := 0.0
	switch o := f.(type) {
	case valueNumber:
		n = float64(o)
	case valueBasicObject:
		n = o.odata.(*numberObjectData).primitiveData
	default:
		panic(fmt.Sprintf("Not a number! %s", f)) // ### throw
	}

	return newString(fmt.Sprintf("%g", n))
}

// ### toLocaleString
// ### valueOf
// ### toFixed
// ### toExponential
// ### toPrecision
