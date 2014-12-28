package govaluate

import (
	"testing"
)

/*
	Represents a test of expression evaluation
*/
type EvaluationTest struct {

	Name string
	Input string
	Parameters []EvaluationParameter
	Expected interface{}
}

type EvaluationParameter struct {

	Name string
	Value interface{}
}

func TestNoParameterEvaluation(test *testing.T) {

	evaluationTests := []EvaluationTest {

		EvaluationTest {

			Name: "Single PLUS",
			Input: "51 + 49",
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Single MINUS",
			Input: "100 - 51",
			Expected: 49.0,
		},
		EvaluationTest {

			Name: "Single MULTIPLY",
			Input: "5 * 20",
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Single DIVIDE",
			Input: "100 / 20",
			Expected: 5.0,
		},
		EvaluationTest {

			Name: "Single even MODULUS",
			Input: "100 % 2",
			Expected: 0.0,
		},
		EvaluationTest {

			Name: "Single odd MODULUS",
			Input: "101 % 2",
			Expected: 1.0,
		},
		EvaluationTest {

			Name: "Compound PLUS",
			Input: "20 + 30 + 50",
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Mutiple operators",
			Input: "20 * 5 - 49",
			Expected: 51.0,
		},
		EvaluationTest {

			Name: "Paren usage",
			Input: "100 - (5 * 10)",
			Expected: 50.0,
		},
		EvaluationTest {

			Name: "Nested parentheses",
			Input: "50 + (5 * (15 - 5))",
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Implicit boolean",
			Input: "2 > 1",
			Expected: true,
		},
		EvaluationTest {

			Name: "Compound boolean",
			Input: "5 < 10 && 1 < 5",
			Expected: true,
		},
		EvaluationTest {

			Name: "Parenthesis boolean",
			Input: "10 < 50 && (1 != 2 && 1 > 0)",
			Expected: true,
		},
		EvaluationTest {

			Name: "Comparison of string constants",
			Input: "'foo' == 'foo'",
			Expected: true,
		},
		EvaluationTest {

			Name: "NEQ comparison of string constants",
			Input: "'foo' != 'bar'",
			Expected: true,
		},
		EvaluationTest {

			Name: "Multiplicative/additive order",
			Input: "5 + 10 * 2",
			Expected: 25.0,
		},
		EvaluationTest {

			Name: "Multiple constant multiplications",
			Input: "10 * 10 * 10",
			Expected: 1000.0,
		},
		EvaluationTest {

			Name: "Multiple adds/multiplications",
			Input: "10 * 10 * 10 + 1 * 10 * 10",
			Expected: 1100.0,
		},
		EvaluationTest {

			Name: "Modulus precedence",
			Input: "1 + 101 % 2 * 5",
			Expected: 2.0,
		},
	}

	runEvaluationTests(evaluationTests, test)
}

func TestParameterizedEvaluation(test *testing.T) {

	evaluationTests := []EvaluationTest {

		EvaluationTest {

			Name: "Single parameter modified by constant",
			Input: "foo + 2",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 2.0,
				},
			},
			Expected: 4.0,
		},
		EvaluationTest {

			Name: "Single parameter modified by variable",
			Input: "foo * bar",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 5.0,
				},
				EvaluationParameter {
					Name: "bar",
					Value: 2.0,
				},
			},
			Expected: 10.0,
		},
		EvaluationTest {

			Name: "Multiple multiplications of the same parameter",
			Input: "foo * foo * foo",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 10.0,
				},
			},
			Expected: 1000.0,
		},
		EvaluationTest {

			Name: "Multiple additions of the same parameter",
			Input: "foo + foo + foo",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 10.0,
				},
			},
			Expected: 30.0,
		},
		EvaluationTest {

			Name: "Parameter name sensitivity",
			Input: "foo + FoO + FOO",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 8.0,
				},
				EvaluationParameter {
					Name: "FoO",
					Value: 4.0,
				},
				EvaluationParameter {
					Name: "FOO",
					Value: 2.0,
				},
			},
			Expected: 14.0,
		},
	}

	runEvaluationTests(evaluationTests, test)
}

func runEvaluationTests(evaluationTests []EvaluationTest, test *testing.T) {

	var expression *EvaluableExpression
	var result interface{}
	var parameters map[string]interface{}
	var err error

	// Run the test cases.
	for _, evaluationTest := range evaluationTests {

		expression, err = NewEvaluableExpression(evaluationTest.Input)

		if(err != nil) {

			test.Logf("Test '%s' failed to parse: '%s'", evaluationTest.Name, err)
			test.Fail()
			continue
		}

		parameters = make(map[string]interface{}, 8)		
		
		for _, parameter := range evaluationTest.Parameters {
			parameters[parameter.Name] = parameter.Value
		}

		result, err = expression.Evaluate(parameters)

		if(err != nil) {

			test.Logf("Test '%s' failed", evaluationTest.Name)
			test.Logf("Encountered error: %s", err.Error())
			test.Fail()
			continue;
		}

		if(result != evaluationTest.Expected) {

			test.Logf("Test '%s' failed", evaluationTest.Name)
			test.Logf("Evaluation result '%v' does not match expected: '%v'", result, evaluationTest.Expected)
			test.Fail()
		}
	}
}

