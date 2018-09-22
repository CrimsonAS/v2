/*
 * Copyright 2018 Crimson AS <info@crimson.no>
 * Author: Robin Burchell <robin.burchell@crimson.no>
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package vm

import (
	"fmt"
	"github.com/CrimsonAS/v2/parser"
	"log"
	"math"
)

type stackFrame struct {
	retAddr int
	vars    []int
	// ### ideally we would reserve space for these inside data_stack
	varValues   []value
	temporaries []value
	outer       *stackFrame
	thisArg     value
}

var stringtable []string

func appendStringtable(name string) int {
	for idx, str := range stringtable {
		if name == str {
			return idx
		}
	}
	stringtable = append(stringtable, name)
	return len(stringtable) - 1
}

type vm struct {
	data_stack    stack
	stack         []stackFrame
	currentFrame  *stackFrame
	code          []opcode
	ip            int
	funcsToDefine []*parser.FunctionExpression // codegen
	returnValue   value
	ignoreReturn  bool
	isNew         int
	canConsume    int
	lastLoadedVar value

	// from codegen
	temporaryIndex int
}

const lookupDebug = false

func (this *vm) setVar(name int, nv value) bool {
	if execDebug {
		log.Printf("Storing %s in %s", nv, stringtable[name])
	}
	sf := this.currentFrame
	for sf != nil {
		for idx, sfvar := range sf.vars {
			if sfvar == name {
				//log.Printf("Set var %d to %+v", name, nv)
				sf.varValues[idx] = nv
				return true
			}
		}
		sf = sf.outer
	}
	return false
}

func (this *vm) findVar(name int) (value, bool) {
	sf := this.currentFrame
	for sf != nil {
		for idx, sfvar := range sf.vars {
			if sfvar == name {
				if execDebug {
					log.Printf("Loading %s gave %s", stringtable[name], sf.varValues[idx])
				}
				return sf.varValues[idx], true
			}
		}
		sf = sf.outer
	}
	if execDebug {
		log.Printf("Loading %s was not found", stringtable[name])
	}
	return nil, false
}

func (this *vm) defineVar(name int, v value) {
	for _, sfvar := range this.currentFrame.vars {
		if sfvar == name {
			//panic(fmt.Sprintf("Var %s already defined", stringtable[name]))
			return
		}
	}

	if execDebug {
		log.Printf("Var %s declared", stringtable[name])
	}
	this.currentFrame.vars = append(this.currentFrame.vars, name)
	this.currentFrame.varValues = append(this.currentFrame.varValues, v)
}

func makeStackFrame(thisArg value, returnAddr int, outer *stackFrame) stackFrame {
	return stackFrame{retAddr: returnAddr, outer: outer, thisArg: thisArg}
}

func New(code string) *vm {
	ast := parser.Parse(code, true /* ignore comments */)

	vm := vm{stack{}, []stackFrame{}, nil, []opcode{}, 0, nil, nil, false, 0, 0, nil, -1}
	vm.stack = []stackFrame{makeStackFrame(newUndefined(), 0, nil)}
	vm.currentFrame = &vm.stack[0]

	il := []tac{}
	vm.generateCodeTAC(ast, &il)
	optimizeTAC(&il)

	if execDebug {
		for idx, op := range il {
			log.Printf("%d: %s", idx, op)
		}
	}

	vm.code = vm.generateBytecode(il)

	if execDebug {
		vm.DumpCode()
	}

	vm.defineVar(appendStringtable("Object"), defineObjectCtor(&vm))
	vm.defineVar(appendStringtable("console"), defineConsoleObject(&vm))
	vm.defineVar(appendStringtable("Math"), defineMathObject(&vm))
	vm.defineVar(appendStringtable("Boolean"), defineBooleanCtor(&vm))
	vm.defineVar(appendStringtable("Array"), defineArrayCtor(&vm))
	vm.defineVar(appendStringtable("String"), defineStringCtor(&vm))

	return &vm
}

const execDebug = false

