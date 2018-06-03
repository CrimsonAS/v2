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

	// ++a
	INCREMENT

	// --a
	DECREMENT

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

	// not used in codegen right now...
	JMP

	// call(a)
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

	// what VM is executing it? (this is mostly for debug purposes, and is kind
	// of dirty...) ### remove
	vm *vm
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
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case MINUS:
		return "SUB"
	case MULTIPLY:
		return "MUL"
	case DIVIDE:
		return "DIV"
	case PUSH_NUMBER:
		return fmt.Sprintf("PUSH number(%f)", this.odata)
	case PUSH_STRING:
		return fmt.Sprintf("PUSH string(%d, \"%s\")", int(this.odata), this.vm.stringtable[int(this.odata)])
	case PUSH_BOOL:
		return fmt.Sprintf("PUSH bool(%f)", this.odata)
	case JMP:
		return fmt.Sprintf("JMP %d", int(this.odata))
	case CALL:
		return "CALL"
	case IN_FUNCTION:
		return fmt.Sprintf("function %s:", this.vm.stringtable[int(this.odata)])
	case JNE:
		return fmt.Sprintf("JNE %d", int(this.odata))
	case RETURN:
		return "RETURN"
	case STORE:
		return fmt.Sprintf("STORE %s", this.vm.stringtable[int(this.odata)])
	case DECLARE:
		return fmt.Sprintf("DECLARE %s", this.vm.stringtable[int(this.odata)])
	case LOAD:
		return fmt.Sprintf("LOAD %s", this.vm.stringtable[int(this.odata)])
	default:
		return fmt.Sprintf("unknown opcode %d", this.otype)
	}
}

// create an opcode with no odata
func (this *vm) simpleOp(o opcode_type) opcode {
	return opcode{o, 0, this}
}

// create an opcode with odata 'i'
func (this *vm) newOpcode(o opcode_type, i float64) opcode {
	return opcode{o, odata(i), this}
}

