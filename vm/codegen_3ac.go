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
)

type tac_address struct {
	valid     bool
	constant  value
	varname   string
	temporary int
}

func (this tac_address) isConstant() bool {
	return !this.isTemp() && !this.isVar()
}

func (this tac_address) isVar() bool {
	return this.varname != "" && this.valid
}

func (this tac_address) isTemp() bool {
	return this.temporary != -1 && this.valid
}

func (this tac_address) String() string {
	if !this.valid {
		return ""
	}
	if this.varname != "" {
		return this.varname
	}

	if this.temporary != -1 {
		return fmt.Sprintf("t_%d", this.temporary)
	}

	return this.constant.String()
}

type tac struct {
	result tac_address
	arg1   tac_address
	op     tac_op_type
	arg2   tac_address
}

func (this tac) String() string {
	if this.op == TAC_ASSIGN {
		return fmt.Sprintf("%s = %s", this.result, this.arg1)
	}
	if this.op == TAC_PUSH_PARAM {
		return fmt.Sprintf("PUSH_PARAM %s", this.arg1)
	}
	if this.op == TAC_LOAD {
		return fmt.Sprintf("%s = LOAD(%s)", this.result, this.arg1)
	}
	if this.op == TAC_CALL {
		return fmt.Sprintf("%s = CALL(%s)", this.result, this.arg1)
	}
	if this.op == TAC_NEW {
		return fmt.Sprintf("%s = NEW(%s)", this.result, this.arg1)
	}
	if this.op == TAC_FUNCTION_PARAMETER {
		return fmt.Sprintf("function_param(%s)", this.arg1)
	}
	if this.op == TAC_FUNCTION {
		return fmt.Sprintf("function(%s)", this.arg1)
	}
	if this.op == TAC_END_FUNCTION {
		return fmt.Sprintf("end function(%s)", this.arg1)
	}
	if this.op == TAC_RETURN {
		return fmt.Sprintf("return(%s)", this.arg1)
	}
	if this.op == TAC_LABEL {
		return fmt.Sprintf("@%s:", this.arg1)
	}
	if this.op == TAC_JMP {
		return fmt.Sprintf("JMP @%s", this.arg1)
	}
	if this.op == TAC_JNE {
		return fmt.Sprintf("JNE %s @%s", this.arg1, this.arg2)
	}
	return fmt.Sprintf("%s = %s %s %s", this.result, this.arg1, this.op, this.arg2)
}

var temporaryIndex = -1

func newTemporary() tac_address {
	temporaryIndex += 1
	return tac_address{true, newUndefined(), "", temporaryIndex}
}

func newConstant(v value) tac_address {
	return tac_address{true, v, "", -1}
}

func newVar(n string) tac_address {
	return tac_address{true, newUndefined(), n, -1}
}

type tac_op_type int

//go:generate stringer -type=tac_op_type
const (
	// Simple math operators
	TAC_ADD tac_op_type = iota
	TAC_SUB
	TAC_MULTIPLY
	TAC_DIVIDE
	TAC_MODULUS
	TAC_LEFT_SHIFT
	TAC_RIGHT_SHIFT
	TAC_UNSIGNED_RIGHT_SHIFT
	TAC_BITWISE_AND // a & b
	TAC_BITWISE_XOR // a ^ b
	TAC_BITWISE_OR  // a | b

	TAC_UPLUS       // +a
	TAC_UMINUS      // -a
	TAC_UNOT        // !a
	TAC_TYPEOF      // typeof a
	TAC_BITWISE_NOT // ~a

	TAC_ASSIGN // =

	TAC_PUSH_PARAM
	TAC_CALL
	TAC_NEW
	TAC_LOAD
	TAC_LOAD_MEMBER

	TAC_LESS_THAN
	TAC_GREATER_THAN
	TAC_EQUALS
	TAC_NOT_EQUALS
	TAC_LESS_THAN_EQ
	TAC_LOGICAL_AND

	TAC_FUNCTION_PARAMETER
	TAC_FUNCTION
	TAC_END_FUNCTION
	TAC_RETURN

	TAC_JNE
	TAC_LABEL
	TAC_JMP
)

func pushVarOrConstant(addr tac_address) []opcode {
	codebuf := []opcode{}

	if addr.isVar() {
		rhsIdx := float64(appendStringtable(addr.varname))
		codebuf = append(codebuf, newOpcode(LOAD, rhsIdx))
	} else if addr.isConstant() {
		codebuf = append(codebuf, pushConstant(addr)...)
	}

	return codebuf
}

