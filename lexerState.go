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
	"errors"
	"fmt"
)

type lexerState struct {
	isEOF          bool
	isNullable     bool
	kind           TokenKind
	validNextKinds []TokenKind
}

// lexer states.
// Constant for all purposes except compiler.
var validLexerStates = []lexerState{

	lexerState{
		kind:       UNKNOWN,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			PATTERN,
			FUNCTION,
			ACCESSOR,
			STRING,
			TIME,
			CLAUSE,
		},
	},

	lexerState{

		kind:       CLAUSE,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			PATTERN,
			FUNCTION,
			ACCESSOR,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},

	lexerState{

		kind:       CLAUSE_CLOSE,
		isEOF:      true,
		isNullable: true,
		validNextKinds: []TokenKind{

			COMPARATOR,
			MODIFIER,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			PATTERN,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
			LOGICALOP,
			TERNARY,
			SEPARATOR,
		},
	},

	lexerState{

		kind:       NUMERIC,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       BOOLEAN,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       STRING,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       TIME,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       PATTERN,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       VARIABLE,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       MODIFIER,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			VARIABLE,
			FUNCTION,
			ACCESSOR,
			STRING,
			BOOLEAN,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState{

		kind:       COMPARATOR,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			ACCESSOR,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
			PATTERN,
		},
	},
	lexerState{

		kind:       LOGICALOP,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			ACCESSOR,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState{

		kind:       PREFIX,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			ACCESSOR,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},

	lexerState{

		kind:       TERNARY,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			STRING,
			TIME,
			VARIABLE,
			FUNCTION,
			ACCESSOR,
			CLAUSE,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       FUNCTION,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{
			CLAUSE,
		},
	},
	lexerState{

		kind:       ACCESSOR,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{
			CLAUSE,
			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:       SEPARATOR,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			STRING,
			TIME,
			VARIABLE,
			FUNCTION,
			ACCESSOR,
			CLAUSE,
		},
	},
}

func (this lexerState) canTransitionTo(kind TokenKind) bool {

	for _, validKind := range this.validNextKinds {

		if validKind == kind {
			return true
		}
	}

	return false
}

func checkExpressionSyntax(tokens []ExpressionToken) error {

	var state lexerState
	var lastToken ExpressionToken
	var err error

	state = validLexerStates[0]

	for _, token := range tokens {

		if !state.canTransitionTo(token.Kind) {

			// call out a specific error for tokens looking like they want to be functions.
			if lastToken.Kind == VARIABLE && token.Kind == CLAUSE {
				return errors.New("Undefined function " + lastToken.Value.(string))
			}

			firstStateName := fmt.Sprintf("%s [%v]", state.kind.String(), lastToken.Value)
			nextStateName := fmt.Sprintf("%s [%v]", token.Kind.String(), token.Value)

			return errors.New("Cannot transition token types from " + firstStateName + " to " + nextStateName)
		}

		state, err = getLexerStateForToken(token.Kind)
		if err != nil {
			return err
		}

		if !state.isNullable && token.Value == nil {

			errorMsg := fmt.Sprintf("Token kind '%v' cannot have a nil value", token.Kind.String())
			return errors.New(errorMsg)
		}

		lastToken = token
	}

	if !state.isEOF {
		return errors.New("Unexpected end of expression")
	}
	return nil
}

func getLexerStateForToken(kind TokenKind) (lexerState, error) {

	for _, possibleState := range validLexerStates {

		if possibleState.kind == kind {
			return possibleState, nil
		}
	}

	errorMsg := fmt.Sprintf("No lexer state found for token kind '%v'\n", kind.String())
	return validLexerStates[0], errors.New(errorMsg)
}
