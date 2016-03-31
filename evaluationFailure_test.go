package govaluate

/*
  Tests to make sure evaluation fails in the expected ways.
*/
import (
	"testing"
	"strings"
)

type DebugStruct struct {
	x int
}

/*
	Represents a test for parsing failures
*/
type EvaluationFailureTest struct {
	Name     string
	Input    string
	Expected string
}

const (
	INVALID_MODIFIER_TYPES string = "cannot be used with the modifier"
	INVALID_COMPARATOR_TYPES = "cannot be used with the comparator"
	INVALID_LOGICALOP_TYPES = "cannot be used with the logical operator"
)

// preset parameter map of types that can be used in an evaluation failure test to check typing.
var EVALUATION_FAILURE_PARAMETERS = map[string]interface{} {
	"number": 1,
	"string": "foo",
	"bool": true,
}

func TestComplexParameter(test *testing.T) {

	expression, _ := NewEvaluableExpression("1")
	parameters := map[string]interface{}{
		"foo": 1i,
	}

	_, err := expression.Evaluate(parameters)
	if err == nil {
		test.Logf("Expected to  fail when giving a complex value, did not")
		test.Fail()
	}
}

func TestStructParameter(test *testing.T) {

	expression, _ := NewEvaluableExpression("1")
	parameters := map[string]interface{}{
		"foo": DebugStruct{},
	}

	_, err := expression.Evaluate(parameters)
	if err == nil {
		test.Logf("Expected to  fail when giving a struct value, did not")
		test.Fail()
	}
}

/*
	While there is type-safe transitions checked at parse-time, tested in the "parsing_test" and "parsingFailure_test" files,
	we also need to make sure that we receive type mismatch errors during evaluation.
*/
func TestOperatorTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest {
		EvaluationFailureTest {

			Name:     "PLUS number to bool",
			Input:    "number + bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest {

			Name:     "MINUS number to bool",
			Input:    "number - bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest {

			Name:     "MINUS number to bool",
			Input:    "number - bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest {

			Name:     "MULTIPLY number to bool",
			Input:    "number * bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest {

			Name:     "DIVIDE number to bool",
			Input:    "number / bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest {

			Name:     "EXPONENT number to bool",
			Input:    "number ^ bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest {

			Name:     "MODULUS number to bool",
			Input:    "number % bool",
			Expected: INVALID_MODIFIER_TYPES,
		},

		EvaluationFailureTest {

			Name:     "AND number to number",
			Input:    "number || number",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest {

			Name:     "OR number to number",
			Input:    "number || number",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest {

			Name:     "AND string to string",
			Input:    "string || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest {

			Name:     "OR string to string",
			Input:    "string || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest {

			Name:     "AND number to string",
			Input:    "number || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest {

			Name:     "OR number to string",
			Input:    "number || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},

		EvaluationFailureTest {

			Name:     "LTE bool to bool",
			Input:    "bool <= bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func runEvaluationFailureTests(evaluationTests []EvaluationFailureTest, test *testing.T) {

	var expression *EvaluableExpression
	var err error

	test.Logf("Running %d parsing test cases", len(evaluationTests))

	for _, testCase := range evaluationTests {

		expression, err = NewEvaluableExpression(testCase.Input)

		if err != nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected evaluation error, but got parsing error: '%s'", err)
			test.Fail()
			continue
		}

		_, err = expression.Evaluate(EVALUATION_FAILURE_PARAMETERS)

		if err == nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected error, received none.")
			test.Fail()
			continue
		}

		if !strings.Contains(err.Error(), testCase.Expected) {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Got error: '%s', expected '%s'", testCase.Expected, err.Error())
			test.Fail()
			continue
		}
	}
}