func (this *vm) pushStack(sf stackFrame) {
	this.stack = append(this.stack, sf)
	this.currentFrame = &this.stack[len(this.stack)-1]

	if execDebug {
		log.Printf("Pushed stack. Stack now: %+v", this.stack)
	}
}

func (this *vm) popStack(rval value) {
	this.stack = this.stack[:len(this.stack)-1]

	if len(this.stack) > 0 {
		this.ip = this.currentFrame.retAddr
		this.currentFrame = &this.stack[len(this.stack)-1]
		this.data_stack.push(rval)
		if execDebug {
			log.Printf("Returning %s up the stack", rval)
			log.Printf("Stack now: %+v", this.stack)
		}
	} else {
		if rval == nil {
			rval = newUndefined()
		}
		this.returnValue = rval
		if execDebug {
			log.Printf("Returning %s from Run()", rval)
		}
	}
}

func (this *vm) DumpCode() {
	log.Printf("String table:")
	for i := 0; i < len(stringtable); i++ {
		log.Printf("%d: %s", i, stringtable[i])
	}
	log.Printf("Program:")
	for i := 0; i < len(this.code); i++ {
		log.Printf("%d: %s", i, this.code[i])
	}
}

func (this *vm) ThrowTypeError(msg string) value {
	if msg != "" {
		panic(fmt.Sprintf("TypeError: %s", msg))
	}
	panic("TypeError")
}

