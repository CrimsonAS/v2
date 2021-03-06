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
	"reflect"
	"strconv"
)

/////////////////////////////////
// constructors
/////////////////////////////////

func newUndefined() valueUndefined {
	return valueUndefined{}
}

func newNull() valueNull {
	return valueNull{}
}

func newBool(val bool) value {
	if val {
		return valueBool(true)
	} else {
		return valueBool(false)
	}
}

func newNumber(val float64) valueNumber {
	return valueNumber(val)
}

func newString(val string) valueString {
	return valueString(val)
}

/////////////////////////////////
// type definitions
/////////////////////////////////

// value represents an abstraction of a JavaScript value.
type value interface {
	ToInteger() int
	ToNumber() float64
	ToBoolean() bool
	ToString() valueString
	ToObject() valueObject
	hasPrimitiveBase() bool
	String() string
}

// ES5 11.9.3
func abstractEqualityComparison(x, y value) bool {
	xt := reflect.TypeOf(x)
	yt := reflect.TypeOf(y)

	if xt == yt {
		switch x.(type) {
		case valueUndefined:
			return true
		case valueNull:
			return true
		case valueNumber:
			xn := x.ToNumber()
			yn := y.ToNumber()
			if math.IsNaN(xn) {
				return false
			}
			if math.IsNaN(yn) {
				return false
			}
			if xn == yn {
				return true
			}
			if xn == +0 && yn == -0 {
				return true
			} else if xn == -0 && yn == +0 {
				return true
			}
			return false
		case valueString:
			return x.ToString() == y.ToString()
		case valueBool:
			xb := x.ToBoolean()
			yb := y.ToBoolean()
			if xb && yb || !xb && !yb {
				return true
			} else {
				return false
			}
		case valueBasicObject:
			if x.(valueBasicObject).odata == y.(valueBasicObject).odata {
				return true
			} else {
				return false
			}
		}
		return false
	}

	switch x.(type) {
	case valueNull:
		switch y.(type) {
		case valueUndefined:
			return true
		}
	case valueUndefined:
		switch y.(type) {
		case valueNull:
			return true
		}
	}

	switch x.(type) {
	case valueNumber:
		switch y.(type) {
		case valueString:
			return x.ToNumber() == y.ToNumber() // ### probably need to handle +/- 0 as above
		}
	case valueString:
		switch y.(type) {
		case valueNumber:
			return x.ToNumber() == y.ToNumber() // ### probably need to handle +/- 0 as above
		}
	}

	switch x.(type) {
	case valueBool:
		return x.ToNumber() == y.ToNumber()
	}

	switch y.(type) {
	case valueBool:
		return x.ToNumber() == y.ToNumber()
	}

	switch x.(type) {
	case valueString:
		switch y.(type) {
		case valueBasicObject:
			return valueToPrimitive(x) == valueToPrimitive(y)
		}
	case valueNumber:
		switch y.(type) {
		case valueBasicObject:
			return valueToPrimitive(x) == valueToPrimitive(y)
		}
	}
	switch x.(type) {
	case valueBasicObject:
		switch y.(type) {
		case valueString:
			return valueToPrimitive(x) == valueToPrimitive(y)
		case valueNumber:
			return valueToPrimitive(x) == valueToPrimitive(y)
		}
	}

	return false
}

// ES5 11.9.6
func strictEqualityComparison(x, y value) bool {
	xt := reflect.TypeOf(x)
	yt := reflect.TypeOf(y)

	if xt != yt {
		return false
	}

	switch x.(type) {
	case valueUndefined:
		return true
	case valueNull:
		return true
	case valueNumber:
		xn := x.ToNumber()
		yn := y.ToNumber()
		if math.IsNaN(xn) {
			return false
		}
		if math.IsNaN(yn) {
			return false
		}
		if xn == yn {
			return true
		}
		if xn == +0 && yn == -0 {
			return true
		} else if xn == -0 && yn == +0 {
			return true
		}

		return false
	case valueString:
		return x.ToString() == y.ToString()
	case valueBool:
		xb := x.ToBoolean()
		yb := y.ToBoolean()
		if xb && yb || !xb && !yb {
			return true
		} else {
			return false
		}
	case valueBasicObject:
		if x.(valueBasicObject).odata == y.(valueBasicObject).odata {
			return true
		} else {
			return false
		}
	}

	return false
}

/////////////////////////////////

type valueUndefined struct {
}

