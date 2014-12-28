package govaluate

import (
	"testing"
)

/*
	Represents a test of correctly creating a SQL query string from an expression.
*/
type QueryTest struct {

	Name string
	Input string
	Expected string
}

func TestSQLSerialization(test *testing.T) {

	testCases := []QueryTest {

		QueryTest {

			Name: "Single GT",
			Input: "1 > 0",
			Expected: "1 > 0",
		},
		QueryTest {

			Name: "Single LT",
			Input: "0 < 1",
			Expected: "0 < 1",
		},
		QueryTest {

			Name: "Single GTE",
			Input: "1 >= 0",
			Expected: "1 >= 0",
		},
		QueryTest {

			Name: "Single LTE",
			Input: "0 <= 1",
			Expected: "0 <= 1",
		},
		QueryTest {

			Name: "Single EQ",
			Input: "1 == 0",
			Expected: "1 = 0",
		},
		QueryTest {

			Name: "Single NEQ",
			Input: "1 != 0",
			Expected: "1 <> 0",
		},

		QueryTest {

			Name: "Parameter names",
			Input: "foo == bar",
			Expected: "[foo] = [bar]",
		},
		QueryTest {

			Name: "Strings",
			Input: "'foo'",
			Expected: "'foo'",
		},
		QueryTest {

			Name: "Date format",
			Input: "'2014-07-04'",
			Expected: "'2014-07-04T00:00:00Z'",
		},
	}

	runQueryTests(testCases, test);
}

func runQueryTests(testCases []QueryTest, test *testing.T) {

	var expression *EvaluableExpression
	var actualQuery string;
	var err error;

	// Run the test cases.
	for _, testCase := range testCases {

		expression, err = NewEvaluableExpression(testCase.Input)

		if(err != nil) {

			test.Logf("Test '%s' failed to parse: %s", testCase.Name, err)
			test.Logf("Expression: '%s'", testCase.Input);
			test.Fail()
			continue
		}

		actualQuery, err = expression.ToSQLQuery();

		if(err != nil) {

			test.Logf("Test '%s' failed to create query: %s", testCase.Name, err)
			test.Logf("Expression: '%s'", testCase.Input);
			test.Fail()
			continue
		}

		if(actualQuery != testCase.Expected) {
			
			test.Logf("Test '%s' did not create expected query. Actual: '%s'", testCase.Name, actualQuery)
			test.Logf("Actual: '%s', expected '%s'", actualQuery, testCase.Expected);
			test.Fail()
			continue
		}
	}
}
