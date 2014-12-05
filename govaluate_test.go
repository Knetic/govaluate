package govaluation

import (
	"testing"
)

/*
	Represents the contents of a test of token parsing
*/
type TokenParsingTest struct {

	Name string
	Input string
	Expected []ExpressionToken
}

func TestConstantParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest {

		TokenParsingTest {

			Name: "Single numeric",
			Input: "1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Single string",
			Input: "'foo'",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: STRING,
						Value: "foo",
					},
			},
		},
		TokenParsingTest {

			Name: "Single boolean",
			Input: "true",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: BOOLEAN,
						Value: true,
					},
			},
		},
		TokenParsingTest {

			Name: "Single numeric",
			Input: "1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Single large numeric",
			Input: "1234567890",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1234567890,
					},
			},
		},
		TokenParsingTest {

			Name: "Single floating-point",
			Input: "0.5",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 0.5,
					},
			},
		},
		TokenParsingTest {

			Name: "Single large floating point",
			Input: "3.14567471",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 3.14567471,
					},
			},
		},
		TokenParsingTest {

			Name: "Single false boolean",
			Input: "false",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: BOOLEAN,
						Value: false,
					},
			},
		},
	}

	runTokenParsingTest(tokenParsingTests, test)
}

func TestLogicalOperatorParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest {

		TokenParsingTest {

			Name: "Boolean AND",
			Input: "true && false",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: BOOLEAN,
						Value: true,
					},
					ExpressionToken {
						Kind: LOGICALOP,
						Value: "&&",
					},
					ExpressionToken {
						Kind: BOOLEAN,
						Value: false,
					},
			},
		},
		TokenParsingTest {

			Name: "Boolean OR",
			Input: "true || false",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: BOOLEAN,
						Value: true,
					},
					ExpressionToken {
						Kind: LOGICALOP,
						Value: "||",
					},
					ExpressionToken {
						Kind: BOOLEAN,
						Value: false,
					},
			},
		},
		TokenParsingTest {

			Name: "Multiple logical operators",
			Input: "true || false && true",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: BOOLEAN,
						Value: true,
					},
					ExpressionToken {
						Kind: LOGICALOP,
						Value: "||",
					},
					ExpressionToken {
						Kind: BOOLEAN,
						Value: false,
					},
					ExpressionToken {
						Kind: LOGICALOP,
						Value: "&&",
					},
					ExpressionToken {
						Kind: BOOLEAN,
						Value: true,
					},
			},
		},
	}

	runTokenParsingTest(tokenParsingTests, test)
}

func TestComparatorParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest {

		TokenParsingTest {

			Name: "Numeric EQ",
			Input: "1 == 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "==",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric NEQ",
			Input: "1 != 2",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "!=",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 2,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric GT",
			Input: "1 > 0",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: ">",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric LT",
			Input: "1 < 2",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "<",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 2,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric GTE",
			Input: "1 >= 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: ">=",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric LTE",
			Input: "1 <= 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "<=",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
	}

	runTokenParsingTest(tokenParsingTests, test)
}

func TestModifierParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest {

		TokenParsingTest {

			Name: "Numeric PLUS",
			Input: "1 + 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "+",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric MINUS",
			Input: "1 - 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "-",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric MULTIPLY",
			Input: "1 * 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "*",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric DIVIDE",
			Input: "1 / 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "/",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric MODULUS",
			Input: "1 % 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "%",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1,
					},
			},
		},
	}

	runTokenParsingTest(tokenParsingTests, test)
}

func TestNoParameterEvaluation(test *testing.T) {
}

func TestParameterizedEvaluation(test *testing.T) {
}

func runTokenParsingTest(tokenParsingTests []TokenParsingTest, test *testing.T) {

	var expression *EvaluableExpression
	var expectedToken ExpressionToken
	var expectedTokenLength, actualTokenLength int

	// Run the test cases.
	for _, parsingTest := range tokenParsingTests {

		expression = NewEvaluableExpression(parsingTest.Input)

		expectedTokenLength = len(parsingTest.Expected)
		actualTokenLength = len(expression.Tokens)

		if(actualTokenLength != expectedTokenLength) {

			test.Log("Test '",parsingTest.Name,"' failed:")
			test.Log("Expected ", expectedTokenLength, " tokens, actually found '", actualTokenLength, "'")
			test.Fail()
		}

		for i, token := range expression.Tokens {

			expectedToken = parsingTest.Expected[i]
			if(token.Kind != expectedToken.Kind) {

				test.Log("Test '", parsingTest.Name, "' failed:")
				test.Log("Expected token kind '", expectedToken.Kind, "' does not match '", token.Kind, "'")
				test.Fail()
			}

			if(token.Value != expectedToken.Value) {

				test.Log("Test '", parsingTest.Name, "' failed:")
				test.Log("Expected token value '", expectedToken.Kind, "' does not match '", token.Kind, "'")
				test.Fail()
			}
		}
	}
}