func pushConstant(addr tac_address) []opcode {
	codebuf := []opcode{}
	switch c := addr.constant.(type) {
	case valueNumber:
		codebuf = append(codebuf, newOpcode(PUSH_NUMBER, float64(c)))
	default:
		panic(fmt.Sprintf("Unknown constant type %T", addr))
	}
	return codebuf
}

func maybePushStore(result tac_address) []opcode {
	codebuf := []opcode{}
	if result.isVar() {
		varIdx := appendStringtable(result.varname)
		codebuf = append(codebuf, newOpcode(STORE, float64(varIdx)))
	}
	return codebuf
}

func callJsFunction(this *vm, params []valueString, addr int) func(vm *vm, f value, args []value) value {
	// Small optimisation: intern strings at codegen time, so we don't have to
	// hash at runtime.
	intArgs := []int{}
	for _, arg := range params {
		intArgs = append(intArgs, appendStringtable(arg.String()))
	}

	return func(vm *vm, f value, args []value) value {
		if execDebug {
			//log.Printf("Calling func! IP %d going to %d, %s", vm.ip, addr, args)
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
				//log.Printf("Defining var %s %s", stringtable[arg], v)
			}
			vm.defineVar(arg, v)
		}

		return newUndefined()
	}
}

func (this *vm) generateBytecode(in []tac) []opcode {
	codebuf := []opcode{}

	type labelInfo struct {
		bytecodeOffset int
	}
	labels := make(map[tac_address]labelInfo)
	type jumpInfo struct {
		label          tac_address
		bytecodeOffset int
	}
	jumps := []jumpInfo{}
	paramNames := []valueString{}
	paramCount := 0

	for idx, op := range in {
		switch op.op {
		case TAC_PUSH_PARAM:
			paramCount++
		case TAC_CALL:
			codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			codebuf = append(codebuf, newOpcode(CALL, float64(paramCount)))
			paramCount = 0
		case TAC_NEW:
			codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			codebuf = append(codebuf, newOpcode(NEW, float64(paramCount)))
			paramCount = 0
		case TAC_FUNCTION_PARAMETER:
			if !op.arg1.isConstant() {
				panic("Not a constant???")
			}
			paramNames = append(paramNames, op.arg1.constant.(valueString))
		case TAC_FUNCTION:
			// Gather all local declarations
			funcIdx := appendStringtable(op.arg1.constant.String())

			runBuiltin := callJsFunction(this, paramNames, len(codebuf))
			callFn := newFunctionObject(runBuiltin, runBuiltin)
			this.defineVar(funcIdx, callFn)
			paramNames = paramNames[:]

			codebuf = append(codebuf, newOpcode(IN_FUNCTION, float64(funcIdx)))
			declaredVars := make(map[int]bool)
			for _, nop := range in[idx:] {
				if nop.op == TAC_END_FUNCTION && nop.arg1 == op.arg1 {
					break
				}

				if nop.result.isVar() {
					varIdx := appendStringtable(nop.result.varname)
					if _, ok := declaredVars[varIdx]; !ok {
						declaredVars[varIdx] = true
						codebuf = append(codebuf, newOpcode(DECLARE, float64(varIdx)))
					}
				}
			}

		case TAC_END_FUNCTION:
			// ignore for now
		case TAC_RETURN:
			if op.arg1.valid {
				codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			}
			codebuf = append(codebuf, simpleOp(RETURN))
		case TAC_LABEL:
			labels[op.arg1] = labelInfo{bytecodeOffset: len(codebuf)}
		case TAC_ASSIGN:
			codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			codebuf = append(codebuf, maybePushStore(op.result)...)
		case TAC_LESS_THAN:
			codebuf = append(codebuf, pushVarOrConstant(op.arg2)...)
			codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			codebuf = append(codebuf, simpleOp(LESS_THAN))
			codebuf = append(codebuf, maybePushStore(op.result)...)
		case TAC_JNE:
			codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			jumps = append(jumps, jumpInfo{label: op.arg2, bytecodeOffset: len(codebuf)})
			codebuf = append(codebuf, newOpcode(JNE, 0))
		case TAC_JMP:
			jumps = append(jumps, jumpInfo{label: op.arg1, bytecodeOffset: len(codebuf)})
			codebuf = append(codebuf, newOpcode(JMP, 0))
		case TAC_ADD:
			codebuf = append(codebuf, pushVarOrConstant(op.arg2)...)
			codebuf = append(codebuf, pushVarOrConstant(op.arg1)...)
			codebuf = append(codebuf, simpleOp(ADD))
			codebuf = append(codebuf, maybePushStore(op.result)...)
		default:
			panic(fmt.Sprintf("unknown tac %s", op))
		}
	}

	for _, jmp := range jumps {
		codebuf[jmp.bytecodeOffset].opdata = opdata(labels[jmp.label].bytecodeOffset - jmp.bytecodeOffset - 1)
	}

	return codebuf
}

