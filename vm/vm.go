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
	vars    map[int]value
	outer   *stackFrame
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
	funcsToDefine []*parser.FunctionExpression
	returnValue   value
	ignoreReturn  bool
	isNew         int
	canConsume    int
	thisArg       value
}

const lookupDebug = false

func (this *vm) setVar(name int, nv value) bool {
	if execDebug {
		//log.Printf("Storing %s in %s", v, stringtable[op.odata.asInt()])
	}
	sf := this.currentFrame
	for sf != nil {
		if _, ok := sf.vars[name]; ok {
			sf.vars[name] = nv
			//log.Printf("Set var %d to %+v", name, nv)
			return true
		}
		sf = sf.outer
	}
	return false
}

func (this *vm) findVar(name int) (value, bool) {
	if execDebug {
		//log.Printf("Loading %s from %d gave %s", stringtable[op.odata.asInt()], op.odata.asInt(), sv)
	}
	sf := this.currentFrame
	for sf != nil {
		if v, ok := sf.vars[name]; ok {
			return v, true
		}
		sf = sf.outer
	}
	return nil, false
}

func (this *vm) defineVar(name int, v value) {
	// ### ensure it doesn't exist
	if execDebug {
		//log.Printf("Var %s declared (%d)", stringtable[op.odata.asInt()], op.odata.asInt())
	}
	this.currentFrame.vars[name] = v
}

func makeStackFrame(returnAddr int, outer *stackFrame) stackFrame {
	return stackFrame{vars: make(map[int]value), retAddr: returnAddr, outer: outer}
}

func New(code string) *vm {
	ast := parser.Parse(code, true /* ignore comments */)

	vm := vm{stack{}, []stackFrame{}, nil, []opcode{}, 0, nil, nil, false, 0, 0, nil}
	vm.stack = []stackFrame{makeStackFrame(0, nil)}
	vm.currentFrame = &vm.stack[0]
	vm.code = vm.generateCode(ast)

	vm.defineVar(appendStringtable("Object"), defineObjectCtor(&vm))
	vm.defineVar(appendStringtable("console"), defineConsoleObject(&vm))
	vm.defineVar(appendStringtable("Math"), defineMathObject(&vm))
	vm.defineVar(appendStringtable("Boolean"), defineBooleanCtor(&vm))
	vm.defineVar(appendStringtable("String"), defineStringCtor(&vm))

	return &vm
}

const codegenDebug = false
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
			rval = newUndefined() // ### are we missing a push somewhere?
		}
		this.returnValue = rval
		if execDebug {
			log.Printf("Returning %s from Run()", rval)
		}

	}
}

func (this *vm) Run() value {
	if codegenDebug {
		log.Printf("String table:")
		for i := 0; i < len(stringtable); i++ {
			log.Printf("%d: %s", i, stringtable[i])
		}
		log.Printf("Program:")
		for i := 0; i < len(this.code); i++ {
			log.Printf("%d: %s", i, this.code[i])
		}
		log.Printf("Starting execution")
	}

	for ; len(this.stack) > 0 && this.ip < len(this.code); this.ip++ {
		op := this.code[this.ip]
		if execDebug {
			log.Printf("Op %d: %s", this.ip, op)
		}
		switch op.otype {
		case PUSH_BOOL:
			b := false
			if op.odata != 0 {
				b = true
			}
			this.data_stack.push(newBool(b))
		case NEW_OBJECT:
			o := newObject()
			o.odata.prototype = &objectProto // ### belongs in newObject?
			this.data_stack.push(o)
		case DEFINE_PROPERTY:
			val := this.data_stack.pop()
			key := this.data_stack.pop()
			obj := this.data_stack.peek().(valueObject)
			pd := &propertyDescriptor{name: key.String(), value: val}
			obj.defineOwnProperty(this, pd.name, pd, false)
		case END_OBJECT:
			this.data_stack.pop()
		case PUSH_UNDEFINED:
			this.data_stack.push(newUndefined())
		case PUSH_NULL:
			this.data_stack.push(newNull())
		case PUSH_NUMBER:
			this.data_stack.push(newNumber(op.odata.asFloat64()))
		case PUSH_STRING:
			this.data_stack.push(newString(stringtable[op.odata.asInt()]))
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
		case LESS_THAN:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() < vals[0].ToNumber()))
		case GREATER_THAN:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() > vals[0].ToNumber()))
		case EQUALS:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() == vals[0].ToNumber()))
		case NOT_EQUALS:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() != vals[0].ToNumber()))
		case LESS_THAN_EQ:
			vals := this.data_stack.popSlice(2)
			this.data_stack.push(newBool(vals[1].ToNumber() <= vals[0].ToNumber()))
		case JMP:
			this.ip += op.odata.asInt()
		case JNE:
			test := this.data_stack.pop()

			if op.odata.asInt() == 0 {
				panic("JNE 0 is an infinite loop")
			}
			if !test.ToBoolean() {
				if execDebug {
					log.Printf("IP is at %d jump rel %d code length %d", this.ip, op.odata.asInt(), len(this.code))
				}
				this.ip += op.odata.asInt()

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
			this.defineVar(op.odata.asInt(), nil)
		case STORE:
			v := this.data_stack.pop()
			ok := this.setVar(op.odata.asInt(), v)
			if !ok {
				panic("var not found")
			}
		case LOAD_MEMBER:
			v := this.data_stack.pop()
			var vo valueObject
			if v.hasPrimitiveBase() {
				// Would be nice if we could do this at codegen time...
				vo = v.ToObject()
			} else {
				vo = v.(valueObject)
			}
			memb := vo.get(this, stringtable[op.odata.asInt()])
			if execDebug {
				log.Printf("LOAD_MEMBER %s.%s got %+v", vo, stringtable[op.odata.asInt()], memb)
			}
			this.data_stack.push(memb)
		case LOAD:
			sv, ok := this.findVar(op.odata.asInt())
			if !ok {
				panic("var not found")
			}
			this.thisArg = sv
			this.data_stack.push(sv)
		default:
			panic(fmt.Sprintf("unhandled opcode %+v", op))
		}
	}

	return this.returnValue
}

func (this *vm) handleCall(op opcode, isNew bool) {
	// my, this is inefficient
	builtinArgs := this.data_stack.popSlice(op.odata.asInt() + 1)

	fn := builtinArgs[len(builtinArgs)-1]
	if len(builtinArgs) > 1 {
		builtinArgs = builtinArgs[0 : len(builtinArgs)-1]
	} else {
		builtinArgs = builtinArgs[:0]
	}

	fo := fn.(valueObject)

	sf := makeStackFrame(this.ip, this.currentFrame)
	this.pushStack(sf)

	var rval value
	if isNew {
		rval = fo.construct(this, this.thisArg, builtinArgs)
	} else {
		rval = fo.call(this, this.thisArg, builtinArgs)
	}

	if this.ignoreReturn {
		this.ignoreReturn = false
	} else {
		this.popStack(rval)
	}
}
