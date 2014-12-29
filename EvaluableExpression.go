package govaluate

import (
	"fmt"
	"errors"
	"bytes"
	"math"
	"time"
)

const isoDateFormat string = "2006-01-02T15:04:05.999999999Z0700";

/*
	EvaluableExpression represents a set of ExpressionTokens which, taken together,
	represent an arbitrary expression that can be evaluated down into a single value.
*/
type EvaluableExpression struct {

	/*
		Represents the query format used to output dates. Typically only used when creating SQL or Mongo queries from an expression.
		Defaults to the complete ISO8601 format, including nanoseconds.
	*/
	QueryDateFormat string;

	tokens []ExpressionToken
	inputExpression string
}

/*
	Creates a new EvaluableExpression from the given [expression] string.
	Returns an error if the given expression has invalid syntax.
*/
func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	var ret *EvaluableExpression;
	var err error

	ret = new(EvaluableExpression)
	ret.QueryDateFormat = isoDateFormat;
	ret.inputExpression = expression;
	ret.tokens, err = parseTokens(expression)

	if(err != nil) {
		return nil, err
	}
	return ret, nil
}

/*
	Evaluate runs the entire expression using the given [parameters]. 
	Each parameter is mapped from a string to a value, such as "foo" = 1.0. 
	If the expression contains a reference to the variable "foo", it will be taken from parameters["foo"].

	This function returns errors if the combination of expression and parameters cannot be run,
	such as if a string parameter is given in an expression that expects it to be a boolean. 
	e.g., "foo == true", where foo is any string.
	These errors are almost exclusively returned for parameters not being present, or being of the wrong type.
	Structural problems with the expression (unexpected tokens, unexpected end of expression, etc) are discovered
	during parsing of the expression in NewEvaluableExpression.

	In all non-error circumstances, this returns the single value result of the expression and parameters given.
	e.g., if the expression is "1 + 1", Evaluate will return 2.0.
	e.g., if the expression is "foo + 1" and parameters contains "foo" = 2, Evaluate will return 3.0
*/
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

			case LT		:	return (value.(float64) < rightValue.(float64)), nil;
			case LTE	:	return (value.(float64) <= rightValue.(float64)), nil;
			case GT		:	return (value.(float64) > rightValue.(float64)), nil;
			case GTE	:	return (value.(float64) >= rightValue.(float64)), nil;
			case EQ		:	return (value == rightValue), nil;
			case NEQ	:	return (value != rightValue), nil;
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

		switch(symbol) {

			case PLUS	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						value = value.(float64) + rightValue.(float64);

			case MINUS	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}

						return value.(float64) - rightValue.(float64), nil;

			default		:	stream.rewind();
						return value, nil;
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

	value, err = evaluateExponentialModifier(stream, parameters);

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

		switch(symbol) {

			case MULTIPLY	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						return value.(float64) * rightValue.(float64), nil;

			case DIVIDE	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						return value.(float64) / rightValue.(float64), nil;

			case MODULUS	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						return math.Mod(value.(float64), rightValue.(float64)), nil;

			default		:	stream.rewind();
						return value, nil;
		}
	}

	stream.rewind();	
	return value, nil;
}

func evaluateExponentialModifier(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

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

		switch(symbol) {

			case EXPONENT	:	rightValue, err = evaluateExponentialModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						return math.Pow(value.(float64), rightValue.(float64)), nil;

			default		:	stream.rewind();
						return value, nil;
		}
	}

	stream.rewind();	
	return value, nil;
}

func evaluateValue(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value interface{};
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
					value = parameters[variableName];

					if(value == nil) {
						errorMessage = "No parameter '"+ variableName +"' found."
						return nil, errors.New(errorMessage);
					}

					return value, nil;

		case NUMERIC	:	fallthrough
		case STRING	:	fallthrough
		case BOOLEAN	:	return token.Value, nil;
		case TIME	:	return float64(token.Value.(time.Time).Unix()), nil;
		default		:	break;
	}

	stream.rewind();
	return nil, errors.New("Unable to evaluate token kind: " + GetTokenKindString(token.Kind));
}

