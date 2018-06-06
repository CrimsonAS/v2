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

func (this *value) set(prop string, v value) {
	switch this.vtype {
	case OBJECT:
		this.odata.set(prop, v)
		return
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this *value) get(prop string) value {
	switch this.vtype {
	case OBJECT:
		return this.odata.get(prop)
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

type objectData struct {
	objectType   objectType
	objects      []value
	callPtr      foFn // used for function object
	constructPtr foFn // used for function object
}

const objectDebug = false

func (this *objectData) get(prop string) value {
	if objectDebug {
		log.Printf("Get Len %d", len(this.objects))
	}
	for i := 0; i < len(this.objects); {
		key := this.objects[i]
		i++
		val := this.objects[i]
		i++

		if objectDebug {
			log.Printf("Searching for %s have %s", prop, key.asString())
		}
		if key.asString() == prop {
			return val
		}
	}

	return newUndefined()
}

func (this *objectData) set(prop string, v value) {
	if objectDebug {
		log.Printf("Set Len %d", len(this.objects))
	}
	for i := 0; i < len(this.objects); {
		key := this.objects[i]
		i++
		//val := this.objects[i]
		i++

		if objectDebug {
			log.Printf("Searching for %s have %s", prop, key.asString())
		}
		if key.asString() == prop {
			this.objects[i-1] = v
			return
		}
	}

	nk := newString(prop)
	this.objects = append(this.objects, nk, v)
	if objectDebug {
		log.Printf("Appended, Len now %d", len(this.objects))
	}
}

type foFn func(vm *vm, f value, args []value) value

func newObject() value {
	v := value{OBJECT, nil, &objectData{OBJECT_PLAIN, nil, nil, nil}}
	return v
}

func newBooleanObject(b bool) value {
	v := value{OBJECT, make([]byte, unsafe.Sizeof(b)), &objectData{BOOLEAN_OBJECT, nil, nil, nil}}
	*(*bool)(unsafe.Pointer(&v.vdata[0])) = b
	return v
}

func newNumberObject(n float64) value {
	v := value{OBJECT, make([]byte, unsafe.Sizeof(n)), &objectData{NUMBER_OBJECT, nil, nil, nil}}
	*(*float64)(unsafe.Pointer(&v.vdata[0])) = n
	return v
}

func newStringObject(s string) value {
	v := value{OBJECT, []byte(s), &objectData{STRING_OBJECT, nil, nil, nil}}
	return v
}

func newFunctionObject(call foFn, construct foFn) value {
	v := value{OBJECT, nil, &objectData{FUNCTION_OBJECT, nil, call, construct}}
	return v
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

func (this value) call(vm *vm, args []value) value {
	if this.vtype != OBJECT {
		panic(fmt.Sprintf("can't convert! %s", this.vtype))
	}
	if this.odata.objectType != FUNCTION_OBJECT {
		panic(fmt.Sprintf("can't call non-function! %d", this.odata.objectType))
	}
	return this.odata.callPtr(vm, this, args)
}

func (this value) construct(vm *vm, args []value) value {
	if this.vtype != OBJECT {
		panic(fmt.Sprintf("can't convert! %s", this.vtype))
	}
	if this.odata.objectType != FUNCTION_OBJECT {
		panic(fmt.Sprintf("can't call non-function! %d", this.odata.objectType))
	}
	return this.odata.constructPtr(vm, this, args)
}
