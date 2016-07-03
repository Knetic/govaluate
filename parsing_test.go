package govaluate

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	"unicode"
)

/*
	Represents a test of parsing all tokens correctly from a string
*/
type TokenParsingTest struct {
	Name     string
	Input    string
	Expected []ExpressionToken
}

func TestConstantParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Single numeric",
			Input: "1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Single two-digit numeric",
			Input: "50",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 50.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Single string",
			Input: "'foo'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  STRING,
					Value: "foo",
				},
			},
		},
		TokenParsingTest{

			Name:  "Single time, RFC3339, only date",
			Input: "'2014-01-02'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  TIME,
					Value: time.Date(2014, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		TokenParsingTest{

			Name:  "Single time, RFC3339, with hh:mm",
			Input: "'2014-01-02 14:12'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  TIME,
					Value: time.Date(2014, time.January, 2, 14, 12, 0, 0, time.UTC),
				},
			},
		},
		TokenParsingTest{

			Name:  "Single time, RFC3339, with hh:mm:ss",
			Input: "'2014-01-02 14:12:22'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  TIME,
					Value: time.Date(2014, time.January, 2, 14, 12, 22, 0, time.UTC),
				},
			},
		},
		TokenParsingTest{

			Name:  "Single boolean",
			Input: "true",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
			},
		},
		TokenParsingTest{

			Name:  "Single large numeric",
			Input: "1234567890",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1234567890.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Single floating-point",
			Input: "0.5",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 0.5,
				},
			},
		},
		TokenParsingTest{

			Name:  "Single large floating point",
			Input: "3.14567471",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 3.14567471,
				},
			},
		},
		TokenParsingTest{

			Name:  "Single false boolean",
			Input: "false",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: false,
				},
			},
		},
		TokenParsingTest{
			Name:  "Single internationalized string",
			Input: "'ÆŦǽഈᚥஇคٸ'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  STRING,
					Value: "ÆŦǽഈᚥஇคٸ",
				},
			},
		},
		TokenParsingTest{
			Name:  "Single internationalized parameter",
			Input: "ÆŦǽഈᚥஇคٸ",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "ÆŦǽഈᚥஇคٸ",
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestLogicalOperatorParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Boolean AND",
			Input: "true && false",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
				ExpressionToken{
					Kind:  LOGICALOP,
					Value: "&&",
				},
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: false,
				},
			},
		},
		TokenParsingTest{

			Name:  "Boolean OR",
			Input: "true || false",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
				ExpressionToken{
					Kind:  LOGICALOP,
					Value: "||",
				},
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: false,
				},
			},
		},
		TokenParsingTest{

			Name:  "Multiple logical operators",
			Input: "true || false && true",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
				ExpressionToken{
					Kind:  LOGICALOP,
					Value: "||",
				},
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: false,
				},
				ExpressionToken{
					Kind:  LOGICALOP,
					Value: "&&",
				},
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestComparatorParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Numeric EQ",
			Input: "1 == 2",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "==",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 2.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric NEQ",
			Input: "1 != 2",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "!=",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 2.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric GT",
			Input: "1 > 0",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: ">",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 0.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric LT",
			Input: "1 < 2",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "<",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 2.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric GTE",
			Input: "1 >= 2",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: ">=",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 2.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric LTE",
			Input: "1 <= 2",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "<=",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 2.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "String REQ",
			Input: "'foobar' =~ 'bar'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  STRING,
					Value: "foobar",
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "=~",
				},

				// it's not particularly clean to test for the contents of a pattern, (since it means modifying the harness below)
				// so pattern contents are left untested.
				ExpressionToken{
					Kind: PATTERN,
				},
			},
		},
		TokenParsingTest{

			Name:  "String NREQ",
			Input: "'foobar' !~ 'bar'",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  STRING,
					Value: "foobar",
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "!~",
				},
				ExpressionToken{
					Kind: PATTERN,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestModifierParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Numeric PLUS",
			Input: "1 + 1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  MODIFIER,
					Value: "+",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric MINUS",
			Input: "1 - 1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  MODIFIER,
					Value: "-",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric MULTIPLY",
			Input: "1 * 1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  MODIFIER,
					Value: "*",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric DIVIDE",
			Input: "1 / 1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  MODIFIER,
					Value: "/",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Numeric MODULUS",
			Input: "1 % 1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  MODIFIER,
					Value: "%",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestPrefixParsing(test *testing.T) {

	testCases := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Sign prefix",
			Input: "-1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  PREFIX,
					Value: "-",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Sign prefix on variable",
			Input: "-foo",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  PREFIX,
					Value: "-",
				},
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo",
				},
			},
		},
		TokenParsingTest{

			Name:  "Boolean prefix",
			Input: "!true",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  PREFIX,
					Value: "!",
				},
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
			},
		},
		TokenParsingTest{

			Name:  "Boolean prefix on variable",
			Input: "!foo",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  PREFIX,
					Value: "!",
				},
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo",
				},
			},
		},
	}

	testCases = combineWhitespaceExpressions(testCases)
	runTokenParsingTest(testCases, test)
}

