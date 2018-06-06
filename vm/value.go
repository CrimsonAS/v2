package vm

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"unsafe"
)

type value_type uint8

func (this value_type) String() string {
	switch this {
	case UNDEFINED:
		return "UNDEFINED"
	case NULL:
		return "NULL"
	case BOOL:
		return "BOOL"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case OBJECT:
		return "OBJECT"
	}
	panic("unreachable")
}

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
	odata *objectData
}

func newUndefined() value {
	return value{}
}

func newNull() value {
	return value{NULL, nil, nil}
}

// A simple allocation pool, to help reduce value allocation time.
type vdataPool struct {
	pool []byte
	head int
}

const maxPoolSize = 1024 * 10

func (this *vdataPool) reallocate() {
	this.pool = make([]byte, maxPoolSize)
	this.head = 0
}

const poolDebug = false

func (this *vdataPool) alloc(sz int) []byte {
	if sz >= maxPoolSize {
		panic("allocation exceeds vdata pool bounds; how?!")
	}
	if this.head+sz >= len(this.pool) {
		if poolDebug {
			log.Printf("Reallocating pool")
		}
		this.reallocate()
	}
	ret := this.pool[this.head : this.head+sz]
	this.head += sz
	return ret
}

var allocPool = vdataPool{}

// allocPool is great, but let's avoid even allocating for booleans that are
// used everywhere and absolutely never change...
var trueVData []byte = make([]byte, 1)
var falseVData []byte = make([]byte, 1)

func init() {
	*(*bool)(unsafe.Pointer(&trueVData[0])) = true
	*(*bool)(unsafe.Pointer(&falseVData[0])) = false
}

func newBool(val bool) value {
	if val {
		return value{BOOL, trueVData, nil}
	} else {
		return value{BOOL, falseVData, nil}
	}
}

func newNumber(val float64) value {
	v := value{NUMBER, allocPool.alloc(8), nil}
	*(*float64)(unsafe.Pointer(&v.vdata[0])) = val
	return v
}

func newString(val string) value {
	v := value{STRING, allocPool.alloc(len(val)), nil}
	for i, c := range []byte(val) {
		v.vdata[i] = c
	}
	return v
}

func (this value) asUndefined() value {
	switch this.vtype {
	case UNDEFINED:
		return newUndefined()
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this value) asNull() value {
	switch this.vtype {
	case NULL:
		return newNull()
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this value) asBool() bool {
	switch this.vtype {
	case BOOL:
		return *(*bool)(unsafe.Pointer(&this.vdata[0]))
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this value) asNumber() float64 {
	switch this.vtype {
	case NUMBER:
		return *(*float64)(unsafe.Pointer(&this.vdata[0]))
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
}

func (this value) asString() string {
	switch this.vtype {
	case STRING:
		return string(this.vdata)
	}
	panic(fmt.Sprintf("can't convert! %s", this.vtype))
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
		// not according to ES spec...
		return "[object]"
		//v := this.toPrimitive()
		//return v.toString()
	}

	panic("unreachable")
}

func (this value) String() string {
	return this.toString()
}
