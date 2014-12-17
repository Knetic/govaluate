package govaluate

import (
	"testing"
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
	}

	runParsingFailureTests(parsingTests, test)
}

func runParsingFailureTests(parsingTests []ParsingFailureTest, test *testing.T) {

	
}
