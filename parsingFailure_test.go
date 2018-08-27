package govaluate

import (
	"fmt"
	"regexp/syntax"
	"strings"
	"testing"
)

const (
	UNEXPECTED_END           string = "Unexpected end of expression"
	INVALID_TOKEN_TRANSITION        = "Cannot transition token types"
	INVALID_TOKEN_KIND              = "Invalid token"
	UNCLOSED_QUOTES                 = "Unclosed string literal"
	UNCLOSED_BRACKETS               = "Unclosed parameter bracket"
	UNBALANCED_PARENTHESIS          = "Unbalanced parenthesis"
	INVALID_NUMERIC                 = "Unable to parse numeric value"
	UNDEFINED_FUNCTION              = "Undefined function"
	HANGING_ACCESSOR                = "Hanging accessor on token"
	UNEXPORTED_ACCESSOR             = "Unable to access unexported"
	INVALID_HEX                     = "Unable to parse hex value"
)

/*
	Represents a test for parsing failures
*/
type ParsingFailureTest struct {
	Name     string
	Input    string
	Expected string
}

func TestParsingFailure(test *testing.T) {

	parsingTests := []ParsingFailureTest{

		{

			Name:     "Invalid equality comparator",
			Input:    "1 = 1",
			Expected: INVALID_TOKEN_KIND,
		},
		{

			Name:     "Invalid equality comparator",
			Input:    "1 === 1",
			Expected: INVALID_TOKEN_KIND,
		},
		{

			Name:     "Too many characters for logical operator",
			Input:    "true &&& false",
			Expected: INVALID_TOKEN_KIND,
		},
		{

			Name:     "Too many characters for logical operator",
			Input:    "true ||| false",
			Expected: INVALID_TOKEN_KIND,
		},
		{

			Name:     "Premature end to expression, via modifier",
			Input:    "10 > 5 +",
			Expected: UNEXPECTED_END,
		},
		{

			Name:     "Premature end to expression, via comparator",
			Input:    "10 + 5 >",
			Expected: UNEXPECTED_END,
		},
		{

			Name:     "Premature end to expression, via logical operator",
			Input:    "10 > 5 &&",
			Expected: UNEXPECTED_END,
		},
		{

			Name:     "Premature end to expression, via ternary operator",
			Input:    "true ?",
			Expected: UNEXPECTED_END,
		},
		{

			Name:     "Hanging REQ",
			Input:    "'wat' =~",
			Expected: UNEXPECTED_END,
		},
		{

			Name:     "Invalid operator change to REQ",
			Input:    " / =~",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Invalid starting token, comparator",
			Input:    "> 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Invalid starting token, modifier",
			Input:    "+ 5",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Invalid starting token, logical operator",
			Input:    "&& 5 < 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Invalid NUMERIC transition",
			Input:    "10 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Invalid STRING transition",
			Input:    "'foo' 'foo'",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Invalid operator transition",
			Input:    "10 > < 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Starting with unbalanced parens",
			Input:    " ) ( arg2",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{

			Name:     "Unclosed bracket",
			Input:    "[foo bar",
			Expected: UNCLOSED_BRACKETS,
		},
		{

			Name:     "Unclosed quote",
			Input:    "foo == 'responseTime",
			Expected: UNCLOSED_QUOTES,
		},
		{

			Name:     "Constant regex pattern fail to compile",
			Input:    "foo =~ '[abc'",
			Expected: string(syntax.ErrMissingBracket),
		},
		{

			Name:     "Unbalanced parenthesis",
			Input:    "10 > (1 + 50",
			Expected: UNBALANCED_PARENTHESIS,
		},
		{

			Name:     "Multiple radix",
			Input:    "127.0.0.1",
			Expected: INVALID_NUMERIC,
		},
		{

			Name:     "Undefined function",
			Input:    "foobar()",
			Expected: UNDEFINED_FUNCTION,
		},
		{

			Name:     "Hanging accessor",
			Input:    "foo.Bar.",
			Expected: HANGING_ACCESSOR,
		},
		{

			// this is expected to change once there are structtags in place that allow aliasing of fields
			Name:     "Unexported parameter access",
			Input:    "foo.bar",
			Expected: UNEXPORTED_ACCESSOR,
		},
		{
			Name:     "Incomplete Hex",
			Input:    "0x",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{
			Name:     "Invalid Hex literal",
			Input:    "0x > 0",
			Expected: INVALID_HEX,
		},
		{
			Name:     "Hex float (Unsupported)",
			Input:    "0x1.1",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		{
			Name:     "Hex invalid letter",
			Input:    "0x12g1",
			Expected: INVALID_TOKEN_TRANSITION,
		},
	}

	runParsingFailureTests(parsingTests, test)
}

func runParsingFailureTests(parsingTests []ParsingFailureTest, test *testing.T) {

	var err error

	fmt.Printf("Running %d parsing test cases...\n", len(parsingTests))

	for _, testCase := range parsingTests {

		_, err = NewEvaluableExpression(testCase.Input)

		if err == nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected a parsing error, found no error.")
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
