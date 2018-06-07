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
	"unsafe"
)

type objectType uint8

const (
	NOT_AN_OBJECT objectType = iota
	OBJECT_PLAIN
	BOOLEAN_OBJECT
	NUMBER_OBJECT
	STRING_OBJECT
	FUNCTION_OBJECT
)

func (this value) defineDefaultProperty(vm *vm, prop string, v value, lt int) value {
	if objectDebug {
		log.Printf("Defining %s on %s = %s", prop, this, v)
	}
	if pdata := this.getOwnProperty(vm, prop); pdata != nil {
		panic(fmt.Sprintf("property already exists: %s", prop))
	}

	pd := propertyDescriptor{name: prop, get: object_get, set: object_set, length: lt, enumerable: true, configurable: false, value: v}
	this.odata.properties = append(this.odata.properties, pd)
	return newUndefined()
}

func (this value) defineReadonlyProperty(vm *vm, prop string, v value, lt int) value {
	if objectDebug {
		log.Printf("Defining %s on %s = %s", prop, this, v)
	}
	if pdata := this.getOwnProperty(vm, prop); pdata != nil {
		panic(fmt.Sprintf("property already exists: %s", prop))
	}

	pd := propertyDescriptor{name: prop, get: object_get, set: object_set, length: lt, enumerable: true, configurable: false, value: v}
	this.odata.properties = append(this.odata.properties, pd)
	return newUndefined()
}

func (this value) set(vm *vm, prop string, v value) value {
	desc := this.getProperty(vm, prop)
	if desc == nil {
		return newUndefined()
	}

	return desc.set(vm, this, prop, desc, v)
}

func (this value) get(vm *vm, prop string) value {
	if objectDebug {
		log.Printf("Getting %s proto %s", prop, this.odata.prototype)
	}
	desc := this.getProperty(vm, prop)
	if desc == nil {
		return newUndefined()
	}

	return desc.get(vm, this, prop, desc)
}

func (this value) getOwnProperty(vm *vm, prop string) *propertyDescriptor {
	if objectDebug {
		log.Printf("GetOwnProperty %s.%s", this, prop)
	}
	for idx, _ := range this.odata.properties {
		//log.Printf("Looking for %s found %s", prop, this.odata.properties[idx].name)
		if this.odata.properties[idx].name == prop {
			return &this.odata.properties[idx]
		}
	}

	// ### STRING_OBJECT GetOwnProperty (es5 15.5.5.2)
	return nil
}

func (this value) getProperty(vm *vm, prop string) *propertyDescriptor {
	pd := this.getOwnProperty(vm, prop)
	if pd != nil {
		return pd
	}

	if this.odata.prototype.vtype != OBJECT {
		return nil
	}

	return this.odata.prototype.getProperty(vm, prop)
}

type foFn func(vm *vm, f value, args []value) value
type getFn func(vm *vm, f value, prop string, pd *propertyDescriptor) value
type setFn func(vm *vm, f value, prop string, pd *propertyDescriptor, v value) value

type propertyDescriptor struct {
	name         string
	get          getFn // [[Get]]
	set          setFn // [[Set]]
	value        value // [[Value]] convenience
	length       int
	propIdx      int
	enumerable   bool // [[Enumerable]]
	configurable bool // [[Configurable]]
}

type objectData struct {
	objectType   objectType
	prototype    value
	properties   []propertyDescriptor
	callPtr      foFn // used for function object
	constructPtr foFn // used for function object
}

const objectDebug = false

func newObject() value {
	v := value{OBJECT, nil, &objectData{OBJECT_PLAIN, value{}, nil, nil, nil}}
	return v
}

func newNumberObject(n float64) value {
	v := value{OBJECT, make([]byte, unsafe.Sizeof(n)), &objectData{NUMBER_OBJECT, value{}, nil, nil, nil}}
	*(*float64)(unsafe.Pointer(&v.vdata[0])) = n
	return v
}

func newFunctionObject(call foFn, construct foFn) value {
	v := value{OBJECT, nil, &objectData{FUNCTION_OBJECT, value{}, nil, call, construct}}
	return v
}

func (this value) checkObjectCoercible(vm *vm) {
	switch this.vtype {
	case UNDEFINED:
		panic("TypeError")
	case NULL:
		panic("TypeError")
	case BOOL:
	case NUMBER:
	case STRING:
	case OBJECT:
	}
}

func (this value) toObject() value {
	switch this.vtype {
	case UNDEFINED:
		panic("TypeError")
	case NULL:
		panic("TypeError")
	case BOOL:
		return newBooleanObject(this.asBool())
	case NUMBER:
		return newNumberObject(this.asNumber())
	case STRING:
		return newStringObject(this.asString())
	case OBJECT:
		return this
	}

	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this value) call(vm *vm, thisArg value, args []value) value {
	if this.vtype != OBJECT {
		panic(fmt.Sprintf("can't convert! %s", this.vtype))
	}
	if this.odata.objectType != FUNCTION_OBJECT {
		panic(fmt.Sprintf("can't call non-function! %d", this.odata.objectType))
	}
	return this.odata.callPtr(vm, thisArg, args)
}

func (this value) construct(vm *vm, thisArg value, args []value) value {
	if this.vtype != OBJECT {
		panic(fmt.Sprintf("can't convert! %s", this.vtype))
	}
	if this.odata.objectType != FUNCTION_OBJECT {
		panic(fmt.Sprintf("can't call non-function! %d", this.odata.objectType))
	}
	return this.odata.constructPtr(vm, thisArg, args)
}
