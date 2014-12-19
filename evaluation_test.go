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

			Name: "Single parameter modified by constant",
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

			test.Log("Test '",evaluationTest.Name,"' failed to parse: ", err)
			test.Fail()
			continue
		}

		parameters = make(map[string]interface{}, 8)		
		
		for _, parameter := range evaluationTest.Parameters {
			parameters[parameter.Name] = parameter.Value
		}

		result, err = expression.Evaluate(parameters)

		if(err != nil) {

			test.Log("Test '", evaluationTest.Name, "' failed:")
			test.Log("Encountered error: " + err.Error())
			test.Fail()
			continue;
		}

		if(result != evaluationTest.Expected) {

			test.Log("Test '", evaluationTest.Name, "' failed:")
			test.Log("Expected evaluation result '", evaluationTest.Expected, "' does not match '", result, "'")
			test.Fail()
		}
	}
}

