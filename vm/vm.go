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
	//log.Printf("Pushing %s onto stack", v)
	this.values = append(this.values, v)
}

func (this *stack) peek() value {
	return this.values[len(this.values)-1]
}

func (this *stack) pop() value {
	v := this.peek()
	//log.Printf("Popping %s from stack", v)
	this.values = this.values[:len(this.values)-1]
	return v
}

type stackFrame struct {
	data_stack stack
	retAddr    int
	vars       map[string]value
}

type vm struct {
	stack         []stackFrame
	currentFrame  *stackFrame
	code          []opcode
	ip            int
	funcsToDefine []*parser.FunctionExpression
	stringtable   []string
}

func makeStackFrame() stackFrame {
	return stackFrame{vars: make(map[string]value)}
}

func NewVM(ast parser.Node) *vm {
	vm := vm{[]stackFrame{}, nil, []opcode{}, 0, nil, nil}
	vm.stack = []stackFrame{makeStackFrame()}
	vm.currentFrame = &vm.stack[0]
	vm.code = vm.generateCode(ast)
	return &vm
}

func (this *vm) Run() value {
	log.Printf("String table:")
	for i := 0; i < len(this.stringtable); i++ {
		log.Printf("%d: %s", i, this.stringtable[i])
	}
	log.Printf("Program:")
	for i := 0; i < len(this.code); i++ {
		log.Printf("%d: %s", i, this.code[i])
	}
	log.Printf("Starting execution")

	for ; this.ip < len(this.code); this.ip++ {
		op := this.code[this.ip]
		log.Printf("Op %d: %s", this.ip, op)
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
			this.currentFrame.data_stack.push(newString(this.stringtable[op.odata.asInt()]))
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
			// ### incomplete (es5 11.4.4)
			val := this.currentFrame.data_stack.pop()
			oldValue := val.toNumber()
			this.currentFrame.data_stack.push(newNumber(oldValue + 1))
		case DECREMENT:
			// ### incomplete (es5 11.4.4)
			val := this.currentFrame.data_stack.pop()
			oldValue := val.toNumber()
			this.currentFrame.data_stack.push(newNumber(oldValue - 1))
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
		case JMP:
			this.ip = op.odata.asInt()
		case JNE:
			test := this.currentFrame.data_stack.pop()

			if !test.toBoolean() {
				this.ip = op.odata.asInt()
			}
		case CALL:
			fn := this.currentFrame.data_stack.pop()
			log.Printf("CALL %s", fn)
			// ### hack: don't store functions as an address, haha.
			this.ip = int(fn.toNumber())
			sf := makeStackFrame()
			sf.retAddr = this.ip + 1
			this.stack = append(this.stack, sf)
			this.currentFrame = &this.stack[len(this.stack)-1]
			log.Printf("Stack now: %+v", this.stack)
		case RETURN:
			log.Printf("Returning from stack: %+v", this.stack)
			this.stack = this.stack[:len(this.stack)-1]

			rval := value{}
			if len(this.currentFrame.data_stack.values) > 0 {
				log.Printf("Returning %s\n", this.currentFrame.data_stack.peek())
				rval = this.currentFrame.data_stack.pop()
			}

			if len(this.stack) > 0 {
				this.ip = this.currentFrame.retAddr
				this.currentFrame = &this.stack[len(this.stack)-1]
				this.currentFrame.data_stack.push(rval)
			} else {
				// return from main
				return rval
			}

			log.Printf("Stack now: %+v", this.stack)
		case IN_FUNCTION:
			// no-op, just for informative/debug purposes
		case DECLARE:
			// ### ensure it doesn't exist
			log.Printf("Var %s declared (%d)", this.stringtable[op.odata.asInt()], op.odata.asInt())
			this.currentFrame.vars[this.stringtable[op.odata.asInt()]] = value{}
		case STORE:
			v := this.currentFrame.data_stack.pop()
			log.Printf("Storing %s in %s (%d)", v, this.stringtable[op.odata.asInt()], op.odata.asInt())
			this.currentFrame.vars[this.stringtable[op.odata.asInt()]] = v
		case LOAD:
			v := this.currentFrame.vars[this.stringtable[op.odata.asInt()]]
			log.Printf("Loading %s from %d gave %s", this.stringtable[op.odata.asInt()], op.odata.asInt(), v)
			this.currentFrame.data_stack.push(v)
		default:
			panic(fmt.Sprintf("unhandled opcode %+v", op))
		}
	}

	panic("unreachable")
}
