package vm

type numberObject struct {
	valueBasicObject
	primitiveData float64
}

func (this *numberObject) Prototype() *valueBasicObject {
	return nil
}

func newNumberObject(n float64) numberObject {
	return numberObject{valueBasicObject: newBasicObject(), primitiveData: n}
}
