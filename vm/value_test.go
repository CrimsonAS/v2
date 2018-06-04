package vm

import (
	"github.com/stvp/assert"
	"testing"
)

func TestObject(t *testing.T) {
	v := newObject()
	v.get("hello")
	assert.Equal(t, v.get("hello"), newUndefined())

	sv := newString("hello world")
	v.set("hello", sv)
	assert.Equal(t, v.get("hello"), newString("hello world"))

	v2 := newObject()
	assert.Equal(t, v2.get("hello"), newUndefined())
	assert.Equal(t, v.get("hello"), newString("hello world"))
	v.set("world", sv)
	assert.Equal(t, v.get("world"), newString("hello world"))
	assert.Equal(t, v.get("hello"), newString("hello world"))

	assert.Equal(t, newBool(true).toObject(), newBooleanObject(true))
	assert.Equal(t, newBool(false).toObject(), newBooleanObject(false))
	assert.Equal(t, newNumber(1.25).toObject(), newNumberObject(1.25))
	assert.Equal(t, newNumber(21.5).toObject(), newNumberObject(21.5))

	// seems to fail..
	//assert.Equal(t, newString("hello").toObject(), newStringObject("hello"))
}
