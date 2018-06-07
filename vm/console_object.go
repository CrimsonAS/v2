package vm

import (
	"log"
)

func defineConsoleObject(vm *vm) value {
	consoleO := newObject()
	consoleO.defineDefaultProperty(vm, "log", newFunctionObject(console_log, nil), 0)
	return consoleO
}

func console_log(vm *vm, f value, args []value) value {
	log.Printf("console.log: %+v", args)
	return newUndefined()
}
