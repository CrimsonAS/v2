package vm

func defineBooleanCtor() value {
	boolO := newFunctionObject(boolean_call, boolean_ctor)
	return boolO
}

func boolean_call(vm *vm, f value, args []value) value {
	return newBool(args[0].toBoolean())
}

func boolean_ctor(vm *vm, f value, args []value) value {
	return newBooleanObject(args[0].toBoolean())
}