func TestEscapedParameters(test *testing.T) {

	testCases := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Single escaped parameter",
			Input: "[foo]",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo",
				},
			},
		},
		TokenParsingTest{

			Name:  "Single escaped parameter with whitespace",
			Input: "[foo bar]",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo bar",
				},
			},
		},
		TokenParsingTest{

			Name:  "Single escaped parameter with escaped closing bracket",
			Input: "[foo[bar\\]]",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo[bar]",
				},
			},
		},
		TokenParsingTest{

			Name:  "Escaped parameters and unescaped parameters",
			Input: "[foo] > bar",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo",
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: ">",
				},
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "bar",
				},
			},
		},
		TokenParsingTest{

			Name:  "Unescaped parameter with space",
			Input: "foo\\ bar > bar",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo bar",
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: ">",
				},
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "bar",
				},
			},
		},
		TokenParsingTest{

			Name:  "Unescaped parameter with space",
			Input: "response\\-time > bar",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "response-time",
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: ">",
				},
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "bar",
				},
			},
		},
		TokenParsingTest{

			Name:  "Parameters with snake_case",
			Input: "foo_bar > baz_quux",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "foo_bar",
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: ">",
				},
				ExpressionToken{
					Kind:  VARIABLE,
					Value: "baz_quux",
				},
			},
		},
		TokenParsingTest{

			Name:  "String literal uses backslash to escape",
			Input: "\"foo\\'bar\"",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  STRING,
					Value: "foo'bar",
				},
			},
		},
	}

	runTokenParsingTest(testCases, test)
}

func TestTernaryParsing(test *testing.T) {
	tokenParsingTests := []TokenParsingTest{

		TokenParsingTest{

			Name:  "Ternary after Boolean",
			Input: "true ? 1",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
				ExpressionToken{
					Kind:  TERNARY,
					Value: "?",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
			},
		},
		TokenParsingTest{

			Name:  "Ternary after Comperator",
			Input: "1 == 0 ? true",
			Expected: []ExpressionToken{
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 1.0,
				},
				ExpressionToken{
					Kind:  COMPARATOR,
					Value: "==",
				},
				ExpressionToken{
					Kind:  NUMERIC,
					Value: 0.0,
				},
				ExpressionToken{
					Kind:  TERNARY,
					Value: "?",
				},
				ExpressionToken{
					Kind:  BOOLEAN,
					Value: true,
				},
			},
		},
	}

	runTokenParsingTest(tokenParsingTests, test)
}

/*
	Tests to make sure that the String() reprsentation of an expression exactly matches what is given to the parse function.
*/
func TestOriginalString(test *testing.T) {

	// include all the token types, to be sure there's no shenaniganery going on.
	expressionString := "2 > 1 &&" +
		"'something' != 'nothing' || " +
		"'2014-01-20' < 'Wed Jul  8 23:07:35 MDT 2015' && " +
		"[escapedVariable name with spaces] <= unescaped\\-variableName &&" +
		"modifierTest + 1000 / 2 > (80 * 100 % 2) && true ? true : false"

	expression, err := NewEvaluableExpression(expressionString)
	if(err != nil) {

		test.Logf("failed to parse original string test: %v", err)
		test.Fail()
	}

	if(expression.String() != expressionString) {
		test.Logf("String() did not give the same expression as given to parse")
		test.Fail()
	}
}

func combineWhitespaceExpressions(testCases []TokenParsingTest) []TokenParsingTest {

	var currentCase, strippedCase TokenParsingTest
	var caseLength int

	caseLength = len(testCases)

	for i := 0; i < caseLength; i++ {

		currentCase = testCases[i]

		strippedCase = TokenParsingTest{

			Name:     (currentCase.Name + " (without whitespace)"),
			Input:    stripUnquotedWhitespace(currentCase.Input),
			Expected: currentCase.Expected,
		}

		testCases = append(testCases, strippedCase, currentCase)
	}

	return testCases
}

func stripUnquotedWhitespace(expression string) string {

	var expressionBuffer bytes.Buffer
	var quoted bool

	for _, character := range expression {

		if !quoted && unicode.IsSpace(character) {
			continue
		}

		if character == '\'' {
			quoted = !quoted
		}

		expressionBuffer.WriteString(string(character))
	}

	return expressionBuffer.String()
}

func runTokenParsingTest(tokenParsingTests []TokenParsingTest, test *testing.T) {

	var expression *EvaluableExpression
	var actualTokens []ExpressionToken
	var actualToken ExpressionToken
	var expectedTokenKindString, actualTokenKindString string
	var expectedTokenLength, actualTokenLength int
	var err error

	fmt.Printf("Running %d parsing test cases...\n", len(tokenParsingTests))

	// Run the test cases.
	for _, parsingTest := range tokenParsingTests {

		expression, err = NewEvaluableExpression(parsingTest.Input)

		if err != nil {

			test.Logf("Test '%s' failed to parse: %s", parsingTest.Name, err)
			test.Logf("Expression: '%s'", parsingTest.Input)
			test.Fail()
			continue
		}

		actualTokens = expression.Tokens()

		expectedTokenLength = len(parsingTest.Expected)
		actualTokenLength = len(actualTokens)

		if actualTokenLength != expectedTokenLength {

			test.Logf("Test '%s' failed:", parsingTest.Name)
			test.Logf("Expected %d tokens, actually found %d", expectedTokenLength, actualTokenLength)
			test.Fail()
			continue
		}

		for i, expectedToken := range parsingTest.Expected {

			actualToken = actualTokens[i]
			if actualToken.Kind != expectedToken.Kind {

				actualTokenKindString = GetTokenKindString(actualToken.Kind)
				expectedTokenKindString = GetTokenKindString(expectedToken.Kind)

				test.Logf("Test '%s' failed:", parsingTest.Name)
				test.Logf("Expected token kind '%v' does not match '%v'", expectedTokenKindString, actualTokenKindString)
				test.Fail()
				continue
			}

			if expectedToken.Value != nil && actualToken.Value != expectedToken.Value {

				test.Logf("Test '%s' failed:", parsingTest.Name)
				test.Logf("Expected token value '%v' does not match '%v'", expectedToken.Value, actualToken.Value)
				test.Fail()
				continue
			}
		}
	}
}
