package vm

import (
	"log"
)

func defineConsoleObject() value {
	consoleO := newObject()
	consoleO.set("log", newFunctionObject(console_log, nil))
	return consoleO
}

func console_log(vm *vm, f value, args []value) value {
	log.Printf("console.log: %+v", args)
	return newUndefined()
}
