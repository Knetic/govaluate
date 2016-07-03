package govaluate

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"
)

const isoDateFormat string = "2006-01-02T15:04:05.999999999Z0700"

var DUMMY_PARAMETERS = MapParameters(map[string]interface{}{})

/*
	EvaluableExpression represents a set of ExpressionTokens which, taken together,
	are an expression that can be evaluated down into a single value.
*/
type EvaluableExpression struct {

	/*
		Represents the query format used to output dates. Typically only used when creating SQL or Mongo queries from an expression.
		Defaults to the complete ISO8601 format, including nanoseconds.
	*/
	QueryDateFormat string

	tokens          []ExpressionToken
	inputExpression string
}

/*
	Parses a new EvaluableExpression from the given [expression] string.
	Returns an error if the given expression has invalid syntax.
*/
func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	var ret *EvaluableExpression
	var err error

	ret = new(EvaluableExpression)
	ret.QueryDateFormat = isoDateFormat
	ret.inputExpression = expression
	ret.tokens, err = parseTokens(expression)

	if err != nil {
		return nil, err
	}
	return ret, nil
}

/*
	Same as `Eval`, but automatically wraps a map of parameters into a `govalute.Parameters` structure.
*/
func (this EvaluableExpression) Evaluate(parameters map[string]interface{}) (interface{}, error) {

	if parameters == nil {
		return this.Eval(nil)
	}
	return this.Eval(MapParameters(parameters))
}

/*
	Runs the entire expression using the given [parameters].
	e.g., If the expression contains a reference to the variable "foo", it will be taken from `parameters.Get("foo")`.

	This function returns errors if the combination of expression and parameters cannot be run,
	such as if a variable in the expression is not present in [parameters].

	In all non-error circumstances, this returns the single value result of the expression and parameters given.
	e.g., if the expression is "1 + 1", this will return 2.0.
	e.g., if the expression is "foo + 1" and parameters contains "foo" = 2, this will return 3.0
*/
func (this EvaluableExpression) Eval(parameters Parameters) (interface{}, error) {

	var stream *tokenStream
	var err error

	if parameters != nil {
		parameters = &sanitizedParameters{parameters}
	}

	if err != nil {
		return nil, err
	}

	stream = newTokenStream(this.tokens)
	return evaluateTokens(stream, parameters)
}

func evaluateTokens(stream *tokenStream, parameters Parameters) (interface{}, error) {

	if stream.hasNext() {
		return evaluateTernary(stream, parameters)
	}
	return nil, nil
}

