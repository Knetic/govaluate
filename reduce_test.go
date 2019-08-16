package govaluate

import "testing"

func TestReduce(test *testing.T) {
	runTest(test, "x + y * z", map[string]interface{}{
		"x": 1.0,
		"z": 7.0,
	}, "1 + [y] * 7")

	runTest(test, "x + y * z", map[string]interface{}{
		"y": 3.0,
		"z": 7.0,
	}, "[x] + 21")

	runTest(test, "x + y * z", map[string]interface{}{
		"x": 0.0,
		"y": 1.0,
	}, "[z]")

	runTest(test, "x - y - x", map[string]interface{}{
		"x": 0.0,
	}, "-[y]")

	runTest(test, "x != 0 && y == x - 1", map[string]interface{}{
		"x": 10.0,
	}, "[y] = 9")

	runTest(test, "x != 0 && y == x - 1", map[string]interface{}{
		"x": 0.0,
	}, "0")

	runTest(test, "x > 0 && y > 0 && z > 0", map[string]interface{}{
		"x": 1.0,
	}, "[y] > 0 AND [z] > 0")

	runTest(test, "x > 0 && y > 0 && z > 0", map[string]interface{}{
		"y": 1.0,
	}, "[x] > 0 AND [z] > 0")

	runTest(test, "x > 0 && y > 0 && z > 0", map[string]interface{}{
		"z": 1.0,
	}, "[x] > 0 AND [y] > 0")

	runTest(test, "x > 0 && y > 0 && z > 0", map[string]interface{}{
		"y": 0.0,
	}, "0")

	runTest(test, "x > 0 || y > 0 && z > 0", map[string]interface{}{
		"x": 1.0,
	}, "1")

	runTest(test, "x > 0 || y > 0 && z > 0", map[string]interface{}{
		"y": 1.0,
		"z": 1.0,
	}, "1")

	runTest(test, "!x && (z == y/2)", map[string]interface{}{
		"x": false,
		"y": 12.0,
	}, "[z] = 6")

	runTest(test, "!x && (z == y/2)", map[string]interface{}{
		"x": true,
	}, "0")

	runTest(test, "x && ((y - 1) == z)", map[string]interface{}{
		"x": true,
		"y": 12.0,
	}, "11 = [z]")

	runTest(test, "y > 4 && (y - 2) == z", map[string]interface{}{
		"y": 12.0,
	}, "10 = [z]")

	runTest(test, "y > 4 && (y - 2) == z", map[string]interface{}{
		"y": 4.0,
	}, "0")

	runTest(test, "x ? (y > 0.15 && y < 0.5) : (z < -0.15 && z > -0.5)", map[string]interface{}{
		"x": true,
	}, "[y] > 0.15 AND [y] < 0.5")

	runTest(test, "x ? (y > 0.15 && y < 0.5) : (z < -0.15 && z > -0.5)", map[string]interface{}{
		"x": false,
	}, "[z] < -0.15 AND [z] > -0.5")
}

func runTest(test *testing.T, input string, parameters map[string]interface{}, expectedOutput string) {
	expr, err := NewEvaluableExpression(input)
	if err != nil {
		test.Errorf("error while creating expression: '%s'; error=%v", input, err)
		return
	}

	reduced, err := expr.Reduce(parameters)
	if err != nil {
		test.Errorf("error while reducing expression: '%s'; error=%v", input, err)
		return
	}

	output, err := reduced.ToSQLQuery()
	if err != nil {
		test.Errorf("error while printing expression: '%s'; error=%v", input, err)
		return
	}

	if output != expectedOutput {
		test.Errorf("Error while testing input: '%v'", input)
		test.Errorf("Expected: '%v'", expectedOutput)
		test.Errorf("Actual:   '%v'", output)
		return
	}
}
