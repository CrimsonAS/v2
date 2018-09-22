/*
 * Copyright 2018 Crimson AS <info@crimson.no>
 * Author: Robin Burchell <robin@crimson.no>
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
	"github.com/stvp/assert"
	"testing"
)

type simpleVMTest struct {
	in  string
	out value
}

func runSimpleVMTestHelper(t *testing.T, tests []simpleVMTest) {
	for _, test := range tests {
		t.Logf("Testing: %s", test.in)
		vm := New(test.in)
		assert.Equal(t, vm.Run(), test.out)
		t.Logf("** Passed %s == %s", test.in, test.out)
	}
}

func TestStrings(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "return \"hello\"",
			out: newString("hello"),
		},
		simpleVMTest{
			in:  "return \"hello\"+\"world\"",
			out: newString("helloworld"),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestPostfixOperators(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; var b; b = a++; return b",
			out: newNumber(0),
		},
		simpleVMTest{
			in:  "var a = 1; var b; b = a--; return b",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; var b = 1; a = b++; return a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; var b = 1; a = b++; return b",
			out: newNumber(2),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestAssignmentOperators(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; a = 1; return a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; a += 1; return a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 1; a /= 2; return a",
			out: newNumber(0.5),
		},
		simpleVMTest{
			in:  "var a = 1; a *= 2; return a",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var a = 1; a *= 2; return a",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var a = 5; a %= 2; return a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 5; a <<= 2; return a",
			out: newNumber(20),
		},
		simpleVMTest{
			in:  "var a = 5; a >>= 1; return a",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var a = 5; a >>>= 1; return a",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var a = 55; a &= 123124; return a",
			out: newNumber(52),
		},
		simpleVMTest{
			in:  "var a = 55; a ^= 123124; return a",
			out: newNumber(123075),
		},
		simpleVMTest{
			in:  "var a = 55; a |= 123124; return a",
			out: newNumber(123127),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestPrefixOperators(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; var b; b = ++a; return b",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; var b; b = --a; return b",
			out: newNumber(-1),
		},
		simpleVMTest{
			in:  "return !false",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return !1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return !!1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return !!!1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return +3",
			out: newNumber(3),
		},
		simpleVMTest{
			in:  "return -3",
			out: newNumber(-3),
		},
		simpleVMTest{
			in:  "return ~500",
			out: newNumber(-501),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestSimple(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "return undefined",
			out: newUndefined(),
		},
		simpleVMTest{
			in:  "return null",
			out: newNull(),
		},
		simpleVMTest{
			in:  "return 2+3",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "return (2+3)",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "return (1+1)",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "return (2+2)*(2+2)",
			out: newNumber(16),
		},
		simpleVMTest{
			in:  "return 10/2",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "if (true) { return 10/2 }",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "if (false) { return 10/2 }",
			out: newUndefined(),
		},
		simpleVMTest{
			in:  "if (true) { return 10/2 }",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "return 2+2*2+2",
			out: newNumber(8),
		},
		simpleVMTest{
			in:  "return 1<2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 1<=2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 2<=1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 2<1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 2>1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 1>2",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 1>=2",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 2>=1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 2>=2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 1==2",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 1==1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 1!=2",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 1!=1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 1&&0",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 0&&0",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 0&&1",
			out: newBool(false),
		},
		simpleVMTest{
			in:  "return 1&&1",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "return 1,2",
			out: newNumber(2),
		},
		simpleVMTest{
			in:  "var a; var b; return a=0, b=2",
			out: newNumber(2), // this will fail with the old compiler.
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestCall(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "return f() function f() { return 5 }",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestThis(t *testing.T) {
	// Roundabout way of checking that 'this' actually works, since returning
	// 'this' gives us no easy way to check it's the right thing...
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "function f() { return this.a }; f.a = 42; return f();",
			out: newNumber(42),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestTypeof(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var v = undefined; return typeof v",
			out: newString("undefined"),
		},
		simpleVMTest{
			in:  "var v = null; return typeof v",
			out: newString("object"),
		},
		simpleVMTest{
			in:  "var v = true; return typeof v",
			out: newString("boolean"),
		},
		simpleVMTest{
			in:  "var v = 5; return typeof v",
			out: newString("number"),
		},
		simpleVMTest{
			in:  "var v = 'test'; return typeof v",
			out: newString("string"),
		},
		simpleVMTest{
			in:  "var v = new Array(); return typeof v",
			out: newString("object"),
		},
		simpleVMTest{
			// ### these should be equivilent, but are not.
			//in:  "var v = function() {}; return typeof v",
			in:  "function v() {}; return typeof v",
			out: newString("function"),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestVar(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 5, b; b = a + 10; return b",
			out: newNumber(15),
		},
		simpleVMTest{
			in:  "var a = 5; return a",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "var a = 5; var a; return a",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestWhile(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; while (a > 10) a = a + 1; return a",
			out: newNumber(0),
		},
		simpleVMTest{
			in:  "var a = 0; while (a < 10) a = a + 1; return a",
			out: newNumber(10),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestDoWhile(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; do { a = a + 1 } while (a > 10); return a",
			out: newNumber(1),
		},
		simpleVMTest{
			in:  "var a = 0; do { a = a + 1 } while (a < 5); return a",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestForStatement(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 0; for (a = 0; a < 5; a = a + 1) { }; return a",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "var a = 10; for (; a < 5; a = a + 1) { }; return a",
			out: newNumber(10),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestReturnStatement(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "function f() { return; } var a; a = f(); return a;",
			out: newUndefined(),
		},
		simpleVMTest{
			in:  "function f() { return 10; } var a; a = f(); return a;",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function a() { return 10; } function b() { return 5; } var c = a(); return c;",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function a() { return 10; } function b() { return 5; } var c = b(); return c;",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "function a() { return 10; } var b = new a(); return b;",
			out: newNumber(10),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestBuiltinFunction(t *testing.T) {
	{
		testFunc := func(vm *vm, f value, args []value) value {
			t.Logf("I'm a message from a builtin function: %+v", args)
			return newString("Hello world")
		}

		vm := New("return testFunc()")
		pf := newFunctionObject(testFunc, nil)
		vm.defineVar(appendStringtable("testFunc"), pf)
		assert.Equal(t, vm.Run(), newString("Hello world"))
	}
	t.Logf("Test one passed")

	// test multiple arguments, and their order
	{
		testFunc := func(vm *vm, f value, args []value) value {
			return newString(args[0].String() + args[1].String())
		}

		vm := New("return testFunc(\"Hello\", \"World\")")
		pf := newFunctionObject(testFunc, nil)
		vm.defineVar(appendStringtable("testFunc"), pf)
		assert.Equal(t, vm.Run(), newString("HelloWorld"))
	}
	t.Logf("Test two passed")

	// test call vs construct
	{
		testCall := func(vm *vm, f value, args []value) value {
			return newNumber(10)
		}
		testConstruct := func(vm *vm, f value, args []value) value {
			return newNumber(20)
		}

		{
			vm := New("return testFunc()")
			pf := newFunctionObject(testCall, testConstruct)
			vm.defineVar(appendStringtable("testFunc"), pf)
			assert.Equal(t, vm.Run(), newNumber(10))
		}
		t.Logf("Call passed")
		{
			vm := New("return new testFunc()")
			pf := newFunctionObject(testCall, testConstruct)
			vm.defineVar(appendStringtable("testFunc"), pf)
			assert.Equal(t, vm.Run(), newNumber(20))
		}
		t.Logf("New passed")
	}
	t.Logf("Test call/construct passed")
}

func TestJSFunction(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "function f(a) { return a; } var b; b = f(10); return b",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a) { return a; } var a; a = f(10); return a",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a, b) { return a; } var a; a = f(10, 5); return a",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a, b) { return b; } var a; a = f(5, 10); return a",
			out: newNumber(10),
		},
		simpleVMTest{
			in:  "function f(a, b) { return a; } var a; a = f(5, 10); return a",
			out: newNumber(5),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestRecursiveLookups(t *testing.T) {
	vm := New("function f(a) { if (a > 3) return a; a = a + 1; return f(a); } var n = f(0); return n")
	ret := vm.Run()
	assert.Equal(t, ret, newNumber(4))
}

func TestFibonnaci(t *testing.T) {
	f := "function fibonacci(n) {\n"
	f += "	var a = 0, b = 1, f = 1;\n"
	f += "	for(var i = 2; i <= n; i++) {\n"
	f += "		f = a + b;\n"
	f += "		a = b;\n"
	f += "		b = f;\n"
	f += "	}\n"
	f += "	return f;\n"
	f += "};\n"
	f += "return fibonacci(10)\n"

	var iterative value
	{
		vm := New(f)
		iterative = vm.Run()
	}

	f = ""
	f = "function fibonacci(n) {\n"
	f += "    if (n < 1) {\n"
	f += "        return 0\n"
	f += "    } else if (n <= 2) {\n"
	f += "        return 1\n"
	f += "    } else {\n"
	f += "        return fibonacci(n - 1) + fibonacci(n - 2)\n"
	f += "    }\n"
	f += "}\n"
	f += "return fibonacci(10)\n"

	var recursive value
	{
		vm := New(f)
		recursive = vm.Run()
	}

	assert.Equal(t, iterative, recursive)
	assert.Equal(t, iterative, newNumber(55))
}

func TestConditionalExpression(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = 1; var b = a == 1 ? 2 : 3; return b",
			out: newNumber(2),
		},
		// ### failing, which seems odd...
		//simpleVMTest{
		//	in:  "var a = 1; var b = a == 1 ? 2 : 3; return b",
		//	out: newNumber(2),
		//},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestObjectProperties(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = {}; return a.b;",
			out: newUndefined(),
		},
		simpleVMTest{
			in:  "var a = {b: 5}; return a.b;",
			out: newNumber(5),
		},
		simpleVMTest{
			in:  "var a = {b: 5}; a.b = 6; return a.b;",
			out: newNumber(6),
		},
	}

	runSimpleVMTestHelper(t, tests)
}

func TestValueStrictEquality(t *testing.T) {
	tests := []simpleVMTest{
		simpleVMTest{
			in:  "var a = true; var b = true; return a === b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = false; var b = false; return a === b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = false; var b = true; return a !== b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = 1.5; var b = 1.5; return a === b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = 1.5; var b = 1.6; return a !== b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = 'hello'; var b = 'hello'; return a === b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = 'hello'; var b = 'world'; return a !== b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = {}; var b = {}; return a !== b",
			out: newBool(true),
		},
		simpleVMTest{
			in:  "var a = []; var b = []; return a !== b",
			out: newBool(true),
		},
	}

	runSimpleVMTestHelper(t, tests)
}
