package govaluate

import (
	"testing"
)

func TestMultipleOperators(test *testing.T) {

	tests := []EvaluationTest {
		EvaluationTest{
			Name:     "Incorrect subtract behavior",
			Input:    "2 - 6 - 10 - 2",
			Expected: -16.0,
		},
	}

	runEvaluationTests(tests, test)
}
