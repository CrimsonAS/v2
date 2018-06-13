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
)

// opcode instructions...
type opcode_type uint8

const (
	// Simple math operators
	ADD opcode_type = iota
	SUB
	MULTIPLY
	DIVIDE
	LEFT_SHIFT
	RIGHT_SHIFT
	UNSIGNED_RIGHT_SHIFT
	BITWISE_AND // a & b
	BITWISE_XOR // a ^ b
	BITWISE_OR  // a | b

	UPLUS       // +a
	UMINUS      // -a
	UNOT        // !a
	TYPEOF      // typeof a
	BITWISE_NOT // ~a

	// a % b
	MODULUS

	// These all push a given value to the stack.
	LOAD_THIS      // 'this'
	PUSH_UNDEFINED // undefined
	PUSH_NULL      // null
	PUSH_NUMBER    // 5
	PUSH_ARRAY     // [a, b, c...]
	PUSH_BOOL      // true
	PUSH_STRING    // "hello" (note: the string index is given via the opdata)

	// LOAD identifier
	// Pushes a variable onto the stack.
	// Note that this also sets the 'this' arg for calls.
	LOAD

	// Loads a member from the topmost stack item, and pushes it to the stack frame.
	LOAD_MEMBER
	STORE_MEMBER
	LOAD_INDEXED
	STORE_INDEXED

	// Jump ip, relative to the current position.
	JMP

	// call/new call the topmost function object on the stack.
	// the opdata gives the argument count
	// arguments are on the stack as well.
	CALL
	NEW

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
	LOGICAL_AND

	INCREMENT // a++
	DECREMENT // a--

	// discard item from the stack
	POP

	// duplicate the top of the stack
	DUP

	// Start an object definition. It will be followed by property definition,
	// and an END_OBJECT.
	NEW_OBJECT

	// Define property of the NEW_OBJECT on the stack, with the arg on the stack.
	DEFINE_PROPERTY

	// End object definition.
	END_OBJECT
)

// 'opdata' is a piece of information attached to an opcode. It can be nothing,
// like in the case of an instruction like ADD (as the operands are pushed onto
// the stack earlier), or an actual value, like when pushing numbers onto the
// stack.
//
// note that it is only ever valid in the context of the opcode_type for the
// opcode -- e.g. PUSH_STRING uses it as an index into the string table, not as
// a number
type opdata float64

func (this opdata) asFloat64() float64 {
	return float64(this)
}

func (this opdata) asInt() int {
	return int(this)
}

// an opcode for the VM to execute.
type opcode struct {
	// what type of instruction?
	otype opcode_type

	// what data is attached to it?
	opdata opdata
}

