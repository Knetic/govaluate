package govaluate

/*
  Contains test cases for all the expression examples given in the README.
  While all of the functionality for these cases should be covered in other tests,
  this is really just a sanity check.
*/
import (
	"testing"
)

func TestBasicEvaluation(test *testing.T) {

	expression, _ := NewEvaluableExpression("10 > 0")
	result, _ := expression.Evaluate(nil)

	if result != true {
		test.Logf("Expected 'true', got '%v'\n", result)
		test.Fail()
	}
}

func TestParameterEvaluation(test *testing.T) {

	expression, _ := NewEvaluableExpression("foo > 0")

	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = -1

	result, _ := expression.Evaluate(parameters)

	if result != false {
		test.Logf("Expected 'false', got '%v'\n", result)
		test.Fail()
	}
}

func TestModifierEvaluation(test *testing.T) {

	expression, _ := NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90")

	parameters := make(map[string]interface{}, 8)
	parameters["requests_made"] = 100
	parameters["requests_succeeded"] = 80

	result, _ := expression.Evaluate(parameters)

	if result != false {
		test.Logf("Expected 'false', got '%v'\n", result)
		test.Fail()
	}
}

func TestStringEvaluation(test *testing.T) {

	expression, _ := NewEvaluableExpression("http_response_body == 'service is ok'")

	parameters := make(map[string]interface{}, 8)
	parameters["http_response_body"] = "service is ok"

	result, _ := expression.Evaluate(parameters)

	if result != true {
		test.Logf("Expected 'false', got '%v'\n", result)
		test.Fail()
	}
}

func TestFloatEvaluation(test *testing.T) {

	expression, _ := NewEvaluableExpression("(mem_used / total_mem) * 100")

	parameters := make(map[string]interface{}, 8)
	parameters["total_mem"] = 1024
	parameters["mem_used"] = 512

	result, _ := expression.Evaluate(parameters)

	if result != 50.0 {
		test.Logf("Expected '50.0', got '%v'\n", result)
		test.Fail()
	}
}

func TestDateComparison(test *testing.T) {

	expression, _ := NewEvaluableExpression("'2014-01-02' > '2014-01-01 23:59:59'")
	result, _ := expression.Evaluate(nil)

	if result != true {
		test.Logf("Expected 'true', got '%v'\n", result)
		test.Fail()
	}
}

func TestMultipleEvaluation(test *testing.T) {
	expression, _ := NewEvaluableExpression("response_time <= 100")
	parameters := make(map[string]interface{}, 8)

	for i := 0; i < 64; i++ {
		parameters["response_time"] = i
		result, _ := expression.Evaluate(parameters)

		if result != true {
			test.Logf("Expected 'true', got '%v'\n", result)
			test.Fail()
		}
	}
}