func evaluateTernary(stream *tokenStream, parameters Parameters) (interface{}, error) {
	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateTernary2(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = TERNARY_SYMBOLS[token.Value.(string)]
		if !keyFound || symbol != TERNARY_FALSE {
			break
		}

		rightValue, err = evaluateTernary2(stream, parameters)
		if err != nil {
			return nil, err
		}

		if value == nil || value == false {
			return rightValue, nil
		} else {
			return value, nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluateTernary2(stream *tokenStream, parameters Parameters) (interface{}, error) {
	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateLogical(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = TERNARY_SYMBOLS[token.Value.(string)]
		if !keyFound || symbol != TERNARY_TRUE {
			break
		}

		rightValue, err = evaluateLogical(stream, parameters)
		if err != nil {
			return nil, err
		}

		// make sure that we're only operating on the appropriate types
		if !isBool(value) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the ternary operator '%v', it is not a bool", value, token.Value))
		}

		if value.(bool) {
			return rightValue, nil
		} else {
			return nil, nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluateLogical(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value, newValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateComparator(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = LOGICAL_SYMBOLS[token.Value.(string)]
		if !keyFound {
			break
		}

		switch symbol {

		case OR:
			if value != nil {
				newValue, err = evaluateLogical(stream, parameters)
			}
		case AND:
			if value != nil {
				newValue, err = evaluateComparator(stream, parameters)
			}
		}

		if err != nil {
			return nil, err
		}

		// make sure that we're only operating on the appropriate types
		if !isBool(value) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the logical operator '%v', it is not a bool", value, token.Value))
		}
		if !isBool(newValue) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the logical operator '%v', it is not a bool", newValue, token.Value))
		}

		if symbol == OR {
			return value.(bool) || newValue.(bool), nil
		} else {
			return value.(bool) && newValue.(bool), nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluateComparator(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var pattern *regexp.Regexp
	var err error
	var keyFound bool

	value, err = evaluateAdditiveModifier(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = COMPARATOR_SYMBOLS[token.Value.(string)]
		if !keyFound {
			break
		}

		rightValue, err = evaluateAdditiveModifier(stream, parameters)
		if err != nil {
			return nil, err
		}

		// make sure that we're only operating on the appropriate types
		if symbol.IsModifierType(NUMERIC_COMPARATORS) {
			if !isFloat64(value) {
				return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the comparator '%v', it is not a number", value, token.Value))
			}
			if !isFloat64(rightValue) {
				return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the comparator '%v', it is not a number", rightValue, token.Value))
			}
		}

		if symbol.IsModifierType(STRING_COMPARATORS) {
			if !isString(value) {
				return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the comparator '%v', it is not a string", value, token.Value))
			}
			if !isRegexOrString(rightValue) {
				return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the comparator '%v', it is not a string", rightValue, token.Value))
			}
		}

		switch symbol {

		case LT:
			return (value.(float64) < rightValue.(float64)), nil
		case LTE:
			return (value.(float64) <= rightValue.(float64)), nil
		case GT:
			return (value.(float64) > rightValue.(float64)), nil
		case GTE:
			return (value.(float64) >= rightValue.(float64)), nil
		case EQ:
			return (value == rightValue), nil
		case NEQ:
			return (value != rightValue), nil
		case REQ:

			switch rightValue.(type) {
			case string:
				pattern, err = regexp.Compile(rightValue.(string))
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Unable to compile regexp pattern '%v': %v", rightValue, err))
				}
			case *regexp.Regexp:
				pattern = rightValue.(*regexp.Regexp)
			}

			return pattern.Match([]byte(value.(string))), nil
		case NREQ:

			switch rightValue.(type) {
			case string:
				pattern, err = regexp.Compile(rightValue.(string))
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Unable to compile regexp pattern '%v': %v", rightValue, err))
				}
			case *regexp.Regexp:
				pattern = rightValue.(*regexp.Regexp)
			}

			return !pattern.Match([]byte(value.(string))), nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluateAdditiveModifier(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateMultiplicativeModifier(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)]
		if !keyFound {
			break
		}

		// short-circuit if this is, in fact, not an additive modifier
		if !symbol.IsModifierType(ADDITIVE_MODIFIERS) {
			stream.rewind()
			return value, nil
		}

		rightValue, err = evaluateAdditiveModifier(stream, parameters)
		if err != nil {
			return nil, err
		}

		// short-circuit to check if we're supposed to do a concat on strings
		if symbol == PLUS {
			if isString(value) || isString(rightValue) {
				return fmt.Sprintf("%v%v", value, rightValue), nil
			}
		}

		// make sure that we're only operating on the appropriate types
		if !isFloat64(value) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the modifier '%v', it is not a number", value, token.Value))
		}
		if !isFloat64(rightValue) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the modifier '%v', it is not a number", rightValue, token.Value))
		}

		switch symbol {

		case PLUS:
			value = value.(float64) + rightValue.(float64)
		case MINUS:
			return value.(float64) - rightValue.(float64), nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluateMultiplicativeModifier(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateExponentialModifier(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)]
		if !keyFound {
			break
		}

		// short circuit if this is, in fact, not multiplicative.
		if !symbol.IsModifierType(MULTIPLICATIVE_MODIFIERS) {
			stream.rewind()
			return value, nil
		}

		rightValue, err = evaluateMultiplicativeModifier(stream, parameters)
		if err != nil {
			return nil, err
		}

		// make sure that we're only operating on the appropriate types
		if !isFloat64(value) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the modifier '%v', it is not a number", value, token.Value))
		}
		if !isFloat64(rightValue) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the modifier '%v', it is not a number", rightValue, token.Value))
		}

		switch symbol {

		case MULTIPLY:
			return value.(float64) * rightValue.(float64), nil
		case DIVIDE:
			return value.(float64) / rightValue.(float64), nil
		case MODULUS:
			return math.Mod(value.(float64), rightValue.(float64)), nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluateExponentialModifier(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateValue(stream, parameters)

	if err != nil {
		return nil, err
	}

	for stream.hasNext() {

		token = stream.next()

		if !isString(token.Value) {
			break
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)]
		if !keyFound {
			break
		}

		// if this isn't actually an exponential modifier, rewind and return.
		if !symbol.IsModifierType(EXPONENTIAL_MODIFIERS) {
			stream.rewind()
			return value, nil
		}

		rightValue, err = evaluateExponentialModifier(stream, parameters)
		if err != nil {
			return nil, err
		}

		// make sure that we're only operating on the appropriate types
		if !isFloat64(value) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the modifier '%v', it is not a number", value, token.Value))
		}
		if !isFloat64(rightValue) {
			return nil, errors.New(fmt.Sprintf("Value '%v' cannot be used with the modifier '%v', it is not a number", rightValue, token.Value))
		}

		switch symbol {
		case EXPONENT:
			return math.Pow(value.(float64), rightValue.(float64)), nil
		}
	}

	stream.rewind()
	return value, nil
}

func evaluatePrefix(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	for stream.hasNext() {

		token = stream.next()

		if token.Kind != PREFIX || !isString(token.Value) {
			break
		}

		symbol, keyFound = PREFIX_SYMBOLS[token.Value.(string)]
		if !keyFound {
			break
		}

		// TODO: need a prefix operator symbol check?

		value, err = evaluateValue(stream, parameters)
		if err != nil {
			return nil, err
		}

		switch symbol {

		case INVERT:
			return !value.(bool), nil

		case NEGATE:
			return -value.(float64), nil

		default:
			stream.rewind()
			return value, nil
		}
	}

	stream.rewind()
	return nil, nil
}

func evaluateValue(stream *tokenStream, parameters Parameters) (interface{}, error) {

	var token ExpressionToken
	var value interface{}
	var errorMessage, variableName string
	var err error

	token = stream.next()

	switch token.Kind {

	case CLAUSE:
		value, err = evaluateTokens(stream, parameters)
		if err != nil {
			return nil, err
		}

		// advance past the CLAUSE_CLOSE token. We know that it's a CLAUSE_CLOSE, because at parse-time we check for unbalanced parens.
		stream.next()
		return value, nil

	case VARIABLE:
		variableName = token.Value.(string)
		value, err = parameters.Get(variableName)
		if err != nil {
			return nil, err
		}

		if value == nil {
			errorMessage = "No parameter '" + variableName + "' found."
			return nil, errors.New(errorMessage)
		}

		return value, nil

	case NUMERIC:
		fallthrough
	case STRING:
		fallthrough
	case PATTERN:
		fallthrough
	case BOOLEAN:
		return token.Value, nil
	case TIME:
		return float64(token.Value.(time.Time).Unix()), nil

	case PREFIX:
		stream.rewind()

		value, err = evaluatePrefix(stream, parameters)
		if err != nil {
			return nil, err
		}
		if value == nil {
			break
		}

		return value, nil
	default:
		break
	}

	stream.rewind()
	return nil, errors.New("Unable to evaluate token kind: " + GetTokenKindString(token.Kind))
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

	var stream *tokenStream
	var token ExpressionToken
	var retBuffer bytes.Buffer
	var toWrite, ret string

	stream = newTokenStream(this.tokens)

	for stream.hasNext() {

		token = stream.next()

		switch token.Kind {

		case STRING:
			toWrite = fmt.Sprintf("'%v' ", token.Value)
		case TIME:
			toWrite = fmt.Sprintf("'%s' ", token.Value.(time.Time).Format(this.QueryDateFormat))

		case LOGICALOP:
			switch LOGICAL_SYMBOLS[token.Value.(string)] {

			case AND:
				toWrite = "AND "
			case OR:
				toWrite = "OR "
			}

		case BOOLEAN:
			if token.Value.(bool) {
				toWrite = "1 "
			} else {
				toWrite = "0 "
			}

		case VARIABLE:
			toWrite = fmt.Sprintf("[%s] ", token.Value.(string))

		case NUMERIC:
			toWrite = fmt.Sprintf("%g ", token.Value.(float64))

		case COMPARATOR:
			switch COMPARATOR_SYMBOLS[token.Value.(string)] {

			case EQ:
				toWrite = "= "
			case NEQ:
				toWrite = "<> "
			default:
				toWrite = fmt.Sprintf("%s ", token.Value.(string))
			}

		case PREFIX:
			toWrite = fmt.Sprintf("%s", token.Value.(string))
		case MODIFIER:
			toWrite = fmt.Sprintf("%s ", token.Value.(string))
		case CLAUSE:
			toWrite = "( "
		case CLAUSE_CLOSE:
			toWrite = ") "

		default:
			toWrite = fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind)
			return "", errors.New(toWrite)
		}

		retBuffer.WriteString(toWrite)
	}

	// trim last space.
	ret = retBuffer.String()
	ret = ret[:len(ret)-1]

	return ret, nil
}

/*
	Returns an array representing the ExpressionTokens that make up this expression.
*/
func (this EvaluableExpression) Tokens() []ExpressionToken {

	return this.tokens
}

/*
	Returns the original expression used to create this EvaluableExpression.
*/
func (this EvaluableExpression) String() string {

	return this.inputExpression
}

func isString(value interface{}) bool {

	switch value.(type) {
	case string:
		return true
	}
	return false
}

func isRegexOrString(value interface{}) bool {

	switch value.(type) {
	case string:
		return true
	case *regexp.Regexp:
		return true
	}
	return false
}

func isBool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return true
	}
	return false
}

func isFloat64(value interface{}) bool {
	switch value.(type) {
	case float64:
		return true
	}
	return false
}
