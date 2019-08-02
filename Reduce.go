package govaluate

import "fmt"

// Reduce does a partial parameter evaluation and returns simplified expression
func (node ExprNode) Reduce(params EvalParams, optimizers map[string]Optimizer) (ExprNode, error) {
	switch node.Type {
	case NodeTypeLiteral:
		// literal can not be reduced
		return node, nil

	case NodeTypeVariable:
		if value, ok := params.Variables[node.Name]; ok {
			// variable is known, replace it with value literal
			return NewExprNodeLiteral(value), nil
		}
		// variable is unknown, return as is
		return node, nil

	case NodeTypeOperator:
		// reduce arguments
		reducedArgs := make([]ExprNode, len(node.Args))
		allArgsKnown := true
		for idx, arg := range node.Args {
			reducedArg, err := arg.Reduce(params, optimizers)
			if err != nil {
				return node, err
			}
			reducedArgs[idx] = reducedArg
			if reducedArg.Type != NodeTypeLiteral {
				allArgsKnown = false
			}
		}

		_, operatorKnown := params.Operators[node.Name]
		if allArgsKnown && operatorKnown {
			// all arguments are known, perform the operation
			value, err := node.Eval(params)
			if err != nil {
				return node, err
			}
			return NewExprNodeLiteral(value), nil
		}

		reducedNode := NewExprNodeOperator(node.Name, reducedArgs...)

		if optimizer, ok := optimizers[node.Name]; ok {
			return optimizer(reducedNode), nil
		}

		return reducedNode, nil
	}
	return node, fmt.Errorf("bad node type: %v", node)
}
