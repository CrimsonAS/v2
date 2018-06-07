package vm

import (
	"math"
	"math/rand"
)

func defineMathObject(vm *vm) value {
	mathO := newObject()
	mathO.odata.prototype = objectProto

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
	return newNumber(math.Abs(args[0].toNumber()))
}

func math_acos(vm *vm, f value, args []value) value {
	return newNumber(math.Acos(args[0].toNumber()))
}

func math_asin(vm *vm, f value, args []value) value {
	return newNumber(math.Asin(args[0].toNumber()))
}

func math_atan(vm *vm, f value, args []value) value {
	return newNumber(math.Atan(args[0].toNumber()))
}

func math_ceil(vm *vm, f value, args []value) value {
	return newNumber(math.Ceil(args[0].toNumber()))
}

func math_cos(vm *vm, f value, args []value) value {
	return newNumber(math.Cos(args[0].toNumber()))
}

func math_exp(vm *vm, f value, args []value) value {
	return newNumber(math.Exp(args[0].toNumber()))
}

func math_floor(vm *vm, f value, args []value) value {
	return newNumber(math.Floor(args[0].toNumber()))
}

func math_log(vm *vm, f value, args []value) value {
	return newNumber(math.Log(args[0].toNumber()))
}

func math_max(vm *vm, f value, args []value) value {
	ret := math.Inf(-1)
	for _, a := range args {
		n := a.toNumber()
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
		n := a.toNumber()
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
	return newNumber(math.Round(args[0].toNumber()))
}

func math_sin(vm *vm, f value, args []value) value {
	return newNumber(math.Sin(args[0].toNumber()))
}

func math_sqrt(vm *vm, f value, args []value) value {
	return newNumber(math.Sqrt(args[0].toNumber()))
}

func math_tan(vm *vm, f value, args []value) value {
	return newNumber(math.Tan(args[0].toNumber()))
}
