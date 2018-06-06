package vm

func defineStringCtor() value {
	stringO := newFunctionObject(string_call, string_ctor)
	return stringO
}

func string_call(vm *vm, f value, args []value) value {
	return newString(args[0].toString())
}

func string_ctor(vm *vm, f value, args []value) value {
	return newStringObject(args[0].toString())
}
