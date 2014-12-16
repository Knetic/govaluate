package govaluate

import (
	"errors"
)

type EvaluableExpression struct {

	tokens []ExpressionToken
	inputExpression string
}

func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	var ret *EvaluableExpression;
	var err error

	ret = new(EvaluableExpression)
	ret.inputExpression = expression;
	ret.tokens, err = parseTokens(expression)

	if(err != nil) {
		return nil, err
	}
	return ret, nil
}

func (this EvaluableExpression) Evaluate(parameters map[string]interface{}) (interface{}, error) {

	var stream *tokenStream;

	stream = newTokenStream(this.tokens);
	return evaluateTokens(stream, parameters);
}

func evaluateTokens(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	if(stream.hasNext()) {
		return evaluateLogical(stream, parameters);
	}
	return nil, nil;
}

func evaluateLogical(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateComparator(stream, parameters);	

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {

		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = LOGICAL_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break;
		}

		switch(symbol) {

			case OR		:	if(value != nil) {
							return evaluateLogical(stream, parameters);
						} else {
							value, err = evaluateComparator(stream, parameters);
						}
			case AND	:	if(value == nil) {
							return evaluateLogical(stream, parameters);
						} else {
							value, err = evaluateComparator(stream, parameters);
						}
		}

		if(err != nil) {
			return nil, err;
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateComparator(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, rightValue interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateAdditiveModifier(stream, parameters);

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {

		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = COMPARATOR_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break
		}

		rightValue, err = evaluateAdditiveModifier(stream, parameters);
		if(err != nil) {
			return nil, err;
		}

		switch(symbol) {

			case LT		:	value  = (value.(float64) < rightValue.(float64));
			case LTE	:	value  = (value.(float64) <= rightValue.(float64));
			case GT		:	value  = (value.(float64) > rightValue.(float64));
			case GTE	:	value  = (value.(float64) >= rightValue.(float64));
			case EQ		:	value  = (value == rightValue);
			case NEQ	:	value  = (value != rightValue);
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateAdditiveModifier(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, rightValue interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateMultiplicativeModifier(stream, parameters);

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {
		
		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break;
		}

		rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
		if(err != nil) {
			return nil, err;
		}

		switch(symbol) {

			case PLUS	:	value = value.(float64) + rightValue.(float64);
			case MINUS	:	value = value.(float64) - rightValue.(float64);
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateMultiplicativeModifier(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, rightValue interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateValue(stream, parameters);

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {

		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break;
		}

		rightValue, err = evaluateValue(stream, parameters);
		if(err != nil) {
			return nil, err;
		}

		switch(symbol) {

			case MULTIPLY	:	value = value.(float64) * rightValue.(float64);
			case DIVIDE	:	value = value.(float64) / rightValue.(float64);
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateValue(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, parameterValue interface{};
	var errorMessage, variableName string;
	var err error;

	token = stream.next();

	switch(token.Kind) {

		case CLAUSE	:	value, err = evaluateTokens(stream, parameters);
					if(err != nil) {
						return nil, err;
					}

					token = stream.next();
					if(token.Kind != CLAUSE_CLOSE) {

						return nil, errors.New("Unbalanced parenthesis");
					}

					return value, nil;

		case VARIABLE	:	variableName = token.Value.(string);
					parameterValue = parameters[variableName];

					if(parameterValue == nil) {
						errorMessage = "No parameter '"+ variableName +"' found."
						return nil, errors.New(errorMessage);
					}

		case NUMERIC	:	fallthrough
		case STRING	:	fallthrough
		case BOOLEAN	:	return token.Value, nil;
		default		:	break;
	}

	stream.rewind();
	return value, nil;
}

func (this EvaluableExpression) Tokens() []ExpressionToken {

	return this.tokens;
}

func (this EvaluableExpression) String() string {

	return this.inputExpression;
}

func isString(value interface{}) bool {

	switch value.(type) {
		case string	:	return true;
		default		:	return false;
	}
}
