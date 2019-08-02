package govaluate

type Optimizer func(ExprNode) ExprNode

func BuiltinOptimizers() map[string]Optimizer {
	return map[string]Optimizer{
		"&&": func(expr ExprNode) ExprNode {
			left := expr.Args[0]
			right := expr.Args[1]
			if left.IsLiteral(false) || right.IsLiteral(false) {
				// false && x -> false, x && false -> false
				return NewExprNodeLiteral(false)
			}
			if left.IsLiteral(true) {
				// true && x -> x
				return right
			}
			if right.IsLiteral(true) {
				// x && true -> x
				return left
			}
			return expr
		},
		"||": func(expr ExprNode) ExprNode {
			left := expr.Args[0]
			right := expr.Args[1]
			if left.IsLiteral(true) || right.IsLiteral(true) {
				// true || x -> true, x || true -> true
				return NewExprNodeLiteral(true)
			}
			if left.IsLiteral(false) {
				// false || x -> x
				return right
			}
			if right.IsLiteral(false) {
				// x || false -> x
				return left
			}
			return expr
		},
		"+": func(expr ExprNode) ExprNode {
			left := expr.Args[0]
			right := expr.Args[1]
			if left.IsLiteral(0.0) {
				// 0 + x -> x
				return right
			}
			if right.IsLiteral(0.0) {
				// x + 0 -> x
				return left
			}
			return expr
		},
		"-": func(expr ExprNode) ExprNode {
			if len(expr.Args) != 2 {
				return expr
			}
			left := expr.Args[0]
			right := expr.Args[1]
			if left.IsLiteral(0.0) {
				// 0 - x -> -x
				return NewExprNodeOperator("-", right)
			}
			if right.IsLiteral(0.0) {
				// x - 0 -> x
				return left
			}
			return expr
		},
		"*": func(expr ExprNode) ExprNode {
			left := expr.Args[0]
			right := expr.Args[1]
			if left.IsLiteral(0.0) || right.IsLiteral(0.0) {
				// 0 * x -> 0, x * 0 -> 0
				return NewExprNodeLiteral(0.0)
			}
			if left.IsLiteral(1.0) {
				// 1 * x -> x
				return right
			}
			if right.IsLiteral(1.0) {
				// x * 1 -> x
				return left
			}
			return expr
		},
		"/": func(expr ExprNode) ExprNode {
			left := expr.Args[0]
			right := expr.Args[1]
			if right.IsLiteral(1.0) {
				// x / 1 -> x
				return left
			}
			return expr
		},
		"?:": func(expr ExprNode) ExprNode {
			condition := expr.Args[0]
			if condition.IsLiteral(true) {
				// true ? x : y -> x
				return expr.Args[1]
			}
			if condition.IsLiteral(false) {
				// false ? x : y -> y
				return expr.Args[2]
			}
			return expr
		},
	}
}
