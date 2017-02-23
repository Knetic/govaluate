package govaluate

/*
  Tests to make sure evaluation fails in the expected ways.
*/
import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type DebugStruct struct {
	x int
}

/*
	Represents a test for parsing failures
*/
type EvaluationFailureTest struct {
	Name       string
	Input      string
	Functions  map[string]ExpressionFunction
	Parameters map[string]interface{}
	Expected   string
}

const (
	INVALID_MODIFIER_TYPES   string = "cannot be used with the modifier"
	INVALID_COMPARATOR_TYPES        = "cannot be used with the comparator"
	INVALID_LOGICALOP_TYPES         = "cannot be used with the logical operator"
	INVALID_TERNARY_TYPES           = "cannot be used with the ternary operator"
	ABSENT_PARAMETER                = "No parameter"
	INVALID_REGEX                   = "Unable to compile regexp pattern"
)

// preset parameter map of types that can be used in an evaluation failure test to check typing.
var EVALUATION_FAILURE_PARAMETERS = map[string]interface{}{
	"number": 1,
	"string": "foo",
	"bool":   true,
}

func TestComplexParameter(test *testing.T) {

	var expression *EvaluableExpression
	var err error

	parameters := map[string]interface{}{
		"complex64":  complex64(0),
		"complex128": complex128(0),
	}

	expression, _ = NewEvaluableExpression("complex64")
	_, err = expression.Evaluate(parameters)
	if err == nil {
		test.Logf("Expected to fail when giving a complex64 value, did not")
		test.Fail()
	}

	expression, _ = NewEvaluableExpression("complex128")
	_, err = expression.Evaluate(parameters)
	if err == nil {
		test.Logf("Expected to fail when giving a complex128 value, did not")
		test.Fail()
	}
}

func TestStructParameter(test *testing.T) {

	expression, _ := NewEvaluableExpression("foo")
	parameters := map[string]interface{}{
		"foo": DebugStruct{},
	}

	_, err := expression.Evaluate(parameters)
	if err == nil {
		test.Logf("Expected to  fail when giving a struct value, did not")
		test.Fail()
	}
}

func TestNilParameterUsage(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:     "Absent parameter used",
			Input:    "foo > 1",
			Expected: ABSENT_PARAMETER,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestModifierTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{

			Name:     "PLUS literal number to literal bool",
			Input:    "1 + true",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "PLUS number to bool",
			Input:    "number + bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "MINUS number to bool",
			Input:    "number - bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "MINUS number to bool",
			Input:    "number - bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "MULTIPLY number to bool",
			Input:    "number * bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "DIVIDE number to bool",
			Input:    "number / bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "EXPONENT number to bool",
			Input:    "number ** bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "MODULUS number to bool",
			Input:    "number % bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "XOR number to bool",
			Input:    "number % bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "BITWISE_OR number to bool",
			Input:    "number | bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "BITWISE_AND number to bool",
			Input:    "number & bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "BITWISE_XOR number to bool",
			Input:    "number ^ bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "BITWISE_LSHIFT number to bool",
			Input:    "number << bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
		EvaluationFailureTest{

			Name:     "BITWISE_RSHIFT number to bool",
			Input:    "number >> bool",
			Expected: INVALID_MODIFIER_TYPES,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestLogicalOperatorTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{

			Name:     "AND number to number",
			Input:    "number && number",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "OR number to number",
			Input:    "number || number",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "AND string to string",
			Input:    "string && string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "OR string to string",
			Input:    "string || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "AND number to string",
			Input:    "number && string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "OR number to string",
			Input:    "number || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "AND bool to string",
			Input:    "bool && string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
		EvaluationFailureTest{

			Name:     "OR bool to string",
			Input:    "bool || string",
			Expected: INVALID_LOGICALOP_TYPES,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

/*
	While there is type-safe transitions checked at parse-time, tested in the "parsing_test" and "parsingFailure_test" files,
	we also need to make sure that we receive type mismatch errors during evaluation.
*/
func TestComparatorTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{

			Name:     "GT literal bool to literal bool",
			Input:    "true > true",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "GT bool to bool",
			Input:    "bool > bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "GTE bool to bool",
			Input:    "bool >= bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "LT bool to bool",
			Input:    "bool < bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "LTE bool to bool",
			Input:    "bool <= bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},

		EvaluationFailureTest{

			Name:     "GT number to string",
			Input:    "number > string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "GTE number to string",
			Input:    "number >= string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "LT number to string",
			Input:    "number < string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "REQ number to string",
			Input:    "number =~ string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "REQ number to bool",
			Input:    "number =~ bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "REQ bool to number",
			Input:    "bool =~ number",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "REQ bool to string",
			Input:    "bool =~ string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "NREQ number to string",
			Input:    "number !~ string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "NREQ number to bool",
			Input:    "number !~ bool",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "NREQ bool to number",
			Input:    "bool !~ number",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "NREQ bool to string",
			Input:    "bool !~ string",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "IN non-array numeric",
			Input:    "1 in 2",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "IN non-array string",
			Input:    "1 in 'foo'",
			Expected: INVALID_COMPARATOR_TYPES,
		},
		EvaluationFailureTest{

			Name:     "IN non-array boolean",
			Input:    "1 in true",
			Expected: INVALID_COMPARATOR_TYPES,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestTernaryTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{

			Name:     "Ternary with number",
			Input:    "10 ? true",
			Expected: INVALID_TERNARY_TYPES,
		},
		EvaluationFailureTest{

			Name:     "Ternary with string",
			Input:    "'foo' ? true",
			Expected: INVALID_TERNARY_TYPES,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestRegexParameterCompilation(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{

			Name:  "Regex equality runtime parsing",
			Input: "'foo' =~ foo",
			Parameters: map[string]interface{}{
				"foo": "[foo",
			},
			Expected: INVALID_REGEX,
		},
		EvaluationFailureTest{

			Name:  "Regex inequality runtime parsing",
			Input: "'foo' =~ foo",
			Parameters: map[string]interface{}{
				"foo": "[foo",
			},
			Expected: INVALID_REGEX,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestFunctionExecution(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{

			Name:  "Function error bubbling",
			Input: "error()",
			Functions: map[string]ExpressionFunction{
				"error": func(arguments ...interface{}) (interface{}, error) {
					return nil, errors.New("Huge problems")
				},
			},
			Expected: "Huge problems",
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func runEvaluationFailureTests(evaluationTests []EvaluationFailureTest, test *testing.T) {

	var expression *EvaluableExpression
	var err error

	fmt.Printf("Running %d negative parsing test cases...\n", len(evaluationTests))

	for _, testCase := range evaluationTests {

		if len(testCase.Functions) > 0 {
			expression, err = NewEvaluableExpressionWithFunctions(testCase.Input, testCase.Functions)
		} else {
			expression, err = NewEvaluableExpression(testCase.Input)
		}

		if err != nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected evaluation error, but got parsing error: '%s'", err)
			test.Fail()
			continue
		}

		if testCase.Parameters == nil {
			testCase.Parameters = EVALUATION_FAILURE_PARAMETERS
		}

		_, err = expression.Evaluate(testCase.Parameters)

		if err == nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected error, received none.")
			test.Fail()
			continue
		}

		if !strings.Contains(err.Error(), testCase.Expected) {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Got error: '%s', expected '%s'", err.Error(), testCase.Expected)
			test.Fail()
			continue
		}
	}
}
