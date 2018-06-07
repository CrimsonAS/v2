package vm

import (
	"fmt"
)

var booleanProto value

func newBooleanObject(b bool) value {
	v := newBool(b)
	v.vtype = OBJECT
	v.odata = &objectData{BOOLEAN_OBJECT, value{}, nil, nil, nil}
	v.odata.prototype = booleanProto
	return v
}

func defineBooleanCtor(vm *vm) value {
	booleanProto = newObject()
	booleanProto.defineDefaultProperty(vm, "toString", newFunctionObject(boolean_prototype_toString, nil), 0)
	booleanProto.defineDefaultProperty(vm, "valueOf", newFunctionObject(boolean_prototype_valueOf, nil), 0)

	boolO := newFunctionObject(boolean_call, boolean_ctor)
	boolO.odata.prototype = booleanProto

	booleanProto.defineDefaultProperty(vm, "constructor", boolO, 0)
	return boolO
}

func boolean_call(vm *vm, f value, args []value) value {
	return newBool(args[0].toBoolean())
}

func boolean_ctor(vm *vm, f value, args []value) value {
	return newBooleanObject(args[0].toBoolean())
}

func boolean_prototype_toString(vm *vm, f value, args []value) value {
	switch f.vtype {
	case BOOL:
		break
	case OBJECT:
		if f.odata.objectType == BOOLEAN_OBJECT {
			break
		}
		fallthrough
	default:
		panic(fmt.Sprintf("Not a boolean! %s", f)) // ### throw
	}

	if f.asBool() {
		return newString("true")
	} else {
		return newString("false")
	}
}

func boolean_prototype_valueOf(vm *vm, f value, args []value) value {
	switch f.vtype {
	case BOOL:
		break
	case OBJECT:
		if f.odata.objectType == BOOLEAN_OBJECT {
			break
		}
		fallthrough
	default:
		panic(fmt.Sprintf("Not a boolean! %s", f)) // ### throw
	}
	return newBool(f.asBool())
}
