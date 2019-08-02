package govaluate

import "fmt"

type EvalContext struct {
	params EvalParams
	args   []ExprNode
}

func (ctx EvalContext) ArgCount() int {
	return len(ctx.args)
}

func (ctx EvalContext) Arg(idx int) (interface{}, error) {
	if idx < len(ctx.args) {
		return ctx.args[idx].Eval(ctx.params)
	}
	return nil, fmt.Errorf("requested argument %d, but argument count is %d", idx, len(ctx.args))
}

func (ctx EvalContext) BooleanArg(idx int) (bool, error) {
	val, err := ctx.Arg(idx)
	if err != nil {
		return false, err
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	return false, fmt.Errorf("argument at %d is not boolean: %v", idx, val)
}

func (ctx EvalContext) NumericArg(idx int) (float64, error) {
	val, err := ctx.Arg(idx)
	if err != nil {
		return 0.0, err
	}
	if floatVal, ok := val.(float64); ok {
		return floatVal, nil
	}
	return 0.0, fmt.Errorf("argument at %d is not numeric: %v", idx, val)
}

func (ctx EvalContext) SliceArg(idx int) ([]interface{}, error) {
	val, err := ctx.Arg(idx)
	if err != nil {
		return []interface{}{}, err
	}
	if sliceVal, ok := val.([]interface{}); ok {
		return sliceVal, nil
	}
	return []interface{}{}, fmt.Errorf("argument at %d is not a slice: %v", idx, val)
}

type Operator func(ctx EvalContext) (interface{}, error)

type EvalParams struct {
	Variables map[string]interface{}
	Operators map[string]Operator
}

func (expr ExprNode) Eval(params EvalParams) (interface{}, error) {
	switch expr.Type {
	case NodeTypeLiteral:
		return expr.Value, nil
	case NodeTypeVariable:
		value, ok := params.Variables[expr.Name]
		if !ok {
			return nil, fmt.Errorf("variable undefined: %v", expr.Name)
		}
		return value, nil
	case NodeTypeOperator:
		operator, ok := params.Operators[expr.Name]
		if !ok {
			return nil, fmt.Errorf("operator undefined: %v", expr.Name)
		}
		return operator(EvalContext{params: params, args: expr.Args})
	}
	return nil, fmt.Errorf("bad expr type: %v", expr)
}

var builtinOperators = BuiltinOperators()

func NewEvalParams(variables map[string]interface{}) EvalParams {
	return EvalParams{
		Variables: variables,
		Operators: builtinOperators,
	}
}