var funcsToDefine []*parser.FunctionExpression

func generateCodeTAC(node parser.Node, retcodebuf *[]tac) tac_address {
	codebuf := []tac{}
	retaddr := tac_address{}

	switch n := node.(type) {
	case *parser.Program:
		codebuf = append(codebuf, tac{arg1: newConstant(newString("%main")), op: TAC_FUNCTION})
		for _, s := range n.Body() {
			generateCodeTAC(s, &codebuf)
		}
		codebuf = append(codebuf, tac{op: TAC_RETURN})
		codebuf = append(codebuf, tac{arg1: newConstant(newString("%main")), op: TAC_END_FUNCTION})

		for _, afunc := range funcsToDefine {
			for _, p := range afunc.Parameters {
				codebuf = append(codebuf, tac{arg1: newConstant(newString(p.String())), op: TAC_FUNCTION_PARAMETER})
			}
			codebuf = append(codebuf, tac{arg1: newConstant(newString(afunc.Identifier.String())), op: TAC_FUNCTION})
			generateCodeTAC(afunc.Body, &codebuf)
			codebuf = append(codebuf, tac{op: TAC_RETURN})
			codebuf = append(codebuf, tac{arg1: newConstant(newString(afunc.Identifier.String())), op: TAC_END_FUNCTION})
		}
	case *parser.VariableStatement:
		for idx, _ := range n.Vars {
			v := n.Vars[idx]
			i := n.Initializers[idx]

			if i != nil {
				exp := generateCodeTAC(i, &codebuf)
				codebuf = append(codebuf, tac{result: newVar(v.String()), arg1: exp, op: TAC_ASSIGN})
			} else {
				codebuf = append(codebuf, tac{result: newVar(v.String()), arg1: newConstant(newUndefined()), op: TAC_ASSIGN})
			}
		}
	case *parser.ExpressionStatement:
		generateCodeTAC(n.X, &codebuf)
	case *parser.ReturnStatement:
		if n.X != nil {
			retaddr = generateCodeTAC(n.X, &codebuf)
			codebuf = append(codebuf, tac{arg1: retaddr, op: TAC_RETURN})
		}
	case *parser.ForStatement:
		if n.Initializer != nil {
			generateCodeTAC(n.Initializer, &codebuf)
		}
		lbl := newTemporary()
		endLbl := newTemporary()
		codebuf = append(codebuf, tac{arg1: lbl, op: TAC_LABEL})
		if n.Test != nil {
			test := generateCodeTAC(n.Test, &codebuf)
			codebuf = append(codebuf, tac{op: TAC_JNE, arg1: test, arg2: endLbl})
		}
		if n.Update != nil {
			generateCodeTAC(n.Update, &codebuf)
		}
		generateCodeTAC(n.Body, &codebuf)
		codebuf = append(codebuf, tac{arg1: lbl, op: TAC_JMP})
		codebuf = append(codebuf, tac{arg1: endLbl, op: TAC_LABEL})
	//case *parser.DoWhileStatement:
	//case *parser.WhileStatement:
	//case *parser.ConditionalExpression:
	//case *parser.IfStatement:
	case *parser.BlockStatement:
		for _, s := range n.Body {
			generateCodeTAC(s, &codebuf)
		}
	case *parser.EmptyStatement:

	//case *parser.ArrayLiteral:
	//case *parser.ObjectLiteral:
	//case *parser.ThisLiteral:
	case *parser.StringLiteral:
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, arg1: newConstant(newString(n.String())), op: TAC_ASSIGN})
	case *parser.IdentifierLiteral:
		return newVar(n.String())
	case *parser.NumericLiteral:
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, arg1: newConstant(newNumber(n.Float64Value())), op: TAC_ASSIGN})
	case *parser.TrueLiteral:
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, arg1: newConstant(newBool(true)), op: TAC_ASSIGN})
	case *parser.FalseLiteral:
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, arg1: newConstant(newBool(false)), op: TAC_ASSIGN})
	case *parser.NullLiteral:
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, arg1: newConstant(newNull()), op: TAC_ASSIGN})

	//case *parser.SequenceExpression:
	case *parser.FunctionExpression:
		funcsToDefine = append(funcsToDefine, n)
	case *parser.NewExpression:
		c := n.X.(*parser.CallExpression)
		for _, arg := range c.Arguments {
			param := generateCodeTAC(arg, &codebuf)
			codebuf = append(codebuf, tac{op: TAC_PUSH_PARAM, arg1: param})
		}
		fid := generateCodeTAC(c.X, &codebuf)
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, op: TAC_NEW, arg1: fid})
	case *parser.CallExpression:
		for _, arg := range n.Arguments {
			param := generateCodeTAC(arg, &codebuf)
			codebuf = append(codebuf, tac{op: TAC_PUSH_PARAM, arg1: param})
		}
		fid := generateCodeTAC(n.X, &codebuf)
		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, op: TAC_CALL, arg1: fid})
	case *parser.UnaryExpression:
		if n.IsPrefix() {
			uref := generateCodeTAC(n.X, &codebuf)
			switch n.Operator() {
			case parser.PLUS:
				retaddr = newTemporary()
				codebuf = append(codebuf, tac{result: retaddr, arg1: uref, op: TAC_ASSIGN})
			case parser.MINUS:
				retaddr = newTemporary()
				codebuf = append(codebuf, tac{result: retaddr, arg1: newConstant(newNumber(0)), op: TAC_SUB, arg2: uref})
			case parser.INCREMENT:
				// ### code generation is likely wrong here
				retaddr = newTemporary()
				codebuf = append(codebuf, tac{result: retaddr, arg1: uref, op: TAC_ADD, arg2: newConstant(newNumber(1))})
				codebuf = append(codebuf, tac{result: uref, arg1: retaddr, op: TAC_ASSIGN})
			default:
				panic(fmt.Sprintf("Unhandled prefix op %s", n.Operator()))
			}
		} else {
			panic(fmt.Sprintf("Unhandled postfix op %s", n.Operator()))
		}
	case *parser.AssignmentExpression:
		rhs := generateCodeTAC(n.Right, &codebuf)
		retaddr = newVar(n.Left.(*parser.IdentifierLiteral).String())
		codebuf = append(codebuf, tac{result: retaddr, arg1: rhs, op: TAC_ASSIGN})
		/*
			var realOp tac_op_type
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
			}*/
	case *parser.BinaryExpression:
		leftRef := generateCodeTAC(n.Left, &codebuf)
		rightRef := generateCodeTAC(n.Right, &codebuf)
		retaddr = newTemporary()

		var realOp tac_op_type
		switch n.Operator() {
		case parser.PLUS:
			realOp = TAC_ADD
		case parser.MINUS:
			realOp = TAC_SUB
		case parser.MULTIPLY:
			realOp = TAC_MULTIPLY
		case parser.DIVIDE:
			realOp = TAC_DIVIDE
		case parser.LEFT_SHIFT:
			realOp = TAC_LEFT_SHIFT
		case parser.RIGHT_SHIFT:
			realOp = TAC_RIGHT_SHIFT
		case parser.UNSIGNED_RIGHT_SHIFT:
			realOp = TAC_UNSIGNED_RIGHT_SHIFT
		case parser.BITWISE_AND:
			realOp = TAC_BITWISE_AND
		case parser.BITWISE_XOR:
			realOp = TAC_BITWISE_XOR
		case parser.BITWISE_OR:
			realOp = TAC_BITWISE_OR
		case parser.MODULUS:
			realOp = TAC_MODULUS
		case parser.LESS_THAN:
			realOp = TAC_LESS_THAN
		case parser.GREATER_THAN:
			realOp = TAC_GREATER_THAN
		case parser.EQUALS:
			realOp = TAC_EQUALS
		case parser.NOT_EQUALS:
			realOp = TAC_NOT_EQUALS
		case parser.LESS_EQ:
			realOp = TAC_LESS_THAN_EQ
		case parser.LOGICAL_AND:
			realOp = TAC_LOGICAL_AND
		default:
			panic(fmt.Sprintf("unknown operator %s", n.Operator()))
		}

		codebuf = append(codebuf, tac{result: retaddr, arg1: leftRef, op: realOp, arg2: rightRef})

	case *parser.DotMemberExpression:
		fid := generateCodeTAC(n.X, &codebuf)
		base := newTemporary()
		codebuf = append(codebuf, tac{result: base, op: TAC_LOAD, arg1: fid})

		retaddr = newTemporary()
		codebuf = append(codebuf, tac{result: retaddr, op: TAC_LOAD_MEMBER, arg1: base, arg2: newConstant(newString(n.Name.String()))})
	//case *parser.BracketMemberExpression:

	default:
		panic(fmt.Sprintf("unknown node %T", node))
	}

	*retcodebuf = append(*retcodebuf, codebuf...)
	return retaddr
}

