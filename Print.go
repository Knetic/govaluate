package govaluate

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// ExprNodePrinter is an output builder for ExprNode.
// Use AppendString or AppendNode from custom handlers to append to output.
type ExprNodePrinter struct {
	nodeHandler func(ExprNode, *ExprNodePrinter) error
	output      strings.Builder
	err         error
}

// PrintConfig is used to override default behavior when printing an expression with a default node handler.
type PrintConfig struct {
	// FormatBoolLiteral overrides boolean literal output.
	// Default handler returns "true" or "false".
	FormatBoolLiteral func(bool) string

	// FormatNumberLiteral overrides number literal output.
	// Default handler formats number with strconv.FormatFloat(value, 'f', -1, 64).
	FormatNumberLiteral func(float64) string

	// FormatStringLiteral overrides string literal output.
	// Default handler returns quoted value with other quotes and newlines escaped with backslash (\).
	FormatStringLiteral func(string) string

	// FormatVariable overrides variable output.
	// Default handler simply returns identifier as is.
	// This can be used to map variables to different names.
	FormatVariable func(string) string

	// OperatorMap contains a mapping for operator name overrides.
	// For example, ** -> pow mapping will change output from x ** y to pow(x, y).
	OperatorMap map[string]string

	// OperatorMapper is similar to OperatorMap, but allows mapping by name and arity (number of arguments).
	// For example, this can override unary and binary minus in a different way.
	OperatorMapper func(name string, arity int) string

	// InfixOperators contains overrides of what operators are printed in infix notation.
	// By default, all operators written in special symbols and "in" operator are considered infix.
	// For example, overriding pow -> true will change output from pow(x, y) to x pow y.
	// This only applies if an operator is binary (two arguments).
	InfixOperators map[string]bool

	// PrecedenceFn overrides precedence of operators.
	// Higher precedence means that the operation should performed first.
	// See defaultPrecedence(string, int) for defaults.
	PrecedenceFn func(name string, arity int) int

	// Operators overrides default behavior when printing a particular operator.
	// By default, special symbol unary operators are printed in prefix notation: !x, ~x.
	// Infix binary operators (see InfixOperators) are printed in infix notation: x + y, x && y.
	// Ternary if (?:) is printed like this: condition ? then : else.
	// All other operators are printed as function calls: square(x), now(), pow(x, y).
	Operators map[string]func(args []ExprNode, output *ExprNodePrinter) error
}

// AppendString appends a token to output as is.
func (b *ExprNodePrinter) AppendString(token string) {
	if b.err == nil {
		b.output.WriteString(token)
	}
}

// AppendNode invokes node handler that will print node to output.
func (b *ExprNodePrinter) AppendNode(node ExprNode) {
	if b.err == nil {
		err := b.nodeHandler(node, b)
		if b.err == nil {
			b.err = err
		}
	}
}

// Print converts ExprNode to string with default node handler.
// PrintConfig can be used to configure output. Use empty PrintConfig for default behavior.
func (expr ExprNode) Print(config PrintConfig) (string, error) {
	return expr.PrintWithHandler(defaultNodeHandler(config))
}

// PrintWithHandler converts ExprNode to string using the specified node handler.
// Node handler takes an ExprNode and feeds output to ExprNodePrinter.
func (expr ExprNode) PrintWithHandler(nodeHandler func(ExprNode, *ExprNodePrinter) error) (string, error) {
	builder := &ExprNodePrinter{nodeHandler: nodeHandler}
	builder.AppendNode(expr)
	return builder.output.String(), builder.err
}

func defaultNodeHandler(config PrintConfig) func(ExprNode, *ExprNodePrinter) error {
	return func(node ExprNode, output *ExprNodePrinter) error {
		switch node.Type {
		case NodeTypeLiteral:
			return literal(node.Value, output, &config)
		case NodeTypeVariable:
			return variable(node.Name, output, &config)
		case NodeTypeOperator:
			return operator(node.Name, node.Args, output, &config)
		}
		return fmt.Errorf("unexpected node: %v", node)
	}
}

func literal(value interface{}, output *ExprNodePrinter, config *PrintConfig) error {
	var literal string
	switch value.(type) {
	case bool:
		literal = boolLiteral(value.(bool), config)
	case float64:
		literal = numberLiteral(value.(float64), config)
	case string:
		literal = stringLiteral(value.(string), config)
	default:
		return fmt.Errorf("unsupported literal type: %v", value)
	}
	output.AppendString(literal)
	return nil
}

func boolLiteral(value bool, config *PrintConfig) string {
	if config.FormatBoolLiteral != nil {
		return config.FormatBoolLiteral(value)
	}
	if value {
		return "true"
	}
	return "false"
}

