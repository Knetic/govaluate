package govaluate

import (

	"errors"
)

type EvaluableExpression struct {

	Tokens []ExpressionToken
	inputExpression string
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

func (this EvaluableExpression) Evaluate(parameters map[string]interface{}) (interface{}, error) {

	var ret interface{};
	var err error;

	ret, _, err = this.evaluateClause(nil, this.Tokens, parameters);

	return ret, err;
}

func (this EvaluableExpression) evaluateClause(leftValue interface{}, tokens []ExpressionToken, parameters map[string]interface{}) (interface{}, int, error) {

	var token ExpressionToken;
	var nextOperation func(interface{}, interface{})(interface{});
	var variableName, errorMessage string;
	var parameterValue interface{};
	var tokensLength int;

	tokensLength = len(tokens);

	for i := 0; i < tokensLength; i++ {

		token = tokens[i];

		if(nextOperation == nil) {

			switch(token.Kind) {
				case NUMERIC 		:	fallthrough
				case STRING		:	fallthrough
				case BOOLEAN		:	leftValue = token.Value
				case VARIABLE		:	variableName = token.Value.(string);
								leftValue = parameters[variableName];

								if(leftValue == nil) {
									errorMessage = "No parameter '"+ variableName +"' found."
									return nil, tokensLength, errors.New(errorMessage);
								}
				
				case COMPARATOR	:	
				case LOGICALOP	:	
				case MODIFIER	:	nextOperation = determineOperator(token.Value.(string), MODIFIER_SYMBOLS);
			}
		} else { 

			switch(token.Kind) {
				case NUMERIC 		:	fallthrough
				case STRING		:	fallthrough
				case BOOLEAN		:	parameterValue = token.Value;
				case VARIABLE		:	variableName = token.Value.(string);
								parameterValue = parameters[variableName];
						
								if(parameterValue == nil) {
									errorMessage = "No parameter '"+ variableName +"' found."
									return nil, tokensLength, errors.New(errorMessage);
								}
			}

			if(parameterValue == nil) {
				errorMessage = "No right value for expression";
				return nil, tokensLength, errors.New(errorMessage);
			}

			leftValue = nextOperation(leftValue, parameterValue);
			nextOperation = nil;
		}
	}
	return leftValue, tokensLength, nil;
}

func determineOperator(operator string, operators map[string]OperatorSymbol) func(interface{}, interface{})interface{} {

	var operatorKind OperatorSymbol;

	operatorKind = operators[operator];

	switch(operatorKind) {

		case	EQ	:	return comparatorEq;
		//case	NEQ	:	return comparatorNeq;
		case	GT	:	return comparatorGt;
		case	LT	:	return comparatorLt;
		case	GTE	:	return comparatorGte;
		case	LTE	:	return comparatorLte;

		case	AND	:	return logicalAnd;
		case	OR	:	return logicalOr;
		
		case	PLUS	:	return modifierPlus;
		case	MINUS	:	return modifierMinus;
		case	MULTIPLY:	return modifierMultiply;
		case	DIVIDE	:	return modifierDivide;
	}

	return comparatorEq;
}

func comparatorEq(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue == rightValue;
}
func comparatorNeq(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue != rightValue;
}
func comparatorGt(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) > rightValue.(float64);
}
func comparatorLt(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) < rightValue.(float64);
}
func comparatorGte(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) >= rightValue.(float64);
}
func comparatorLte(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) <= rightValue.(float64);
}
func logicalAnd(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(bool) && rightValue.(bool);
}
func logicalOr(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(bool) || rightValue.(bool);
}
func modifierPlus(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) + rightValue.(float64);
}
func modifierMinus(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) - rightValue.(float64);
}
func modifierMultiply(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) * rightValue.(float64);
}
func modifierDivide(leftValue interface{}, rightValue interface{}) interface{} {
	return leftValue.(float64) / rightValue.(float64);
}

func (this EvaluableExpression) String() string {

	return this.inputExpression;
}
