package vm

import (
	"fmt"
	"log"
	"unsafe"
)

type objectType int

const (
	OBJECT_PLAIN objectType = iota
	BOOLEAN_OBJECT
	NUMBER_OBJECT
	STRING_OBJECT
	FUNCTION_OBJECT
)

func (this value) set(prop string, v value) {
	switch this.vtype {
	case OBJECT:
		os := *(*objectData)(unsafe.Pointer(&this.vdata[0]))
		os.set(prop, v)
		*(*objectData)(unsafe.Pointer(&this.vdata[0])) = os
		return
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this value) get(prop string) value {
	switch this.vtype {
	case OBJECT:
		os := *(*objectData)(unsafe.Pointer(&this.vdata[0]))
		return os.get(prop)
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

type objectData struct {
	objectType objectType
	objects    []value
}

func (this *objectData) get(prop string) value {
	log.Printf("Get Len %d", len(this.objects))
	for i := 0; i < len(this.objects); {
		key := this.objects[i]
		i++
		val := this.objects[i]
		i++

		log.Printf("Searching for %s have %s", prop, key.asString())
		if key.asString() == prop {
			return val
		}
	}

	return newUndefined()
}

func (this *objectData) set(prop string, v value) {
	log.Printf("Set Len %d", len(this.objects))
	for i := 0; i < len(this.objects); {
		key := this.objects[i]
		i++
		//val := this.objects[i]
		i++

		log.Printf("Searching for %s have %s", prop, key.asString())
		if key.asString() == prop {
			this.objects[i-1] = v
			return
		}
	}

	nk := newString(prop)
	this.objects = append(this.objects, nk, v)
	log.Printf("Appended, Len now %d", len(this.objects))
}

type booleanObjectData struct {
	objectData
	primitiveValue bool
}
type numberObjectData struct {
	objectData
	primitiveValue float64
}
type stringObjectData struct {
	objectData
	primitiveValue string
}
type functionObjectData struct {
	objectData
	primitiveValue foFn
}
type foFn func(vm *vm, f value, args []value) value

func newObject() value {
	val := objectData{}

	v := value{OBJECT, make([]byte, unsafe.Sizeof(val))}
	*(*objectData)(unsafe.Pointer(&v.vdata[0])) = val
	return v
}

func newBooleanObject(b bool) value {
	val := booleanObjectData{}
	val.objectType = BOOLEAN_OBJECT
	val.primitiveValue = b

	v := value{OBJECT, make([]byte, unsafe.Sizeof(val))}
	*(*booleanObjectData)(unsafe.Pointer(&v.vdata[0])) = val
	return v
}

func newNumberObject(n float64) value {
	val := numberObjectData{}
	val.objectType = NUMBER_OBJECT
	val.primitiveValue = n

	v := value{OBJECT, make([]byte, unsafe.Sizeof(val))}
	*(*numberObjectData)(unsafe.Pointer(&v.vdata[0])) = val
	return v
}

func newStringObject(s string) value {
	val := stringObjectData{}
	val.objectType = STRING_OBJECT
	val.primitiveValue = s

	log.Printf("Created string with val %s", s)
	v := value{OBJECT, make([]byte, unsafe.Sizeof(val))}
	*(*stringObjectData)(unsafe.Pointer(&v.vdata[0])) = val

	log.Printf("val %s", v.toString())
	return v
}

func newFunctionObject(call foFn) value {
	val := functionObjectData{}
	val.objectType = FUNCTION_OBJECT
	val.primitiveValue = call

	log.Printf("Created function with val %+v", call)
	v := value{OBJECT, make([]byte, unsafe.Sizeof(val))}
	*(*functionObjectData)(unsafe.Pointer(&v.vdata[0])) = val

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
	fod := *(*functionObjectData)(unsafe.Pointer(&this.vdata[0]))
	if fod.objectType != FUNCTION_OBJECT {
		panic(fmt.Sprintf("can't call non-function! %d", fod.objectType))
	}

	return fod.primitiveValue(vm, this, args)
}
