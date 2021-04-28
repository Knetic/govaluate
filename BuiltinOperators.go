package govaluate

import (
	"math"
)

func BuiltinOperators() map[string]Operator {
	return map[string]Operator{
		"==": builtinEq,
		"!=": builtinNeq,
		"<":  builtinLt,
		"<=": builtinLte,
		">":  builtinGt,
		">=": builtinGte,

		"&&": builtinLogicalAnd,
		"||": builtinLogicalOr,
		"!":  builtinLogicalNot,

		"+":  builtinSum,
		"-":  builtinMinus,
		"*":  builtinMul,
		"/":  builtinDiv,
		"%":  builtinMod,
		"**": builtinPow,

		"&":  builtinBitwiseAnd,
		"|":  builtinBitwiseOr,
		"^":  builtinBitwiseXor,
		"<<": builtinBitwiseLShift,
		">>": builtinBitwiseRShift,
		"~":  builtinBitwiseInverse,

		"?:": builtinTernaryIf,
		"??": builtinCoalesce,

		"array": builtinArray,
		"in":    builtinContains,
		"[]":    builtinIndexer,

		"floor": builtinFloor,
		"ceil":  builtinCeil,
		"round": builtinRound,
		"sqrt":  builtinSqrt,
		"sin":   builtinSin,
		"cos":   builtinCos,
		"max":   builtinMax,
	}
}

func builtinEq(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryArgs(ctx)
	return a == b, err
}

func builtinNeq(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryArgs(ctx)
	return a != b, err
}

func builtinLt(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a < b, err
}

func builtinLte(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a <= b, err
}

func builtinGt(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a > b, err
}

func builtinGte(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a >= b, err
}

func builtinLogicalAnd(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return nil, err
	}
	left, err := ctx.BooleanArg(0)
	if err != nil || !left {
		return false, err
	}
	return ctx.BooleanArg(1)
}

func builtinLogicalOr(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return nil, err
	}
	left, err := ctx.BooleanArg(0)
	if err != nil || left {
		return true, err
	}
	return ctx.BooleanArg(1)
}

func builtinLogicalNot(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(1); err != nil {
		return nil, err
	}
	arg, err := ctx.BooleanArg(0)
	return !arg, err
}

func builtinSum(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a + b, err
}

func builtinMinus(ctx EvalContext) (interface{}, error) {
	if ctx.ArgCount() == 1 {
		right, err := ctx.NumericArg(0)
		return -right, err
	}

	return builtinSub(ctx)
}

func builtinSub(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a - b, err
}

func builtinMul(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a * b, err
}

func builtinDiv(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return a / b, err
}

func builtinMod(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryIntegerArgs(ctx)
	return float64(a % b), err
}

func builtinPow(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return math.Pow(a, b), err
}

func builtinBitwiseAnd(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryIntegerArgs(ctx)
	return float64(a & b), err
}

func builtinBitwiseOr(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryIntegerArgs(ctx)
	return float64(a | b), err
}

func builtinBitwiseXor(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryIntegerArgs(ctx)
	return float64(a ^ b), err
}

func builtinBitwiseLShift(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryIntegerArgs(ctx)
	if b < 0 {
		return 0.0, ctx.FormatError("shift count is negative: %d", b)
	}
	return float64(a << uint(b)), err
}

func builtinBitwiseRShift(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryIntegerArgs(ctx)
	if b < 0 {
		return 0.0, ctx.FormatError("shift count is negative: %d", b)
	}
	return float64(a >> uint(b)), err
}

func builtinBitwiseInverse(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(1); err != nil {
		return nil, err
	}
	arg, err := ctx.IntegerArg(0)
	return float64(^arg), err
}

func builtinTernaryIf(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(3); err != nil {
		return nil, err
	}
	condition, err := ctx.BooleanArg(0)
	if err != nil {
		return nil, err
	}
	if condition {
		return ctx.Arg(1)
	}
	return ctx.Arg(2)
}

func builtinCoalesce(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return nil, err
	}
	left, err := ctx.Arg(0)
	if err != nil || left != nil {
		return left, err
	}
	return ctx.Arg(1)
}

func builtinArray(ctx EvalContext) (interface{}, error) {
	items := make([]interface{}, ctx.ArgCount())
	for i := 0; i < len(items); i++ {
		item, err := ctx.Arg(i)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}
	return items, nil
}

func builtinContains(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return nil, err
	}
	item, err := ctx.Arg(0)
	if err != nil {
		return nil, err
	}
	slice, err := ctx.SliceArg(1)
	if err != nil {
		return nil, err
	}
	for _, v := range slice {
		if item == v {
			return true, nil
		}
	}
	return false, nil
}

func builtinIndexer(ctx EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return nil, err
	}
	slice, err := ctx.SliceArg(0)
	if err != nil {
		return nil, err
	}
	index, err := ctx.IntegerArg(1)
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= len(slice) {
		return nil, ctx.FormatError("index out of bounds: %d, len: %d", index, len(slice))
	}
	return slice[index], nil
}

func builtinFloor(ctx EvalContext) (interface{}, error) {
	arg, err := unaryNumericArg(ctx)
	return math.Floor(arg), err
}

func builtinCeil(ctx EvalContext) (interface{}, error) {
	arg, err := unaryNumericArg(ctx)
	return math.Ceil(arg), err
}

func builtinRound(ctx EvalContext) (interface{}, error) {
	arg, err := unaryNumericArg(ctx)
	return math.Round(arg), err
}

func builtinSqrt(ctx EvalContext) (interface{}, error) {
	arg, err := unaryNumericArg(ctx)
	return math.Sqrt(arg), err
}

func builtinSin(ctx EvalContext) (interface{}, error) {
	arg, err := unaryNumericArg(ctx)
	return math.Sin(arg), err
}

func builtinCos(ctx EvalContext) (interface{}, error) {
	arg, err := unaryNumericArg(ctx)
	return math.Cos(arg), err
}

func builtinMax(ctx EvalContext) (interface{}, error) {
	a, b, err := binaryNumericArgs(ctx)
	return math.Max(a, b), err
}

func binaryArgs(ctx EvalContext) (interface{}, interface{}, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return nil, nil, err
	}
	left, err := ctx.Arg(0)
	if err != nil {
		return nil, nil, err
	}
	right, err := ctx.Arg(1)
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func binaryNumericArgs(ctx EvalContext) (float64, float64, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return 0.0, 0.0, err
	}
	left, err := ctx.NumericArg(0)
	if err != nil {
		return 0.0, 0.0, err
	}
	right, err := ctx.NumericArg(1)
	if err != nil {
		return 0.0, 0.0, err
	}
	return left, right, nil
}

func unaryNumericArg(ctx EvalContext) (float64, error) {
	if err := ctx.CheckArgCount(1); err != nil {
		return 0.0, err
	}
	return ctx.NumericArg(0)
}

func binaryIntegerArgs(ctx EvalContext) (int, int, error) {
	if err := ctx.CheckArgCount(2); err != nil {
		return 0.0, 0.0, err
	}
	left, err := ctx.IntegerArg(0)
	if err != nil {
		return 0.0, 0.0, err
	}
	right, err := ctx.IntegerArg(1)
	if err != nil {
		return 0.0, 0.0, err
	}
	return left, right, nil
}
