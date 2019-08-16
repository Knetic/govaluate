package govaluate

import "fmt"

// Reduce does a partial parameter evaluation and returns simplified expression
func (expr ExprNode) Reduce(params EvalParams, optimizers map[string]Optimizer) (ExprNode, error) {
	switch expr.Type {
	case NodeTypeLiteral:
		// literal can not be reduced
		return expr, nil

	case NodeTypeVariable:
		if value, ok := params.Variables[expr.Name]; ok {
			// variable is known, replace it with value literal
			return NewExprNodeLiteral(value, expr.SourcePos, expr.SourceLen), nil
		}
		// variable is unknown, return as is
		return expr, nil

	case NodeTypeOperator:
		// reduce arguments
		reducedArgs := make([]ExprNode, len(expr.Args))
		allArgsKnown := true
		for idx, arg := range expr.Args {
			reducedArg, err := arg.Reduce(params, optimizers)
			if err != nil {
				return expr, err
			}
			reducedArgs[idx] = reducedArg
			if reducedArg.Type != NodeTypeLiteral {
				allArgsKnown = false
			}
		}

		_, operatorKnown := params.Operators[expr.Name]
		if allArgsKnown && operatorKnown {
			// all arguments are known, perform the operation
			value, err := expr.Eval(params)
			if err != nil {
				return expr, err
			}
			return NewExprNodeLiteral(value, expr.SourcePos, expr.SourceLen), nil
		}

		expr.Args = reducedArgs

		if optimizer, ok := optimizers[expr.Name]; ok {
			return optimizer(expr), nil
		}

		return expr, nil
	}
	return expr, fmt.Errorf("bad node type: %v", expr)
}
