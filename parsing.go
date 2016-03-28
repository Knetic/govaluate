package govaluate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func parseTokens(expression string) ([]ExpressionToken, error) {

	var ret []ExpressionToken
	var token, lastToken ExpressionToken
	var state lexerState
	var stream *lexerStream
	var err error
	var found bool

	state = validLexerStates[0]
	stream = newLexerStream(expression)

	for stream.canRead() {

		token, err, found = readToken(stream, state)

		if err != nil {
			return ret, err
		}

		if !found {
			break
		}

		if !state.canTransitionTo(token.Kind) {

			firstStateName := fmt.Sprintf("%s [%v]", GetTokenKindString(state.kind), lastToken.Value)
			nextStateName := fmt.Sprintf("%s [%v]", GetTokenKindString(token.Kind), token.Value)

			return ret, errors.New("Cannot transition token types from " + firstStateName + " to " + nextStateName)
		}

		// append this valid token, find new lexer state.
		ret = append(ret, token)

		for _, possibleState := range validLexerStates {

			if possibleState.kind == token.Kind {

				state = possibleState
				break
			}
		}

		lastToken = token
	}

	if !state.isEOF {
		return ret, errors.New("Unexpected end of expression")
	}

	return ret, nil
}

func readToken(stream *lexerStream, state lexerState) (ExpressionToken, error, bool) {

	var ret ExpressionToken
	var tokenValue interface{}
	var tokenTime time.Time
	var tokenString string
	var kind TokenKind
	var character rune
	var found bool
	var completed bool

	// numeric is 0-9, or .
	// string starts with '
	// variable is alphanumeric, always starts with a letter
	// bracket always means variable
	// symbols are anything non-alphanumeric
	// all others read into a buffer until they reach the end of the stream
	for stream.canRead() {

		character = stream.readCharacter()

		if unicode.IsSpace(character) {
			continue
		}

		kind = UNKNOWN

		// numeric constant
		if isNumeric(character) {

			tokenString = readTokenUntilFalse(stream, isNumeric)
			tokenValue, _ = strconv.ParseFloat(tokenString, 64)
			kind = NUMERIC
			break
		}

		// escaped variable
		if character == '[' {

			tokenValue, completed = readUntilFalse(stream, true, false, isNotClosingBracket)
			kind = VARIABLE

			if !completed {
				return ExpressionToken{}, errors.New("Unclosed parameter bracket"), false
			}

			mapVal := []interface{}{}
			err := json.Unmarshal([]byte(`[`+tokenValue.(string)+`]`), &mapVal)
			if err == nil {
				kind = ARRAY
				tokenValue = mapVal
			}

			// above method normally rewinds us to the closing bracket, which we want to skip.
			stream.rewind(-1)
			break
		}

		// regex pattern
		if character == '/' {
			if state.kind == COMPARATOR {

				tokenValue, completed = readUntilFalseNoEscape(stream, true, false, isNotClosingSlash)
				kind = PATTERN

				if !completed {
					return ExpressionToken{}, errors.New("Unclosed parameter /"), false
				}

				// above method normally rewinds us to the closing bracket, which we want to skip.
				stream.rewind(-1)
				break
			}
		}

		// regular variable
		if unicode.IsLetter(character) {

			tokenValue = readTokenUntilFalse(stream, isVariableName)
			kind = VARIABLE

			if tokenValue == "true" {

				kind = BOOLEAN
				tokenValue = true
			} else if tokenValue == "false" {

				kind = BOOLEAN
				tokenValue = false
			} else if true == isLogicalOp(tokenValue.(string)) {
				tokenValue = strings.ToUpper(tokenValue.(string))
				kind = LOGICALOP
			} else if true == isComparator(tokenValue.(string)) {
				tokenValue = strings.ToUpper(tokenValue.(string))
				kind = COMPARATOR
			}
			break
		}

		if !isNotQuote(character) {
			tokenValue, completed = readUntilFalse(stream, true, false, isNotQuote)

			if !completed {
				return ExpressionToken{}, errors.New("Unclosed string literal"), false
			}

			// advance the stream one position, since reading until false assumes the terminator is a real token
			stream.rewind(-1)

			// check to see if this can be parsed as a time.
			tokenTime, found = tryParseTime(tokenValue.(string))
			if found {
				kind = TIME
				tokenValue = tokenTime
			} else {
				kind = STRING
			}
			break
		}

		if character == '(' {
			tokenValue = character
			kind = CLAUSE
			break
		}

		if character == ')' {
			tokenValue = character
			kind = CLAUSE_CLOSE
			break
		}

		// must be a known symbol
		tokenString = readTokenUntilFalse(stream, isNotAlphanumeric)
		tokenValue = tokenString

		// quick hack for the case where "-" can mean "prefixed negation" or "minus", which are used
		// very differently.
		if state.canTransitionTo(PREFIX) {
			_, found = PREFIX_SYMBOLS[tokenString]
			if found {

				kind = PREFIX
				break
			}
		}
		_, found = MODIFIER_SYMBOLS[tokenString]
		if found {

			kind = MODIFIER
			break
		}

		_, found = LOGICAL_SYMBOLS[tokenString]
		if found {

			kind = LOGICALOP
			break
		}

		_, found = COMPARATOR_SYMBOLS[tokenString]
		if found {
			kind = COMPARATOR
			break
		}

		errorMessage := fmt.Sprintf("Invalid token: '%s'", tokenString)
		return ret, errors.New(errorMessage), false
	}

	ret.Kind = kind
	ret.Value = tokenValue

	return ret, nil, (kind != UNKNOWN)
}

