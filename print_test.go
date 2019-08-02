package govaluate

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintOverrideLiterals(t *testing.T) {
	expr, err := Parse("1 > x || true || false")
	assert.Nil(t, err)

	output, err := expr.Print(PrintConfig{
		FormatBoolLiteral: func(value bool) string {
			if value {
				return "TRUE"
			}
			return "FALSE"
		},
		FormatNumberLiteral: func(value float64) string {
			return strconv.FormatFloat(value, 'f', 2, 64)
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "1.00 > x || TRUE || FALSE", output)
}

func TestPrintOverrideVariables(t *testing.T) {
	expr, err := Parse("x ? y : z")
	assert.Nil(t, err)

	output, err := expr.Print(PrintConfig{
		FormatVariable: func(name string) string {
			switch name {
			case "x":
				return "condition"
			case "y":
				return "valueIfTrue"
			case "z":
				return "valueIfFalse"
			}
			return name
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "condition ? valueIfTrue : valueIfFalse", output)
}

func TestPrintOverrideOperators(t *testing.T) {
	expr, err := Parse("n > 0 ? 2**n : -n - 1")
	assert.Nil(t, err)

	output, err := expr.Print(PrintConfig{
		OperatorMap: map[string]string{
			"?:": "IF",
			"**": "POW",
			">":  "GT",
		},
		OperatorMapper: func(operator string, arity int) string {
			if operator == "-" {
				if arity == 1 {
					return "NEGATE"
				} else if arity == 2 {
					return "SUB"
				}
			}
			return operator
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "IF(GT(n, 0), POW(2, n), SUB(NEGATE(n), 1))", output)
}

func TestPrintOverrideInfix(t *testing.T) {
	expr, err := Parse("2 ** n")
	assert.Nil(t, err)

	output, err := expr.Print(PrintConfig{
		OperatorMap: map[string]string{
			"**": "pow",
		},
		InfixOperators: map[string]bool{
			"pow": true,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "2 pow n", output)
}
