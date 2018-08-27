package govaluate

import (
	"fmt"
	"strings"
	"testing"
)

const (
	EXPERR_NIL_VALUE string = "cannot have a nil value"
)

/*
	Contains a single test case for the EvaluableExpression.NewEvaluableExpressionFromTokens() method.

	These tests, and the ones in `tokenExpressionFailure_test` will be fairly incomplete.
	Creating an expression from a string and from tokens _must_ both perform the same syntax checks.
	So all the checks in `parsing_test` will follow the same logic as the ones here.

	These tests check some corner cases - such as tokens having nil values when they must have something.
	Cases that cannot occur through the normal parser, but may occur in other parsers.
*/
type ExpressionTokenSyntaxTest struct {
	Name     string
	Input    []ExpressionToken
	Expected string
}

func TestNilValues(test *testing.T) {

	cases := []ExpressionTokenSyntaxTest{
		{
			Name: "Nil numeric",
			Input: []ExpressionToken{
				{
					Kind: NUMERIC,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil string",
			Input: []ExpressionToken{
				{
					Kind: STRING,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil bool",
			Input: []ExpressionToken{
				{
					Kind: BOOLEAN,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil time",
			Input: []ExpressionToken{
				{
					Kind: TIME,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil pattern",
			Input: []ExpressionToken{
				{
					Kind: PATTERN,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil variable",
			Input: []ExpressionToken{
				{
					Kind: VARIABLE,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil prefix",
			Input: []ExpressionToken{
				{
					Kind: PREFIX,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil comparator",
			Input: []ExpressionToken{
				{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				{
					Kind: COMPARATOR,
				},
				{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil logicalop",
			Input: []ExpressionToken{
				{
					Kind:  BOOLEAN,
					Value: true,
				},
				{
					Kind: LOGICALOP,
				},
				{
					Kind:  BOOLEAN,
					Value: true,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil modifer",
			Input: []ExpressionToken{
				{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				{
					Kind: MODIFIER,
				},
				{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
		{
			Name: "Nil ternary",
			Input: []ExpressionToken{
				{
					Kind:  BOOLEAN,
					Value: true,
				},
				{
					Kind: TERNARY,
				},
				{
					Kind:  BOOLEAN,
					Value: true,
				},
			},
			Expected: EXPERR_NIL_VALUE,
		},
	}

	runExpressionFromTokenTests(cases, true, test)
}

func runExpressionFromTokenTests(cases []ExpressionTokenSyntaxTest, expectFail bool, test *testing.T) {

	var err error

	fmt.Printf("Running %d expression from expression token tests...\n", len(cases))

	for _, testCase := range cases {

		_, err = NewEvaluableExpressionFromTokens(testCase.Input)

		if err != nil {
			if expectFail {

				if !strings.Contains(err.Error(), testCase.Expected) {

					test.Logf("Test '%s' failed", testCase.Name)
					test.Logf("Got error: '%s', expected '%s'", err.Error(), testCase.Expected)
					test.Fail()
				}
				continue
			}

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Got error: '%s'", err)
			test.Fail()
			continue
		} else {
			if expectFail {

				test.Logf("Test '%s' failed", testCase.Name)
				test.Logf("Expected error, found none\n")
				test.Fail()
				continue
			}
		}
	}
}
