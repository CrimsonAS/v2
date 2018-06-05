package vm

import (
	"fmt"
	"log"
	"unsafe"
)

type objectType int

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
	objectType objectType
	objects    []value
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
	v := value{OBJECT, nil, nil, &objectData{OBJECT_PLAIN, nil}}
	return v
}

func newBooleanObject(b bool) value {
	v := value{OBJECT, make([]byte, unsafe.Sizeof(b)), nil, &objectData{BOOLEAN_OBJECT, nil}}
	*(*bool)(unsafe.Pointer(&v.vdata[0])) = b
	return v
}

func newNumberObject(n float64) value {
	v := value{OBJECT, make([]byte, unsafe.Sizeof(n)), nil, &objectData{NUMBER_OBJECT, nil}}
	*(*float64)(unsafe.Pointer(&v.vdata[0])) = n
	return v
}

func newStringObject(s string) value {
	v := value{OBJECT, []byte(s), nil, &objectData{STRING_OBJECT, nil}}
	return v
}

func newFunctionObject(call foFn) value {
	v := value{OBJECT, nil, call, &objectData{FUNCTION_OBJECT, nil}}
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
	return this.fptr(vm, this, args)
}
