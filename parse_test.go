package govaluate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrecedence(t *testing.T) {
	expr, err := Parse("x + y * 2")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("+", []ExprNode{
			NewExprNodeVariable("x", 0, 1),
			NewExprNodeOperator("*", []ExprNode{
				NewExprNodeVariable("y", 4, 1),
				NewExprNodeLiteral(2.0, 8, 1),
			}, 4, 5, OperatorTypeInfix),
		}, 0, 9, OperatorTypeInfix),
		expr,
	)

	expr, err = Parse("(x + y) * 2")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("*", []ExprNode{
			NewExprNodeOperator("+", []ExprNode{
				NewExprNodeVariable("x", 1, 1),
				NewExprNodeVariable("y", 5, 1),
			}, 0, 7, OperatorTypeInfix),
			NewExprNodeLiteral(2.0, 10, 1),
		}, 0, 11, OperatorTypeInfix),
		expr,
	)
}

func TestParseFunction(t *testing.T) {
	expr, err := Parse("now()")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("now", []ExprNode{}, 0, 5, OperatorTypeCall),
		expr,
	)

	expr, err = Parse("pow(x, 2)")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("pow", []ExprNode{
			NewExprNodeVariable("x", 4, 1),
			NewExprNodeLiteral(2.0, 7, 1),
		}, 0, 9, OperatorTypeCall),
		expr,
	)

	expr, err = Parse("max(3, 4)")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("max", []ExprNode{
			NewExprNodeLiteral(3.0, 4, 1),
			NewExprNodeLiteral(4.0, 7, 1),
		}, 0, 9, OperatorTypeCall),
		expr,
	)
}

func TestParseUnary(t *testing.T) {
	expr, err := Parse("-x**2")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("-", []ExprNode{
			NewExprNodeOperator("**", []ExprNode{
				NewExprNodeVariable("x", 1, 1),
				NewExprNodeLiteral(2.0, 4, 1),
			}, 1, 4, OperatorTypeInfix),
		}, 0, 5, OperatorTypePrefix),
		expr,
	)

	expr, err = Parse("-x - y")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("-", []ExprNode{
			NewExprNodeOperator("-", []ExprNode{
				NewExprNodeVariable("x", 1, 1),
			}, 0, 2, OperatorTypePrefix),
			NewExprNodeVariable("y", 5, 1),
		}, 0, 6, OperatorTypeInfix),
		expr,
	)
}

func TestParseTernary(t *testing.T) {
	expr, err := Parse("x > 0 ? x**2 : -1")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("?:", []ExprNode{
			NewExprNodeOperator(">", []ExprNode{
				NewExprNodeVariable("x", 0, 1),
				NewExprNodeLiteral(0.0, 4, 1),
			}, 0, 5, OperatorTypeInfix),
			NewExprNodeOperator("**", []ExprNode{
				NewExprNodeVariable("x", 8, 1),
				NewExprNodeLiteral(2.0, 11, 1),
			}, 8, 4, OperatorTypeInfix),
			NewExprNodeOperator("-", []ExprNode{
				NewExprNodeLiteral(1.0, 16, 1),
			}, 15, 2, OperatorTypePrefix),
		}, 0, 17, OperatorTypeTernary),
		expr,
	)
}

func TestParseInArray(t *testing.T) {
	expr, err := Parse("x in [1, 2, 3]")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("in", []ExprNode{
			NewExprNodeVariable("x", 0, 1),
			NewExprNodeOperator("array", []ExprNode{
				NewExprNodeLiteral(1.0, 6, 1),
				NewExprNodeLiteral(2.0, 9, 1),
				NewExprNodeLiteral(3.0, 12, 1),
			}, 5, 9, OperatorTypeArray),
		}, 0, 14, OperatorTypeInfix),
		expr,
	)
}

func TestParseArrays(t *testing.T) {
	expr, err := Parse("[[x], []]")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("array", []ExprNode{
			NewExprNodeOperator("array", []ExprNode{
				NewExprNodeVariable("x", 2, 1),
			}, 1, 3, OperatorTypeArray),
			NewExprNodeOperator("array", []ExprNode{}, 6, 2, OperatorTypeArray),
		}, 0, 9, OperatorTypeArray),
		expr,
	)
}

func TestParseError(t *testing.T) {
	_, err := Parse("(1 + 2(")
	assert.EqualError(t, err, "unmatched bracket: '(', expecting ')', pos: 6")

	_, err = Parse("(1 + 2)(")
	assert.EqualError(t, err, "unexpected token Bracket{'('}, expecting operator, pos: 7")

	_, err = Parse("f(x ? y)")
	assert.EqualError(t, err, "unexpected token Bracket{')'}, expecting ':', pos: 7")

	_, err = Parse("f(x ? y : )")
	assert.EqualError(t, err, "unexpected token Bracket{')'}, expecting value, pos: 10")

	_, err = Parse("x +")
	assert.EqualError(t, err, "unexpected eof, expecting value")

	_, err = Parse("x ? y ? 1 : 0 : 0")
	assert.EqualError(t, err, "unexpected token Operator{?}, expecting ':', pos: 6")

	_, err = Parse("2 in [a, b, c")
	assert.EqualError(t, err, "unexpected eof, expecting ']', ','")
}
