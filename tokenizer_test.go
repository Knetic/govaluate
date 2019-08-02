package govaluate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeWhitespace(t *testing.T) {
	valid := []string{" ", "\n", "\r\n", "\t \r\n"}
	values := []interface{}{nil, nil, nil, nil}
	suffix := []string{"", "+ 1", ".square()", ")"}
	testTokenizerSuccess(t, tokenizeWhitespace, TokenKindWhitespace, valid, suffix, values)
	testTokenizerFail(t, tokenizeWhitespace, suffix)
}

func TestTokenizeIdentifier(t *testing.T) {
	valid := []string{"x", "data", "is_valid", "TestIdent", "var123"}
	values := []interface{}{"x", "data", "is_valid", "TestIdent", "var123"}
	suffix := []string{"", " + 1", "+1", ", ", ".square()", ")"}
	testTokenizerSuccess(t, tokenizeIdentifier, TokenKindIdentifier, valid, suffix, values)
	testTokenizerFail(t, tokenizeIdentifier, append(suffix, "123var"))
}

func TestTokenizeHexNumber(t *testing.T) {
	valid := []string{"0x0011", "0x1", "0xc0ffee", "0xABCDEF"}
	values := []interface{}{float64(0x11), 1.0, float64(0xc0ffee), float64(0xabcdef)}
	suffix := []string{"", "g", "+1", ".square()"}
	testTokenizerSuccess(t, tokenizeHexNumber, TokenKindNumber, valid, suffix, values)
	testTokenizerFail(t, tokenizeHexNumber, append(suffix, "0123", "abcd", "0x"))
}

func TestTokenizeDecimalNumber(t *testing.T) {
	valid := []string{"1", "15", "123.24", ".7", "17.43"}
	values := []interface{}{1.0, 15.0, 123.24, 0.7, 17.43}
	suffix := []string{"", "abc", ".", " + 1", ".square()"}
	testTokenizerSuccess(t, tokenizeDecimalNumber, TokenKindNumber, valid, suffix, values)
	testTokenizerFail(t, tokenizeDecimalNumber, suffix)
}

func TestTokenizeString(t *testing.T) {
	valid := []string{"\"\"", "\"hello\"", "\"hello, \\\"quoted\\\"\"", "''", "'\"hey\"'"}
	values := []interface{}{"", "hello", "hello, \"quoted\"", "", "\"hey\""}
	suffix := []string{"", "abc", ".", " + 1", ")", "["}
	testTokenizerSuccess(t, tokenizeString, TokenKindString, valid, suffix, values)
	testTokenizerFail(t, tokenizeString, append(suffix, "\"unterminated", "'"))
}

func TestTokenizeOperator(t *testing.T) {
	valid := []string{"+", "-", "<=", "**", "|>", "&&", "||"}
	values := []interface{}{"+", "-", "<=", "**", "|>", "&&", "||"}
	suffix := []string{"", "abc", "7", "\"str\"", " x", ")", "("}
	testTokenizerSuccess(t, tokenizeOperator, TokenKindOperator, valid, suffix, values)
	testTokenizerFail(t, tokenizeOperator, suffix)
}

func TestTokenizeBracket(t *testing.T) {
	valid := []string{"(", ")", "[", "]", "{", "}"}
	values := []interface{}{'(', ')', '[', ']', '{', '}'}
	invalid := []string{"", "abc", "7", "\"str\"", " x"}
	suffix := append(invalid, ")", "(")
	testTokenizerSuccess(t, tokenizeBracket, TokenKindBracket, valid, suffix, values)
	testTokenizerFail(t, tokenizeBracket, invalid)
}

func testTokenizerSuccess(t *testing.T, tokenizer Tokenizer, kind ExprTokenKind, prefix, suffix []string, expectedValues []interface{}) {
	for idx, a := range prefix {
		for _, b := range suffix {
			token := tokenizer(a + b)
			assert.Equal(t, NewExprToken(kind, expectedValues[idx], len(a)), token, "prefix='%s' suffix='%s'", a, b)
		}
	}
}

func testTokenizerFail(t *testing.T, tokenizer Tokenizer, input []string) {
	for _, a := range input {
		token := tokenizer(a)
		assert.Equal(t, 0, token.SourceLen, "input='%s'", a)
	}
}

func testTokenizer(t *testing.T, tokenizer Tokenizer, kind ExprTokenKind, valid, invalid []string, values []interface{}) {

	for idx, a := range valid {
		for _, b := range invalid {
			token := tokenizer(a + b)
			assert.Equal(t, NewExprToken(kind, values[idx], len(a)), token, "valid='%s' invalid='%s'", a, b)
		}
	}

	for _, b := range invalid {
		token := tokenizer(b)
		assert.Equal(t, 0, token.SourceLen, "invalid='%s'", b)
	}
}

func TestTokenize(t *testing.T) {
	tokens, err := Tokenize("x + 1")
	expected := []ExprToken{
		ExprToken{TokenKindIdentifier, "x", 1, 0},
		ExprToken{TokenKindOperator, "+", 1, 2},
		ExprToken{TokenKindNumber, 1.0, 1, 4},
	}
	assert.Equal(t, expected, tokens)
	assert.Nil(t, err)

	tokens, err = Tokenize("data.x < -17.0 && (x || y)")
	expected = []ExprToken{
		ExprToken{TokenKindIdentifier, "data", 4, 0},
		ExprToken{TokenKindOperator, ".", 1, 4},
		ExprToken{TokenKindIdentifier, "x", 1, 5},
		ExprToken{TokenKindOperator, "<", 1, 7},
		ExprToken{TokenKindOperator, "-", 1, 9},
		ExprToken{TokenKindNumber, 17.0, 4, 10},
		ExprToken{TokenKindOperator, "&&", 2, 15},
		ExprToken{TokenKindBracket, '(', 1, 18},
		ExprToken{TokenKindIdentifier, "x", 1, 19},
		ExprToken{TokenKindOperator, "||", 2, 21},
		ExprToken{TokenKindIdentifier, "y", 1, 24},
		ExprToken{TokenKindBracket, ')', 1, 25},
	}
	assert.Equal(t, expected, tokens)
	assert.Nil(t, err)

	tokens, err = Tokenize("x in \"str")
	assert.EqualError(t, err, "unable to parse input at pos=5")
}
