package govaluate

import (
	"fmt"
	"strings"
)

// Grammar:
// expr    = ternary | binary | indexer ;
// ternary = indexer, "?", expr, ":", expr ;
// binary  = indexer, operator, expr
//         | indexer, ident, expr ;
// indexer = value, { "[", expr, "]" } ;
// value   = literal | call | boolean | ident | "(", expr, ")" | array | prefix ;
// call    = ident, "(", args, ")" ;
// array   = "[", args, "]" ;
// args    = [ expr, { ",", expr }, [ "," ] ] ;
// prefix  = operator, expr ;
// boolean = "true" | "false" ;

// Parse converts expression string to an AST, which can be evaluated.
func Parse(input string) (ExprNode, error) {
	s := NewTokenStream(input)
	expr, err := parseExpr(s, 0)
	if err == nil && !s.Peek().Is(TokenKindEOF, nil) {
		return ExprNode{}, unexpectedToken(s.Peek(), "operator")
	}
	if tokenizerErr := s.Error(); tokenizerErr != nil {
		return ExprNode{}, tokenizerErr
	}
	return expr, err
}

// MustParse returns an AST or panics if string cannot be parsed.
func MustParse(input string) ExprNode {
	expr, err := Parse(input)
	if err != nil {
		panic(fmt.Errorf("MustParse error: %v", err))
	}
	return expr
}

// TryParse returns an ExprNodeLiteral if the string cannot be parsed
func TryParse(input string) ExprNode {
	expr, err := Parse(input)
	if err != nil {
		return NewExprNodeLiteral(input, 0, len(input))
	}
	return expr
}

func parseExpr(s *TokenStream, minPrecedence int) (ExprNode, error) {
	lhs, err := parseIndexer(s)
	if err != nil {
		return lhs, err
	}
	return parseExprInner(s, lhs, minPrecedence)
}

func parseExprInner(s *TokenStream, lhs ExprNode, minPrecedence int) (ExprNode, error) {
	operator, precedence, ok := peekOperator(s)
	for ok && precedence >= minPrecedence {
		s.Next()
		if operator == "?" {
			return parseTernaryIf(s, lhs)
		}
		rhs, err := parseIndexer(s)
		if err != nil {
			return ExprNode{}, err
		}
		innerOperator, innerPrecedence, innerOk := peekOperator(s)
		for innerOk && innerPrecedence > precedence {
			rhs, err = parseExprInner(s, rhs, innerPrecedence)
			if err != nil {
				return ExprNode{}, err
			}
			innerOperator, innerPrecedence, innerOk = peekOperator(s)
		}
		pos, len := lhs.SourcePos, rhs.SourcePos+rhs.SourceLen-lhs.SourcePos
		lhs = NewExprNodeOperator(operator, []ExprNode{lhs, rhs}, pos, len, OperatorTypeInfix)
		operator, precedence, ok = innerOperator, innerPrecedence, innerOk
	}
	return lhs, nil
}

func parseIndexer(s *TokenStream) (ExprNode, error) {
	value, err := parseValue(s)
	if err != nil {
		return ExprNode{}, err
	}
	res := value
	for s.Peek().Is(TokenKindBracket, '[') {
		s.Next()
		index, err := parseExpr(s, 0)
		if err != nil {
			return ExprNode{}, err
		}
		bracket, err := consumeBracket(s, ']')
		if err != nil {
			return ExprNode{}, err
		}
		pos, len := value.SourcePos, bracket.SourcePos+bracket.SourceLen-value.SourcePos
		res = NewExprNodeOperator("[]", []ExprNode{res, index}, pos, len, OperatorTypeIndexer)
	}
	return res, nil
}

