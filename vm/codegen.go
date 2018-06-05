package vm

import (
	"fmt"
	"github.com/CrimsonAS/v2/parser"
	"log"
)

// opcode instructions...
type opcode_type int

const (
	// a + b
	PLUS opcode_type = iota

	// +a
	UPLUS

	// -a
	UMINUS

	// !a
	UNOT

	// a - b
	MINUS

	// a * b
	MULTIPLY

	// a / b
	DIVIDE

	// 5
	PUSH_NUMBER

	// true
	PUSH_BOOL

	// "hello"
	PUSH_STRING

	// LOAD hello (from the stack frame)
	LOAD

	// LOAD member (from property of the object on the stack frame)
	LOAD_MEMBER

	// not used in codegen right now...
	JMP

	// call(a)
	// the odata gives the argument count
	// arguments are on the stack, as is the thing to call.
	CALL

	// used to tell the VM which function it's inside, for debug printing
	// purposes.
	IN_FUNCTION

	// jump if false (misnamed ###)
	JNE

	// return from function
	RETURN

	// declare var
	DECLARE

	// store var
	STORE

	LESS_THAN
	GREATER_THAN
	EQUALS
	NOT_EQUALS
	LESS_THAN_EQ

	INCREMENT // a++
	DECREMENT // a--

	// discard item from the stack
	POP

	// duplicate the top of the stack
	DUP
)

// 'odata' is a piece of information attached to an opcode. It can be nothing,
// like in the case of an instruction like ADD (as the operands are pushed onto
// the stack earlier), or an actual value, like when pushing numbers onto the
// stack.
//
// note that it is only ever valid in the context of the opcode_type for the
// opcode -- e.g. PUSH_STRING uses it as an index into the string table, not as
// a number
type odata float64

func (this odata) asFloat64() float64 {
	return float64(this)
}

func (this odata) asInt() int {
	return int(this)
}

// an opcode for the VM to execute.
type opcode struct {
	// what type of instruction?
	otype opcode_type

	// what data is attached to it?
	odata odata
}

func (this opcode) String() string {
	switch this.otype {
	case PLUS:
		return "ADD"
	case UPLUS:
		return "UPLUS"
	case UMINUS:
		return "UMINUS"
	case UNOT:
		return "UNOT"
	case MINUS:
		return "SUB"
	case MULTIPLY:
		return "MUL"
	case DIVIDE:
		return "DIV"
	case DUP:
		return "DUP"
	case POP:
		return "POP"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case LESS_THAN:
		return "<"
	case LESS_THAN_EQ:
		return "<="
	case GREATER_THAN:
		return "<"
	case EQUALS:
		return "=="
	case NOT_EQUALS:
		return "!="
	case PUSH_NUMBER:
		return fmt.Sprintf("PUSH number(%f)", this.odata)
	case PUSH_STRING:
		return fmt.Sprintf("PUSH string(%d, \"%s\")", int(this.odata), stringtable[int(this.odata)])
	case PUSH_BOOL:
		return fmt.Sprintf("PUSH bool(%f)", this.odata)
	case JMP:
		return fmt.Sprintf("JMP %d", int(this.odata))
	case CALL:
		return fmt.Sprintf("CALL(argc: %d)", int(this.odata))
	case IN_FUNCTION:
		return fmt.Sprintf("function %s:", stringtable[int(this.odata)])
	case JNE:
		return fmt.Sprintf("JNE %d", int(this.odata))
	case RETURN:
		return "RETURN"
	case STORE:
		return fmt.Sprintf("STORE %s", stringtable[int(this.odata)])
	case DECLARE:
		return fmt.Sprintf("DECLARE %s", stringtable[int(this.odata)])
	case LOAD:
		return fmt.Sprintf("LOAD %s", stringtable[int(this.odata)])
	case LOAD_MEMBER:
		return fmt.Sprintf("LOAD_MEMBER %s", stringtable[int(this.odata)])
	default:
		return fmt.Sprintf("unknown opcode %d", this.otype)
	}
}

