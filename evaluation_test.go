package govaluate

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"
)

/*
	Represents a test of expression evaluation
*/
type EvaluationTest struct {
	Name       string
	Input      string
	Functions  map[string]ExpressionFunction
	Parameters []EvaluationParameter
	Expected   interface{}
}

type EvaluationParameter struct {
	Name  string
	Value interface{}
}

func TestNoParameterEvaluation(test *testing.T) {

	evaluationTests := []EvaluationTest{

		{

			Name:     "Single PLUS",
			Input:    "51 + 49",
			Expected: 100.0,
		},
		{

			Name:     "Single MINUS",
			Input:    "100 - 51",
			Expected: 49.0,
		},
		{

			Name:     "Single BITWISE AND",
			Input:    "100 & 50",
			Expected: 32.0,
		},
		{

			Name:     "Single BITWISE OR",
			Input:    "100 | 50",
			Expected: 118.0,
		},
		{

			Name:     "Single BITWISE XOR",
			Input:    "100 ^ 50",
			Expected: 86.0,
		},
		{

			Name:     "Single shift left",
			Input:    "2 << 1",
			Expected: 4.0,
		},
		{

			Name:     "Single shift right",
			Input:    "2 >> 1",
			Expected: 1.0,
		},
		{

			Name:     "Single BITWISE NOT",
			Input:    "~10",
			Expected: -11.0,
		},
		{

			Name:     "Single MULTIPLY",
			Input:    "5 * 20",
			Expected: 100.0,
		},
		{

			Name:     "Single DIVIDE",
			Input:    "100 / 20",
			Expected: 5.0,
		},
		{

			Name:     "Single even MODULUS",
			Input:    "100 % 2",
			Expected: 0.0,
		},
		{

			Name:     "Single odd MODULUS",
			Input:    "101 % 2",
			Expected: 1.0,
		},
		{

			Name:     "Single EXPONENT",
			Input:    "10 ** 2",
			Expected: 100.0,
		},
		{

			Name:     "Compound PLUS",
			Input:    "20 + 30 + 50",
			Expected: 100.0,
		},
		{

			Name:     "Compound BITWISE AND",
			Input:    "20 & 30 & 50",
			Expected: 16.0,
		},
		{

			Name:     "Mutiple operators",
			Input:    "20 * 5 - 49",
			Expected: 51.0,
		},
		{

			Name:     "Parenthesis usage",
			Input:    "100 - (5 * 10)",
			Expected: 50.0,
		},
		{

			Name:     "Nested parentheses",
			Input:    "50 + (5 * (15 - 5))",
			Expected: 100.0,
		},
		{

			Name:     "Nested parentheses with bitwise",
			Input:    "100 ^ (23 * (2 | 5))",
			Expected: 197.0,
		},
		{

			Name:     "Logical OR operation of two clauses",
			Input:    "(1 == 1) || (true == true)",
			Expected: true,
		},
		{

			Name:     "Logical AND operation of two clauses",
			Input:    "(1 == 1) && (true == true)",
			Expected: true,
		},
		{

			Name:     "Implicit boolean",
			Input:    "2 > 1",
			Expected: true,
		},
		{

			Name:     "Compound boolean",
			Input:    "5 < 10 && 1 < 5",
			Expected: true,
		},
		{

			Name:     "Evaluated true && false operation (for issue #8)",
			Input:    "1 > 10 && 11 > 10",
			Expected: false,
		},
		{

			Name:     "Evaluated true && false operation (for issue #8)",
			Input:    "true == true && false == true",
			Expected: false,
		},
		{

			Name:     "Parenthesis boolean",
			Input:    "10 < 50 && (1 != 2 && 1 > 0)",
			Expected: true,
		},
		{

			Name:     "Comparison of string constants",
			Input:    "'foo' == 'foo'",
			Expected: true,
		},
		{

			Name:     "NEQ comparison of string constants",
			Input:    "'foo' != 'bar'",
			Expected: true,
		},
		{

			Name:     "REQ comparison of string constants",
			Input:    "'foobar' =~ 'oba'",
			Expected: true,
		},
		{

			Name:     "NREQ comparison of string constants",
			Input:    "'foo' !~ 'bar'",
			Expected: true,
		},
		{

			Name:     "Multiplicative/additive order",
			Input:    "5 + 10 * 2",
			Expected: 25.0,
		},
		{

			Name:     "Multiple constant multiplications",
			Input:    "10 * 10 * 10",
			Expected: 1000.0,
		},
		{

			Name:     "Multiple adds/multiplications",
			Input:    "10 * 10 * 10 + 1 * 10 * 10",
			Expected: 1100.0,
		},
		{

			Name:     "Modulus precedence",
			Input:    "1 + 101 % 2 * 5",
			Expected: 6.0,
		},
		{

			Name:     "Exponent precedence",
			Input:    "1 + 5 ** 3 % 2 * 5",
			Expected: 6.0,
		},
		{

			Name:     "Bit shift precedence",
			Input:    "50 << 1 & 90",
			Expected: 64.0,
		},
		{

			Name:     "Bit shift precedence",
			Input:    "90 & 50 << 1",
			Expected: 64.0,
		},
		{

			Name:     "Bit shift precedence amongst non-bitwise",
			Input:    "90 + 50 << 1 * 5",
			Expected: 4480.0,
		},
		{
			Name:     "Order of non-commutative same-precedence operators (additive)",
			Input:    "1 - 2 - 4 - 8",
			Expected: -13.0,
		},
		{
			Name:     "Order of non-commutative same-precedence operators (multiplicative)",
			Input:    "1 * 4 / 2 * 8",
			Expected: 16.0,
		},
		{
			Name:     "Null coalesce precedence",
			Input:    "true ?? true ? 100 + 200 : 400",
			Expected: 300.0,
		},
		{

			Name:     "Identical date equivalence",
			Input:    "'2014-01-02 14:12:22' == '2014-01-02 14:12:22'",
			Expected: true,
		},
		{

			Name:     "Positive date GT",
			Input:    "'2014-01-02 14:12:22' > '2014-01-02 12:12:22'",
			Expected: true,
		},
		{

			Name:     "Negative date GT",
			Input:    "'2014-01-02 14:12:22' > '2014-01-02 16:12:22'",
			Expected: false,
		},
		{

			Name:     "Positive date GTE",
			Input:    "'2014-01-02 14:12:22' >= '2014-01-02 12:12:22'",
			Expected: true,
		},
		{

			Name:     "Negative date GTE",
			Input:    "'2014-01-02 14:12:22' >= '2014-01-02 16:12:22'",
			Expected: false,
		},
		{

			Name:     "Positive date LT",
			Input:    "'2014-01-02 14:12:22' < '2014-01-02 16:12:22'",
			Expected: true,
		},
		{

			Name:     "Negative date LT",
			Input:    "'2014-01-02 14:12:22' < '2014-01-02 11:12:22'",
			Expected: false,
		},
		{

			Name:     "Positive date LTE",
			Input:    "'2014-01-02 09:12:22' <= '2014-01-02 12:12:22'",
			Expected: true,
		},
		{

			Name:     "Negative date LTE",
			Input:    "'2014-01-02 14:12:22' <= '2014-01-02 11:12:22'",
			Expected: false,
		},
		{

			Name:     "Sign prefix comparison",
			Input:    "-1 < 0",
			Expected: true,
		},
		{

			Name:     "Lexicographic LT",
			Input:    "'ab' < 'abc'",
			Expected: true,
		},
		{

			Name:     "Lexicographic LTE",
			Input:    "'ab' <= 'abc'",
			Expected: true,
		},
		{

			Name:     "Lexicographic GT",
			Input:    "'aba' > 'abc'",
			Expected: false,
		},
		{

			Name:     "Lexicographic GTE",
			Input:    "'aba' >= 'abc'",
			Expected: false,
		},
		{

			Name:     "Boolean sign prefix comparison",
			Input:    "!true == false",
			Expected: true,
		},
		{

			Name:     "Inversion of clause",
			Input:    "!(10 < 0)",
			Expected: true,
		},
		{

			Name:     "Negation after modifier",
			Input:    "10 * -10",
			Expected: -100.0,
		},
		{

			Name:     "Ternary with single boolean",
			Input:    "true ? 10",
			Expected: 10.0,
		},
		{

			Name:     "Ternary nil with single boolean",
			Input:    "false ? 10",
			Expected: nil,
		},
		{

			Name:     "Ternary with comparator boolean",
			Input:    "10 > 5 ? 35.50",
			Expected: 35.50,
		},
		{

			Name:     "Ternary nil with comparator boolean",
			Input:    "1 > 5 ? 35.50",
			Expected: nil,
		},
		{

			Name:     "Ternary with parentheses",
			Input:    "(5 * (15 - 5)) > 5 ? 35.50",
			Expected: 35.50,
		},
		{

			Name:     "Ternary precedence",
			Input:    "true ? 35.50 > 10",
			Expected: true,
		},
		{

			Name:     "Ternary-else",
			Input:    "false ? 35.50 : 50",
			Expected: 50.0,
		},
		{

			Name:     "Ternary-else inside clause",
			Input:    "(false ? 5 : 35.50) > 10",
			Expected: true,
		},
		{

			Name:     "Ternary-else (true-case) inside clause",
			Input:    "(true ? 1 : 5) < 10",
			Expected: true,
		},
		{

			Name:     "Ternary-else before comparator (negative case)",
			Input:    "true ? 1 : 5 > 10",
			Expected: 1.0,
		},
		{

			Name:     "Nested ternaries (#32)",
			Input:    "(2 == 2) ? 1 : (true ? 2 : 3)",
			Expected: 1.0,
		},
		{

			Name:     "Nested ternaries, right case (#32)",
			Input:    "false ? 1 : (true ? 2 : 3)",
			Expected: 2.0,
		},
		{

			Name:     "Doubly-nested ternaries (#32)",
			Input:    "true ? (false ? 1 : (false ? 2 : 3)) : (false ? 4 : 5)",
			Expected: 3.0,
		},
		{

			Name:     "String to string concat",
			Input:    "'foo' + 'bar' == 'foobar'",
			Expected: true,
		},
		{

			Name:     "String to float64 concat",
			Input:    "'foo' + 123 == 'foo123'",
			Expected: true,
		},
		{

			Name:     "Float64 to string concat",
			Input:    "123 + 'bar' == '123bar'",
			Expected: true,
		},
		{

			Name:     "String to date concat",
			Input:    "'foo' + '02/05/1970' == 'foobar'",
			Expected: false,
		},
		{

			Name:     "String to bool concat",
			Input:    "'foo' + true == 'footrue'",
			Expected: true,
		},
		{

			Name:     "Bool to string concat",
			Input:    "true + 'bar' == 'truebar'",
			Expected: true,
		},
		{

			Name:     "Null coalesce left",
			Input:    "1 ?? 2",
			Expected: 1.0,
		},
		{

			Name:     "Array membership literals",
			Input:    "1 in (1, 2, 3)",
			Expected: true,
		},
		{

			Name:     "Array membership literal with inversion",
			Input:    "!(1 in (1, 2, 3))",
			Expected: false,
		},
		{

			Name:     "Logical operator reordering (#30)",
			Input:    "(true && true) || (true && false)",
			Expected: true,
		},
		{

			Name:     "Logical operator reordering without parens (#30)",
			Input:    "true && true || true && false",
			Expected: true,
		},
		{

			Name:     "Logical operator reordering with multiple OR (#30)",
			Input:    "false || true && true || false",
			Expected: true,
		},
		{

			Name:     "Left-side multiple consecutive (should be reordered) operators",
			Input:    "(10 * 10 * 10) > 10",
			Expected: true,
		},
		{

			Name:     "Three-part non-paren logical op reordering (#44)",
			Input:    "false && true || true",
			Expected: true,
		},
		{

			Name:     "Three-part non-paren logical op reordering (#44), second one",
			Input:    "true || false && true",
			Expected: true,
		},
		{

			Name:     "Logical operator reordering without parens (#45)",
			Input:    "true && true || false && false",
			Expected: true,
		},
		{

			Name:  "Single function",
			Input: "foo()",
			Functions: map[string]ExpressionFunction{
				"foo": func(arguments ...interface{}) (interface{}, error) {
					return true, nil
				},
			},

			Expected: true,
		},
		{

			Name:  "Function with argument",
			Input: "passthrough(1)",
			Functions: map[string]ExpressionFunction{
				"passthrough": func(arguments ...interface{}) (interface{}, error) {
					return arguments[0], nil
				},
			},

			Expected: 1.0,
		},

		{

			Name:  "Function with arguments",
			Input: "passthrough(1, 2)",
			Functions: map[string]ExpressionFunction{
				"passthrough": func(arguments ...interface{}) (interface{}, error) {
					return arguments[0].(float64) + arguments[1].(float64), nil
				},
			},

			Expected: 3.0,
		},
		{

			Name:  "Nested function with precedence",
			Input: "sum(1, sum(2, 3), 2 + 2, true ? 4 : 5)",
			Functions: map[string]ExpressionFunction{
				"sum": func(arguments ...interface{}) (interface{}, error) {

					sum := 0.0
					for _, v := range arguments {
						sum += v.(float64)
					}
					return sum, nil
				},
			},

			Expected: 14.0,
		},
		{

			Name:  "Empty function and modifier, compared",
			Input: "numeric()-1 > 0",
			Functions: map[string]ExpressionFunction{
				"numeric": func(arguments ...interface{}) (interface{}, error) {
					return 2.0, nil
				},
			},

			Expected: true,
		},
		{

			Name:  "Empty function comparator",
			Input: "numeric() > 0",
			Functions: map[string]ExpressionFunction{
				"numeric": func(arguments ...interface{}) (interface{}, error) {
					return 2.0, nil
				},
			},

			Expected: true,
		},
		{

			Name:  "Empty function logical operator",
			Input: "success() && !false",
			Functions: map[string]ExpressionFunction{
				"success": func(arguments ...interface{}) (interface{}, error) {
					return true, nil
				},
			},

			Expected: true,
		},
		{

			Name:  "Empty function ternary",
			Input: "nope() ? 1 : 2.0",
			Functions: map[string]ExpressionFunction{
				"nope": func(arguments ...interface{}) (interface{}, error) {
					return false, nil
				},
			},

			Expected: 2.0,
		},
		{

			Name:  "Empty function null coalesce",
			Input: "null() ?? 2",
			Functions: map[string]ExpressionFunction{
				"null": func(arguments ...interface{}) (interface{}, error) {
					return nil, nil
				},
			},

			Expected: 2.0,
		},
		{

			Name:  "Empty function with prefix",
			Input: "-ten()",
			Functions: map[string]ExpressionFunction{
				"ten": func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				},
			},

			Expected: -10.0,
		},
		{

			Name:  "Empty function as part of chain",
			Input: "10 - numeric() - 2",
			Functions: map[string]ExpressionFunction{
				"numeric": func(arguments ...interface{}) (interface{}, error) {
					return 5.0, nil
				},
			},

			Expected: 3.0,
		},
		{

			Name:  "Empty function near separator",
			Input: "10 in (1, 2, 3, ten(), 8)",
			Functions: map[string]ExpressionFunction{
				"ten": func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				},
			},

			Expected: true,
		},
		{

			Name:  "Enclosed empty function with modifier and comparator (#28)",
			Input: "(ten() - 1) > 3",
			Functions: map[string]ExpressionFunction{
				"ten": func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				},
			},

			Expected: true,
		},
		{

			Name:  "Ternary/Java EL ambiguity",
			Input: "false ? foo:length()",
			Functions: map[string]ExpressionFunction{
				"length": func(arguments ...interface{}) (interface{}, error) {
					return 1.0, nil
				},
			},
			Expected: 1.0,
		},
	}

	runEvaluationTests(evaluationTests, test)
}