func parseValue(s *TokenStream) (ExprNode, error) {
	token := s.Next()

	switch token.Kind {
	case TokenKindNumber, TokenKindString:
		return NewExprNodeLiteral(token.Value, token.SourcePos, token.SourceLen), nil

	case TokenKindIdentifier:
		// function call
		if s.Peek().Is(TokenKindBracket, '(') {
			return parseCall(s, token)
		}

		// boolean literal
		switch token.Value {
		case "true":
			return NewExprNodeLiteral(true, token.SourcePos, token.SourceLen), nil
		case "false":
			return NewExprNodeLiteral(false, token.SourcePos, token.SourceLen), nil
		}

		// variable
		return NewExprNodeVariable(token.Value.(string), token.SourcePos, token.SourceLen), nil

	case TokenKindBracket:
		switch token.Value {
		case '(':
			// expression in brackets
			expr, err := parseExpr(s, 0)
			if err != nil {
				return ExprNode{}, err
			}
			bracket, err := consumeBracket(s, ')')
			if err != nil {
				return ExprNode{}, err
			}
			expr.SourcePos = token.SourcePos
			expr.SourceLen = bracket.SourcePos + bracket.SourceLen - token.SourcePos
			return expr, nil

		case '[':
			// array
			items, err := parseArgs(s, ']')
			if err != nil {
				return ExprNode{}, err
			}
			bracket, err := consumeBracket(s, ']')
			if err != nil {
				return ExprNode{}, err
			}
			pos, len := token.SourcePos, bracket.SourcePos+bracket.SourceLen-token.SourcePos
			return NewExprNodeOperator("array", items, pos, len, OperatorTypeArray), nil
		}

	case TokenKindOperator:
		// prefix operator
		// consume all operators with higher precedence
		expr, err := parseExpr(s, defaultPrecedence(token.Value.(string), 1))
		if err != nil {
			return ExprNode{}, err
		}
		// then apply prefix operator
		pos, len := token.SourcePos, expr.SourcePos+expr.SourceLen-token.SourcePos
		return NewExprNodeOperator(token.Value.(string), []ExprNode{expr}, pos, len, OperatorTypePrefix), nil
	}

	return ExprNode{}, unexpectedToken(token, "value")
}

func parseCall(s *TokenStream, nameToken ExprToken) (ExprNode, error) {
	if _, err := consumeBracket(s, '('); err != nil {
		return ExprNode{}, err
	}
	args, err := parseArgs(s, ')')
	if err != nil {
		return ExprNode{}, err
	}
	bracket, err := consumeBracket(s, ')')
	if err != nil {
		return ExprNode{}, err
	}
	name := nameToken.Value.(string)
	pos, len := nameToken.SourcePos, bracket.SourcePos+bracket.SourceLen-nameToken.SourcePos
	return NewExprNodeOperator(name, args, pos, len, OperatorTypeCall), nil
}

func parseArgs(s *TokenStream, until rune) ([]ExprNode, error) {
	args := []ExprNode{}
	for !s.Peek().Is(TokenKindBracket, until) {
		arg, err := parseExpr(s, defaultPrecedence(",", 2)+1)
		if err != nil {
			return args, err
		}
		args = append(args, arg)
		if s.Peek().Is(TokenKindOperator, ",") {
			s.Next()
		} else if !s.Peek().Is(TokenKindBracket, until) {
			return args, unexpectedToken(s.Peek(), "'"+string(until)+"'", "','")
		}
	}
	return args, nil
}

func parseTernaryIf(s *TokenStream, condition ExprNode) (ExprNode, error) {
	precedence := defaultPrecedence("?:", 3)
	valueIfTrue, err := parseExpr(s, precedence+1)
	if err != nil {
		return ExprNode{}, err
	}
	if !s.Peek().Is(TokenKindOperator, ":") {
		return ExprNode{}, unexpectedToken(s.Peek(), "':'")
	}
	s.Next()
	valueIfFalse, err := parseExpr(s, precedence)
	if err != nil {
		return ExprNode{}, err
	}
	args := []ExprNode{condition, valueIfTrue, valueIfFalse}
	pos, len := condition.SourcePos, valueIfFalse.SourcePos+valueIfFalse.SourceLen-condition.SourcePos
	return NewExprNodeOperator("?:", args, pos, len, OperatorTypeTernary), nil
}

func peekOperator(s *TokenStream) (string, int, bool) {
	if token := s.Peek(); token.Kind == TokenKindOperator || token.Kind == TokenKindIdentifier {
		name := token.Value.(string)
		return name, defaultPrecedence(name, 2), true
	}
	return "", 0, false
}

func consumeBracket(s *TokenStream, bracket rune) (ExprToken, error) {
	token := s.Next()
	if token.Kind != TokenKindBracket {
		return token, unexpectedToken(token, "'"+string(bracket)+"'")
	}
	if token.Value != bracket {
		return token, fmt.Errorf("unmatched bracket: '%v', expecting '%v', pos: %d", string(token.Value.(rune)), string(bracket), token.SourcePos)
	}
	return token, nil
}

func unexpectedToken(token ExprToken, expected ...string) error {
	if token.Is(TokenKindEOF, nil) {
		return fmt.Errorf("unexpected eof, expecting %s", strings.Join(expected, ", "))
	}
	return fmt.Errorf("unexpected token %v, expecting %s, pos: %d", token, strings.Join(expected, ", "), token.SourcePos)
}
