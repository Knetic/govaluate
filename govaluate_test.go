package govaluate

import (
	"testing"
)

/*
	Represents a test of parsing all tokens correctly from a string
*/
type TokenParsingTest struct {

	Name string
	Input string
	Expected []ExpressionToken
}

/*
	Represents a test of expression evaluation
*/
type EvaluationTest struct {

	Name string
	Input string
	Parameters []EvaluationParameter
	Expected interface{}
}

type EvaluationParameter struct {

	Name string
	Value interface{}
}

/*
	Represents a test for parsing failures
*/
type ParsingFailureTest struct {

	Name string
	Input string
	Expected string
}

func TestConstantParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest {

		TokenParsingTest {

			Name: "Single numeric",
			Input: "1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
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

			Name: "Single large numeric",
			Input: "1234567890",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1234567890.0,
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
						Value: 1.0,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "==",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric NEQ",
			Input: "1 != 2",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "!=",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 2.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric GT",
			Input: "1 > 0",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: ">",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 0.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric LT",
			Input: "1 < 2",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "<",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 2.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric GTE",
			Input: "1 >= 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: ">=",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric LTE",
			Input: "1 <= 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: COMPARATOR,
						Value: "<=",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
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
						Value: 1.0,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "+",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric MINUS",
			Input: "1 - 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "-",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric MULTIPLY",
			Input: "1 * 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "*",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
			},
		},
		TokenParsingTest {

			Name: "Numeric DIVIDE",
			Input: "1 / 1",
			Expected: []ExpressionToken {
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
					},
					ExpressionToken {
						Kind: MODIFIER,
						Value: "/",
					},
					ExpressionToken {
						Kind: NUMERIC,
						Value: 1.0,
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
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Single MINUS",
			Input: "100 - 51",
			Expected: 49.0,
		},
		EvaluationTest {

			Name: "Single MULTIPLY",
			Input: "5 * 20",
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Single DIVIDE",
			Input: "100 / 20",
			Expected: 5.0,
		},
		EvaluationTest {

			Name: "Compound PLUS",
			Input: "20 + 30 + 50",
			Expected: 100.0,
		},
		EvaluationTest {

			Name: "Mutiple operators",
			Input: "20 * 5 - 49",
			Expected: 51.0,
		},
		EvaluationTest {

			Name: "Paren usage",
			Input: "100 - (5 * 10)",
			Expected: 50.0,
		},
		EvaluationTest {

			Name: "Nested parentheses",
			Input: "50 + (5 * (5 - 3))",
			Expected: 100.0,
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

	evaluationTests := []EvaluationTest {

		EvaluationTest {

			Name: "Single parameter modified by constant",
			Input: "foo + 2",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 2.0,
				},
			},
			Expected: 4.0,
		},
		EvaluationTest {

			Name: "Single parameter modified by constant",
			Input: "foo * bar",
			Parameters: []EvaluationParameter {

				EvaluationParameter {
					Name: "foo",
					Value: 5.0,
				},
				EvaluationParameter {
					Name: "bar",
					Value: 2.0,
				},
			},
			Expected: 10.0,
		},
	}

	runEvaluationTests(evaluationTests, test)
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

func runTokenParsingTest(tokenParsingTests []TokenParsingTest, test *testing.T) {

	var expression *EvaluableExpression
	var actualTokens []ExpressionToken;
	var actualToken ExpressionToken
	var expectedTokenLength, actualTokenLength int
	var err error

	// Run the test cases.
	for _, parsingTest := range tokenParsingTests {

		expression, err = NewEvaluableExpression(parsingTest.Input)

		if(err != nil) {

			test.Log("Test '",parsingTest.Name,"' failed to parse: ", err)
			test.Fail()
			continue
		}

		actualTokens = expression.Tokens();

		expectedTokenLength = len(parsingTest.Expected);
		actualTokenLength = len(actualTokens);

		if(actualTokenLength != expectedTokenLength) {

			test.Log("Test '",parsingTest.Name,"' failed:")
			test.Log("Expected ", expectedTokenLength, " tokens, actually found '", actualTokenLength, "'")
			test.Fail()
			continue
		}

		for i, expectedToken := range parsingTest.Expected {

			actualToken = actualTokens[i]
			if(actualToken.Kind != expectedToken.Kind) {

				test.Log("Test '", parsingTest.Name, "' failed:")
				test.Log("Expected token kind '", expectedToken.Kind, "' does not match '", actualToken.Kind, "'")
				test.Fail()
				continue
			}

			if(actualToken.Value != expectedToken.Value) {

				test.Log("Test '", parsingTest.Name, "' failed:")
				test.Log("Expected token value '", expectedToken.Value, "' does not match '", actualToken.Value, "'")
				test.Fail()
				continue
			}
		}
	}
}

func runEvaluationTests(evaluationTests []EvaluationTest, test *testing.T) {

	var expression *EvaluableExpression
	var result interface{}
	var parameters map[string]interface{}
	var err error

	// Run the test cases.
	for _, evaluationTest := range evaluationTests {

		expression, err = NewEvaluableExpression(evaluationTest.Input)

		if(err != nil) {

			test.Log("Test '",evaluationTest.Name,"' failed to parse: ", err)
			test.Fail()
			continue
		}

		parameters = make(map[string]interface{}, 8)		
		
		for _, parameter := range evaluationTest.Parameters {
			parameters[parameter.Name] = parameter.Value
		}

		result, err = expression.Evaluate(parameters)

		if(err != nil) {

			test.Log("Test '", evaluationTest.Name, "' failed:")
			test.Log("Encountered error: " + err.Error())
			test.Fail()
			continue;
		}

		if(result != evaluationTest.Expected) {

			test.Log("Test '", evaluationTest.Name, "' failed:")
			test.Log("Expected evaluation result '", evaluationTest.Expected, "' does not match '", result, "'")
			test.Fail()
		}
	}
}

func runParsingFailureTests(parsingTests []ParsingFailureTest, test *testing.T) {

	
}
