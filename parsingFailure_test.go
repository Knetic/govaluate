package govaluate

import (
	"testing"
	"strings"
)

const (

	UNEXPECTED_END string = "Unexpected end of expression"
	INVALID_TOKEN_TRANSITION = "Cannot transition token types"
	INVALID_TOKEN_KIND = "Invalid token"
)

/*
	Represents a test for parsing failures
*/
type ParsingFailureTest struct {

	Name string
	Input string
	Expected string
}

func TestParsingFailure(test *testing.T) {

	parsingTests := []ParsingFailureTest {

		ParsingFailureTest {

			Name: "Invalid equality comparator",
			Input: "1 = 1",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Invalid equality comparator",
			Input: "1 === 1",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Half of a logical operator",
			Input: "true & false",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Half of a logical operator",
			Input: "true | false",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Too many characters for logical operator",
			Input: "true &&& false",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Too many characters for logical operator",
			Input: "true ||| false",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Premature end to expression, via modifier",
			Input: "10 > 5 +",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Premature end to expression, via comparator",
			Input: "10 + 5 >",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Premature end to expression, via logical operator",
			Input: "10 > 5 &&",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Invalid starting token, comparator",
			Input: "> 10",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Invalid starting token, modifier",
			Input: "+ 5",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Invalid starting token, logical operator",
			Input: "&& 5 < 10",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Non-equality comparison (>) with string",
			Input: "100 > 'foo'",
			Expected: "",
		},
		ParsingFailureTest {

			Name: "Non-equality comparison (<) with string",
			Input: "100 < 'foo'",
			Expected: "",
		},
	}

	runParsingFailureTests(parsingTests, test)
}

func runParsingFailureTests(parsingTests []ParsingFailureTest, test *testing.T) {

	var err error;

	for _, testCase := range parsingTests {

		_, err = NewEvaluableExpression(testCase.Input);

		if(err == nil) {
			
			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected a parsing error, found no error.")
			test.Fail()
			continue;
		}

		if(!strings.Contains(err.Error(), testCase.Expected)) {
			
			
			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Got error: '%s', expected '%s'", testCase.Expected, err.Error())
			test.Fail()
			continue;
		} 
	}
}