func (this opcode) String() string {
	switch this.otype {
	case ADD:
		return "ADD"
	case UPLUS:
		return "UPLUS"
	case UMINUS:
		return "UMINUS"
	case UNOT:
		return "UNOT"
	case TYPEOF:
		return "TYPEOF"
	case BITWISE_NOT:
		return "BITWISE_NOT"
	case SUB:
		return "SUB"
	case MULTIPLY:
		return "MUL"
	case DIVIDE:
		return "DIV"
	case LEFT_SHIFT:
		return "<<"
	case RIGHT_SHIFT:
		return ">>"
	case UNSIGNED_RIGHT_SHIFT:
		return ">>>"
	case BITWISE_AND:
		return "&"
	case BITWISE_XOR:
		return "^"
	case BITWISE_OR:
		return "|"
	case MODULUS:
		return "MOD"
	case NEW_OBJECT:
		return "NEW_OBJECT"
	case DEFINE_PROPERTY:
		return "DEFINE_PROPERTY"
	case END_OBJECT:
		return "END_OBJECT"
	case DUP:
		return "DUP"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case LESS_THAN:
		return "<"
	case LESS_THAN_EQ:
		return "<="
	case LOGICAL_AND:
		return "&&"
	case GREATER_THAN:
		return "<"
	case EQUALS:
		return "=="
	case NOT_EQUALS:
		return "!="
	case LOAD_THIS:
		return "LOAD this"
	case PUSH_UNDEFINED:
		return "PUSH undefined"
	case PUSH_NULL:
		return "PUSH null"
	case PUSH_NUMBER:
		return fmt.Sprintf("PUSH number(%f)", this.opdata)
	case PUSH_ARRAY:
		return fmt.Sprintf("PUSH array(%f)", this.opdata)
	case PUSH_STRING:
		return fmt.Sprintf("PUSH string(%d, \"%s\")", int(this.opdata), stringtable[int(this.opdata)])
	case PUSH_BOOL:
		return fmt.Sprintf("PUSH bool(%f)", this.opdata)
	case JMP:
		return fmt.Sprintf("JMP %d", int(this.opdata))
	case CALL:
		return fmt.Sprintf("CALL(argc: %d)", int(this.opdata))
	case NEW:
		return fmt.Sprintf("NEW(argc: %d)", int(this.opdata))
	case IN_FUNCTION:
		return fmt.Sprintf("function %s:", stringtable[int(this.opdata)])
	case JNE:
		return fmt.Sprintf("JNE %d", int(this.opdata))
	case RETURN:
		return "RETURN"
	case STORE:
		return fmt.Sprintf("STORE %s", stringtable[int(this.opdata)])
	case DECLARE:
		return fmt.Sprintf("DECLARE %s", stringtable[int(this.opdata)])
	case LOAD:
		return fmt.Sprintf("LOAD %s", stringtable[int(this.opdata)])
	case LOAD_MEMBER:
		return fmt.Sprintf("LOAD_MEMBER %s", stringtable[int(this.opdata)])
	case STORE_MEMBER:
		return fmt.Sprintf("STORE_MEMBER %s", stringtable[int(this.opdata)])
	case LOAD_INDEXED:
		return fmt.Sprintf("LOAD_INDEXED")
	case STORE_INDEXED:
		return fmt.Sprintf("STORE_INDEXED")
	default:
		return fmt.Sprintf("unknown opcode %d", this.otype)
	}
}

// create an opcode with no opdata
func simpleOp(o opcode_type) opcode {
	return opcode{o, 0}
}

