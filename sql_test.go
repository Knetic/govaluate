package govaluate

import (
	"testing"
)

/*
	Represents a test of correctly creating a SQL query string from an expression.
*/
type QueryTest struct {
	Name     string
	Input    string
	Expected string
}

func TestSQLSerialization(test *testing.T) {

	testCases := []QueryTest{

		{

			Name:     "Single GT",
			Input:    "1 > 0",
			Expected: "1 > 0",
		},
		{

			Name:     "Single LT",
			Input:    "0 < 1",
			Expected: "0 < 1",
		},
		{

			Name:     "Single GTE",
			Input:    "1 >= 0",
			Expected: "1 >= 0",
		},
		{

			Name:     "Single LTE",
			Input:    "0 <= 1",
			Expected: "0 <= 1",
		},
		{

			Name:     "Single EQ",
			Input:    "1 == 0",
			Expected: "1 = 0",
		},
		{

			Name:     "Single NEQ",
			Input:    "1 != 0",
			Expected: "1 <> 0",
		},

		{

			Name:     "Parameter names",
			Input:    "foo == bar",
			Expected: "[foo] = [bar]",
		},
		{

			Name:     "Strings",
			Input:    "'foo'",
			Expected: "'foo'",
		},
		{

			Name:     "Date format",
			Input:    "'2014-07-04T00:00:00Z'",
			Expected: "'2014-07-04T00:00:00Z'",
		},
		{

			Name:     "Single PLUS",
			Input:    "10 + 10",
			Expected: "10 + 10",
		},
		{

			Name:     "Single MINUS",
			Input:    "10 - 10",
			Expected: "10 - 10",
		},
		{

			Name:     "Single MULTIPLY",
			Input:    "10 * 10",
			Expected: "10 * 10",
		},
		{

			Name:     "Single DIVIDE",
			Input:    "10 / 10",
			Expected: "10 / 10",
		},
		{

			Name:     "Single true bool",
			Input:    "true",
			Expected: "1",
		},
		{

			Name:     "Single false bool",
			Input:    "false",
			Expected: "0",
		},
		{

			Name:     "Single AND",
			Input:    "true && true",
			Expected: "1 AND 1",
		},
		{

			Name:     "Single OR",
			Input:    "true || true",
			Expected: "1 OR 1",
		},
		{

			Name:     "Clauses",
			Input:    "10 + (foo + bar)",
			Expected: "10 + ( [foo] + [bar] )",
		},
		{

			Name:     "Negate prefix",
			Input:    "foo < -1",
			Expected: "[foo] < -1",
		},
		{

			Name:     "Invert prefix",
			Input:    "!(foo > 1)",
			Expected: "NOT ( [foo] > 1 )",
		},
		{

			Name:     "Exponent",
			Input:    "1 ** 2",
			Expected: "POW(1, 2)",
		},
		{

			Name:     "Modulus",
			Input:    "10 % 2",
			Expected: "MOD(10, 2)",
		},
		{

			Name:     "Membership operator",
			Input:    "foo IN (1, 2, 3)",
			Expected: "[foo] in ( 1 , 2 , 3 )",
		},
		{

			Name:     "Null coalescence",
			Input:    "foo ?? bar",
			Expected: "COALESCE([foo], [bar])",
		},
		/*
			// Ternaries don't work yet, because the outputter is not yet sophisticated enough to produce them.
			QueryTest{

				Name:     "Full ternary",
				Input:    "[foo] == 5 ? 1 : 2",
				Expected: "IF([foo] = 5, 1, 2)",
			},
			QueryTest{

				Name:     "Half ternary",
				Input:    "[foo] == 5 ? 1",
				Expected: "IF([foo] = 5, 1)",
			},
			QueryTest{

				Name:     "Full ternary with implicit bool",
				Input:    "[foo] ? 1 : 2",
				Expected: "IF([foo] = 0, 1, 2)",
			},
			QueryTest{

				Name:     "Half ternary with implicit bool",
				Input:    "[foo] ? 1",
				Expected: "IF([foo] = 0, 1)",
			},*/
		{

			Name:     "Regex equals",
			Input:    "'foo' =~ '[fF][oO]+'",
			Expected: "'foo' RLIKE '[fF][oO]+'",
		},
		{

			Name:     "Regex not-equals",
			Input:    "'foo' !~ '[fF][oO]+'",
			Expected: "'foo' NOT RLIKE '[fF][oO]+'",
		},
	}

	runQueryTests(testCases, test)
}

func runQueryTests(testCases []QueryTest, test *testing.T) {

	var expression *EvaluableExpression
	var actualQuery string
	var err error

	test.Logf("Running %d SQL translation test cases", len(testCases))

	// Run the test cases.
	for _, testCase := range testCases {

		expression, err = NewEvaluableExpression(testCase.Input)

		if err != nil {

			test.Logf("Test '%s' failed to parse: %s", testCase.Name, err)
			test.Logf("Expression: '%s'", testCase.Input)
			test.Fail()
			continue
		}

		actualQuery, err = expression.ToSQLQuery()

		if err != nil {

			test.Logf("Test '%s' failed to create query: %s", testCase.Name, err)
			test.Logf("Expression: '%s'", testCase.Input)
			test.Fail()
			continue
		}

		if actualQuery != testCase.Expected {

			test.Logf("Test '%s' did not create expected query.", testCase.Name)
			test.Logf("Actual: '%s', expected '%s'", actualQuery, testCase.Expected)
			test.Fail()
			continue
		}
	}
}
