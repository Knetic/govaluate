package govaluate

// ExprNode is a structured representation of an expression.
// There are three types of nodes: literal, variable and operator. The latter
// can have child nodes. They form a tree, where each node is an expression itself.
type ExprNode struct {
	Type  ExprNodeType
	Name  string
	Value interface{}
	Args  []ExprNode

	SourcePos, SourceLen int
	OperatorType         OperatorType
}

// ExprNodeType is a type of ExprNode.
type ExprNodeType int

const (
	// NodeTypeLiteral is just a constant literal, e.g. a boolean, a number, or a string.
	// ExprNode.Value contains the actual value.
	NodeTypeLiteral ExprNodeType = iota

	// NodeTypeVariable is a variable.
	// ExprNode.Name contains the name of the variable.
	NodeTypeVariable

	// NodeTypeOperator is an operation over the arguments, the other nodes.
	// It can be a function call, a binary operator (+, -, *, etc) with two arguments,
	// unary (!, -, ~), ternary (?:), etc. There are no restrictions on operator name,
	// it just needs to be defined at the evaluation phase.
	// ExprNode.Name is the name of the operation.
	// ExprNode.Args are the arguments.
	NodeTypeOperator
)

type OperatorType int

const (
	OperatorTypeCall OperatorType = iota
	OperatorTypeInfix
	OperatorTypePrefix
	OperatorTypeTernary
	OperatorTypeArray
	OperatorTypeIndexer
)

// NewExprNodeLiteral constructs a literal node.
func NewExprNodeLiteral(value interface{}, sourcePos, sourceLen int) ExprNode {
	return ExprNode{
		Type:      NodeTypeLiteral,
		Value:     value,
		SourcePos: sourcePos,
		SourceLen: sourceLen,
	}
}

// NewExprNodeVariable constructs a variable node.
func NewExprNodeVariable(name string, sourcePos, sourceLen int) ExprNode {
	return ExprNode{
		Type:      NodeTypeVariable,
		Name:      name,
		SourcePos: sourcePos,
		SourceLen: sourceLen,
	}
}

// NewExprNodeOperator constructs an operator node.
func NewExprNodeOperator(name string, args []ExprNode, sourcePos, sourceLen int, operatorType OperatorType) ExprNode {
	return ExprNode{
		Type:         NodeTypeOperator,
		Name:         name,
		Args:         args,
		SourcePos:    sourcePos,
		SourceLen:    sourceLen,
		OperatorType: operatorType,
	}
}

// IsOperator returns true if this expression is an operator with matching name.
func (expr ExprNode) IsOperator(name string) bool {
	return expr.Type == NodeTypeOperator && expr.Name == name
}

// IsLiteral returns true if this expression is a literal with matching value.
func (expr ExprNode) IsLiteral(value interface{}) bool {
	return expr.Type == NodeTypeLiteral && expr.Value == value
}

// GetValue returns expression value, if it's a constant.
func (expr ExprNode) GetValue() (interface{}, bool) {
	if expr.Type == NodeTypeLiteral {
		return expr.Value, true
	}
	return nil, false
}

// VarsCount returns a map where keys are the variable names in the expression,
// and values are how many times they are referenced.
func (expr ExprNode) VarsCount() map[string]int {
	vars := map[string]int{}
	collectVars(expr, vars)
	return vars
}

// Vars returns a list of variables referenced in the expression.
func (expr ExprNode) Vars() []string {
	vars := expr.VarsCount()
	res := make([]string, 0, len(vars))
	for key := range vars {
		res = append(res, key)
	}
	return res
}

func collectVars(expr ExprNode, output map[string]int) {
	switch expr.Type {
	case NodeTypeVariable:
		output[expr.Name]++
	case NodeTypeOperator:
		for _, arg := range expr.Args {
			collectVars(arg, output)
		}
	}
}