// create an opcode with opdata 'i'
func newOpcode(o opcode_type, i float64) opcode {
	return opcode{o, opdata(i)}
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

func (this *vm) generateCodeForLiteral(node parser.Node) []opcode {
	codebuf := []opcode{}
	switch n := node.(type) {
	case *parser.ArrayLiteral:
		for _, elem := range n.Elements {
			codebuf = append(codebuf, this.generateCode(elem)...)
		}
		codebuf = append(codebuf, newOpcode(PUSH_ARRAY, float64(len(n.Elements))))
	case *parser.NumericLiteral:
		codebuf = append(codebuf, newOpcode(PUSH_NUMBER, n.Float64Value()))
	case *parser.TrueLiteral:
		codebuf = append(codebuf, newOpcode(PUSH_BOOL, 1))
	case *parser.FalseLiteral:
		codebuf = append(codebuf, newOpcode(PUSH_BOOL, 0))
	case *parser.ThisLiteral:
		codebuf = append(codebuf, newOpcode(LOAD_THIS, 0))
	case *parser.NullLiteral:
		codebuf = append(codebuf, simpleOp(PUSH_NULL))
	case *parser.ObjectLiteral:
		codebuf = append(codebuf, simpleOp(NEW_OBJECT))

		for _, prop := range n.Properties {
			switch pk := prop.Key.(type) {
			case *parser.IdentifierLiteral:
				kl := appendStringtable(pk.String())
				codebuf = append(codebuf, newOpcode(PUSH_STRING, float64(kl)))
			case *parser.NumericLiteral:
				kl := appendStringtable(pk.String())
				codebuf = append(codebuf, newOpcode(PUSH_STRING, float64(kl)))
			case *parser.StringLiteral:
				kl := appendStringtable(pk.String())
				codebuf = append(codebuf, newOpcode(PUSH_STRING, float64(kl)))
			default:
				panic("unknown object key")
			}
			codebuf = append(codebuf, this.generateCode(prop.X)...)
			codebuf = append(codebuf, simpleOp(DEFINE_PROPERTY))
		}

		if this.canConsume > 0 {
			codebuf = append(codebuf, simpleOp(DUP))
		}

		codebuf = append(codebuf, simpleOp(END_OBJECT))
	default:
		panic(fmt.Sprintf("unknown literal %T", node))
	}

	return codebuf
}

func (this *vm) generateCodeForStatement(node parser.Node) []opcode {
	codebuf := []opcode{}
	switch n := node.(type) {
	case *parser.VariableStatement:
		// ### these should come at the start of a function
		this.canConsume++
		defer func() { this.canConsume = this.canConsume - 1 }()

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
	case *parser.ExpressionStatement:
		codebuf = append(codebuf, this.generateCode(n.X)...)
	case *parser.ReturnStatement:
		this.canConsume++
		defer func() { this.canConsume = this.canConsume - 1 }()
		if n.X != nil {
			codebuf = append(codebuf, this.generateCode(n.X)...)
		} else {
			codebuf = append(codebuf, simpleOp(PUSH_UNDEFINED))
		}
		codebuf = append(codebuf, simpleOp(RETURN))
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
		jumpLen := 0
		if len(test) > 0 {
			codebuf = append(codebuf, test...)
			codebuf = append(codebuf, newOpcode(JNE, float64(len(update)+len(body)+1))) // jump over the update, body and JMP
			jumpLen = 2                                                                 // back over the JNE
		} else {
			jumpLen = 1 // back over the JMP only
		}
		codebuf = append(codebuf, body...)
		codebuf = append(codebuf, update...)
		codebuf = append(codebuf, newOpcode(JMP, float64(-(len(body)+len(update)+len(test)+jumpLen)))) // back to the test start
	case *parser.DoWhileStatement:
		codebuf = append(codebuf, this.generateCode(n.Body)...)               // do { n.Body }
		codebuf = append(codebuf, this.generateCode(n.X)...)                  // while (X)
		codebuf = append(codebuf, newOpcode(JNE, float64(1)))                 // jump over the following JMP
		codebuf = append(codebuf, newOpcode(JMP, float64(-(len(codebuf)+1)))) // back to the trueBranch start
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
	case *parser.BlockStatement:
		for _, s := range n.Body {
			codebuf = append(codebuf, this.generateCode(s)...)
		}
	case *parser.EmptyStatement:
	default:
		panic(fmt.Sprintf("unknown statement %T", node))
	}

	return codebuf
}

func (this *vm) generateCodeForExpression(node parser.Node) []opcode {
	codebuf := []opcode{}
	switch n := node.(type) {
	case *parser.FunctionExpression:
		this.funcsToDefine = append(this.funcsToDefine, n)
	case *parser.NewExpression:
		this.isNew++
		codebuf = append(codebuf, this.generateCode(n.X)...)
		this.isNew--
	case *parser.CallExpression:
		for _, arg := range n.Arguments {
			codebuf = append(codebuf, this.generateCode(arg)...)
		}
		codebuf = append(codebuf, this.generateCode(n.X)...)

		if this.isNew > 0 {
			codebuf = append(codebuf, newOpcode(NEW, float64(len(n.Arguments))))
		} else {
			codebuf = append(codebuf, newOpcode(CALL, float64(len(n.Arguments))))
		}
	case *parser.StringLiteral:
		sl := appendStringtable(n.String())
		codebuf = append(codebuf, newOpcode(PUSH_STRING, float64(sl)))
	case *parser.IdentifierLiteral:
		if n.String() == "undefined" {
			// I have no idea why this is an identifier, but there you go
			codebuf = append(codebuf, simpleOp(PUSH_UNDEFINED))
		} else {
			il := appendStringtable(n.String())
			codebuf = append(codebuf, newOpcode(LOAD, float64(il)))
		}
	case *parser.UnaryExpression:
		if n.IsPrefix() {
			switch n.Operator() {
			case parser.PLUS:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(UPLUS))
			case parser.MINUS:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(UMINUS))
			case parser.LOGICAL_NOT:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(UNOT))
			case parser.TYPEOF:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(TYPEOF))
			case parser.BITWISE_NOT:
				codebuf = append(codebuf, this.generateCode(n.X)...)
				codebuf = append(codebuf, simpleOp(BITWISE_NOT))

			// See the comment for postfix INCREMENT/DECREMENT.
			case parser.INCREMENT:
				varIdx := float64(0.0)
				switch lhs := n.X.(type) {
				case *parser.IdentifierLiteral:
					varIdx = float64(appendStringtable(lhs.String()))
					codebuf = append(codebuf, this.generateCode(n.X)...)
					codebuf = append(codebuf, simpleOp(INCREMENT))
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, newOpcode(STORE, varIdx))
				case *parser.DotMemberExpression:
					varIdx = float64(appendStringtable(lhs.Name.String()))
					codebuf = append(codebuf, simpleOp(DUP)) // so we can store back to it
					codebuf = append(codebuf, this.generateCode(lhs.X)...)
					codebuf = append(codebuf, newOpcode(LOAD_MEMBER, float64(varIdx)))
					codebuf = append(codebuf, simpleOp(INCREMENT))
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, newOpcode(STORE_MEMBER, float64(varIdx)))
				}
			case parser.DECREMENT:
				varIdx := float64(0.0)
				switch lhs := n.X.(type) {
				case *parser.IdentifierLiteral:
					varIdx = float64(appendStringtable(lhs.String()))
					codebuf = append(codebuf, this.generateCode(n.X)...)
					codebuf = append(codebuf, simpleOp(DECREMENT))
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, newOpcode(STORE, varIdx))
				case *parser.DotMemberExpression:
					varIdx = float64(appendStringtable(lhs.Name.String()))
					codebuf = append(codebuf, this.generateCode(lhs.X)...)
					codebuf = append(codebuf, simpleOp(DUP)) // so we can store back to it
					codebuf = append(codebuf, newOpcode(LOAD_MEMBER, float64(varIdx)))
					codebuf = append(codebuf, simpleOp(DECREMENT))
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, newOpcode(STORE_MEMBER, float64(varIdx)))
				}
			default:
				panic(fmt.Sprintf("unknown prefix unary operator %s", n.Operator()))
			}
		} else {
			switch lhs := n.X.(type) {
			case *parser.IdentifierLiteral:
				varIdx := float64(appendStringtable(lhs.String()))

				switch n.Operator() {
				case parser.INCREMENT:
					codebuf = append(codebuf, this.generateCode(n.X)...)
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, simpleOp(INCREMENT))
					codebuf = append(codebuf, newOpcode(STORE, varIdx))
				case parser.DECREMENT:
					codebuf = append(codebuf, this.generateCode(n.X)...)
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, simpleOp(DECREMENT))
					codebuf = append(codebuf, newOpcode(STORE, varIdx))
				default:
					panic(fmt.Sprintf("unknown postfix unary operator %s", n.Operator()))
				}
			case *parser.DotMemberExpression:
				switch n.Operator() {
				case parser.INCREMENT:
					codebuf = append(codebuf, this.generateCode(n.X)...)
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, simpleOp(INCREMENT))
					varIdx := float64(appendStringtable(lhs.Name.String()))
					codebuf = append(codebuf, newOpcode(STORE_MEMBER, varIdx))
				case parser.DECREMENT:
					codebuf = append(codebuf, this.generateCode(n.X)...)
					if this.canConsume > 0 {
						codebuf = append(codebuf, simpleOp(DUP))
					}
					codebuf = append(codebuf, simpleOp(DECREMENT))
					varIdx := float64(appendStringtable(lhs.Name.String()))
					codebuf = append(codebuf, newOpcode(STORE_MEMBER, varIdx))
				}
			}
		}
	case *parser.AssignmentExpression:
		this.canConsume++
		defer func() { this.canConsume = this.canConsume - 1 }()

		var realOp opcode_type
		codebuf = append(codebuf, this.generateCode(n.Right)...)

		switch n.Operator() {
		case parser.ASSIGNMENT:
			realOp = STORE
		case parser.PLUS_EQ:
			realOp = ADD
		case parser.MINUS_EQ:
			realOp = SUB
		case parser.MULTIPLY_EQ:
			realOp = MULTIPLY
		case parser.DIVIDE_EQ:
			realOp = DIVIDE
		case parser.MODULUS_EQ:
			realOp = MODULUS
		case parser.LEFT_SHIFT_EQ:
			realOp = LEFT_SHIFT
		case parser.RIGHT_SHIFT_EQ:
			realOp = RIGHT_SHIFT
		case parser.UNSIGNED_RIGHT_SHIFT_EQ:
			realOp = UNSIGNED_RIGHT_SHIFT
		case parser.AND_EQ:
			realOp = BITWISE_AND
		case parser.XOR_EQ:
			realOp = BITWISE_XOR
		case parser.OR_EQ:
			realOp = BITWISE_OR
		default:
			panic(fmt.Sprintf("unknown operator %s", n.Operator()))
		}

		if realOp != STORE {
			// If it isn't a direct assignment, load the left hand side, perform
			// the op.
			switch lhs := n.Left.(type) {
			case *parser.IdentifierLiteral:
				varIdx := float64(appendStringtable(lhs.String()))
				codebuf = append(codebuf, newOpcode(LOAD, varIdx))
			case *parser.DotMemberExpression:
				codebuf = append(codebuf, this.generateCode(lhs.X)...)
				varIdx := appendStringtable(lhs.Name.String())
				codebuf = append(codebuf, newOpcode(LOAD_MEMBER, float64(varIdx)))
			case *parser.BracketMemberExpression:
				codebuf = append(codebuf, this.generateCode(lhs.X)...)
				codebuf = append(codebuf, this.generateCode(lhs.Y)...)
				codebuf = append(codebuf, simpleOp(LOAD_INDEXED))
			default:
				panic(fmt.Sprintf("unknown left hand side for assignment %T", n.Left))
			}
			codebuf = append(codebuf, simpleOp(realOp))
		}

		// Now store the result back to the left hand side.
		switch lhs := n.Left.(type) {
		case *parser.IdentifierLiteral:
			varIdx := float64(appendStringtable(lhs.String()))
			codebuf = append(codebuf, newOpcode(STORE, varIdx))
		case *parser.DotMemberExpression:
			codebuf = append(codebuf, this.generateCode(lhs.X)...)
			varIdx := appendStringtable(lhs.Name.String())
			codebuf = append(codebuf, newOpcode(STORE_MEMBER, float64(varIdx)))
		case *parser.BracketMemberExpression:
			codebuf = append(codebuf, this.generateCode(lhs.Y)...)
			codebuf = append(codebuf, this.generateCode(lhs.X)...)
			codebuf = append(codebuf, simpleOp(STORE_INDEXED))
		default:
			panic(fmt.Sprintf("unknown left hand side for assignment %T", n.Left))
		}
	case *parser.BinaryExpression:
		this.canConsume++
		defer func() { this.canConsume = this.canConsume - 1 }()
		codebuf = append(codebuf, this.generateCode(n.Right)...)
		codebuf = append(codebuf, this.generateCode(n.Left)...)
		switch n.Operator() {
		case parser.PLUS:
			codebuf = append(codebuf, simpleOp(ADD))
		case parser.MINUS:
			codebuf = append(codebuf, simpleOp(SUB))
		case parser.MULTIPLY:
			codebuf = append(codebuf, simpleOp(MULTIPLY))
		case parser.DIVIDE:
			codebuf = append(codebuf, simpleOp(DIVIDE))
		case parser.LEFT_SHIFT:
			codebuf = append(codebuf, simpleOp(LEFT_SHIFT))
		case parser.RIGHT_SHIFT:
			codebuf = append(codebuf, simpleOp(RIGHT_SHIFT))
		case parser.UNSIGNED_RIGHT_SHIFT:
			codebuf = append(codebuf, simpleOp(UNSIGNED_RIGHT_SHIFT))
		case parser.BITWISE_AND:
			codebuf = append(codebuf, simpleOp(BITWISE_AND))
		case parser.BITWISE_XOR:
			codebuf = append(codebuf, simpleOp(BITWISE_XOR))
		case parser.BITWISE_OR:
			codebuf = append(codebuf, simpleOp(BITWISE_OR))
		case parser.MODULUS:
			codebuf = append(codebuf, simpleOp(MODULUS))
		case parser.LESS_THAN:
			codebuf = append(codebuf, simpleOp(LESS_THAN))
		case parser.GREATER_THAN:
			codebuf = append(codebuf, simpleOp(GREATER_THAN))
		case parser.EQUALS:
			codebuf = append(codebuf, simpleOp(EQUALS))
		case parser.NOT_EQUALS:
			codebuf = append(codebuf, simpleOp(NOT_EQUALS))
		case parser.LESS_EQ:
			codebuf = append(codebuf, simpleOp(LESS_THAN_EQ))
		case parser.LOGICAL_AND:
			codebuf = append(codebuf, simpleOp(LOGICAL_AND))
		default:
			panic(fmt.Sprintf("unknown operator %s", n.Operator()))
		}
	case *parser.DotMemberExpression:
		codebuf = append(codebuf, this.generateCode(n.X)...)
		varIdx := appendStringtable(n.Name.String())
		codebuf = append(codebuf, newOpcode(LOAD_MEMBER, float64(varIdx)))
	case *parser.BracketMemberExpression:
		codebuf = append(codebuf, this.generateCode(n.X)...)
		codebuf = append(codebuf, this.generateCode(n.Y)...)
		codebuf = append(codebuf, simpleOp(LOAD_INDEXED))
	default:
		panic(fmt.Sprintf("unknown expression %T", node))
	}

	return codebuf
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
			callFn := newFunctionObject(runBuiltin, nil)
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
	case *parser.ArrayLiteral:
		return this.generateCodeForLiteral(n)
	case *parser.NumericLiteral:
		return this.generateCodeForLiteral(n)
	case *parser.TrueLiteral:
		return this.generateCodeForLiteral(n)
	case *parser.FalseLiteral:
		return this.generateCodeForLiteral(n)
	case *parser.NullLiteral:
		return this.generateCodeForLiteral(n)
	case *parser.ObjectLiteral:
		return this.generateCodeForLiteral(n)
	case *parser.ThisLiteral:
		return this.generateCodeForLiteral(n)

	case *parser.VariableStatement:
		return this.generateCodeForStatement(n)
	case *parser.ExpressionStatement:
		return this.generateCodeForStatement(n)
	case *parser.ReturnStatement:
		return this.generateCodeForStatement(n)
	case *parser.ForStatement:
		return this.generateCodeForStatement(n)
	case *parser.DoWhileStatement:
		return this.generateCodeForStatement(n)
	case *parser.WhileStatement:
		return this.generateCodeForStatement(n)
	case *parser.ConditionalExpression:
		return this.generateCodeForStatement(n)
	case *parser.IfStatement:
		return this.generateCodeForStatement(n)
	case *parser.BlockStatement:
		return this.generateCodeForStatement(n)
	case *parser.EmptyStatement:
		return this.generateCodeForStatement(n)

	case *parser.FunctionExpression:
		return this.generateCodeForExpression(n)
	case *parser.NewExpression:
		return this.generateCodeForExpression(n)
	case *parser.CallExpression:
		return this.generateCodeForExpression(n)
	case *parser.StringLiteral:
		return this.generateCodeForExpression(n)
	case *parser.IdentifierLiteral:
		return this.generateCodeForExpression(n)
	case *parser.UnaryExpression:
		return this.generateCodeForExpression(n)
	case *parser.AssignmentExpression:
		return this.generateCodeForExpression(n)
	case *parser.BinaryExpression:
		return this.generateCodeForExpression(n)
	case *parser.DotMemberExpression:
		return this.generateCodeForExpression(n)
	case *parser.BracketMemberExpression:
		return this.generateCodeForExpression(n)

	default:
		panic(fmt.Sprintf("unknown node %T", node))
	}

	panic(fmt.Sprintf("unreachable %T", node))
}
