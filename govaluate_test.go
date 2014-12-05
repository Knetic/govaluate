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

type EvaluationTest struct {

	Name string
	Input string
	Parameters map[string]interface{}
	Expected interface{}
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

	evaluationTests := []EvaluationTest {

		EvaluationTest {

			Name: "Single PLUS",
			Input: "51 + 49",
			Expected: 100,
		},
		EvaluationTest {

			Name: "Single MINUS",
			Input: "100 - 51",
			Expected: 49,
		},
		EvaluationTest {

			Name: "Single MULTIPLY",
			Input: "5 * 20",
			Expected: 100,
		},
		EvaluationTest {

			Name: "Single DIVIDE",
			Input: "100 / 20",
			Expected: 5,
		},
		EvaluationTest {

			Name: "Single MODULUS",
			Input: "100 % 2",
			Expected: 0,
		},
		EvaluationTest {

			Name: "Compound PLUS",
			Input: "20 + 30 + 50",
			Expected: 100,
		},
		EvaluationTest {

			Name: "Mutiple operators",
			Input: "20 * 5 - 49",
			Expected: 51,
		},
		EvaluationTest {

			Name: "Paren usage",
			Input: "100 - (5 * 10)",
			Expected: 50,
		},
		EvaluationTest {

			Name: "Nested parentheses",
			Input: "50 + (5 * (5 - 3))",
			Expected: 100,
		},
		EvaluationTest {

			Name: "Implicit boolean",
			Input: "2 > 1",
			Expected: true,
		},
		EvaluationTest {

			Name: "Compound boolean",
			Input: "5 < 10 && 1 < 5",
			Expected: true,
		},
		EvaluationTest {

			Name: "Parenthesis boolean",
			Input: "10 < 50 && (1 != 2 && 1 == 1)",
			Expected: true,
		},
		EvaluationTest {

			Name: "Comparison of string constants",
			Input: "'foo' == 'foo'",
			Expected: true,
		},
		EvaluationTest {

			Name: "NEQ comparison of string constants",
			Input: "'foo' != 'bar'",
			Expected: true,
		},
	}

	runEvaluationTests(evaluationTests, test)
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

func runEvaluationTests(evaluationTests []EvaluationTest, test *testing.T) {

	var expression *EvaluableExpression
	var result interface{}

	// Run the test cases.
	for _, evaluationTest := range evaluationTests {

		expression = NewEvaluableExpression(evaluationTest.Input)
		result = expression.Evaluate(evaluationTest.Parameters)

		if(result != evaluationTest.Expected) {

			test.Log("Test '", evaluationTest.Name, "' failed:")
			test.Log("Expected evaluation result '", evaluationTest.Expected, "' does not match '", result, "'")
			test.Fail()
		}
	}
}
