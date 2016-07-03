package govaluate

import (
	"fmt"
	"regexp"
	"testing"
)

/*
	Represents a test of expression evaluation
*/
type EvaluationTest struct {
	Name       string
	Input      string
	Parameters []EvaluationParameter
	Expected   interface{}
}

type EvaluationParameter struct {
	Name  string
	Value interface{}
}

func TestNoParameterEvaluation(test *testing.T) {

	evaluationTests := []EvaluationTest{

		EvaluationTest{

			Name:     "Single PLUS",
			Input:    "51 + 49",
			Expected: 100.0,
		},
		EvaluationTest{

			Name:     "Single MINUS",
			Input:    "100 - 51",
			Expected: 49.0,
		},
		EvaluationTest{

			Name:     "Single MULTIPLY",
			Input:    "5 * 20",
			Expected: 100.0,
		},
		EvaluationTest{

			Name:     "Single DIVIDE",
			Input:    "100 / 20",
			Expected: 5.0,
		},
		EvaluationTest{

			Name:     "Single even MODULUS",
			Input:    "100 % 2",
			Expected: 0.0,
		},
		EvaluationTest{

			Name:     "Single odd MODULUS",
			Input:    "101 % 2",
			Expected: 1.0,
		},
		EvaluationTest{

			Name:     "Compound PLUS",
			Input:    "20 + 30 + 50",
			Expected: 100.0,
		},
		EvaluationTest{

			Name:     "Mutiple operators",
			Input:    "20 * 5 - 49",
			Expected: 51.0,
		},
		EvaluationTest{

			Name:     "Parenthesis usage",
			Input:    "100 - (5 * 10)",
			Expected: 50.0,
		},
		EvaluationTest{

			Name:     "Nested parentheses",
			Input:    "50 + (5 * (15 - 5))",
			Expected: 100.0,
		},
		EvaluationTest{

			Name:     "Logical OR operation of two clauses",
			Input:    "(1 == 1) || (true == true)",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Logical AND operation of two clauses",
			Input:    "(1 == 1) && (true == true)",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Implicit boolean",
			Input:    "2 > 1",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Compound boolean",
			Input:    "5 < 10 && 1 < 5",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Evaluated true && false operation (for issue #8)",
			Input:    "1 > 10 && 11 > 10",
			Expected: false,
		},
		EvaluationTest{

			Name:     "Evaluated true && false operation (for issue #8)",
			Input:    "true == true && false == true",
			Expected: false,
		},
		EvaluationTest{

			Name:     "Parenthesis boolean",
			Input:    "10 < 50 && (1 != 2 && 1 > 0)",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Comparison of string constants",
			Input:    "'foo' == 'foo'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "NEQ comparison of string constants",
			Input:    "'foo' != 'bar'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "REQ comparison of string constants",
			Input:    "'foobar' =~ 'oba'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "NREQ comparison of string constants",
			Input:    "'foo' !~ 'bar'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Multiplicative/additive order",
			Input:    "5 + 10 * 2",
			Expected: 25.0,
		},
		EvaluationTest{

			Name:     "Multiple constant multiplications",
			Input:    "10 * 10 * 10",
			Expected: 1000.0,
		},
		EvaluationTest{

			Name:     "Multiple adds/multiplications",
			Input:    "10 * 10 * 10 + 1 * 10 * 10",
			Expected: 1100.0,
		},
		EvaluationTest{

			Name:     "Modulus precedence",
			Input:    "1 + 101 % 2 * 5",
			Expected: 2.0,
		},
		EvaluationTest{

			Name:     "Exponent precedence",
			Input:    "1 + 5 ^ 3 % 2 * 5",
			Expected: 6.0,
		},
		EvaluationTest{

			Name:     "Identical date equivalence",
			Input:    "'2014-01-02 14:12:22' == '2014-01-02 14:12:22'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Positive date GT",
			Input:    "'2014-01-02 14:12:22' > '2014-01-02 12:12:22'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Negative date GT",
			Input:    "'2014-01-02 14:12:22' > '2014-01-02 16:12:22'",
			Expected: false,
		},
		EvaluationTest{

			Name:     "Positive date GTE",
			Input:    "'2014-01-02 14:12:22' >= '2014-01-02 12:12:22'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Negative date GTE",
			Input:    "'2014-01-02 14:12:22' >= '2014-01-02 16:12:22'",
			Expected: false,
		},
		EvaluationTest{

			Name:     "Positive date LT",
			Input:    "'2014-01-02 14:12:22' < '2014-01-02 16:12:22'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Negative date LT",
			Input:    "'2014-01-02 14:12:22' < '2014-01-02 11:12:22'",
			Expected: false,
		},
		EvaluationTest{

			Name:     "Positive date LTE",
			Input:    "'2014-01-02 09:12:22' <= '2014-01-02 12:12:22'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Negative date LTE",
			Input:    "'2014-01-02 14:12:22' <= '2014-01-02 11:12:22'",
			Expected: false,
		},
		EvaluationTest{

			Name:     "Sign prefix comparison",
			Input:    "-1 < 0",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Boolean sign prefix comparison",
			Input:    "!true == false",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Inversion of clause",
			Input:    "!(10 < 0)",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Negation after modifier",
			Input:    "10 * -10",
			Expected: -100.0,
		},
		EvaluationTest{

			Name:     "Ternary with single boolean",
			Input:    "true ? 10",
			Expected: 10.0,
		},
		EvaluationTest{

			Name:     "Ternary nil with single boolean",
			Input:    "false ? 10",
			Expected: nil,
		},
		EvaluationTest{

			Name:     "Ternary with comparator boolean",
			Input:    "10 > 5 ? 35.50",
			Expected: 35.50,
		},
		EvaluationTest{

			Name:     "Ternary nil with comparator boolean",
			Input:    "1 > 5 ? 35.50",
			Expected: nil,
		},
		EvaluationTest{

			Name:     "Ternary with parentheses",
			Input:    "(5 * (15 - 5)) > 5 ? 35.50",
			Expected: 35.50,
		},
		EvaluationTest{

			Name:     "Ternary precedence",
			Input:    "true ? 35.50 > 10",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Ternary-else",
			Input:    "false ? 35.50 : 50",
			Expected: 50.0,
		},
		EvaluationTest{

			Name:     "Ternary-else inside clause",
			Input:    "(false ? 5 : 35.50) > 10",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Ternary-else (true-case) inside clause",
			Input:    "(true ? 35.50 : 5) > 10",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Ternary-else before comparator (negative case)",
			Input:    "true ? 35.50 : 5 > 10",
			Expected: 35.50,
		},
		EvaluationTest{

			Name:     "String to string concat",
			Input:    "'foo' + 'bar' == 'foobar'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "String to float64 concat",
			Input:    "'foo' + 123 == 'foo123'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Float64 to string concat",
			Input:    "123 + 'bar' == '123bar'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "String to date concat",
			Input:    "'foo' + '02/05/1970' == 'foobar'",
			Expected: false,
		},
		EvaluationTest{

			Name:     "String to bool concat",
			Input:    "'foo' + true == 'footrue'",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Bool to string concat",
			Input:    "true + 'bar' == 'truebar'",
			Expected: true,
		},
	}

	runEvaluationTests(evaluationTests, test)
}

func TestParameterizedEvaluation(test *testing.T) {

	evaluationTests := []EvaluationTest{

		EvaluationTest{

			Name:  "Single parameter modified by constant",
			Input: "foo + 2",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: 2.0,
				},
			},
			Expected: 4.0,
		},
		EvaluationTest{

			Name:  "Single parameter modified by variable",
			Input: "foo * bar",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: 5.0,
				},
				EvaluationParameter{
					Name:  "bar",
					Value: 2.0,
				},
			},
			Expected: 10.0,
		},
		EvaluationTest{

			Name:  "Multiple multiplications of the same parameter",
			Input: "foo * foo * foo",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: 10.0,
				},
			},
			Expected: 1000.0,
		},
		EvaluationTest{

			Name:  "Multiple additions of the same parameter",
			Input: "foo + foo + foo",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: 10.0,
				},
			},
			Expected: 30.0,
		},
		EvaluationTest{

			Name:  "Parameter name sensitivity",
			Input: "foo + FoO + FOO",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: 8.0,
				},
				EvaluationParameter{
					Name:  "FoO",
					Value: 4.0,
				},
				EvaluationParameter{
					Name:  "FOO",
					Value: 2.0,
				},
			},
			Expected: 14.0,
		},
		EvaluationTest{

			Name:  "Sign prefix comparison against prefixed variable",
			Input: "-1 < -foo",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: -8.0,
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Fixed-point parameter",
			Input: "foo > 1",
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "foo",
					Value: 2,
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:     "Modifier after closing clause",
			Input:    "(2 + 2) + 2 == 6",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Comparator after closing clause",
			Input:    "(2 + 2) >= 4",
			Expected: true,
		},
		EvaluationTest{

			Name:  "Two-boolean logical operation (for issue #8)",
			Input: "(foo == true) || (bar == true)",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: true,
				},
				EvaluationParameter{
					Name:  "bar",
					Value: false,
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Two-variable integer logical operation (for issue #8)",
			Input: "foo > 10 && bar > 10",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: 1,
				},
				EvaluationParameter{
					Name:  "bar",
					Value: 11,
				},
			},
			Expected: false,
		},
		EvaluationTest{

			Name:  "Regex against right-hand parameter",
			Input: "'foobar' =~ foo",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "obar",
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Not-regex against right-hand paramter",
			Input: "'foobar' !~ foo",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "baz",
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Regex against two parameter",
			Input: "foo =~ bar",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "foobar",
				},
				EvaluationParameter{
					Name:  "bar",
					Value: "oba",
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Not-regex against two paramter",
			Input: "foo !~ bar",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "foobar",
				},
				EvaluationParameter{
					Name:  "bar",
					Value: "baz",
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Pre-compiled regex",
			Input: "foo =~ bar",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "foobar",
				},
				EvaluationParameter{
					Name:  "bar",
					Value: regexp.MustCompile("[fF][oO]+"),
				},
			},
			Expected: true,
		},
		EvaluationTest{

			Name:  "Pre-compiled not-regex",
			Input: "foo !~ bar",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "foobar",
				},
				EvaluationParameter{
					Name:  "bar",
					Value: regexp.MustCompile("[fF][oO]+"),
				},
			},
			Expected: false,
		},
		EvaluationTest{

			Name:  "Single boolean parameter",
			Input: "commission ? 10",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "commission",
					Value: true,
				},
			},
			Expected: 10.0,
		},
		EvaluationTest{

			Name:  "True comparator with a parameter",
			Input: "partner == 'amazon' ? 10",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "partner",
					Value: "amazon",
				},
			},
			Expected: 10.0,
		},
		EvaluationTest{

			Name:  "False comparator with a parameter",
			Input: "partner == 'amazon' ? 10",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "partner",
					Value: "ebay",
				},
			},
			Expected: nil,
		},
		EvaluationTest{

			Name:  "True comparator with multiple parameters",
			Input: "theft && period == 24 ? 60",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "theft",
					Value: true,
				},
				EvaluationParameter{
					Name:  "period",
					Value: 24,
				},
			},
			Expected: 60.0,
		},
		EvaluationTest{

			Name:  "False comparator with multiple parameters",
			Input: "theft && period == 24 ? 60",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "theft",
					Value: false,
				},
				EvaluationParameter{
					Name:  "period",
					Value: 24,
				},
			},
			Expected: nil,
		},
		EvaluationTest{

			Name:  "String concat with single string parameter",
			Input: "foo + 'bar'",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "baz",
				},
			},
			Expected: "bazbar",
		},
		EvaluationTest{

			Name:  "String concat with multiple string parameter",
			Input: "foo + bar",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "baz",
				},
				EvaluationParameter{
					Name:  "bar",
					Value: "quux",
				},
			},
			Expected: "bazquux",
		},
		EvaluationTest{

			Name:  "String concat with float parameter",
			Input: "foo + bar",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "baz",
				},
				EvaluationParameter{
					Name:  "bar",
					Value: 123.0,
				},
			},
			Expected: "baz123",
		},
		EvaluationTest{

			Name:  "Mixed multiple string concat",
			Input: "foo + 123 + 'bar' + true",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "foo",
					Value: "baz",
				},
			},
			Expected: "baz123bartrue",
		},
		EvaluationTest{

			Name:  "Integer width spectrum",
			Input: "uint8 + uint16 + uint32 + uint64 + int8 + int16 + int32 + int64",
			Parameters: []EvaluationParameter{
				EvaluationParameter{
					Name:  "uint8",
					Value: uint8(0),
				},
				EvaluationParameter{
					Name:  "uint16",
					Value: uint16(0),
				},
				EvaluationParameter{
					Name:  "uint32",
					Value: uint32(0),
				},
				EvaluationParameter{
					Name:  "uint64",
					Value: uint64(0),
				},
				EvaluationParameter{
					Name:  "int8",
					Value: int8(0),
				},
				EvaluationParameter{
					Name:  "int16",
					Value: int16(0),
				},
				EvaluationParameter{
					Name:  "int32",
					Value: int32(0),
				},
				EvaluationParameter{
					Name:  "int64",
					Value: int64(0),
				},
			},
			Expected: 0.0,
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

func runEvaluationTests(evaluationTests []EvaluationTest, test *testing.T) {

	var expression *EvaluableExpression
	var result interface{}
	var parameters map[string]interface{}
	var err error

	fmt.Printf("Running %d evaluation test cases...\n", len(evaluationTests))

	// Run the test cases.
	for _, evaluationTest := range evaluationTests {

		expression, err = NewEvaluableExpression(evaluationTest.Input)

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
