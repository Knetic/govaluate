package govaluation

import (

	"fmt"
)

type EvaluableExpression struct {

	Tokens []*ExpressionToken
	inputExpression string
}

/*
	Represents a single parsed token.
*/
type ExpressionToken struct {

	Kind TokenKind
	Value string
}

/*
	Represents all valid types of tokens that a token can be.
*/
type TokenKind int {

	COMMENT iota
	NUMBER
	STRING
	LT
	GT
	LTE
	GTE
	EQ
	NEQ
	AND
	OR
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULUS
	EOF
}

/*
	Represents a lexer state, which can parse a given set of TokenKinds and transition to a set of other lexer states
*/
type lexerState interface {

	validKinds []TokenKind
	nextStates []lexerState
	evaluate (*func(*chan)(*lexerState))
}

func NewEvaluableExpression(expression string) *EvaluableExpression {

	var ret *EvaluableExpression;

	ret = new(EvaluableExpression)
	ret.inputExpression = expression;
	ret.Tokens = parseTokens(expression)

	return ret
}

func parseTokens(expression string) []*ExpressionToken {

	var ret []*ExpressionToken

	ret = new([]*ExpressionToken)

	return ret
}

func (self EvaluableExpression) String() string {

	return self.inputExpression;
}