// create an opcode with no odata
func simpleOp(o opcode_type) opcode {
	return opcode{o, 0}
}

// create an opcode with odata 'i'
func newOpcode(o opcode_type, i float64) opcode {
	return opcode{o, odata(i)}
}

func callBuiltinAddr(this *vm, params []*parser.IdentifierLiteral, addr int) func(vm *vm, f value, args []value) value {
	// Small optimisation: intern strings at codegen time, so we don't have to
	// hash at runtime.
	intArgs := []int{}
	for _, arg := range params {
		intArgs = append(intArgs, appendStringtable(arg.String()))
	}

	return func(vm *vm, f value, args []value) value {
		if execDebug {
			log.Printf("Calling func! IP %d going to %d, %s", vm.ip, addr, args)
		}
		// alter the IP of the new stack frame the CALL set up to be in
		// the function's code.
		vm.ip = addr

		// bit of a dirty hack here. we tell the VM to ignore the return
		// value of the builtin function, and instead, wait for the
		// return instruction to pop the stack.
		vm.ignoreReturn = true

		for idx, arg := range intArgs {
			v := args[idx]
			if execDebug {
				log.Printf("Defining var %s %s", stringtable[arg], v)
			}
			vm.defineVar(arg, v)
		}

		return newUndefined()
	}
}