func optimizeTAC(codebuf *[]tac) {
	return
	for i := 0; i < 50; i++ {
		removeDeadCode(codebuf)
		simplifyExpression(codebuf)
		copyPropagation(codebuf)
	}
}

func removeDeadCode(retcodebuf *[]tac) {
	codebuf := *retcodebuf

	// look for temporaries that aren't subsequently used as an RHS
	// and remove them.
	for idx, op := range codebuf {
		if op.result.isTemp() {
			// Don't eliminate temps that have side effects (function calls)
			if op.op == TAC_CALL || op.op == TAC_NEW {
				continue
			}

			found := false
			for _, nop := range codebuf[idx:] {
				if nop.arg1 == op.result || nop.arg2 == op.result {
					found = true
				}
			}
			if !found {
				codebuf = append(codebuf[0:idx], codebuf[idx+1:]...)
			}
		}
	}

	// optimize:
	// return; return -> return
	// this probably won't crop up in user code, but it does crop up in our own
	// code generation, as we surround function expressions with a begin/end
	// function instruction.
	remove := false
	for idx, op := range codebuf {
		if op.op == TAC_RETURN {
			if remove {
				codebuf = append(codebuf[0:idx], codebuf[idx+1:]...)
				remove = false
			} else {
				if idx < len(codebuf)-1 && codebuf[idx+1].op == TAC_RETURN {
					remove = true
				}
			}
		}
	}

	// optimize:
	// t_N = a * 5
	// t_N+1 = t_N
	//
	// to:
	// t_N+1 = a * 5
	for idx, op := range codebuf {
		if op.result.isTemp() {
			// Don't perform this optimization right now, as it breaks with dead code
			// elimination. If we optimize t_N = NEW, a = t_N to a = NEW, then
			// we also need to remove the old new instruction below...
			if op.op == TAC_CALL || op.op == TAC_NEW {
				continue
			}

			foundIdx := -1
			foundCount := 0
			for nidx, nop := range codebuf[idx:] {
				if nop.op == TAC_ASSIGN && nop.arg1 == op.result {
					foundIdx = idx + nidx
					foundCount++
				}
			}

			if foundCount == 1 {
				codebuf[foundIdx].arg1 = op.arg1
				codebuf[foundIdx].arg2 = op.arg2
				codebuf[foundIdx].op = op.op
				// We could remove the original instruction here, but dead code
				// elimination will take it for us.
			}
		}
	}

	*retcodebuf = codebuf
}

