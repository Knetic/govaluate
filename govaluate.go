package govaluation

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
}

/*
	Represents all valid types of tokens that a token can be.
*/
type TokenKind int
type ComparatorToken string
type LogicalOperatorToken string
type ModifierToken string

const (

	NUMERIC TokenKind = iota
	BOOLEAN
	STRING
	COMPARATOR
	LOGICALOP
	MODIFIER
)

const (

	EQ  ComparatorToken = "=" 
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

func NewEvaluableExpression(expression string) *EvaluableExpression {

	var ret *EvaluableExpression;

	ret = new(EvaluableExpression)
	ret.inputExpression = expression;
	ret.Tokens = parseTokens(expression)

	return ret
}

func parseTokens(expression string) []ExpressionToken {

	var ret []ExpressionToken
	return ret
}

func (self EvaluableExpression) String() string {

	return self.inputExpression;
}