func TestParameterizedEvaluation(test *testing.T) {

	evaluationTests := []EvaluationTest{

		{

			Name:  "Single parameter modified by constant",
			Input: "foo + 2",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: 2.0,
				},
			},
			Expected: 4.0,
		},
		{

			Name:  "Single parameter modified by variable",
			Input: "foo * bar",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: 5.0,
				},
				{
					Name:  "bar",
					Value: 2.0,
				},
			},
			Expected: 10.0,
		},
		{

			Name:  "Multiple multiplications of the same parameter",
			Input: "foo * foo * foo",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: 10.0,
				},
			},
			Expected: 1000.0,
		},
		{

			Name:  "Multiple additions of the same parameter",
			Input: "foo + foo + foo",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: 10.0,
				},
			},
			Expected: 30.0,
		},
		{

			Name:  "Parameter name sensitivity",
			Input: "foo + FoO + FOO",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: 8.0,
				},
				{
					Name:  "FoO",
					Value: 4.0,
				},
				{
					Name:  "FOO",
					Value: 2.0,
				},
			},
			Expected: 14.0,
		},
		{

			Name:  "Sign prefix comparison against prefixed variable",
			Input: "-1 < -foo",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: -8.0,
				},
			},
			Expected: true,
		},
		{

			Name:  "Fixed-point parameter",
			Input: "foo > 1",
			Parameters: []EvaluationParameter{

				{
					Name:  "foo",
					Value: 2,
				},
			},
			Expected: true,
		},
		{

			Name:     "Modifier after closing clause",
			Input:    "(2 + 2) + 2 == 6",
			Expected: true,
		},
		{

			Name:     "Comparator after closing clause",
			Input:    "(2 + 2) >= 4",
			Expected: true,
		},
		{

			Name:  "Two-boolean logical operation (for issue #8)",
			Input: "(foo == true) || (bar == true)",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: true,
				},
				{
					Name:  "bar",
					Value: false,
				},
			},
			Expected: true,
		},
		{

			Name:  "Two-variable integer logical operation (for issue #8)",
			Input: "foo > 10 && bar > 10",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: 1,
				},
				{
					Name:  "bar",
					Value: 11,
				},
			},
			Expected: false,
		},
		{

			Name:  "Regex against right-hand parameter",
			Input: "'foobar' =~ foo",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "obar",
				},
			},
			Expected: true,
		},
		{

			Name:  "Not-regex against right-hand parameter",
			Input: "'foobar' !~ foo",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
			},
			Expected: true,
		},
		{

			Name:  "Regex against two parameters",
			Input: "foo =~ bar",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "foobar",
				},
				{
					Name:  "bar",
					Value: "oba",
				},
			},
			Expected: true,
		},
		{

			Name:  "Not-regex against two parameters",
			Input: "foo !~ bar",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "foobar",
				},
				{
					Name:  "bar",
					Value: "baz",
				},
			},
			Expected: true,
		},
		{

			Name:  "Pre-compiled regex",
			Input: "foo =~ bar",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "foobar",
				},
				{
					Name:  "bar",
					Value: regexp.MustCompile("[fF][oO]+"),
				},
			},
			Expected: true,
		},
		{

			Name:  "Pre-compiled not-regex",
			Input: "foo !~ bar",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "foobar",
				},
				{
					Name:  "bar",
					Value: regexp.MustCompile("[fF][oO]+"),
				},
			},
			Expected: false,
		},
		{

			Name:  "Single boolean parameter",
			Input: "commission ? 10",
			Parameters: []EvaluationParameter{
				{
					Name:  "commission",
					Value: true,
				},
			},
			Expected: 10.0,
		},
		{

			Name:  "True comparator with a parameter",
			Input: "partner == 'amazon' ? 10",
			Parameters: []EvaluationParameter{
				{
					Name:  "partner",
					Value: "amazon",
				},
			},
			Expected: 10.0,
		},
		{

			Name:  "False comparator with a parameter",
			Input: "partner == 'amazon' ? 10",
			Parameters: []EvaluationParameter{
				{
					Name:  "partner",
					Value: "ebay",
				},
			},
			Expected: nil,
		},
		{

			Name:  "True comparator with multiple parameters",
			Input: "theft && period == 24 ? 60",
			Parameters: []EvaluationParameter{
				{
					Name:  "theft",
					Value: true,
				},
				{
					Name:  "period",
					Value: 24,
				},
			},
			Expected: 60.0,
		},
		{

			Name:  "False comparator with multiple parameters",
			Input: "theft && period == 24 ? 60",
			Parameters: []EvaluationParameter{
				{
					Name:  "theft",
					Value: false,
				},
				{
					Name:  "period",
					Value: 24,
				},
			},
			Expected: nil,
		},
		{

			Name:  "String concat with single string parameter",
			Input: "foo + 'bar'",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
			},
			Expected: "bazbar",
		},
		{

			Name:  "String concat with multiple string parameter",
			Input: "foo + bar",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
				{
					Name:  "bar",
					Value: "quux",
				},
			},
			Expected: "bazquux",
		},
		{

			Name:  "String concat with float parameter",
			Input: "foo + bar",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
				{
					Name:  "bar",
					Value: 123.0,
				},
			},
			Expected: "baz123",
		},
		{

			Name:  "Mixed multiple string concat",
			Input: "foo + 123 + 'bar' + true",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
			},
			Expected: "baz123bartrue",
		},
		{

			Name:  "Integer width spectrum",
			Input: "uint8 + uint16 + uint32 + uint64 + int8 + int16 + int32 + int64",
			Parameters: []EvaluationParameter{
				{
					Name:  "uint8",
					Value: uint8(0),
				},
				{
					Name:  "uint16",
					Value: uint16(0),
				},
				{
					Name:  "uint32",
					Value: uint32(0),
				},
				{
					Name:  "uint64",
					Value: uint64(0),
				},
				{
					Name:  "int8",
					Value: int8(0),
				},
				{
					Name:  "int16",
					Value: int16(0),
				},
				{
					Name:  "int32",
					Value: int32(0),
				},
				{
					Name:  "int64",
					Value: int64(0),
				},
			},
			Expected: 0.0,
		},
		{

			Name:  "Floats",
			Input: "float32 + float64",
			Parameters: []EvaluationParameter{
				{
					Name:  "float32",
					Value: float32(0.0),
				},
				{
					Name:  "float64",
					Value: float64(0.0),
				},
			},
			Expected: 0.0,
		},
		{

			Name:  "Null coalesce right",
			Input: "foo ?? 1.0",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: nil,
				},
			},
			Expected: 1.0,
		},
		{

			Name:  "Multiple comparator/logical operators (#30)",
			Input: "(foo >= 2887057408 && foo <= 2887122943) || (foo >= 168100864 && foo <= 168118271)",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: 2887057409,
				},
			},
			Expected: true,
		},
		{

			Name:  "Multiple comparator/logical operators, opposite order (#30)",
			Input: "(foo >= 168100864 && foo <= 168118271) || (foo >= 2887057408 && foo <= 2887122943)",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: 2887057409,
				},
			},
			Expected: true,
		},
		{

			Name:  "Multiple comparator/logical operators, small value (#30)",
			Input: "(foo >= 2887057408 && foo <= 2887122943) || (foo >= 168100864 && foo <= 168118271)",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: 168100865,
				},
			},
			Expected: true,
		},
		{

			Name:  "Multiple comparator/logical operators, small value, opposite order (#30)",
			Input: "(foo >= 168100864 && foo <= 168118271) || (foo >= 2887057408 && foo <= 2887122943)",
			Parameters: []EvaluationParameter{
				{
					Name:  "foo",
					Value: 168100865,
				},
			},
			Expected: true,
		},
		{

			Name:  "Incomparable array equality comparison",
			Input: "arr == arr",
			Parameters: []EvaluationParameter{
				{
					Name:  "arr",
					Value: []int{0, 0, 0},
				},
			},
			Expected: true,
		},
		{

			Name:  "Incomparable array not-equality comparison",
			Input: "arr != arr",
			Parameters: []EvaluationParameter{
				{
					Name:  "arr",
					Value: []int{0, 0, 0},
				},
			},
			Expected: false,
		},
		{

			Name:  "Mixed function and parameters",
			Input: "sum(1.2, amount) + name",
			Functions: map[string]ExpressionFunction{
				"sum": func(arguments ...interface{}) (interface{}, error) {

					sum := 0.0
					for _, v := range arguments {
						sum += v.(float64)
					}
					return sum, nil
				},
			},
			Parameters: []EvaluationParameter{
				{
					Name:  "amount",
					Value: .8,
				},
				{
					Name:  "name",
					Value: "awesome",
				},
			},

			Expected: "2awesome",
		},
		{

			Name:  "Short-circuit OR",
			Input: "true || fail()",
			Functions: map[string]ExpressionFunction{
				"fail": func(arguments ...interface{}) (interface{}, error) {
					return nil, errors.New("Did not short-circuit")
				},
			},
			Expected: true,
		},
		{

			Name:  "Short-circuit AND",
			Input: "false && fail()",
			Functions: map[string]ExpressionFunction{
				"fail": func(arguments ...interface{}) (interface{}, error) {
					return nil, errors.New("Did not short-circuit")
				},
			},
			Expected: false,
		},
		{

			Name:  "Short-circuit ternary",
			Input: "true ? 1 : fail()",
			Functions: map[string]ExpressionFunction{
				"fail": func(arguments ...interface{}) (interface{}, error) {
					return nil, errors.New("Did not short-circuit")
				},
			},
			Expected: 1.0,
		},
		{

			Name:  "Short-circuit coalesce",
			Input: "'foo' ?? fail()",
			Functions: map[string]ExpressionFunction{
				"fail": func(arguments ...interface{}) (interface{}, error) {
					return nil, errors.New("Did not short-circuit")
				},
			},
			Expected: "foo",
		},
		{

			Name:       "Simple parameter call",
			Input:      "foo.String",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   fooParameter.Value.(dummyParameter).String,
		},
		{

			Name:       "Simple parameter function call",
			Input:      "foo.Func()",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "funk",
		},
		{

			Name:       "Simple parameter call from pointer",
			Input:      "fooptr.String",
			Parameters: []EvaluationParameter{fooPtrParameter},
			Expected:   fooParameter.Value.(dummyParameter).String,
		},
		{

			Name:       "Simple parameter function call from pointer",
			Input:      "fooptr.Func()",
			Parameters: []EvaluationParameter{fooPtrParameter},
			Expected:   "funk",
		},
		{

			Name:       "Simple parameter function call from pointer",
			Input:      "fooptr.Func3()",
			Parameters: []EvaluationParameter{fooPtrParameter},
			Expected:   "fronk",
		},
		{

			Name:       "Simple parameter call",
			Input:      "foo.String == 'hi'",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   false,
		},
		{

			Name:       "Simple parameter call with modifier",
			Input:      "foo.String + 'hi'",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   fooParameter.Value.(dummyParameter).String + "hi",
		},
		{

			Name:       "Simple parameter function call, two-arg return",
			Input:      "foo.Func2()",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "frink",
		},
		{

			Name:       "Parameter function call with all argument types",
			Input:      "foo.TestArgs(\"hello\", 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1.0, 2.0, true)",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "hello: 33",
		},

		{

			Name:       "Simple parameter function call, one arg",
			Input:      "foo.FuncArgStr('boop')",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "boop",
		},
		{

			Name:       "Simple parameter function call, one arg",
			Input:      "foo.FuncArgStr('boop') + 'hi'",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "boophi",
		},
		{

			Name:       "Nested parameter function call",
			Input:      "foo.Nested.Dunk('boop')",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "boopdunk",
		},
		{

			Name:       "Nested parameter call",
			Input:      "foo.Nested.Funk",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "funkalicious",
		},
		{

			Name:       "Parameter call with + modifier",
			Input:      "1 + foo.Int",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   102.0,
		},
		{

			Name:       "Parameter string call with + modifier",
			Input:      "'woop' + (foo.String)",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   "woopstring!",
		},
		{

			Name:       "Parameter call with && operator",
			Input:      "true && foo.BoolFalse",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   false,
		},
		{

			Name:       "Null coalesce nested parameter",
			Input:      "foo.Nil ?? false",
			Parameters: []EvaluationParameter{fooParameter},
			Expected:   false,
		},
	}

	runEvaluationTests(evaluationTests, test)
}