// generate the code for a given AST node, and return it in bytecode.
func (this *vm) generateCode(node parser.Node) []opcode {
	codebuf := []opcode{}
	switch n := node.(type) {
	case *parser.Program:
		funcidx := appendStringtable("%main")
		codebuf = append(codebuf, newOpcode(IN_FUNCTION, float64(funcidx)))

		for _, s := range n.Body() {
			codebuf = append(codebuf, this.generateCode(s)...)
		}

		codebuf = append(codebuf, simpleOp(RETURN))

		for _, n := range this.funcsToDefine {
			runBuiltin := callBuiltinAddr(this, n.Parameters, len(codebuf))
			callFn := newFunctionObject(runBuiltin)
			varIdx := appendStringtable(n.Identifier.String())
			this.defineVar(varIdx, callFn)

			codebuf = append(codebuf, newOpcode(IN_FUNCTION, float64(varIdx)))
			codebuf = append(codebuf, this.generateCode(n.Body)...)

			// Generate a return if the function didn't
			if codebuf[len(codebuf)-1].otype != RETURN {
				codebuf = append(codebuf, simpleOp(RETURN))
			}
		}
		return codebuf
	case *parser.ExpressionStatement:
		codebuf = append(codebuf, this.generateCode(n.X)...)
		// this seems to break tests... but why?
		//codebuf = append(codebuf, simpleOp(POP))
		return codebuf
	case *parser.FunctionExpression:
		this.funcsToDefine = append(this.funcsToDefine, n)
		return codebuf
	case *parser.CallExpression:
		codebuf = append(codebuf, this.generateCode(n.X)...)

		// We push the arguments in _reverse_ order so that CALL can build an
		// argument list just by popping.
		for idx, _ := range n.Arguments {
			codebuf = append(codebuf, this.generateCode(n.Arguments[len(n.Arguments)-(idx+1)])...)
		}

		codebuf = append(codebuf, newOpcode(CALL, float64(len(n.Arguments))))
		return codebuf
	case *parser.StringLiteral:
		sl := appendStringtable(n.String())
		codebuf = append(codebuf, newOpcode(PUSH_STRING, float64(sl)))
		return codebuf
	case *parser.IdentifierLiteral:
		il := appendStringtable(n.String())
		codebuf = append(codebuf, newOpcode(LOAD, float64(il)))
		return codebuf
	case *parser.NumericLiteral:
		codebuf = append(codebuf, newOpcode(PUSH_NUMBER, n.Float64Value()))
		return codebuf
	case *parser.TrueLiteral:
		codebuf = append(codebuf, newOpcode(PUSH_BOOL, 1))
		return codebuf
	case *parser.FalseLiteral:
		codebuf = append(codebuf, newOpcode(PUSH_BOOL, 0))
		return codebuf
	case *parser.UnaryExpression:
		if n.IsPrefix() {
			codebuf = append(codebuf, this.generateCode(n.X)...)
			switch n.Operator() {
			case parser.PLUS:
				codebuf = append(codebuf, simpleOp(UPLUS))
			case parser.MINUS:
				codebuf = append(codebuf, simpleOp(UMINUS))
			case parser.LOGICAL_NOT:
				codebuf = append(codebuf, simpleOp(UNOT))

			// See the comment for postfix INCREMENT/DECREMENT.
			case parser.INCREMENT:
				lhs := n.X.(*parser.IdentifierLiteral)
				varIdx := float64(appendStringtable(lhs.String()))
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(INCREMENT))
				codebuf = append(codebuf, simpleOp(DUP))
				codebuf = append(codebuf, newOpcode(STORE, varIdx))
			case parser.DECREMENT:
				lhs := n.X.(*parser.IdentifierLiteral)
				varIdx := float64(appendStringtable(lhs.String()))
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(DECREMENT))
				codebuf = append(codebuf, simpleOp(DUP))
				codebuf = append(codebuf, newOpcode(STORE, varIdx))
			default:
				panic(fmt.Sprintf("unknown prefix unary operator %s", n.Operator()))
			}
		} else {
			lhs := n.X.(*parser.IdentifierLiteral)
			varIdx := float64(appendStringtable(lhs.String()))

			// These are pretty ugly.
			// The DUP is needed so that subsequent assignment operations can
			// pop the original value, but that's pretty disgusting.
			switch n.Operator() {
			case parser.INCREMENT:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(DUP))
				codebuf = append(codebuf, simpleOp(INCREMENT))
				codebuf = append(codebuf, newOpcode(STORE, varIdx))
			case parser.DECREMENT:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(DUP))
				codebuf = append(codebuf, simpleOp(DECREMENT))
				codebuf = append(codebuf, newOpcode(STORE, varIdx))
			default:
				panic(fmt.Sprintf("unknown postfix unary operator %s", n.Operator()))
			}
		}
		return codebuf
	case *parser.BinaryExpression:
		switch n.Operator() {
		case parser.PLUS:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(PLUS))
			return codebuf
		case parser.MINUS:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(MINUS))
			return codebuf
		case parser.MULTIPLY:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(MULTIPLY))
			return codebuf
		case parser.DIVIDE:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(DIVIDE))
			return codebuf
		case parser.LESS_THAN:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(LESS_THAN))
			return codebuf
		case parser.GREATER_THAN:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(GREATER_THAN))
			return codebuf
		case parser.EQUALS:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(EQUALS))
			return codebuf
		case parser.NOT_EQUALS:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(NOT_EQUALS))
			return codebuf
		case parser.LESS_EQ:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, simpleOp(LESS_THAN_EQ))
			return codebuf
		case parser.ASSIGNMENT:
			lhs := n.Left.(*parser.IdentifierLiteral)
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			varIdx := float64(appendStringtable(lhs.String()))
			codebuf = append(codebuf, simpleOp(DUP)) // duplicate so it's available as a return value too...
			codebuf = append(codebuf, newOpcode(STORE, varIdx))
			return codebuf
		default:
			panic(fmt.Sprintf("unknown operator %s", n.Operator()))
		}
	case *parser.ReturnStatement:
		if n.X != nil {
			codebuf = append(codebuf, this.generateCode(n.X)...)
		} else {
			panic("should LOAD(undefined)")
		}
		codebuf = append(codebuf, simpleOp(RETURN))
		return codebuf
	case *parser.ForStatement:
		// for (init; test; update) { body }
		if n.Initializer != nil {
			codebuf = append(codebuf, this.generateCode(n.Initializer)...) // init
		}
		test := []opcode{}
		update := []opcode{}
		if n.Test != nil {
			test = this.generateCode(n.Test)
		}
		if n.Update != nil {
			update = this.generateCode(n.Update)
		}
		body := this.generateCode(n.Body)
		codebuf = append(codebuf, test...)
		codebuf = append(codebuf, newOpcode(JNE, float64(len(update)+len(body)+1))) // jump over the update, body and JMP
		codebuf = append(codebuf, body...)
		codebuf = append(codebuf, update...)
		codebuf = append(codebuf, newOpcode(JMP, float64(-(len(body)+len(update)+len(test)+2)))) // back to the test start
		return codebuf
	case *parser.DoWhileStatement:
		codebuf = append(codebuf, this.generateCode(n.Body)...)               // do { n.Body }
		codebuf = append(codebuf, this.generateCode(n.X)...)                  // while (X)
		codebuf = append(codebuf, newOpcode(JNE, float64(1)))                 // jump over the following JMP
		codebuf = append(codebuf, newOpcode(JMP, float64(-(len(codebuf)+1)))) // back to the trueBranch start
		return codebuf
	case *parser.WhileStatement:
		test := this.generateCode(n.X)

		trueBranch := this.generateCode(n.Body)
		codebuf = append(codebuf, test...)
		// if (!test) -> skip trueBranch
		// the '1's here are for the instructions we're inserting ourselves
		// (JMP/JNE)
		codebuf = append(codebuf, newOpcode(JNE, float64(len(trueBranch)+1)))
		codebuf = append(codebuf, trueBranch...)
		// jmp back to the test
		codebuf = append(codebuf, newOpcode(JMP, float64(-(len(codebuf)+1))))
		return codebuf
	case *parser.ConditionalExpression:
		// ### duplicates IfStatement
		test := this.generateCode(n.X)
		trueBranch := this.generateCode(n.Then)
		falseBranch := []opcode{}
		if n.Else != nil {
			falseBranch = this.generateCode(n.Else)
		}
		codebuf = append(codebuf, test...)
		codebuf = append(codebuf, newOpcode(JNE, float64(len(trueBranch))))
		codebuf = append(codebuf, trueBranch...)
		codebuf = append(codebuf, falseBranch...)
		return codebuf
	case *parser.IfStatement:
		// ### duplicates ConditionalExpression
		test := this.generateCode(n.ConditionExpr)
		trueBranch := this.generateCode(n.ThenStmt)
		falseBranch := []opcode{}
		if n.ElseStmt != nil {
			falseBranch = this.generateCode(n.ElseStmt)
		}
		codebuf = append(codebuf, test...)
		codebuf = append(codebuf, newOpcode(JNE, float64(len(trueBranch))))
		codebuf = append(codebuf, trueBranch...)
		codebuf = append(codebuf, falseBranch...)
		return codebuf
	case *parser.VariableStatement:
		// ### these should come at the start of a function
		for idx, _ := range n.Vars {
			v := n.Vars[idx]
			i := n.Initializers[idx]

			varIdx := float64(appendStringtable(v.String()))
			codebuf = append(codebuf, newOpcode(DECLARE, varIdx))

			if i != nil {
				codebuf = append(codebuf, this.generateCode(i)...)
				codebuf = append(codebuf, newOpcode(STORE, varIdx))
			}
		}
		return codebuf
	case *parser.BlockStatement:
		for _, s := range n.Body {
			codebuf = append(codebuf, this.generateCode(s)...)
		}
		return codebuf
	case *parser.EmptyStatement:
		// ### do we need to generate anything?
		return codebuf
	case *parser.DotMemberExpression:
		codebuf = append(codebuf, this.generateCode(n.X)...)
		varIdx := appendStringtable(n.Name.String())
		codebuf = append(codebuf, newOpcode(LOAD_MEMBER, float64(varIdx)))
		return codebuf
	default:
		panic(fmt.Sprintf("unknown node %T", node))
	}

	panic(fmt.Sprintf("unreachable %T", node))
}
