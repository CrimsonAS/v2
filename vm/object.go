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

func (this value) defineDefaultProperty(vm *vm, prop string, v value, lt int) bool {
	pd := &propertyDescriptor{name: prop, length: lt, enumerable: true, configurable: false, value: v}
	return this.defineOwnProperty(vm, prop, pd, true)
}

func (this value) defineReadonlyProperty(vm *vm, prop string, v value, lt int) bool {
	pd := &propertyDescriptor{name: prop, length: lt, enumerable: true, configurable: false, value: v}
	return this.defineOwnProperty(vm, prop, pd, true)
}

func (this value) canPut(vm *vm, prop string) bool {
	desc := this.getOwnProperty(vm, prop)
	if desc != nil {
		if desc.isAccessorDescriptor() {
			if desc.set != nil {
				return false
			} else {
				return true
			}
		} else {
			return desc.writable
		}
	}

	proto := this.odata.prototype
	if proto.vtype == NULL {
		return this.odata.extensible
	}

	inherited := proto.getProperty(vm, prop)
	if inherited == nil {
		return this.odata.extensible
	}

	if inherited.isAccessorDescriptor() {
		if inherited.set == nil {
			return false
		} else {
			return true
		}
	} else {
		if !this.odata.extensible {
			return false
		} else {
			return inherited.writable
		}
	}
	return true
}

// ### ArrayObject [[DefineOwnProperty]] 15.4.5.1
func (this value) defineOwnProperty(vm *vm, prop string, desc *propertyDescriptor, throw bool) bool {
	if objectDebug {
		//log.Printf("Defining %s on %s %t", prop, this, this.odata.extensible)
	}
	current := this.getOwnProperty(vm, prop)
	extensible := this.odata.extensible
	if current == nil {
		if !extensible {
			if throw {
				panic("TypeError")
			} else {
				return false
			}
		}

		var pd *propertyDescriptor
		if desc.isGenericDescriptor() || desc.isDataDescriptor() {
			pd = &propertyDescriptor{name: prop, value: desc.value, writable: desc.writable, enumerable: desc.enumerable, configurable: desc.configurable}
		} else {
			pd = &propertyDescriptor{name: prop, get: desc.get, set: desc.set, enumerable: desc.enumerable, configurable: desc.configurable}
		}
		this.odata.properties = append(this.odata.properties, pd)
		//log.Printf("Added new property %s %+v", prop, pd)
		return true
	}

	// ### 8.12.9 5/6
	// "return true if every field in 'desc' is absent
	// return true if every field in desc also occurs in current, and the value
	// of every field in desc is the same value as the corresponding field in
	// current using the SameValue algorithm (9.12)

	if !current.configurable {
		if desc.configurable {
			if throw {
				panic("TypeError")
			} else {
				return false
			}
		}

		if desc.enumerable && !current.enumerable {
			if throw {
				panic("TypeError")
			} else {
				return false
			}
		}
	}

	if desc.isGenericDescriptor() {
		// no validation needed (es5 8.12.9 8)
	} else if current.isDataDescriptor() != desc.isDataDescriptor() {
		if !current.configurable {
			if throw {
				panic("TypeError")
			} else {
				return false
			}
		}

		if current.isDataDescriptor() {
			// ###
			// Convert the property named P of object O from a data property to
			// an accessor property. Preserve the existing values of
			// [[Configurable]] and [[Enumerable]] and set the rest of the
			// property's attributes to their default values.
		} else {
			// ###
			// Convert the property named P of object O from an accessor
			// property to a data property. Preserve the existing values of the
			// converted property's [[Configurable]] and [[Enumerable]
			// attributes, and set the rest of the property's attributes to
			// their default values.
		}
	} else if current.isDataDescriptor() && desc.isDataDescriptor() {
		if !current.configurable {
			if !current.writable && desc.writable {
				if throw {
					panic("TypeError")
				} else {
					return false
				}
			}

			if !current.writable {
				// ### reject if the value field of desc is present and
				// SameValue(desc.value, current.value) is false
			}
		}

		// If it's configurable, any change is OK.
	} else if current.isAccessorDescriptor() && desc.isAccessorDescriptor() {
		if !current.configurable {
			// ### reject if desc.set is present and samevalue(desc.set,
			// current.set) is false
			// reject if desc.get is present and samevalue(desc.get,
			// current.get) is false
		}
	}

	current.name = desc.name
	current.get = desc.get
	current.set = desc.set
	current.value = desc.value
	current.length = desc.length
	current.propIdx = desc.propIdx
	current.writable = desc.writable
	current.enumerable = desc.enumerable
	current.configurable = desc.configurable
	//log.Printf("Defined property %s on %s", prop, this)
	return true
}

func (this value) put(vm *vm, prop string, v value, throw bool) {
	if objectDebug {
		log.Printf("Setting %s = %s on %s", prop, v, this)
	}
	if !this.canPut(vm, prop) {
		if throw {
			panic("TypeError")
		}
		return
	}

	ownDesc := this.getOwnProperty(vm, prop)
	if ownDesc.isDataDescriptor() {
		valueDesc := &propertyDescriptor{value: v}
		this.defineOwnProperty(vm, prop, valueDesc, throw)
		return
	}

	desc := this.getProperty(vm, prop)
	if desc.isAccessorDescriptor() {
		desc.set(vm, this, prop, desc, v)
	} else {
		newDesc := &propertyDescriptor{value: v, writable: true, enumerable: true, configurable: true}
		this.defineOwnProperty(vm, prop, newDesc, throw)
	}
}

func (this value) get(vm *vm, prop string) value {
	desc := this.getProperty(vm, prop)
	if desc == nil {
		return newUndefined()
	}

	if desc.isDataDescriptor() {
		return desc.value
	} else if desc.isAccessorDescriptor() {
		return desc.get(vm, this, prop, desc)
	}

	panic("unreachable")
}

func (this value) getOwnProperty(vm *vm, prop string) *propertyDescriptor {
	if objectDebug {
		//log.Printf("GetOwnProperty %s.%s", this, prop)
	}
	for idx, _ := range this.odata.properties {
		//log.Printf("Looking for %s found %s", prop, this.odata.properties[idx].name)
		if this.odata.properties[idx].name == prop {
			return this.odata.properties[idx]
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
	writable     bool // [[Writable]]
	enumerable   bool // [[Enumerable]]
	configurable bool // [[Configurable]]
}

func (this *propertyDescriptor) isDataDescriptor() bool {
	// ### does this match the spec? hmm..
	return this.get == nil && this.set == nil
}

func (this *propertyDescriptor) isGenericDescriptor() bool {
	if !this.isAccessorDescriptor() && !this.isDataDescriptor() {
		return true
	}

	return false
}

func (this *propertyDescriptor) isAccessorDescriptor() bool {
	return this.get != nil || this.set != nil
}

type objectData struct {
	objectType   objectType
	prototype    value
	properties   []*propertyDescriptor
	callPtr      foFn // used for function object
	constructPtr foFn // used for function object
	extensible   bool // ### is this needed?
}

const objectDebug = true

func newNumberObject(n float64) value {
	v := value{OBJECT, make([]byte, unsafe.Sizeof(n)), &objectData{NUMBER_OBJECT, value{}, nil, nil, nil, true}}
	*(*float64)(unsafe.Pointer(&v.vdata[0])) = n
	return v
}

func newFunctionObject(call foFn, construct foFn) value {
	v := value{OBJECT, nil, &objectData{FUNCTION_OBJECT, value{}, nil, call, construct, true}}
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
