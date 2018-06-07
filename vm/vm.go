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

type stack struct {
	values []value
}

func (this *stack) push(v value) {
	this.values = append(this.values, v)
	//log.Printf("Pushed %s len %d", v, len(this.values))
}

func (this *stack) peek() value {
	return this.values[len(this.values)-1]
}

func (this *stack) popSlice(length int) []value {
	//log.Printf("popSlice %s %d", this.values, length)
	to := len(this.values)
	from := to - length
	ret := this.values[from:to]
	this.values = this.values[:from]
	//log.Printf("popSlice from %d to %d slice len %s values now %s", from, to, ret, this.values)
	return ret
}
func (this *stack) pop() value {
	v := this.peek()
	this.values = this.values[:len(this.values)-1]
	//log.Printf("Popped %s len %d", v, len(this.values))
	return v
}

type stackFrame struct {
	data_stack stack
	retAddr    int
	vars       map[int]*value
	outer      *stackFrame
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

func (this *vm) findVar(name int) *value {
	sf := this.currentFrame
	for sf != nil {
		if v, ok := sf.vars[name]; ok {
			return v
		}
		sf = sf.outer
	}
	return nil
}

func (this *vm) defineVar(name int, v value) {
	this.currentFrame.vars[name] = &v
}

func makeStackFrame(preparedData stack, returnAddr int, outer *stackFrame) stackFrame {
	return stackFrame{data_stack: preparedData, vars: make(map[int]*value), retAddr: returnAddr, outer: outer}
}

func New(code string) *vm {
	ast := parser.Parse(code, true /* ignore comments */)

	vm := vm{[]stackFrame{}, nil, []opcode{}, 0, nil, value{}, false, 0, 0, value{}}
	vm.stack = []stackFrame{makeStackFrame(stack{}, 0, nil)}
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
		this.currentFrame.data_stack.push(rval)
		if execDebug {
			log.Printf("Returning %s up the stack", rval)
			log.Printf("Stack now: %+v", this.stack)
		}
	} else {
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
			this.currentFrame.data_stack.push(newBool(b))
		case NEW_OBJECT:
			o := newObject()
			o.odata.prototype = objectProto
			this.currentFrame.data_stack.push(o)
		case DEFINE_PROPERTY:
			val := this.currentFrame.data_stack.pop()
			key := this.currentFrame.data_stack.pop()
			obj := this.currentFrame.data_stack.peek()
			pd := &propertyDescriptor{name: key.toString(), value: val}
			obj.defineOwnProperty(this, pd.name, pd, false)
		case END_OBJECT:
			this.currentFrame.data_stack.pop()
		case PUSH_UNDEFINED:
			this.currentFrame.data_stack.push(newUndefined())
		case PUSH_NULL:
			this.currentFrame.data_stack.push(newNull())
		case PUSH_NUMBER:
			this.currentFrame.data_stack.push(newNumber(op.odata.asFloat64()))
		case PUSH_STRING:
			this.currentFrame.data_stack.push(newString(stringtable[op.odata.asInt()]))
		case UPLUS:
			val := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newNumber(val.toNumber()))
		case UMINUS:
			expr := this.currentFrame.data_stack.pop()
			oldVal := expr.toNumber()
			if math.IsNaN(oldVal) {
				this.currentFrame.data_stack.push(newNumber(math.NaN()))
			} else {
				this.currentFrame.data_stack.push(newNumber(oldVal * -1))
			}
		case UNOT:
			expr := this.currentFrame.data_stack.pop()
			oldVal := expr.toBoolean()
			if oldVal {
				this.currentFrame.data_stack.push(newBool(false))
			} else {
				this.currentFrame.data_stack.push(newBool(true))
			}
		case INCREMENT:
			v := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newNumber(v.toNumber() + 1))
		case DECREMENT:
			v := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newNumber(v.toNumber() - 1))
		case ADD:
			// ### could (should) specialize this in codegen for numeric types
			// vs unknown types?
			vals := this.currentFrame.data_stack.popSlice(2)
			vals[0] = vals[0].toPrimitive()
			vals[1] = vals[1].toPrimitive()
			if vals[0].vtype == STRING || vals[1].vtype == STRING {
				this.currentFrame.data_stack.push(newString(vals[1].toString() + vals[0].toString()))
			} else {
				this.currentFrame.data_stack.push(newNumber(vals[1].toNumber() + vals[0].toNumber()))
			}
		case SUB:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newNumber(vals[1].toNumber() - vals[0].toNumber()))
		case MULTIPLY:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newNumber(vals[1].toNumber() * vals[0].toNumber()))
		case DIVIDE:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newNumber(vals[1].toNumber() / vals[0].toNumber()))
		case LESS_THAN:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newBool(vals[1].toNumber() < vals[0].toNumber()))
		case GREATER_THAN:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newBool(vals[1].toNumber() > vals[0].toNumber()))
		case EQUALS:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newBool(vals[1].toNumber() == vals[0].toNumber()))
		case NOT_EQUALS:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newBool(vals[1].toNumber() != vals[0].toNumber()))
		case LESS_THAN_EQ:
			vals := this.currentFrame.data_stack.popSlice(2)
			this.currentFrame.data_stack.push(newBool(vals[1].toNumber() <= vals[0].toNumber()))
		case JMP:
			this.ip += op.odata.asInt()
		case JNE:
			test := this.currentFrame.data_stack.pop()

			if op.odata.asInt() == 0 {
				panic("JNE 0 is an infinite loop")
			}
			if !test.toBoolean() {
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
			rval := value{}
			if len(this.currentFrame.data_stack.values) > 0 {
				rval = this.currentFrame.data_stack.pop()
			}

			this.popStack(rval)
		case DUP:
			cv := this.currentFrame.data_stack.peek()
			this.currentFrame.data_stack.push(cv)
		case IN_FUNCTION:
			// no-op, just for informative/debug purposes
		case DECLARE:
			// ### ensure it doesn't exist
			if execDebug {
				log.Printf("Var %s declared (%d)", stringtable[op.odata.asInt()], op.odata.asInt())
			}
			this.currentFrame.vars[op.odata.asInt()] = &value{}
		case STORE:
			v := this.currentFrame.data_stack.pop()
			if execDebug {
				log.Printf("Storing %s in %s (%d)", v, stringtable[op.odata.asInt()], op.odata.asInt())
			}
			sv := this.findVar(op.odata.asInt())
			if sv == nil {
				panic("var not found")
			}
			*sv = v
		case LOAD_MEMBER:
			v := this.currentFrame.data_stack.pop()
			if v.hasPrimitiveBase() {
				// Would be nice if we could do this at codegen time...
				v = v.toObject()
			}
			memb := v.get(this, stringtable[op.odata.asInt()])
			if execDebug {
				log.Printf("LOAD_MEMBER %s.%s got %+v", v, stringtable[op.odata.asInt()], memb)
			}
			this.currentFrame.data_stack.push(memb)
		case LOAD:
			sv := this.findVar(op.odata.asInt())
			if sv == nil {
				panic("var not found")
			}
			if execDebug {
				log.Printf("Loading %s from %d gave %s", stringtable[op.odata.asInt()], op.odata.asInt(), *sv)
			}
			this.thisArg = *sv
			this.currentFrame.data_stack.push(*sv)
		default:
			panic(fmt.Sprintf("unhandled opcode %+v", op))
		}
	}

	return this.returnValue
}

func (this *vm) handleCall(op opcode, isNew bool) {
	// my, this is inefficient
	builtinArgs := this.currentFrame.data_stack.popSlice(op.odata.asInt() + 1)
	fn := builtinArgs[len(builtinArgs)-1]
	if len(builtinArgs) > 1 {
		builtinArgs = builtinArgs[0 : len(builtinArgs)-1]
	} else {
		builtinArgs = builtinArgs[:0]
	}

	if fn.vtype != OBJECT {
		panic(fmt.Sprintf("CALL without a function: %s", fn))
	}

	sf := makeStackFrame(stack{}, this.ip, this.currentFrame)
	this.pushStack(sf)

	var rval value
	if isNew {
		rval = fn.construct(this, this.thisArg, builtinArgs)
	} else {
		rval = fn.call(this, this.thisArg, builtinArgs)
	}

	if this.ignoreReturn {
		this.ignoreReturn = false
	} else {
		this.popStack(rval)
	}
}
