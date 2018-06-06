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

const stackDebug = false

func (this *stack) push(v value) {
	if stackDebug {
		log.Printf("Pushing %s onto stack", v)
	}
	this.values = append(this.values, v)
}

func (this *stack) peek() value {
	return this.values[len(this.values)-1]
}

func (this *stack) pop() value {
	v := this.peek()
	if stackDebug {
		log.Printf("Popping %s from stack", v)
	}
	this.values = this.values[:len(this.values)-1]
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

func NewVM(ast parser.Node) *vm {
	vm := vm{[]stackFrame{}, nil, []opcode{}, 0, nil, value{}, false, 0}
	vm.stack = []stackFrame{makeStackFrame(stack{}, 0, nil)}
	vm.currentFrame = &vm.stack[0]
	vm.code = vm.generateCode(ast)

	vm.defineVar(appendStringtable("console"), defineConsoleObject())
	vm.defineVar(appendStringtable("Math"), defineMathObject())
	vm.defineVar(appendStringtable("Boolean"), defineBooleanCtor())
	vm.defineVar(appendStringtable("String"), defineStringCtor())

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
		case PLUS:
			// ### could (should) specialize this in codegen for numeric types
			// vs unknown types?
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			lprim := lval.toPrimitive()
			rprim := rval.toPrimitive()
			if lprim.vtype == STRING || rprim.vtype == STRING {
				this.currentFrame.data_stack.push(newString(lprim.toString() + rprim.toString()))
			} else {
				this.currentFrame.data_stack.push(newNumber(lprim.toNumber() + rprim.toNumber()))
			}
		case MINUS:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newNumber(lval.toNumber() - rval.toNumber()))
		case MULTIPLY:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newNumber(lval.toNumber() * rval.toNumber()))
		case DIVIDE:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newNumber(lval.toNumber() / rval.toNumber()))
		case LESS_THAN:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newBool(lval.toNumber() < rval.toNumber()))
		case GREATER_THAN:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newBool(lval.toNumber() > rval.toNumber()))
		case EQUALS:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newBool(lval.toNumber() == rval.toNumber()))
		case NOT_EQUALS:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newBool(lval.toNumber() != rval.toNumber()))
		case LESS_THAN_EQ:
			lval := this.currentFrame.data_stack.pop()
			rval := this.currentFrame.data_stack.pop()
			this.currentFrame.data_stack.push(newBool(lval.toNumber() <= rval.toNumber()))
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
		case POP:
			if len(this.currentFrame.data_stack.values) > 0 {
				this.currentFrame.data_stack.pop()
			}
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
			memb := v.get(stringtable[op.odata.asInt()])
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
			this.currentFrame.data_stack.push(*sv)
		default:
			panic(fmt.Sprintf("unhandled opcode %+v", op))
		}
	}

	return this.returnValue
}

func (this *vm) handleCall(op opcode, isNew bool) {
	// my, this is inefficient
	builtinArgs := []value{}
	for i := 0; i < op.odata.asInt(); i++ {
		v := this.currentFrame.data_stack.pop()
		if execDebug {
			log.Printf("Loaded arg %d %s", i, v)
		}
		builtinArgs = append(builtinArgs, v)
	}

	fn := this.currentFrame.data_stack.pop()
	if fn.vtype != OBJECT {
		panic(fmt.Sprintf("CALL without a function: %s", fn))
	}

	sf := makeStackFrame(stack{}, this.ip, this.currentFrame)
	this.pushStack(sf)

	var rval value
	if isNew {
		rval = fn.construct(this, builtinArgs)
	} else {
		rval = fn.call(this, builtinArgs)
	}

	if this.ignoreReturn {
		this.ignoreReturn = false
	} else {
		this.popStack(rval)
	}
}