func numberLiteral(value float64, config *PrintConfig) string {
	if config.FormatNumberLiteral != nil {
		return config.FormatNumberLiteral(value)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func stringLiteral(value string, config *PrintConfig) string {
	if config.FormatStringLiteral != nil {
		return config.FormatStringLiteral(value)
	}
	escapedValue := strings.NewReplacer(
		"\\", "\\\\",
		"\"", "\\\"",
		"\r", "\\r",
		"\n", "\\n",
	).Replace(value)
	return "\"" + escapedValue + "\""
}

func variable(name string, output *ExprNodePrinter, config *PrintConfig) error {
	variable := name
	if config.FormatVariable != nil {
		variable = config.FormatVariable(name)
	}
	output.AppendString(variable)
	return nil
}

func operator(name string, args []ExprNode, output *ExprNodePrinter, config *PrintConfig) error {
	arity := len(args)
	mappedName := config.mappedName(name, arity)

	if fn, ok := config.Operators[mappedName]; ok {
		return fn(args, output)
	}

	// binary operator: x + y
	infix := config.isInfix(name, arity)
	if infix {
		selfPrecedence := config.precedence(name, arity)
		leftPrecedence := config.precedenceForNode(args[0])
		rightPrecedence := config.precedenceForNode(args[1])
		if leftPrecedence < selfPrecedence {
			output.AppendString("(")
		}
		output.AppendNode(args[0])
		if leftPrecedence < selfPrecedence {
			output.AppendString(")")
		}
		output.AppendString(" ")
		output.AppendString(mappedName)
		output.AppendString(" ")
		if rightPrecedence <= selfPrecedence {
			output.AppendString("(")
		}
		output.AppendNode(args[1])
		if rightPrecedence <= selfPrecedence {
			output.AppendString(")")
		}
		return nil
	}

	// prefix operator: !x
	prefix := arity == 1 && isSpecial(mappedName)
	if prefix {
		selfPrecedence := config.precedence(name, arity)
		rightPrecedence := config.precedenceForNode(args[0])
		output.AppendString(mappedName)
		if rightPrecedence < selfPrecedence {
			output.AppendString("(")
		}
		output.AppendNode(args[0])
		if rightPrecedence < selfPrecedence {
			output.AppendString(")")
		}
		return nil
	}

	// ternary if: x ? y : z
	if mappedName == "?:" && arity == 3 {
		selfPrecedence := config.precedence(name, arity)
		conditionPrecedence := config.precedenceForNode(args[0])
		thenPrecedence := config.precedenceForNode(args[1])
		elsePrecedence := config.precedenceForNode(args[2])
		if conditionPrecedence <= selfPrecedence {
			output.AppendString("(")
		}
		output.AppendNode(args[0])
		if conditionPrecedence <= selfPrecedence {
			output.AppendString(")")
		}
		output.AppendString(" ? ")
		if thenPrecedence <= selfPrecedence {
			output.AppendString("(")
		}
		output.AppendNode(args[1])
		if thenPrecedence <= selfPrecedence {
			output.AppendString(")")
		}
		output.AppendString(" : ")
		if elsePrecedence < selfPrecedence {
			output.AppendString("(")
		}
		output.AppendNode(args[2])
		if elsePrecedence < selfPrecedence {
			output.AppendString(")")
		}
		return nil
	}

	// function call: fn(a, b, c)
	output.AppendString(mappedName)
	output.AppendString("(")
	for idx, arg := range args {
		if idx > 0 {
			output.AppendString(", ")
		}
		output.AppendNode(arg)
	}
	output.AppendString(")")
	return nil
}

func isSpecial(name string) bool {
	for _, r := range []rune(name) {
		if unicode.IsLetter(r) {
			return false
		}
	}
	return len(name) > 0
}

func (config *PrintConfig) mappedName(operator string, arity int) string {
	if mappedName, ok := config.OperatorMap[operator]; ok {
		return mappedName
	}
	if config.OperatorMapper != nil {
		if mappedName := config.OperatorMapper(operator, arity); mappedName != "" {
			return mappedName
		}
	}
	return operator
}

func (config *PrintConfig) isInfix(operator string, arity int) bool {
	if arity != 2 {
		return false
	}
	mappedName := config.mappedName(operator, arity)
	if infix, found := config.InfixOperators[mappedName]; found {
		return infix
	}
	return isSpecial(mappedName) || mappedName == "in"
}

func (config *PrintConfig) precedenceForNode(node ExprNode) int {
	if node.Type == NodeTypeOperator {
		return config.precedence(node.Name, len(node.Args))
	}
	// variable and literal have max precedence
	return math.MaxInt32
}

func (config *PrintConfig) precedence(operator string, arity int) int {
	if config.PrecedenceFn != nil {
		mappedName := config.mappedName(operator, arity)
		return config.PrecedenceFn(mappedName, arity)
	}
	return defaultPrecedence(operator, arity)
}

func defaultPrecedence(operator string, arity int) int {
	if arity == 1 {
		return 10
	}
	switch operator {
	case ",":
		return 0
	case "?:", "?", ":":
		return 1
	case "??":
		return 2
	case "||":
		return 3
	case "&&":
		return 4
	case "==", "!=", ">", "<", ">=", "<=", "=~", "!~", "in":
		return 5
	case "&", "|", "^", "<<", ">>":
		return 7
	case "+", "-":
		return 8
	case "*", "/", "%":
		return 9
	case "**":
		return 11
	}
	return 6
}
