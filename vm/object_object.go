package vm

import (
	"log"
)

var objectProto value

func defineObjectCtor(vm *vm) value {
	objectProto = newObject()
	objectProto.defineDefaultProperty(vm, "toString", newFunctionObject(object_prototype_toString, nil), 0)

	objectCtor := newFunctionObject(object_call, object_ctor)
	objectCtor.odata.prototype = objectProto
	objectProto.set(vm, "constructor", objectCtor)

	return objectCtor
}

func object_call(vm *vm, f value, args []value) value {
	return newBool(args[0].toBoolean())
}

func object_ctor(vm *vm, f value, args []value) value {
	o := newObject()
	o.odata.prototype = objectProto
	return o
}

func object_prototype_toString(vm *vm, f value, args []value) value {
	if f.vtype == UNDEFINED {
		return newString("[object Undefined]")
	} else if f.vtype == NULL {
		return newString("[object Null]")
	}

	o := f.toObject()
	switch o.odata.objectType {
	case OBJECT_PLAIN:
		return newString("[object Object]")
	case BOOLEAN_OBJECT:
		return newString("[object Boolean]")
	case NUMBER_OBJECT:
		return newString("[object Number]")
	case STRING_OBJECT:
		return newString("[object String]")
	case FUNCTION_OBJECT:
		return newString("[object Function]")
	}

	panic("unknown object type")
}

func object_get(vm *vm, f value, prop string, pd *propertyDescriptor) value {
	if objectDebug {
		log.Printf("object_get %s", prop)
	}
	return pd.value
}

func object_set(vm *vm, f value, prop string, pd *propertyDescriptor, v value) value {
	if objectDebug {
		log.Printf("object_set %s", prop)
	}
	pd.value = v
	return newUndefined()
}
