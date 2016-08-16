package govaluate

import (
	"bytes"
	"errors"
	"fmt"
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

	tokens           []ExpressionToken
	evaluationStages *evaluationStage
	inputExpression  string
}

/*
	Parses a new EvaluableExpression from the given [expression] string.
	Returns an error if the given expression has invalid syntax.
*/
func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	functions := make(map[string]ExpressionFunction)
	return NewEvaluableExpressionWithFunctions(expression, functions)
}

/*
	Similar to [NewEvaluableExpression], except that instead of a string, an already-tokenized expression is given.
	This is useful in cases where you may be generating an expression automatically, or using some other parser (e.g., to parse from a query language)
*/
func NewEvaluableExpressionFromTokens(tokens []ExpressionToken) (*EvaluableExpression, error) {

	var ret *EvaluableExpression
	var err error

	ret = new(EvaluableExpression)
	ret.QueryDateFormat = isoDateFormat

	err = checkBalance(tokens)
	if(err != nil) {
		return nil, err
	}

	err = checkExpressionSyntax(tokens)
	if(err != nil) {
		return nil, err
	}

	ret.tokens, err = optimizeTokens(tokens)
	if(err != nil) {
		return nil, err
	}

	ret.evaluationStages, err = planStages(ret.tokens)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

/*
	Similar to [NewEvaluableExpression], except enables the use of user-defined functions.
	Functions passed into this will be available to the expression.
*/
func NewEvaluableExpressionWithFunctions(expression string, functions map[string]ExpressionFunction) (*EvaluableExpression, error) {

	var ret *EvaluableExpression
	var err error

	ret = new(EvaluableExpression)
	ret.QueryDateFormat = isoDateFormat
	ret.inputExpression = expression

	ret.tokens, err = parseTokens(expression, functions)
	if err != nil {
		return nil, err
	}

	err = checkExpressionSyntax(ret.tokens)
	if(err != nil) {
		return nil, err
	}

	ret.evaluationStages, err = planStages(ret.tokens)
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

	if this.evaluationStages == nil {
		return nil, nil
	}

	if parameters != nil {
		parameters = &sanitizedParameters{parameters}
	}
	return evaluateStage(this.evaluationStages, parameters)
}

func evaluateStage(stage *evaluationStage, parameters Parameters) (interface{}, error) {

	var left, right interface{}
	var err error

	if stage.leftStage != nil {
		left, err = evaluateStage(stage.leftStage, parameters)
		if err != nil {
			return nil, err
		}
	}

	if stage.rightStage != nil {
		right, err = evaluateStage(stage.rightStage, parameters)
		if err != nil {
			return nil, err
		}
	}

	// type checks
	if stage.typeCheck == nil {

		err = typeCheck(stage.leftTypeCheck, left, stage.symbol, stage.typeErrorFormat)
		if err != nil {
			return nil, err
		}

		err = typeCheck(stage.rightTypeCheck, right, stage.symbol, stage.typeErrorFormat)
		if err != nil {
			return nil, err
		}
	} else {
		// special case where the type check needs to know both sides to determine if the operator can handle it
		if !stage.typeCheck(left, right) {
			errorMsg := fmt.Sprintf(stage.typeErrorFormat, left, stage.symbol.String())
			return nil, errors.New(errorMsg)
		}
	}

	return stage.operator(left, right, parameters)
}

func typeCheck(check stageTypeCheck, value interface{}, symbol OperatorSymbol, format string) error {

	if check == nil {
		return nil
	}

	if check(value) {
		return nil
	}

	errorMsg := fmt.Sprintf(format, value, symbol.String())
	return errors.New(errorMsg)
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
			switch PREFIX_SYMBOLS[token.Value.(string)] {

			case INVERT:
				toWrite = fmt.Sprintf("NOT ")
			default:
				toWrite = fmt.Sprintf("%s", token.Value.(string))
			}
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