func readTokenUntilFalse(stream *lexerStream, condition func(rune) bool) string {

	var ret string

	stream.rewind(1)
	ret, _ = readUntilFalse(stream, false, true, condition)
	return ret
}

/*
	Returns the string that was read until the given [condition] was false, or whitespace was broken.
	Returns false if the stream ended before whitespace was broken or condition was met.
*/
func readUntilFalse(stream *lexerStream, includeWhitespace bool, breakWhitespace bool, condition func(rune) bool) (string, bool) {

	var tokenBuffer bytes.Buffer
	var character rune
	var conditioned bool

	conditioned = false

	for stream.canRead() {

		character = stream.readCharacter()

		// Use backslashes to escape anything
		if character == '\\' {

			character = stream.readCharacter()
			tokenBuffer.WriteString(string(character))
			continue
		}

		if unicode.IsSpace(character) {

			if breakWhitespace && tokenBuffer.Len() > 0 {
				conditioned = true
				break
			}
			if !includeWhitespace {
				continue
			}
		}

		if condition(character) {
			tokenBuffer.WriteString(string(character))
		} else {
			conditioned = true
			stream.rewind(1)
			break
		}
	}

	return tokenBuffer.String(), conditioned
}

func readUntilFalseNoEscape(stream *lexerStream, includeWhitespace bool, breakWhitespace bool, condition func(rune) bool) (string, bool) {

	var tokenBuffer bytes.Buffer
	var character rune
	var conditioned bool

	conditioned = false

	for stream.canRead() {

		character = stream.readCharacter()

		// Use backslashes to escape anything
		// if character == '\\' {

		// 	character = stream.readCharacter()
		// 	tokenBuffer.WriteString(string(character))
		// 	continue
		// }

		if unicode.IsSpace(character) {

			if breakWhitespace && tokenBuffer.Len() > 0 {
				conditioned = true
				break
			}
			if !includeWhitespace {
				continue
			}
		}

		if condition(character) {
			tokenBuffer.WriteString(string(character))
		} else {
			conditioned = true
			stream.rewind(1)
			break
		}
	}

	return tokenBuffer.String(), conditioned
}

func isLogicalOp(txt string) bool {
	txt = strings.ToUpper(txt)
	if _, ok := LOGICAL_SYMBOLS[txt]; ok {
		return true
	}
	return false
}

func isComparator(txt string) bool {
	txt = strings.ToUpper(txt)
	if _, ok := COMPARATOR_SYMBOLS[txt]; ok {
		return true
	}
	return false
}

func isNumeric(character rune) bool {

	return unicode.IsDigit(character) || character == '.'
}

func isNotQuote(character rune) bool {

	return character != '\'' && character != '"'
}

func isNotAlphanumeric(character rune) bool {

	return !(unicode.IsDigit(character) ||
		unicode.IsLetter(character) ||
		character == '(' ||
		character == ')' ||
		!isNotQuote(character))
}

func isVariableName(character rune) bool {

	return unicode.IsLetter(character) ||
		unicode.IsDigit(character) ||
		character == '_'
}

func isNotClosingBracket(character rune) bool {

	return character != ']'
}

func isNotClosingSlash(character rune) bool {

	return character != '/'
}

/*
	Attempts to parse the [candidate] as a Time.
	Tries a series of standardized date formats, returns the Time if one applies,
	otherwise returns false through the second return.
*/
func tryParseTime(candidate string) (time.Time, bool) {

	var ret time.Time
	var found bool

	timeFormats := [...]string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.Kitchen,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",                         // RFC 3339
		"2006-01-02 15:04",                   // RFC 3339 with minutes
		"2006-01-02 15:04:05",                // RFC 3339 with seconds
		"2006-01-02 15:04:05-07:00",          // RFC 3339 with seconds and timezone
		"2006-01-02T15Z0700",                 // ISO8601 with hour
		"2006-01-02T15:04Z0700",              // ISO8601 with minutes
		"2006-01-02T15:04:05Z0700",           // ISO8601 with seconds
		"2006-01-02T15:04:05.999999999Z0700", // ISO8601 with nanoseconds
	}

	for _, format := range timeFormats {

		ret, found = tryParseExactTime(candidate, format)
		if found {
			return ret, true
		}
	}

	return time.Now(), false
}

func tryParseExactTime(candidate string, format string) (time.Time, bool) {

	var ret time.Time
	var err error

	ret, err = time.Parse(format, candidate)
	if err != nil {
		return time.Now(), false
	}

	return ret, true
}
