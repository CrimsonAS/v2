package vm

type functionObjectData struct {
	*valueBasicObjectData
	callPtr      foFn
	constructPtr foFn
	prototype    *valueBasicObject
}

func (this *functionObjectData) Prototype() *valueBasicObject {
	return this.prototype
}

func newFunctionObject(call foFn, construct foFn) valueBasicObject {
	v := valueBasicObject{&functionObjectData{&valueBasicObjectData{extensible: true}, call, construct, nil}}
	return v
}

func (this valueBasicObject) call(vm *vm, thisArg value, args []value) value {
	return this.odata.(*functionObjectData).callPtr(vm, thisArg, args)
}

func (this valueBasicObject) construct(vm *vm, thisArg value, args []value) value {
	return this.odata.(*functionObjectData).constructPtr(vm, thisArg, args)
}