func (this *vm) Run() value {
	for ; len(this.stack) > 0 && this.ip < len(this.code); this.ip++ {
		op := this.code[this.ip]
		if execDebug {
			log.Printf("Op %d: %s (stack: %+v §§ temporaries %+v)", this.ip, op, this.data_stack, this.currentFrame.temporaries)
		}
		switch op.otype {
		case PUSH_BOOL:
			b := false
			if op.opdata != 0 {
				b = true
			}
			this.data_stack.push(newBool(b))
		case NEW_OBJECT:
			o := newBasicObject()
			this.data_stack.push(o)
		case DEFINE_PROPERTY:
			val := this.data_stack.pop()
			key := this.data_stack.pop()
			obj := this.data_stack.peek().(valueObject)
			pn := key.ToString()
			pd := &propertyDescriptor{name: pn.String(), value: val, hasValue: true, writable: true, hasWritable: true, configurable: true, hasConfigurable: true}
			obj.defineOwnProperty(this, pn, pd, false)
		case END_OBJECT:
			this.data_stack.pop()
		case PUSH_UNDEFINED:
			this.data_stack.push(newUndefined())
		case PUSH_NULL:
			this.data_stack.push(newNull())
		case PUSH_ARRAY:
			vals := this.data_stack.popSlice(op.opdata.asInt())
			this.data_stack.push(newArrayObject(vals))
		case PUSH_NUMBER:
			this.data_stack.push(newNumber(op.opdata.asFloat64()))
		case PUSH_STRING:
			this.data_stack.push(newString(stringtable[op.opdata.asInt()]))
		case UPLUS:
			val := this.data_stack.pop()
			this.data_stack.push(newNumber(val.ToNumber()))
		case UMINUS:
			expr := this.data_stack.pop()
			oldVal := expr.ToNumber()
			if math.IsNaN(oldVal) {
				this.data_stack.push(newNumber(math.NaN()))
			} else {
				this.data_stack.push(newNumber(oldVal * -1))
			}
		case UNOT:
			expr := this.data_stack.pop()
			oldVal := expr.ToBoolean()
			if oldVal {
				this.data_stack.push(newBool(false))
			} else {
				this.data_stack.push(newBool(true))
			}
		case INCREMENT:
			v := this.data_stack.pop()
			this.data_stack.push(newNumber(v.ToNumber() + 1))
		case DECREMENT:
			v := this.data_stack.pop()
			this.data_stack.push(newNumber(v.ToNumber() - 1))
		case ADD:
			// ### could (should) specialize this in codegen for numeric types
			vals := this.data_stack.popSlice(2)
			vals[0] = valueToPrimitive(vals[0])
			vals[1] = valueToPrimitive(vals[1])

			oneIsString := false
			switch vals[0].(type) {
			case valueString:
				oneIsString = true
			}
			switch vals[1].(type) {
			case valueString:
				oneIsString = true
			}
			if oneIsString {
				this.data_stack.push(newString(vals[1].String() + vals[0].String()))
			} else {
				this.data_stack.push(newNumber(vals[1].ToNumber() + vals[0].ToNumber()))
			}
		case SUB:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(vals[1].ToNumber() - vals[0].ToNumber()))
		case MULTIPLY:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(vals[1].ToNumber() * vals[0].ToNumber()))
		case DIVIDE:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(vals[1].ToNumber() / vals[0].ToNumber()))
		case MODULUS:
			vals := this.data_stack.popSlice(2)
			// ### using math is probably going to hurt performance?
			this.data_stack.push(newNumber(math.Mod(vals[1].ToNumber(), vals[0].ToNumber())))
		case LEFT_SHIFT:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(float64(vals[1].ToInteger() << uint(vals[0].ToInteger()))))
		case RIGHT_SHIFT:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(float64(vals[1].ToInteger() >> uint(vals[0].ToInteger()))))
		case UNSIGNED_RIGHT_SHIFT:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(float64(uint32(vals[1].ToInteger()) >> uint(vals[0].ToInteger()))))
		case BITWISE_AND:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(float64(vals[1].ToInteger() & vals[0].ToInteger())))
		case BITWISE_XOR:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(float64(vals[1].ToInteger() ^ vals[0].ToInteger())))
		case BITWISE_OR:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newNumber(float64(vals[1].ToInteger() | vals[0].ToInteger())))
		case BITWISE_NOT:
			v := this.data_stack.pop()
			this.data_stack.push(newNumber(float64(^v.ToInteger())))
		case LESS_THAN:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() < vals[0].ToNumber()))
		case GREATER_THAN:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() > vals[0].ToNumber()))
		case GREATER_THAN_EQ:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() >= vals[0].ToNumber()))
		case EQUALS:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() == vals[0].ToNumber()))
		case NOT_EQUALS:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() != vals[0].ToNumber()))
		case STRICT_EQUALS:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(strictEqualityComparison(vals[1], vals[0])))
		case STRICT_NOT_EQUALS:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(!strictEqualityComparison(vals[1], vals[0])))
		case LESS_THAN_EQ:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() <= vals[0].ToNumber()))
		case LOGICAL_AND:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToBoolean() && vals[0].ToBoolean()))
		case POP:
			this.data_stack.pop()
		case JMP:
			this.ip += op.opdata.asInt()
		case JNE:
			test := this.data_stack.pop()

			if op.opdata.asInt() == 0 {
				panic("JNE 0 is an infinite loop")
			}
			if !test.ToBoolean() {
				if execDebug {
					log.Printf("IP is at %d jump rel %d code length %d", this.ip, op.opdata.asInt(), len(this.code))
				}
				this.ip += op.opdata.asInt()

				if this.ip >= len(this.code) {
					panic("JNE blew over opcode length")
				}
			}
		case CALL:
			this.handleCall(op, false)
		case NEW:
			this.handleCall(op, true)
		case RETURN:
			// can't inline this to popStack, because the builtin case doesn't
			// have a value pushed onto the data_stack.
			var rval value
			if len(this.data_stack.values) > 0 {
				rval = this.data_stack.pop()
			}

			this.popStack(rval)
		case DUP:
			cv := this.data_stack.peek()
			this.data_stack.push(cv)
		case IN_FUNCTION:
			// no-op, just for informative/debug purposes
		case DECLARE:
			this.defineVar(op.opdata.asInt(), nil)
		case STORE:
			v := this.data_stack.pop()
			ok := this.setVar(op.opdata.asInt(), v)
			if !ok {
				panic(fmt.Sprintf("var %s not found", stringtable[op.opdata.asInt()]))
			}
		case STORE_MEMBER:
			v := this.data_stack.pop()
			nv := this.data_stack.pop()
			var vo valueObject
			if v.hasPrimitiveBase() {
				// Would be nice if we could do this at codegen time...
				vo = v.ToObject()
			} else {
				vo = v.(valueObject)
			}
			vo.put(this, newString(stringtable[op.opdata.asInt()]), nv, true)
		case LOAD_MEMBER:
			v := this.data_stack.pop()
			var vo valueObject
			if v.hasPrimitiveBase() {
				// Would be nice if we could do this at codegen time...
				vo = v.ToObject()
			} else {
				vo = v.(valueObject)
			}
			this.data_stack.push(vo.get(this, newString(stringtable[op.opdata.asInt()])))
		case LOAD_INDEXED:
			v := this.data_stack.pop()
			idx := this.data_stack.pop().ToInteger()
			var vo valueObject
			if v.hasPrimitiveBase() {
				// Would be nice if we could do this at codegen time...
				vo = v.ToObject()
			} else {
				vo = v.(valueObject)
			}

			this.data_stack.push(vo.get(this, newNumber(float64(idx))))
		case STORE_INDEXED:
			v := this.data_stack.pop()
			idx := this.data_stack.pop().ToInteger()
			var vo valueObject
			if v.hasPrimitiveBase() {
				// Would be nice if we could do this at codegen time...
				vo = v.ToObject()
			} else {
				vo = v.(valueObject)
			}

			nv := this.data_stack.pop()
			vo.put(this, newNumber(float64(idx)), nv, true)
		case LOAD:
			sv, ok := this.findVar(op.opdata.asInt())
			if !ok {
				panic("var not found")
			}
			this.lastLoadedVar = sv
			this.data_stack.push(sv)
		case LOAD_THIS:
			// ### 'this' should be valid in global contexts too, but isn't currently.
			if this.currentFrame.thisArg == nil {
				panic("'this' in global context not yet supported...")
			}
			this.data_stack.push(this.currentFrame.thisArg)
		case LOAD_TEMPORARY:
			idx := op.opdata.asInt()
			if idx >= len(this.currentFrame.temporaries) {
				this.data_stack.push(newUndefined())
			} else {
				this.data_stack.push(this.currentFrame.temporaries[idx])
			}
		case STORE_TEMPORARY:
			idx := op.opdata.asInt()
			for len(this.currentFrame.temporaries) <= idx {
				this.currentFrame.temporaries = append(this.currentFrame.temporaries, newUndefined())
			}
			this.currentFrame.temporaries[idx] = this.data_stack.pop()
		case TYPEOF:
			v := this.data_stack.pop()

			// ### does the value interface need another member?
			switch v.(type) {
			case valueUndefined:
				this.data_stack.push(newString("undefined"))
			case valueNull:
				this.data_stack.push(newString("object"))
			case valueBool:
				this.data_stack.push(newString("boolean"))
			case valueNumber:
				this.data_stack.push(newString("number"))
			case valueString:
				this.data_stack.push(newString("string"))
			case functionObject:
				this.data_stack.push(newString("function"))
			case valueBasicObject:
				this.data_stack.push(newString("object"))
			case arrayObject:
				this.data_stack.push(newString("object"))
			default:
				panic("Unknown type")
			}
		default:
			panic(fmt.Sprintf("unhandled opcode %+v", op))
		}
	}

	return this.returnValue
}

func (this *vm) handleCall(op opcode, isNew bool) {
	// my, this is inefficient
	builtinArgs := this.data_stack.popSlice(op.opdata.asInt() + 1)

	fn := builtinArgs[len(builtinArgs)-1]
	if len(builtinArgs) > 1 {
		builtinArgs = builtinArgs[0 : len(builtinArgs)-1]
	} else {
		builtinArgs = builtinArgs[:0]
	}

	fo := fn.(functionObject)

	sf := makeStackFrame(this.lastLoadedVar, this.ip, this.currentFrame)
	this.pushStack(sf)

	var rval value
	if isNew {
		rval = fo.construct(this, this.lastLoadedVar, builtinArgs)
	} else {
		rval = fo.call(this, this.lastLoadedVar, builtinArgs)
	}

	if this.ignoreReturn {
		this.ignoreReturn = false
	} else {
		this.popStack(rval)
	}
}
