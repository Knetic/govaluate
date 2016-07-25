package govaluate

import (
	"testing"
)

func TestMultipleOperators(test *testing.T) {

	tests := []EvaluationTest {
		EvaluationTest{
			Name:     "Incorrect subtract behavior",
			Input:    "1 - 2 - 4 - 8",
			Expected: -13.0,
		},
	}

	runEvaluationTests(tests, test)
}
