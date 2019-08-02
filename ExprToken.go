package govaluate

import "fmt"

type ExprToken struct {
	Kind      ExprTokenKind
	Value     interface{}
	SourceLen int
	SourcePos int
}

type ExprTokenKind int

const (
	TokenKindEOF ExprTokenKind = iota
	TokenKindWhitespace
	TokenKindIdentifier
	TokenKindNumber
	TokenKindString
	TokenKindOperator
	TokenKindBracket
)

func NewExprToken(kind ExprTokenKind, value interface{}, sourceLen int) ExprToken {
	return ExprToken{
		Kind:      kind,
		Value:     value,
		SourceLen: sourceLen,
	}
}

func (token ExprToken) Is(kind ExprTokenKind, value interface{}) bool {
	return token.Kind == kind && token.Value == value
}

func (token ExprToken) String() string {
	switch token.Kind {
	case TokenKindEOF:
		return "EOF"
	case TokenKindWhitespace:
		return "Whitespace{}"
	case TokenKindIdentifier:
		return fmt.Sprintf("Identifier{%v}", token.Value)
	case TokenKindNumber:
		return fmt.Sprintf("Number{%v}", token.Value)
	case TokenKindString:
		return fmt.Sprintf("String{%v}", token.Value)
	case TokenKindOperator:
		return fmt.Sprintf("Operator{%v}", token.Value)
	case TokenKindBracket:
		return fmt.Sprintf("Bracket{'%v'}", string(token.Value.(rune)))
	}
	return fmt.Sprintf("Unknown{%v, %v}", token.Kind, token.Value)
}
