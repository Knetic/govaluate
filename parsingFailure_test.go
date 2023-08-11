/*
The MIT License (MIT)

Copyright (c) 2014-2016 George Lester

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
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

		ParsingFailureTest{

			Name:     "Invalid equality comparator",
			Input:    "1 = 1",
			Expected: INVALID_TOKEN_KIND,
		},
		ParsingFailureTest{

			Name:     "Invalid equality comparator",
			Input:    "1 === 1",
			Expected: INVALID_TOKEN_KIND,
		},
		ParsingFailureTest{

			Name:     "Too many characters for logical operator",
			Input:    "true &&& false",
			Expected: INVALID_TOKEN_KIND,
		},
		ParsingFailureTest{

			Name:     "Too many characters for logical operator",
			Input:    "true ||| false",
			Expected: INVALID_TOKEN_KIND,
		},
		ParsingFailureTest{

			Name:     "Premature end to expression, via modifier",
			Input:    "10 > 5 +",
			Expected: UNEXPECTED_END,
		},
		ParsingFailureTest{

			Name:     "Premature end to expression, via comparator",
			Input:    "10 + 5 >",
			Expected: UNEXPECTED_END,
		},
		ParsingFailureTest{

			Name:     "Premature end to expression, via logical operator",
			Input:    "10 > 5 &&",
			Expected: UNEXPECTED_END,
		},
		ParsingFailureTest{

			Name:     "Premature end to expression, via ternary operator",
			Input:    "true ?",
			Expected: UNEXPECTED_END,
		},
		ParsingFailureTest{

			Name:     "Hanging REQ",
			Input:    "'wat' =~",
			Expected: UNEXPECTED_END,
		},
		ParsingFailureTest{

			Name:     "Invalid operator change to REQ",
			Input:    " / =~",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Invalid starting token, comparator",
			Input:    "> 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Invalid starting token, modifier",
			Input:    "+ 5",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Invalid starting token, logical operator",
			Input:    "&& 5 < 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Invalid NUMERIC transition",
			Input:    "10 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Invalid STRING transition",
			Input:    "'foo' 'foo'",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Invalid operator transition",
			Input:    "10 > < 10",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Starting with unbalanced parens",
			Input:    " ) ( arg2",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{

			Name:     "Unclosed bracket",
			Input:    "[foo bar",
			Expected: UNCLOSED_BRACKETS,
		},
		ParsingFailureTest{

			Name:     "Unclosed quote",
			Input:    "foo == 'responseTime",
			Expected: UNCLOSED_QUOTES,
		},
		ParsingFailureTest{

			Name:     "Constant regex pattern fail to compile",
			Input:    "foo =~ '[abc'",
			Expected: string(syntax.ErrMissingBracket),
		},
		ParsingFailureTest{

			Name:     "Unbalanced parenthesis",
			Input:    "10 > (1 + 50",
			Expected: UNBALANCED_PARENTHESIS,
		},
		ParsingFailureTest{

			Name:     "Multiple radix",
			Input:    "127.0.0.1",
			Expected: INVALID_NUMERIC,
		},
		ParsingFailureTest{

			Name:     "Undefined function",
			Input:    "foobar()",
			Expected: UNDEFINED_FUNCTION,
		},
		ParsingFailureTest{

			Name:     "Hanging accessor",
			Input:    "foo.Bar.",
			Expected: HANGING_ACCESSOR,
		},
		ParsingFailureTest{

			// this is expected to change once there are structtags in place that allow aliasing of fields
			Name:     "Unexported parameter access",
			Input:    "foo.bar",
			Expected: UNEXPORTED_ACCESSOR,
		},
		ParsingFailureTest{
			Name:     "Incomplete Hex",
			Input:    "0x",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{
			Name:     "Invalid Hex literal",
			Input:    "0x > 0",
			Expected: INVALID_HEX,
		},
		ParsingFailureTest{
			Name:     "Hex float (Unsupported)",
			Input:    "0x1.1",
			Expected: INVALID_TOKEN_TRANSITION,
		},
		ParsingFailureTest{
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
