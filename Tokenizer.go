package govaluate

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Tokenizer reads single token from the beginning of input.
// If input can not be parsed, EOF token is returned.
type Tokenizer func(string) ExprToken

// TokenStream produces a stream of tokens from input string, skipping whitespace.
// Stream ends with EOF token, whitespace is skipped.
type TokenStream struct {
	input     string
	pos       int
	hasNext   bool
	nextToken ExprToken
}

// Tokenize converts input string to a list of tokens.
func Tokenize(input string) ([]ExprToken, error) {
	tokens := []ExprToken{}
	s := NewTokenStream(input)
	for token := s.Next(); !token.Is(TokenKindEOF, nil); token = s.Next() {
		tokens = append(tokens, token)
	}
	return tokens, s.Error()
}

// NewTokenStream constructs a new TokenStream with given input.
func NewTokenStream(input string) *TokenStream {
	return &TokenStream{input: input}
}

// Peek returns the next token in stream, without consuming it.
func (s *TokenStream) Peek() ExprToken {
	return s.next(false)
}

// Next returns the next token in stream, consuming it, and advancing stream to the next token.
// When there are no more tokens in input, an EOF token is returned.
func (s *TokenStream) Next() ExprToken {
	return s.next(true)
}

// Error returns an error if there was an error reading input.
func (s *TokenStream) Error() error {
	if s.Peek().Is(TokenKindEOF, nil) && s.pos != len(s.input) {
		return fmt.Errorf("unable to parse input at pos=%d", s.pos)
	}
	return nil
}

func (s *TokenStream) next(advance bool) ExprToken {
	if !s.hasNext {
		s.nextToken = s.readNext()
		s.hasNext = true
	}
	token := s.nextToken
	if advance && !token.Is(TokenKindEOF, nil) {
		s.hasNext = false
	}
	return token
}

func (s *TokenStream) readNext() ExprToken {
	if s.pos == len(s.input) {
		return ExprToken{}
	}
	whitespace := tokenizeWhitespace(s.input[s.pos:])
	if whitespace.SourceLen > 0 {
		s.pos += whitespace.SourceLen
	}
	nextInput := s.input[s.pos:]
	var tokenizers = [...]Tokenizer{
		tokenizeIdentifier,
		tokenizeHexNumber,
		tokenizeDecimalNumber,
		tokenizeString,
		tokenizeBracket,
		tokenizeOperator,
	}
	for _, tokenizer := range tokenizers {
		token := tokenizer(nextInput)
		if token.SourceLen > 0 {
			token.SourcePos = s.pos
			s.pos += token.SourceLen
			return token
		}
	}
	return ExprToken{}
}

func tokenizeWhitespace(input string) ExprToken {
	nextNonWhitespace := len(input)
	for idx, ch := range input {
		if !unicode.IsSpace(ch) {
			nextNonWhitespace = idx
			break
		}
	}
	return NewExprToken(TokenKindWhitespace, nil, nextNonWhitespace)
}

func tokenizeIdentifier(input string) ExprToken {
	nextNonIdent := len(input)
	for idx, ch := range input {
		if !unicode.IsLetter(ch) && ch != '_' && (idx == 0 || !unicode.IsDigit(ch)) {
			nextNonIdent = idx
			break
		}
	}
	if nextNonIdent == 0 {
		return ExprToken{}
	}
	return NewExprToken(TokenKindIdentifier, input[:nextNonIdent], nextNonIdent)
}

func tokenizeHexNumber(input string) ExprToken {
	if len(input) < 3 || input[0] != '0' || input[1] != 'x' {
		return ExprToken{}
	}

	nextNonDigit := len(input)
	for idx, ch := range input[2:] {
		if !('0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F') {
			nextNonDigit = 2 + idx
			break
		}
	}
	if nextNonDigit < 3 {
		return ExprToken{}
	}

	value, err := strconv.ParseUint(input[2:nextNonDigit], 16, 64)
	if err != nil {
		return ExprToken{}
	}
	return NewExprToken(TokenKindNumber, float64(value), nextNonDigit)
}

func tokenizeDecimalNumber(input string) ExprToken {
	type parserState int
	const (
		initial parserState = iota
		integer
		maybeFloat
		float
	)
	nextNonDecimal := len(input)
	state := initial
outer:
	for idx, ch := range input {
		if ch == '.' {
			switch state {
			case initial, integer:
				// first point, maybe it's a float
				state = maybeFloat
			case maybeFloat, float:
				// second point, stop reading
				nextNonDecimal = idx
				break outer
			}
		} else if '0' <= ch && ch <= '9' {
			switch state {
			case initial:
				state = integer
			case maybeFloat:
				state = float
			}
		} else {
			nextNonDecimal = idx
			break outer
		}
	}
	if state == maybeFloat {
		// do not consume trailing point char
		nextNonDecimal--
	}
	if nextNonDecimal == 0 {
		return ExprToken{}
	}
	value, err := strconv.ParseFloat(input[:nextNonDecimal], 64)
	if err != nil {
		return ExprToken{}
	}
	return NewExprToken(TokenKindNumber, value, nextNonDecimal)
}

func tokenizeString(input string) ExprToken {
	// must start with a double or single quote
	if len(input) < 2 || input[0] != '"' && input[0] != '\'' {
		return ExprToken{}
	}
	var builder strings.Builder
	escape := false
	for idx, ch := range input[1:] {
		if escape {
			// second char of escape sequence
			switch ch {
			case 'r':
				builder.WriteRune('\r')
			case 'n':
				builder.WriteRune('\n')
			case 't':
				builder.WriteRune('\t')
			default:
				builder.WriteRune(ch)
			}
			escape = false
		} else if ch == '\\' {
			// first char of escape sequence
			escape = true
		} else if ch == rune(input[0]) {
			// closing quote
			return NewExprToken(TokenKindString, builder.String(), idx+2)
		} else {
			// read string until closing quote
			builder.WriteRune(ch)
		}
	}
	return ExprToken{}
}

func tokenizeBracket(input string) ExprToken {
	if len(input) == 0 {
		return ExprToken{}
	}
	switch input[0] {
	case '(', ')', '[', ']', '{', '}':
		return NewExprToken(TokenKindBracket, rune(input[0]), 1)
	}
	return ExprToken{}
}

func tokenizeOperator(input string) ExprToken {
	if len(input) == 0 {
		return ExprToken{}
	}

	// handle symbols that can not be combined
	switch input[0] {
	case ',':
		return NewExprToken(TokenKindOperator, input[:1], 1)
	}

	// read symbols as a single token, like ==, ||, &&
	operatorSymbols := []rune("~!#$%^&*-+|\\=:./?<>")
	nextNonOperator := len(input)
	for idx, ch := range input {
		found := false
		for _, sym := range operatorSymbols {
			if sym == ch {
				found = true
				break
			}
		}
		if !found {
			nextNonOperator = idx
			break
		}
	}
	if nextNonOperator == 0 {
		return ExprToken{}
	}
	return NewExprToken(TokenKindOperator, input[:nextNonOperator], nextNonOperator)
}
