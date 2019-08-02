package govaluate

import (
	"fmt"
	"math"
)

func BuiltinOperators() map[string]Operator {
	return map[string]Operator{
		"==": binaryOp(func(a, b interface{}) interface{} {
			return a == b
		}),
		"!=": binaryOp(func(a, b interface{}) interface{} {
			return a != b
		}),
		">": binaryNumericOp(func(a, b float64) interface{} {
			return a > b
		}),
		"<": binaryNumericOp(func(a, b float64) interface{} {
			return a < b
		}),
		">=": binaryNumericOp(func(a, b float64) interface{} {
			return a >= b
		}),
		"<=": binaryNumericOp(func(a, b float64) interface{} {
			return a <= b
		}),
		"&&": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() != 2 {
				return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
			}
			left, err := ctx.BooleanArg(0)
			if err != nil || !left {
				return false, err
			}
			return ctx.BooleanArg(1)
		},
		"||": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() != 2 {
				return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
			}
			left, err := ctx.BooleanArg(0)
			if err != nil || left {
				return true, err
			}
			return ctx.BooleanArg(1)
		},
		"+": binaryNumericOp(func(a, b float64) interface{} {
			return a + b
		}),
		"-": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() == 1 {
				right, err := ctx.NumericArg(0)
				if err != nil {
					return nil, err
				}
				return -right, nil
			}
			if ctx.ArgCount() == 2 {
				left, err := ctx.NumericArg(0)
				if err != nil {
					return nil, err
				}
				right, err := ctx.NumericArg(1)
				if err != nil {
					return nil, err
				}
				return left - right, nil
			}
			return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
		},
		"&": binaryNumericOp(func(a, b float64) interface{} {
			return float64(int(a) & int(b))
		}),
		"|": binaryNumericOp(func(a, b float64) interface{} {
			return float64(int(a) | int(b))
		}),
		"^": binaryNumericOp(func(a, b float64) interface{} {
			return float64(int(a) ^ int(b))
		}),
		"<<": binaryNumericOp(func(a, b float64) interface{} {
			return float64(int(a) << uint(b))
		}),
		">>": binaryNumericOp(func(a, b float64) interface{} {
			return float64(int(a) >> uint(b))
		}),
		"*": binaryNumericOp(func(a, b float64) interface{} {
			return a * b
		}),
		"/": binaryNumericOp(func(a, b float64) interface{} {
			return a / b
		}),
		"%": binaryNumericOp(func(a, b float64) interface{} {
			return float64(int(a) % int(b))
		}),
		"**": binaryNumericOp(func(a, b float64) interface{} {
			return math.Pow(a, b)
		}),
		"~": unaryNumericOp(func(a float64) interface{} {
			return float64(^int(a))
		}),
		"!": unaryBooleanOp(func(a bool) interface{} {
			return !a
		}),
		"?:": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() != 3 {
				return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
			}
			condition, err := ctx.BooleanArg(0)
			if err != nil {
				return nil, err
			}
			if condition {
				return ctx.Arg(1)
			}
			return ctx.Arg(2)
		},
		"??": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() != 2 {
				return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
			}
			left, err := ctx.Arg(0)
			if err != nil || left != nil {
				return left, err
			}
			return ctx.Arg(1)
		},
		"array": func(ctx EvalContext) (interface{}, error) {
			items := make([]interface{}, ctx.ArgCount())
			for i := 0; i < len(items); i++ {
				item, err := ctx.Arg(i)
				if err != nil {
					return nil, err
				}
				items[i] = item
			}
			return items, nil
		},
		"in": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() != 2 {
				return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
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
		},
		"[]": func(ctx EvalContext) (interface{}, error) {
			if ctx.ArgCount() != 2 {
				return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
			}
			slice, err := ctx.SliceArg(0)
			if err != nil {
				return nil, err
			}
			index, err := ctx.NumericArg(1)
			if err != nil {
				return nil, err
			}
			indexInt := int(index)
			if float64(indexInt) != index {
				return nil, fmt.Errorf("invalid index: %v", index)
			}
			if indexInt < 0 || indexInt >= len(slice) {
				return nil, fmt.Errorf("index out of bounds: %d; len=%d", indexInt, len(slice))
			}
			return slice[indexInt], nil
		},
		"floor": unaryNumericOp(func(v float64) interface{} {
			return math.Floor(v)
		}),
		"ceil": unaryNumericOp(func(v float64) interface{} {
			return math.Ceil(v)
		}),
		"round": unaryNumericOp(func(v float64) interface{} {
			return math.Round(v)
		}),
		"sqrt": unaryNumericOp(func(v float64) interface{} {
			return math.Sqrt(v)
		}),
		"sin": unaryNumericOp(func(v float64) interface{} {
			return math.Sin(v)
		}),
		"cos": unaryNumericOp(func(v float64) interface{} {
			return math.Cos(v)
		}),
	}
}

func binaryOp(fn func(interface{}, interface{}) interface{}) Operator {
	return func(ctx EvalContext) (interface{}, error) {
		if ctx.ArgCount() != 2 {
			return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
		}
		left, err := ctx.Arg(0)
		if err != nil {
			return nil, err
		}
		right, err := ctx.Arg(1)
		if err != nil {
			return nil, err
		}
		return fn(left, right), nil
	}
}

func binaryNumericOp(fn func(float64, float64) interface{}) Operator {
	return func(ctx EvalContext) (interface{}, error) {
		if ctx.ArgCount() != 2 {
			return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
		}
		left, err := ctx.NumericArg(0)
		if err != nil {
			return nil, err
		}
		right, err := ctx.NumericArg(1)
		if err != nil {
			return nil, err
		}
		return fn(left, right), nil
	}
}

func unaryNumericOp(fn func(float64) interface{}) Operator {
	return func(ctx EvalContext) (interface{}, error) {
		if ctx.ArgCount() != 1 {
			return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
		}
		right, err := ctx.NumericArg(0)
		if err != nil {
			return nil, err
		}
		return fn(right), nil
	}
}

func unaryBooleanOp(fn func(bool) interface{}) Operator {
	return func(ctx EvalContext) (interface{}, error) {
		if ctx.ArgCount() != 1 {
			return nil, fmt.Errorf("wrong number of arguments: %d", ctx.ArgCount())
		}
		right, err := ctx.BooleanArg(0)
		if err != nil {
			return nil, err
		}
		return fn(right), nil
	}
}