func simplifyExpression(retcodebuf *[]tac) {
	codebuf := *retcodebuf
	zero := newConstant(newNumber(0))
	for idx, op := range codebuf {
		// lhs = rhs * 0 -> lhs = 0
		if op.op == TAC_MULTIPLY {
			if op.arg1 == zero {
				op.op = TAC_ASSIGN
				op.arg1 = zero
				op.arg2 = tac_address{}
				codebuf[idx] = op
			} else if op.arg2 == zero {
				op.op = TAC_ASSIGN
				op.arg1 = zero
				op.arg2 = tac_address{}
				codebuf[idx] = op
			}
		}

		// lhs = rhs + 0 -> lhs = rhs
		if op.op == TAC_ADD || op.op == TAC_SUB {
			if op.arg1 == zero {
				op.op = TAC_ASSIGN
				op.arg1 = op.arg2
				op.arg2 = tac_address{}
				codebuf[idx] = op
			}
			if op.arg2 == zero {
				op.op = TAC_ASSIGN
				op.arg2 = tac_address{}
				codebuf[idx] = op
			}
		}
	}
	*retcodebuf = codebuf
}

func copyPropagation(retcodebuf *[]tac) {
	codebuf := *retcodebuf
	for idx, op := range codebuf {
		// t_1 = 5, t_2 = t_1 ADD 0 => t_2 = 5 ADD 0
		if op.op == TAC_ASSIGN && op.result.isTemp() {
			for nidx := 0; nidx < len(codebuf)-idx; nidx++ {
				nop := codebuf[idx+nidx]
				if nop.arg1 == op.result {
					nop.arg1 = op.arg1
					codebuf[idx+nidx] = nop
				}
				if nop.arg2 == op.result {
					nop.arg2 = op.arg1
					codebuf[idx+nidx] = nop
				}
			}
		}
	}
	*retcodebuf = codebuf
}
