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
	"math"
	"math/rand"
)

func defineMathObject(vm *vm) valueBasicObject {
	mathO := valueBasicObject{&rootObjectData{&valueBasicObjectData{extensible: true}}}

	mathO.defineReadonlyProperty(vm, "E", newNumber(2.7182818284590452354), 1)
	mathO.defineReadonlyProperty(vm, "LN10", newNumber(2.302585092994046), 1)
	mathO.defineReadonlyProperty(vm, "LN2", newNumber(0.6931471805599453), 1)
	mathO.defineReadonlyProperty(vm, "LOG2E", newNumber(1.4426950408889634), 1)
	mathO.defineReadonlyProperty(vm, "LO10E", newNumber(0.4342944819032518), 1)
	mathO.defineReadonlyProperty(vm, "PI", newNumber(3.1415926535897932), 1)
	mathO.defineReadonlyProperty(vm, "SQRT1_2", newNumber(0.7071067811865476), 1)
	mathO.defineReadonlyProperty(vm, "SQRT2", newNumber(1.4142135623730951), 1)

	mathO.defineDefaultProperty(vm, "abs", newFunctionObject(math_abs, nil), 1)
	mathO.defineDefaultProperty(vm, "acos", newFunctionObject(math_acos, nil), 1)
	mathO.defineDefaultProperty(vm, "asin", newFunctionObject(math_asin, nil), 1)
	mathO.defineDefaultProperty(vm, "atan", newFunctionObject(math_atan, nil), 1)
	// atan2
	mathO.defineDefaultProperty(vm, "ceil", newFunctionObject(math_ceil, nil), 1)
	mathO.defineDefaultProperty(vm, "cos", newFunctionObject(math_cos, nil), 1)
	mathO.defineDefaultProperty(vm, "exp", newFunctionObject(math_exp, nil), 1)
	mathO.defineDefaultProperty(vm, "floor", newFunctionObject(math_floor, nil), 1)
	mathO.defineDefaultProperty(vm, "log", newFunctionObject(math_log, nil), 1)
	mathO.defineDefaultProperty(vm, "max", newFunctionObject(math_max, nil), 1)
	mathO.defineDefaultProperty(vm, "min", newFunctionObject(math_min, nil), 1)
	// pow
	mathO.defineDefaultProperty(vm, "random", newFunctionObject(math_random, nil), 1)
	mathO.defineDefaultProperty(vm, "round", newFunctionObject(math_round, nil), 1)
	mathO.defineDefaultProperty(vm, "sin", newFunctionObject(math_sin, nil), 1)
	mathO.defineDefaultProperty(vm, "sqrt", newFunctionObject(math_sqrt, nil), 1)
	mathO.defineDefaultProperty(vm, "tan", newFunctionObject(math_tan, nil), 1)
	return mathO
}

func math_abs(vm *vm, f value, args []value) value {
	return newNumber(math.Abs(args[0].ToNumber()))
}

func math_acos(vm *vm, f value, args []value) value {
	return newNumber(math.Acos(args[0].ToNumber()))
}

func math_asin(vm *vm, f value, args []value) value {
	return newNumber(math.Asin(args[0].ToNumber()))
}

func math_atan(vm *vm, f value, args []value) value {
	return newNumber(math.Atan(args[0].ToNumber()))
}

func math_ceil(vm *vm, f value, args []value) value {
	return newNumber(math.Ceil(args[0].ToNumber()))
}

func math_cos(vm *vm, f value, args []value) value {
	return newNumber(math.Cos(args[0].ToNumber()))
}

func math_exp(vm *vm, f value, args []value) value {
	return newNumber(math.Exp(args[0].ToNumber()))
}

func math_floor(vm *vm, f value, args []value) value {
	return newNumber(math.Floor(args[0].ToNumber()))
}

func math_log(vm *vm, f value, args []value) value {
	return newNumber(math.Log(args[0].ToNumber()))
}

func math_max(vm *vm, f value, args []value) value {
	ret := math.Inf(-1)
	for _, a := range args {
		n := a.ToNumber()
		if math.IsNaN(n) {
			return newNumber(math.NaN())
		}
		ret = math.Max(ret, n)
	}
	return newNumber(ret)
}

func math_min(vm *vm, f value, args []value) value {
	ret := math.Inf(+1)
	for _, a := range args {
		n := a.ToNumber()
		if math.IsNaN(n) {
			return newNumber(math.NaN())
		}
		ret = math.Min(ret, n)
	}
	return newNumber(ret)
}

func math_random(vm *vm, f value, args []value) value {
	return newNumber(rand.Float64())
}

func math_round(vm *vm, f value, args []value) value {
	return newNumber(math.Round(args[0].ToNumber()))
}

func math_sin(vm *vm, f value, args []value) value {
	return newNumber(math.Sin(args[0].ToNumber()))
}

func math_sqrt(vm *vm, f value, args []value) value {
	return newNumber(math.Sqrt(args[0].ToNumber()))
}

func math_tan(vm *vm, f value, args []value) value {
	return newNumber(math.Tan(args[0].ToNumber()))
}
