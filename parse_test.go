package govaluate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrecedence(t *testing.T) {
	expr, err := Parse("x + y * 2")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("+",
			NewExprNodeVariable("x"),
			NewExprNodeOperator("*",
				NewExprNodeVariable("y"),
				NewExprNodeLiteral(2.0),
			),
		),
		expr,
	)

	expr, err = Parse("(x + y) * 2")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("*",
			NewExprNodeOperator("+",
				NewExprNodeVariable("x"),
				NewExprNodeVariable("y"),
			),
			NewExprNodeLiteral(2.0),
		),
		expr,
	)
}

func TestParseFunction(t *testing.T) {
	expr, err := Parse("now()")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("now"),
		expr,
	)

	expr, err = Parse("pow(x, 2)")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("pow",
			NewExprNodeVariable("x"),
			NewExprNodeLiteral(2.0),
		),
		expr,
	)
}

func TestParseUnary(t *testing.T) {
	expr, err := Parse("-x**2")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("-",
			NewExprNodeOperator("**",
				NewExprNodeVariable("x"),
				NewExprNodeLiteral(2.0),
			),
		),
		expr,
	)

	expr, err = Parse("-x - y")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("-",
			NewExprNodeOperator("-",
				NewExprNodeVariable("x"),
			),
			NewExprNodeVariable("y"),
		),
		expr,
	)
}

func TestParseTernary(t *testing.T) {
	expr, err := Parse("x > 0 ? x**2 : -1")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("?:",
			NewExprNodeOperator(">",
				NewExprNodeVariable("x"),
				NewExprNodeLiteral(0.0),
			),
			NewExprNodeOperator("**",
				NewExprNodeVariable("x"),
				NewExprNodeLiteral(2.0),
			),
			NewExprNodeOperator("-",
				NewExprNodeLiteral(1.0),
			),
		),
		expr,
	)
}

func TestParseInArray(t *testing.T) {
	expr, err := Parse("x in [1, 2, 3]")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("in",
			NewExprNodeVariable("x"),
			NewExprNodeOperator("array",
				NewExprNodeLiteral(1.0),
				NewExprNodeLiteral(2.0),
				NewExprNodeLiteral(3.0),
			),
		),
		expr,
	)
}

func TestParseArrays(t *testing.T) {
	expr, err := Parse("[[x], []]")
	assert.Nil(t, err)
	assert.Equal(t,
		NewExprNodeOperator("array",
			NewExprNodeOperator("array",
				NewExprNodeVariable("x"),
			),
			NewExprNodeOperator("array"),
		),
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
