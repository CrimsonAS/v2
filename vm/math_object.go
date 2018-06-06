package vm

import (
	"math"
	"math/rand"
)

func defineMathObject() value {
	mathO := newObject()

	mathO.set("E", newNumber(2.7182818284590452354))
	mathO.set("LN10", newNumber(2.302585092994046))
	mathO.set("LN2", newNumber(0.6931471805599453))
	mathO.set("LOG2E", newNumber(1.4426950408889634))
	mathO.set("LO10E", newNumber(0.4342944819032518))
	mathO.set("PI", newNumber(3.1415926535897932))
	mathO.set("SQRT1_2", newNumber(0.7071067811865476))
	mathO.set("SQRT2", newNumber(1.4142135623730951))

	mathO.set("abs", newFunctionObject(math_abs, nil))
	mathO.set("acos", newFunctionObject(math_acos, nil))
	mathO.set("asin", newFunctionObject(math_asin, nil))
	mathO.set("atan", newFunctionObject(math_atan, nil))
	// atan2
	mathO.set("ceil", newFunctionObject(math_ceil, nil))
	mathO.set("cos", newFunctionObject(math_cos, nil))
	mathO.set("exp", newFunctionObject(math_exp, nil))
	mathO.set("floor", newFunctionObject(math_floor, nil))
	mathO.set("log", newFunctionObject(math_log, nil))
	mathO.set("max", newFunctionObject(math_max, nil))
	mathO.set("min", newFunctionObject(math_min, nil))
	// pow
	mathO.set("random", newFunctionObject(math_random, nil))
	mathO.set("round", newFunctionObject(math_round, nil))
	mathO.set("sin", newFunctionObject(math_sin, nil))
	mathO.set("sqrt", newFunctionObject(math_sqrt, nil))
	mathO.set("tan", newFunctionObject(math_tan, nil))
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