/*
	Returns a string representing this expression as if it were written in SQL.
	This function assumes that all parameters exist within the same table, and that the table essentially represents
	a serialized object of some sort (e.g., hibernate).
	If your data model is more normalized, you may need to consider iterating through each actual token given by `Tokens()`
	to create your query.

	Boolean values are considered to be "1" for true, "0" for false.

	Times are formatted according to this.QueryDateFormat.
*/
func (this EvaluableExpression) ToSQLQuery() (string, error) {

	var stream *tokenStream;
	var token ExpressionToken;
	var retBuffer bytes.Buffer;
	var toWrite, ret string;

	stream = newTokenStream(this.tokens);

	for(stream.hasNext()) {

		token = stream.next();

		switch(token.Kind) {

			case STRING		:	toWrite = fmt.Sprintf("'%v' ", token.Value);
			case TIME		:	toWrite = fmt.Sprintf("'%s' ", token.Value.(time.Time).Format(this.QueryDateFormat));
			
			case LOGICALOP		:	switch(LOGICAL_SYMBOLS[token.Value.(string)]) {

								case AND	:	toWrite = "AND ";
								case OR		:	toWrite = "OR ";
							}

			case BOOLEAN		:	if(token.Value.(bool)) {
								toWrite = "1 ";
							} else {
								toWrite = "0 ";
							}

			case VARIABLE		:	toWrite = fmt.Sprintf("[%s] ", token.Value.(string));

			case NUMERIC 		:	toWrite = fmt.Sprintf("%g ", token.Value.(float64));

			case COMPARATOR		:	switch(COMPARATOR_SYMBOLS[token.Value.(string)]) {

								case EQ		:	toWrite = "= ";
								case NEQ	:	toWrite = "<> ";
								default		:	toWrite = fmt.Sprintf("%s ", token.Value.(string));
							}
	
			case MODIFIER		:	toWrite = fmt.Sprintf("%s ", token.Value.(string));
			case CLAUSE		:	toWrite = "( "
			case CLAUSE_CLOSE	:	toWrite = ") "

			default			:	toWrite = fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind);
							return "", errors.New(toWrite);
		}

		retBuffer.WriteString(toWrite);
	}

	// trim last space.
	ret = retBuffer.String();
	ret = ret[:len(ret)-1];

	return ret, nil;
}

/*
	Returns a string representing this expression as if it were written as a Mongo query.
*/
func (this EvaluableExpression) ToMongoQuery() (string, error) {

	var stream *tokenStream;
	var token ExpressionToken;
	var retBuffer bytes.Buffer;
	var toWrite, ret string;

	stream = newTokenStream(this.tokens);

	for(stream.hasNext()) {

		token = stream.next();

		switch(token.Kind) {

			case STRING		:	toWrite = fmt.Sprintf("\"%s\" ", token.Value.(string));
			case TIME		:	toWrite = fmt.Sprintf("ISODate(\"%s\") ", token.Value.(time.Time).Format(isoDateFormat));
			case LOGICALOP		:	
			case BOOLEAN		:	if(token.Value.(bool)) {
								toWrite = "true ";
							} else {
								toWrite = "false ";
							}
			case VARIABLE		:	toWrite = fmt.Sprintf("%s ", token.Value.(string));
			case NUMERIC 		:	toWrite = fmt.Sprintf("%g ", token.Value.(float64));
			case COMPARATOR		:	
			case CLAUSE		:	fallthrough;
			case CLAUSE_CLOSE	:	continue;

			case MODIFIER		:	toWrite = fmt.Sprintf("Unable to use modifiers in Mongo queries (found '%s')", token.Kind);
							return "", errors.New(toWrite);

			default			:	toWrite = fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind);
							return "", errors.New(toWrite);
		}

		retBuffer.WriteString(toWrite);
	}

	// trim last space.
	ret = retBuffer.String();
	ret = ret[:len(ret)-1];

	return ret, nil;
}

/*
	Returns an array representing the ExpressionTokens that make up this expression.
*/
func (this EvaluableExpression) Tokens() []ExpressionToken {

	return this.tokens;
}

/*
	Returns the original expression used to create this EvaluableExpression.
*/
func (this EvaluableExpression) String() string {

	return this.inputExpression;
}

func isString(value interface{}) bool {

	switch value.(type) {
		case string	:	return true;
		default		:	break;
	}
	return false;
}
