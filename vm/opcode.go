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

	LOAD_TEMPORARY
	STORE_TEMPORARY

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
	LESS_THAN_EQ
	GREATER_THAN
	GREATER_THAN_EQ
	EQUALS
	NOT_EQUALS
	STRICT_EQUALS
	STRICT_NOT_EQUALS
	LOGICAL_AND
	LOGICAL_OR

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
	case POP:
		return "POP"
	case LESS_THAN:
		return "<"
	case LESS_THAN_EQ:
		return "<="
	case LOGICAL_AND:
		return "&&"
	case GREATER_THAN:
		return ">"
	case GREATER_THAN_EQ:
		return ">="
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
	case LOAD_TEMPORARY:
		return fmt.Sprintf("LOAD_TEMPORARY %d", int(this.opdata))
	case STORE_TEMPORARY:
		return fmt.Sprintf("STORE_TEMPORARY %d", int(this.opdata))
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
