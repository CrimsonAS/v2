package vm

type numberObjectData struct {
	*valueObjectData
	primitiveData value
}

func (this *numberObjectData) Prototype() *valueObject {
	return nil
}

func newNumberObject(n float64) valueObject {
	v := valueObject{&numberObjectData{&valueObjectData{extensible: true}, newNumber(n)}}
	return v
}
