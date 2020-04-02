package govaluate

import "fmt"

type Operator func(ctx EvalContext) (interface{}, error)

type EvalContext struct {
	params EvalParams
	expr   ExprNode
}

func (ctx EvalContext) ArgCount() int {
	return len(ctx.expr.Args)
}

func (ctx EvalContext) CheckArgCount(count int) error {
	if ctx.ArgCount() != count {
		return ctx.FormatError("wrong number of arguments: %d, expected: %d", ctx.ArgCount(), count)
	}
	return nil
}

func (ctx EvalContext) Arg(idx int) (interface{}, error) {
	args := ctx.expr.Args
	if idx >= len(args) {
		return nil, ctx.FormatError("requested argument #%d, but argument count is %d", idx+1, len(args))
	}

	val, err := args[idx].Eval(ctx.params)
	if err != nil {
		return val, fmt.Errorf("%s / %s", formatArgName(ctx.expr, idx), err.Error())
	}

	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	default:
		return v, nil
	}
}

func (ctx EvalContext) BooleanArg(idx int) (bool, error) {
	val, err := ctx.Arg(idx)
	if err != nil {
		return false, err
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	return false, formatArgError(ctx.expr, idx, "is not boolean: %v", val)
}

func (ctx EvalContext) NumericArg(idx int) (float64, error) {
	val, err := ctx.Arg(idx)
	if err != nil {
		return 0.0, err
	}

	if numVal, ok := val.(float64); ok {
		return numVal, nil
	}

	return 0.0, formatArgError(ctx.expr, idx, "is not numeric: %v", val)
}

func (ctx EvalContext) IntegerArg(idx int) (int, error) {
	val, err := ctx.NumericArg(idx)
	if err != nil {
		return 0, err
	}
	intVal := int(val)
	if float64(intVal) != val {
		return 0.0, formatArgError(ctx.expr, idx, "is not integer: %v", val)
	}
	return intVal, nil
}

func (ctx EvalContext) SliceArg(idx int) ([]interface{}, error) {
	val, err := ctx.Arg(idx)
	if err != nil {
		return []interface{}{}, err
	}
	if sliceVal, ok := val.([]interface{}); ok {
		return sliceVal, nil
	}
	return []interface{}{}, formatArgError(ctx.expr, idx, "is not array: %v", val)
}

func (ctx EvalContext) FormatError(msg string, msgArgs ...interface{}) error {
	return fmt.Errorf("%s [op=%s; pos=%d; len=%d]",
		fmt.Sprintf(msg, msgArgs...),
		ctx.expr.Name, ctx.expr.SourcePos, ctx.expr.SourceLen,
	)
}

func formatArgError(expr ExprNode, idx int, msg string, msgArgs ...interface{}) error {
	return fmt.Errorf("%s %s [%s]",
		formatArgName(expr, idx),
		fmt.Sprintf(msg, msgArgs...),
		formatDebugInfo(expr.Args[idx]),
	)
}

func formatDebugInfo(expr ExprNode) string {
	return fmt.Sprintf("pos=%d; len=%d", expr.SourcePos, expr.SourceLen)
}

func formatArgName(expr ExprNode, idx int) string {
	switch expr.OperatorType {
	case OperatorTypeInfix:
		if idx == 0 {
			return fmt.Sprintf("lhs of %s", expr.Name)
		} else if idx == 1 {
			return fmt.Sprintf("rhs of %s", expr.Name)
		}
	case OperatorTypePrefix:
		if idx == 0 {
			return fmt.Sprintf("argument of %s", expr.Name)
		}
	case OperatorTypeTernary:
		if idx == 0 {
			return "ternary condition"
		} else if idx == 1 {
			return "ternary then"
		} else if idx == 2 {
			return "ternary else"
		}
	case OperatorTypeArray:
		return fmt.Sprintf("array item #%d", idx+1)
	case OperatorTypeIndexer:
		if idx == 0 {
			return "indexer receiver"
		} else if idx == 1 {
			return "index"
		}
	}
	return fmt.Sprintf("argument #%d of %s", idx+1, expr.Name)
}
