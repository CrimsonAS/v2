package vm

type functionObjectData struct {
	*valueObjectData
	callPtr      foFn
	constructPtr foFn
	prototype    *valueObject
}

func (this *functionObjectData) Prototype() *valueObject {
	return this.prototype
}

func newFunctionObject(call foFn, construct foFn) valueObject {
	v := valueObject{&functionObjectData{&valueObjectData{extensible: true}, call, construct, nil}}
	return v
}

func (this valueObject) call(vm *vm, thisArg value, args []value) value {
	return this.odata.(*functionObjectData).callPtr(vm, thisArg, args)
}

func (this valueObject) construct(vm *vm, thisArg value, args []value) value {
	return this.odata.(*functionObjectData).constructPtr(vm, thisArg, args)
}
