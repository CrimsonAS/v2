package vm

type functionObject struct {
	valueBasicObject
	callPtr      foFn
	constructPtr foFn
	prototype    *valueBasicObject
}

func newFunctionObject(call foFn, construct foFn) functionObject {
	return functionObject{valueBasicObject: newBasicObject(), callPtr: call, constructPtr: construct}
}

func (this *functionObject) call(vm *vm, thisArg value, args []value) value {
	return this.callPtr(vm, thisArg, args)
}

func (this *functionObject) construct(vm *vm, thisArg value, args []value) value {
	return this.constructPtr(vm, thisArg, args)
}

//////////////////////////////////////

func (this *functionObject) Prototype() *valueBasicObject {
	return this.prototype
}