func (this valueUndefined) ToInteger() int {
	return 0
}
func (this valueUndefined) ToNumber() float64 {
	return math.NaN()
}
func (this valueUndefined) ToBoolean() bool {
	return false
}
func (this valueUndefined) ToString() valueString {
	return "undefined"
}
func (this valueUndefined) ToObject() valueObject {
	panic("TypeError") // ### can't ThrowTypeError as we have no vm
}
func (this valueUndefined) hasPrimitiveBase() bool {
	return false
}
func (this valueUndefined) String() string {
	return this.ToString().String()
}

/////////////////////////////////

type valueNull struct {
}

func (this valueNull) ToInteger() int {
	return 0
}
func (this valueNull) ToNumber() float64 {
	return +0
}
func (this valueNull) ToBoolean() bool {
	return false
}
func (this valueNull) ToString() valueString {
	return "null"
}
func (this valueNull) ToObject() valueObject {
	panic("TypeError") // ### can't ThrowTypeError as we have no vm
}
func (this valueNull) hasPrimitiveBase() bool {
	return false
}
func (this valueNull) String() string {
	return this.ToString().String()
}

/////////////////////////////////

type valueBool bool

func (this valueBool) ToInteger() int {
	if this {
		return 1
	} else {
		return 0
	}
}

func (this valueBool) ToNumber() float64 {
	if this {
		return 1
	} else {
		return +0
	}
}

func (this valueBool) ToBoolean() bool {
	return bool(this)
}

func (this valueBool) ToString() valueString {
	if this {
		return "true"
	} else {
		return "false"
	}
}

func (this valueBool) ToObject() valueObject {
	return newBooleanObject(bool(this))
}

func (this valueBool) hasPrimitiveBase() bool {
	return true
}

func (this valueBool) String() string {
	return this.ToString().String()
}

/////////////////////////////////

type valueString string

func (this valueString) ToInteger() int {
	panic("not implemented")
}

func (this valueString) ToNumber() float64 {
	// ### toNumber(string) not implemented (es5 9.3.1)
	v, _ := strconv.ParseFloat(this.String(), 64)
	return v
}

func (this valueString) ToBoolean() bool {
	return len(this) > 0
}

func (this valueString) ToString() valueString {
	return this
}

func (this valueString) ToObject() valueObject {
	return newStringObject(this.String())
}

func (this valueString) hasPrimitiveBase() bool {
	return true
}

func (this valueString) String() string {
	return string(this)
}

/////////////////////////////////

type valueNumber float64

func (this valueNumber) ToInteger() int {
	tf := float64(this)
	if math.IsNaN(tf) {
		return +0
	}
	if int(tf) == 0 || math.IsInf(tf, 0) {
		return int(tf)
	}

	if tf > 0 {
		return int(1 * math.Floor(math.Abs(tf)))
	} else {
		return int(-1 * math.Floor(math.Abs(tf)))
	}
}

func (this valueNumber) ToNumber() float64 {
	return float64(this)
}

func (this valueNumber) ToBoolean() bool {
	if int(this) == 0 || math.IsNaN(float64(this)) {
		return false
	} else {
		return true
	}
}

func (this valueNumber) ToString() valueString {
	return newString(fmt.Sprintf("%f", float64(this)))
}

func (this valueNumber) ToObject() valueObject {
	return newNumberObject(float64(this))
}

func (this valueNumber) hasPrimitiveBase() bool {
	return true
}

func (this valueNumber) String() string {
	return this.ToString().String()
}

/////////////////////////////////

func (this valueBasicObject) ToInteger() int {
	return int(this.ToNumber())
}

func (this valueBasicObject) ToNumber() float64 {
	panic("object conversion not implemented")
}

func (this valueBasicObject) ToBoolean() bool {
	return true
}

func (this valueBasicObject) ToString() valueString {
	return "[object]"
}

func (this valueBasicObject) ToObject() valueObject {
	return this
}

func (this valueBasicObject) hasPrimitiveBase() bool {
	return false
}

func (this valueBasicObject) String() string {
	return this.ToString().String()
}

//////////////////////////////////////
//////////////////////////////////////

func checkObjectCoercible(vm *vm, v value) {
	switch v.(type) {
	case valueUndefined:
		vm.ThrowTypeError("")
	case valueNull:
		vm.ThrowTypeError("")
	case valueBool:
	case valueNumber:
	case valueString:
	case valueObject:
	}
}

func valueToPrimitive(v value) value {
	switch v.(type) {
	case valueUndefined:
		return v
	case valueNull:
		return v
	case valueBool:
		return v
	case valueNumber:
		return v
	case valueString:
		return v
	case valueObject:
		panic("object conversion not implemented")
	}
	panic("unreachable")
}