/*
	Tests the behavior of a nil set of parameters.
*/
func TestNilParameters(test *testing.T) {

	expression, _ := NewEvaluableExpression("true")
	_, err := expression.Evaluate(nil)

	if err != nil {
		test.Fail()
	}
}

/*
	Tests functionality related to using functions with a struct method receiver.
	Created to test #54.
*/
func TestStructFunctions(test *testing.T) {

	parseFormat := "2006"
	y2k, _ := time.Parse(parseFormat, "2000")
	y2k1, _ := time.Parse(parseFormat, "2001")

	functions := map[string]ExpressionFunction{
		"func1": func(args ...interface{}) (interface{}, error) {
			return float64(y2k.Year()), nil
		},
		"func2": func(args ...interface{}) (interface{}, error) {
			return float64(y2k1.Year()), nil
		},
	}

	exp, _ := NewEvaluableExpressionWithFunctions("func1() + func2()", functions)
	result, _ := exp.Evaluate(nil)

	if result != 4001.0 {
		test.Logf("Function calling method did not return the right value. Got: %v, expected %d\n", result, 4001)
		test.Fail()
	}
}

func runEvaluationTests(evaluationTests []EvaluationTest, test *testing.T) {

	var expression *EvaluableExpression
	var result interface{}
	var parameters map[string]interface{}
	var err error

	fmt.Printf("Running %d evaluation test cases...\n", len(evaluationTests))

	// Run the test cases.
	for _, evaluationTest := range evaluationTests {

		if evaluationTest.Functions != nil {
			expression, err = NewEvaluableExpressionWithFunctions(evaluationTest.Input, evaluationTest.Functions)
		} else {
			expression, err = NewEvaluableExpression(evaluationTest.Input)
		}

		if err != nil {

			test.Logf("Test '%s' failed to parse: '%s'", evaluationTest.Name, err)
			test.Fail()
			continue
		}

		parameters = make(map[string]interface{}, 8)

		for _, parameter := range evaluationTest.Parameters {
			parameters[parameter.Name] = parameter.Value
		}

		result, err = expression.Evaluate(parameters)

		if err != nil {

			test.Logf("Test '%s' failed", evaluationTest.Name)
			test.Logf("Encountered error: %s", err.Error())
			test.Fail()
			continue
		}

		if result != evaluationTest.Expected {

			test.Logf("Test '%s' failed", evaluationTest.Name)
			test.Logf("Evaluation result '%v' does not match expected: '%v'", result, evaluationTest.Expected)
			test.Fail()
		}
	}
}
