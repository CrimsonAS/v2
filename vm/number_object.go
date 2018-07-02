package vm

type numberObjectData struct {
	*valueBasicObjectData
	primitiveData value
}

func (this *numberObjectData) Prototype() *valueBasicObject {
	return nil
}

func newNumberObject(n float64) valueBasicObject {
	v := valueBasicObject{&numberObjectData{&valueBasicObjectData{extensible: true}, newNumber(n)}}
	return v
}
