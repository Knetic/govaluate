package govaluate

import (
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

			Name:     "Paren usage",
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

			Name:     "Evaluated false NAND false operation",
			Input:    "false NAND false",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Evaluated false nand true operation",
			Input:    "false NAND true",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Evaluated true nand false operation",
			Input:    "true NAND false",
			Expected: true,
		},
		EvaluationTest{

			Name:     "Evaluated true NAND true operation",
			Input:    "true NAND true",
			Expected: false,
		},

		EvaluationTest{
			Name:     "basic logical test true or false",
			Input:    "true OR false",
			Expected: true,
		},
		EvaluationTest{
			Name:     "basic logical test false or true",
			Input:    "false OR true",
			Expected: true,
		},
		EvaluationTest{
			Name:     "basic logical test false or false",
			Input:    "false OR false",
			Expected: false,
		},

		EvaluationTest{
			Name:     "basic logical test false and false",
			Input:    "false AND false",
			Expected: false,
		},
		EvaluationTest{
			Name:     "basic logical test true and true",
			Input:    "true AND true",
			Expected: true,
		},
		EvaluationTest{
			Name:     "basic logical test true and false",
			Input:    "true AND false",
			Expected: false,
		},
		EvaluationTest{
			Name:     "basic logical test false and true",
			Input:    "false AND true",
			Expected: false,
		},

		EvaluationTest{
			Name:     "basic logical test false and true 2",
			Input:    "(false AND true) OR true",
			Expected: true,
		},

		EvaluationTest{
			Name:     "false XOR false",
			Input:    "false XOR false",
			Expected: false,
		},

		EvaluationTest{
			Name:     "false xor true",
			Input:    "false xor true",
			Expected: true,
		},
		EvaluationTest{
			Name:     "true XOR false",
			Input:    "true XOR false",
			Expected: true,
		},
		EvaluationTest{
			Name:     "true xor true",
			Input:    "true xor true",
			Expected: false,
		},
		EvaluationTest{
			Name:     `"500" =~ /5\d\d/`,
			Input:    `"500" =~ /5\d\d/`,
			Expected: true,
		},
		EvaluationTest{
			Name:     `"500" !~ /5\d\d/`,
			Input:    `"500" !~ /4\d\d/`,
			Expected: true,
		},
		EvaluationTest{
			Name:     `string in a array of strings`,
			Input:    `"foo" in ["boo", "bar", "foo", "zob"]`,
			Expected: true,
		},
		EvaluationTest{
			Name:     `number in a array of numbers`,
			Input:    `6 in [4,5,6]`,
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
			Name:  `string in a param array of strings`,
			Input: `"foo" in [bar]`,
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "bar",
					Value: []string{"boo", "bar", "foo", "zob"},
				},
			},
			Expected: true,
		},
		EvaluationTest{
			Name:  `number in a param array of numbers`,
			Input: `4 in [bar]`,
			Parameters: []EvaluationParameter{

				EvaluationParameter{
					Name:  "bar",
					Value: []float64{4, 5, 6},
				},
			},
			Expected: true,
		},
	}

	runEvaluationTests(evaluationTests, test)
}

func runEvaluationTests(evaluationTests []EvaluationTest, test *testing.T) {

	var expression *EvaluableExpression
	var result interface{}
	var parameters map[string]interface{}
	var err error

	test.Logf("Running %d evaluation test cases", len(evaluationTests))

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
