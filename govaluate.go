package govaluation

import (
	"errors"
	"bytes"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type EvaluableExpression struct {

	Tokens []ExpressionToken
	inputExpression string
}

/*
	Represents a single parsed token.
*/
type ExpressionToken struct {

	Kind TokenKind
	Value interface{}
	Evaluate func(left interface{}, right interface{}) interface{}
}

/*
	Represents all valid types of tokens that a token can be.
*/
type TokenKind int
type ComparatorToken string
type LogicalOperatorToken string
type ModifierToken string

const (

	UNKNOWN TokenKind = iota

	NUMERIC 
	BOOLEAN
	STRING
	VARIABLE

	COMPARATOR
	LOGICALOP
	MODIFIER

	CLAUSE
)

const (

	EQ  ComparatorToken = "==" 
	NEQ ComparatorToken = "!="
	GT  ComparatorToken = ">"
	LT  ComparatorToken = "<"
	GTE ComparatorToken = ">="
	LTE ComparatorToken = "<="
)

const (

	AND LogicalOperatorToken = "&&"
	OR  LogicalOperatorToken = "||"
)

const (

	PLUS ModifierToken 	= "+"
	MINUS ModifierToken 	= "-"
	MULTIPLY ModifierToken 	= "*"
	DIVIDE ModifierToken 	= "/"
	MODULUS ModifierToken 	= "%"
)

type lexerState struct {

	isEOF bool
	kind TokenKind
	validNextKinds []TokenKind
}

type lexerStream struct {
	source string
	position int
	length int
}

// lexer states.
// Constant for all purposes except compiler.

var CLAUSESTATE lexerState = lexerState {

	kind: CLAUSE,
	isEOF: false,
	validNextKinds: []TokenKind {

		NUMERIC,
		BOOLEAN,
		VARIABLE,
		STRING,
		CLAUSE,
	},
}

var NUMERICSTATE lexerState = lexerState {

	kind: NUMERIC,
	isEOF: true,
	validNextKinds: []TokenKind {

		MODIFIER,
		COMPARATOR,
		LOGICALOP,
	},
}

var STRINGSTATE lexerState = lexerState {

	kind: STRING,
	isEOF: true,
	validNextKinds: []TokenKind {

		MODIFIER,
		COMPARATOR,
		LOGICALOP,
	},
}

var VARIABLESTATE lexerState = lexerState {

	kind: VARIABLE,
	isEOF: true,
	validNextKinds: []TokenKind {

		MODIFIER,
		COMPARATOR,
		LOGICALOP,
	},
}

var MODIFIERSTATE lexerState = lexerState {

	kind: MODIFIER,
	isEOF: false,
	validNextKinds: []TokenKind {

		NUMERIC,
		VARIABLE,
	},
}

var COMPARATORSTATE lexerState = lexerState {

	kind: COMPARATOR,
	isEOF: false,
	validNextKinds: []TokenKind {

		NUMERIC,
		BOOLEAN,
		VARIABLE,
		STRING,
	},
}

var LOGICALOPSTATE lexerState = lexerState {

	kind: LOGICALOP,
	isEOF: false,
	validNextKinds: []TokenKind {

		NUMERIC,
		BOOLEAN,
		VARIABLE,
		STRING,
	},
}

func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	var ret *EvaluableExpression;
	var err error

	ret = new(EvaluableExpression)
	ret.inputExpression = expression;
	ret.Tokens, err = parseTokens(expression)

	if(err != nil) {
		return nil, err
	}
	return ret, nil
}

func parseTokens(expression string) ([]ExpressionToken, error) {

	var ret []ExpressionToken
	var token ExpressionToken
	var state lexerState
	var stream *lexerStream
	var err error
	var found bool

	state = CLAUSESTATE;
	stream = newLexerStream(expression);

	for ;; {

		token, err, found = readToken(stream);

		if(err != nil) {
			return ret, err	
		}

		if(!found) {
			break;
		}

		if(state.canTransitionTo(token.Kind)) {

			ret = append(ret, token)
		} else {
			return ret, errors.New("Cannot transition token types") // TODO: make this more descriptive.
		}
	}

	return ret, nil
}

func readToken(stream *lexerStream) (ExpressionToken, error, bool) {

	var ret ExpressionToken
	var tokenValue interface{}
	var tokenString string
	var kind TokenKind
	var character rune

	kind = UNKNOWN;

	// numeric is 0-9, or .
	// string starts with '
	// variable is alphanumeric, always starts with a letter
	// symbols are anything non-alphanumeric
	// all others read into a buffer until they reach the end of the stream
	for(stream.canRead()) {

		character = stream.readCharacter()

		if(unicode.IsSpace(character)) {
			continue
		}

		// variable
		if(unicode.IsLetter(character)) {

			stream.rewind(1)

			tokenValue = readUntilFalse(stream, unicode.IsLetter);
			kind = VARIABLE;

			if(tokenValue == "true") {

				kind = BOOLEAN
				tokenValue = true
			} else {

				if(tokenValue == "false") {

					kind = BOOLEAN	
					tokenValue = false
				}
			}
			break;
		}

		// numeric constant
		if(isNumeric(character)) {

			stream.rewind(1)

			tokenString = readUntilFalse(stream, isNumeric);
			tokenValue, _ = strconv.ParseFloat(tokenString, 64)
			kind = NUMERIC;
			break;
		}

		if(isSingleQuote(character)) {
			tokenValue = readUntilFalse(stream, isSingleQuote);
			kind = STRING;
			break;
		}

		// must be a known symbol
		
	}

	ret.Kind = kind;
	ret.Value = tokenValue;

	return ret, nil, (kind != UNKNOWN);
}

func readUntilFalse(stream *lexerStream, condition func(rune)(bool)) string {

	var tokenBuffer bytes.Buffer
	var character rune

	for(stream.canRead()) {

		character = stream.readCharacter()

		if(condition(character)) {
			tokenBuffer.WriteString(string(character));
		} else {
			stream.rewind(1)
			break;
		}
	}

	return tokenBuffer.String();
}

func isNumeric(character rune) bool {
	return unicode.IsDigit(character) || character == '.'
}
func isSingleQuote(character rune) bool {
	return character == '\''
}
func isNotAlphaNumeric(character rune) bool {
	return !(unicode.IsDigit(character) || unicode.IsLetter(character))
}

func newLexerStream(source string) *lexerStream {

	var ret *lexerStream

	ret = new(lexerStream)
	ret.source = source
	ret.length = len(source)
	return ret
}

func (this *lexerStream) readCharacter() rune {

	var character rune

	character, _ = utf8.DecodeRuneInString(this.source[this.position:])
	this.position += 1
	return character
}

func (this *lexerStream) rewind(amount int) {
	this.position -= amount
}

func (this lexerStream) canRead() bool {
	return this.position < this.length;
}

func (this lexerState) canTransitionTo(kind TokenKind) bool {

	for _, validKind := range this.validNextKinds {

		if(validKind == kind) {
			return true
		}
	}

	return false
}

func (this EvaluableExpression) Evaluate(parameters map[string]interface{}) interface{} {

	return false
}

func (this EvaluableExpression) String() string {

	return this.inputExpression;
}
