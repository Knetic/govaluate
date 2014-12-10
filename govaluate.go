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
type OperatorSymbol int

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
	EQ	OperatorSymbol = 1
	NEQ	OperatorSymbol = iota
	GT
	LT
	GTE
	LTE
	AND
	OR
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULUS
)

// map of all valid symbols
var COMPARATOR_SYMBOLS = map[string]OperatorSymbol {

	"==": EQ,
	"!=": NEQ,
	">": GT,
	">=": GTE,
	"<": LT,
	"<=": LTE,
};

var LOGICAL_SYMBOLS = map[string]OperatorSymbol {

	"&&": AND,
	"||": OR,
};

var MODIFIER_SYMBOLS = map[string]OperatorSymbol {
	
	"+": PLUS,
	"-": MINUS,
	"*": MULTIPLY,
	"/": DIVIDE,
	"%": MODULUS,
};

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
// TODO: make this an array, instead of named states
// TODO: then iterate through all valid states to find state for kind.

var VALID_LEXER_STATES = []lexerState {

	lexerState {

		kind: CLAUSE,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			CLAUSE,
		},
	},

	lexerState {

		kind: NUMERIC,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: BOOLEAN,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: STRING,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: VARIABLE,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: MODIFIER,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			VARIABLE,
		},
	},
	lexerState {

		kind: COMPARATOR,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
		},
	},
	lexerState {

		kind: LOGICALOP,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
		},
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

	state = VALID_LEXER_STATES[0];
	stream = newLexerStream(expression);

	for ;; {

		token, err, found = readToken(stream);

		if(err != nil) {
			return ret, err	
		}

		if(!found) {
			break;
		}

		if(!state.canTransitionTo(token.Kind)) {

			return ret, errors.New("Cannot transition token types") // TODO: make this more descriptive.
		}

		// append this valid token, find new lexer state.		
		ret = append(ret, token)
		
		for _, possibleState := range VALID_LEXER_STATES {
			
			if(possibleState.kind == token.Kind) {
				
				state = possibleState
				break;
			}
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

			tokenValue = readUntilFalse(stream, false, unicode.IsLetter);
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

			tokenString = readUntilFalse(stream, false, isNumeric);
			tokenValue, _ = strconv.ParseFloat(tokenString, 64)
			kind = NUMERIC;
			break;
		}

		if(!isNotSingleQuote(character)) {
			tokenValue = readUntilFalse(stream, true, isNotSingleQuote);
			kind = STRING;

			// advance the stream one position, since reading until false assumes the terminator is a real token
			stream.rewind(-1)
			break;
		}

		// must be a known symbol
		stream.rewind(1);
		tokenString = readUntilFalse(stream, false, isNotAlphanumeric);
		stream.rewind(1);

		tokenValue = tokenString

		if(MODIFIER_SYMBOLS[tokenString] != 0) {

			kind = MODIFIER;
			break;
		}

		if(LOGICAL_SYMBOLS[tokenString] != 0) {

			kind = LOGICALOP;
			break;
		}

		if(COMPARATOR_SYMBOLS[tokenString] != 0) {

			kind = COMPARATOR;
			break;
		}

		kind = UNKNOWN
		stream.rewind(-1);
	}

	ret.Kind = kind;
	ret.Value = tokenValue;

	return ret, nil, (kind != UNKNOWN);
}

func readUntilFalse(stream *lexerStream, includeWhitespace bool, condition func(rune)(bool)) string {

	//TODO: eliminate the "includewhitespace", build a separate function which handles whitespace and ends with single quotes
	//TODO: then remove all the "rewind" cruft above.
	var tokenBuffer bytes.Buffer
	var character rune

	for(stream.canRead()) {

		character = stream.readCharacter()

		if(!includeWhitespace && unicode.IsSpace(character)) {
			continue;
		}

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
func isNotSingleQuote(character rune) bool {
	return character != '\''
}
func isNotAlphanumeric(character rune) bool {
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