// generate the code for a given AST node, and return it in bytecode.
func (this *vm) generateCode(node parser.Node) []opcode {
	codebuf := []opcode{}
	switch n := node.(type) {
	case *parser.Program:
		codebuf = append(codebuf, this.newOpcode(IN_FUNCTION, float64(len(this.stringtable))))
		this.stringtable = append(this.stringtable, "%main")

		for _, s := range n.Body() {
			codebuf = append(codebuf, this.generateCode(s)...)
		}

		codebuf = append(codebuf, this.simpleOp(RETURN))

		for _, n := range this.funcsToDefine {
			// ### push and pop stack frames here? need to cooperate with
			// runtime somehow
			//if _, ok := this.currentFrame.vars[n.Identifier.String()]; ok {
			//	panic(fmt.Sprintf("Already defined function %s", n.Identifier.String()))
			//} else {
			this.currentFrame.vars[n.Identifier.String()] = newNumber(float64(len(codebuf)))
			//}

			codebuf = append(codebuf, this.newOpcode(IN_FUNCTION, float64(len(this.stringtable))))
			this.stringtable = append(this.stringtable, n.Identifier.String())
			codebuf = append(codebuf, this.generateCode(n.Body)...)
			codebuf = append(codebuf, this.simpleOp(RETURN))
		}
		return codebuf

		return codebuf
	case *parser.FunctionExpression:
		this.funcsToDefine = append(this.funcsToDefine, n)
		return codebuf
	case *parser.CallExpression:
		codebuf = append(codebuf, this.generateCode(n.X)...)
		codebuf = append(codebuf, this.simpleOp(CALL))
		return codebuf
	case *parser.StringLiteral:
		codebuf = append(codebuf, this.newOpcode(PUSH_STRING, float64(len(this.stringtable))))
		this.stringtable = append(this.stringtable, n.String())
		return codebuf
	case *parser.IdentifierLiteral:
		codebuf = append(codebuf, this.newOpcode(LOAD, float64(len(this.stringtable))))
		this.stringtable = append(this.stringtable, n.String())
		return codebuf
	case *parser.NumericLiteral:
		codebuf = append(codebuf, this.newOpcode(PUSH_NUMBER, n.Float64Value()))
		return codebuf
	case *parser.TrueLiteral:
		codebuf = append(codebuf, this.newOpcode(PUSH_BOOL, 1))
		return codebuf
	case *parser.FalseLiteral:
		codebuf = append(codebuf, this.newOpcode(PUSH_BOOL, 0))
		return codebuf
	case *parser.UnaryExpression:
		if n.IsPrefix() {
			codebuf = append(codebuf, this.generateCode(n.X)...)
			switch n.Operator() {
			case parser.INCREMENT:
				codebuf = append(codebuf, this.simpleOp(INCREMENT))
			case parser.DECREMENT:
				codebuf = append(codebuf, this.simpleOp(DECREMENT))
			case parser.PLUS:
				codebuf = append(codebuf, this.simpleOp(UPLUS))
			case parser.MINUS:
				codebuf = append(codebuf, this.simpleOp(UMINUS))
			case parser.LOGICAL_NOT:
				codebuf = append(codebuf, this.simpleOp(UNOT))
			default:
				panic(fmt.Sprintf("unknown prefix unary operator %s", n.Operator()))
			}
		} else {
			switch n.Operator() {
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
			codebuf = append(codebuf, this.simpleOp(PLUS))
			return codebuf
		case parser.MINUS:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, this.simpleOp(MINUS))
			return codebuf
		case parser.MULTIPLY:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, this.simpleOp(MULTIPLY))
			return codebuf
		case parser.DIVIDE:
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			codebuf = append(codebuf, this.generateCode(n.Left)...)
			codebuf = append(codebuf, this.simpleOp(DIVIDE))
			return codebuf
		case parser.ASSIGNMENT:
			log.Printf("%+v", n.Right)
			lhs := n.Left.(*parser.IdentifierLiteral)
			codebuf = append(codebuf, this.generateCode(n.Right)...)
			varIdx := float64(len(this.stringtable))
			this.stringtable = append(this.stringtable, lhs.String())
			codebuf = append(codebuf, this.newOpcode(STORE, varIdx))
			// LOAD to return... is this right? should STORE push()? is there
			// another way altogether?
			codebuf = append(codebuf, this.newOpcode(LOAD, varIdx))
			return codebuf
		default:
			panic(fmt.Sprintf("unknown operator %s", n.Operator()))
		}
	case *parser.IfStatement:
		test := this.generateCode(n.ConditionExpr)
		trueBranch := this.generateCode(n.ThenStmt)
		falseBranch := []opcode{}
		if n.ElseStmt != nil {
			falseBranch = this.generateCode(n.ElseStmt)
		}
		codebuf = append(codebuf, test...)
		// jump past the true branch, and the JNE
		// ### note that these are relative (no use of global code position), but
		// the VM treats them as absolute. one of these two will need to be
		// fixed, or JNE in functions will break.
		codebuf = append(codebuf, this.newOpcode(JNE, float64(len(test)+len(trueBranch)+1)))
		codebuf = append(codebuf, trueBranch...)
		codebuf = append(codebuf, falseBranch...)
		return codebuf
	case *parser.VariableStatement:
		// ### these should come at the start of a function
		for idx, _ := range n.Vars {
			v := n.Vars[idx]
			i := n.Initializers[idx]

			varIdx := float64(len(this.stringtable))
			codebuf = append(codebuf, this.newOpcode(DECLARE, varIdx))
			this.stringtable = append(this.stringtable, v.String())

			if i != nil {
				codebuf = append(codebuf, this.generateCode(i)...)
				codebuf = append(codebuf, this.newOpcode(STORE, varIdx))
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
	default:
		panic(fmt.Sprintf("unknown node %T", node))
	}

	panic("unreachable")
}
