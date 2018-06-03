package vm

import (
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

type value_type int

const (
	UNDEFINED value_type = iota
	NULL
	BOOL
	NUMBER
	STRING
	OBJECT
)

type value struct {
	vtype value_type
	vdata []byte // horrible un-type-safe voodoo
}

func newUndefined() value {
	return value{}
}

func newNull() value {
	return value{UNDEFINED, nil}
}

func newBool(b bool) value {
	v := value{BOOL, make([]byte, 8)}
	*(*bool)(unsafe.Pointer(&v.vdata[0])) = b
	return v
}

func newNumber(val float64) value {
	v := value{NUMBER, make([]byte, 8)}
	*(*float64)(unsafe.Pointer(&v.vdata[0])) = val
	return v
}

func newString(val string) value {
	v := value{STRING, []byte(val)}
	return v
}

func newObject() value {
	v := value{OBJECT, nil}
	return v
}

func (this value) asUndefined() value {
	switch this.vtype {
	case UNDEFINED:
		return value{}
	}
	panic("can't convert!")
}

func (this value) asNull() value {
	switch this.vtype {
	case NULL:
		return newNull()
	}
	panic("can't convert!")
}

func (this value) asBool() bool {
	switch this.vtype {
	case BOOL:
		return *(*bool)(unsafe.Pointer(&this.vdata[0]))
	}
	panic("can't convert!")
}

func (this value) asNumber() float64 {
	switch this.vtype {
	case NUMBER:
		return *(*float64)(unsafe.Pointer(&this.vdata[0]))
	}
	panic("can't convert!")
}

func (this value) asString() string {
	switch this.vtype {
	case STRING:
		return string(this.vdata)
	}
	panic("can't convert!")
}

/*
func (this value) asObject() float64 {
*/

func (this value) toPrimitive() value {
	switch this.vtype {
	case UNDEFINED:
		fallthrough
	case NULL:
		fallthrough
	case BOOL:
		fallthrough
	case NUMBER:
		fallthrough
	case STRING:
		return this
	case OBJECT:
		panic("object conversion not implemented")
	}
	panic("unreachable")
}

func (this value) toBoolean() bool {
	switch this.vtype {
	case UNDEFINED:
		fallthrough
	case NULL:
		return false
	case BOOL:
		return this.asBool()
	case NUMBER:
		n := this.asNumber()
		if int(n) == 0 || math.IsNaN(n) {
			return false
		} else {
			return true
		}
	case STRING:
		s := this.asString()
		return len(s) > 0
	case OBJECT:
		return true
	}
	panic("unreachable")
}

func (this value) toNumber() float64 {
	switch this.vtype {
	case UNDEFINED:
		return math.NaN()
	case NULL:
		return +0
	case BOOL:
		if this.asBool() {
			return 1
		} else {
			return +0
		}
	case NUMBER:
		return this.asNumber()
	case STRING:
		// ### toNumber(string) not implemented (es5 9.3.1)
		v, _ := strconv.ParseFloat(this.asString(), 64)
		return v
	case OBJECT:
		v := this.toPrimitive()
		return v.toNumber()
	}
	panic("unreachable")
}

func (this value) toInteger() float64 { // ### return type ok?
	number := this.toNumber()
	if math.IsNaN(number) {
		return +0
	}
	if int(number) == 0 || math.IsInf(number, 0) {
		return number
	}

	if number > 0 {
		return 1 * math.Floor(math.Abs(number))
	} else {
		return -1 * math.Floor(math.Abs(number))
	}
}

func (this value) toInt32() int64 {
	panic("not implemented (es5 9.5)")
}

func (this value) toUInt32() int64 {
	panic("not implemented (es5 9.6)")
}

func (this value) toUInt16() int64 {
	panic("not implemented (es5 9.7)")
}

func (this value) toString() string {
	switch this.vtype {
	case UNDEFINED:
		return "undefined"
	case NULL:
		return "null"
	case BOOL:
		if this.asBool() {
			return "true"
		} else {
			return "false"
		}
	case NUMBER:
		// may be wrong, check es5 9.8.1
		return fmt.Sprintf("%f", this.asNumber())
	case STRING:
		return this.asString()
	case OBJECT:
		v := this.toPrimitive()
		return v.toString()
	}

	panic("unreachable")
}

func (this value) toObject() value {
	switch this.vtype {
	case UNDEFINED:
		panic("TypeError")
	case NULL:
		panic("TypeError")
	case BOOL:
	case NUMBER:
	case STRING:
	case OBJECT:
	}
	panic("toObject() not implemented (es5 9.9)")
}

func (this value) String() string {
	return this.toString()
}
