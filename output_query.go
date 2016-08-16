package govaluate

import (

	"fmt"
	"time"
	"errors"
)

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
	var transactions *expressionOutputStream
	var transaction string
	var err error

	stream = newTokenStream(this.tokens)
	transactions = new(expressionOutputStream)

	for stream.hasNext() {

		transaction, err = this.findNextSQLString(stream, transactions)
		if(err != nil) {
			return "", err
		}

		transactions.add(transaction)
	}

	return transactions.createString(" "), nil
}

func (this EvaluableExpression) findNextSQLString(stream *tokenStream, transactions *expressionOutputStream) (string, error) {

	var token ExpressionToken
	var ret string

	token = stream.next()

	switch token.Kind {

	case STRING:
		ret = fmt.Sprintf("'%v'", token.Value)
	case TIME:
		ret = fmt.Sprintf("'%s'", token.Value.(time.Time).Format(this.QueryDateFormat))

	case LOGICALOP:
		switch LOGICAL_SYMBOLS[token.Value.(string)] {

		case AND:
			ret = "AND"
		case OR:
			ret = "OR"
		}

	case BOOLEAN:
		if token.Value.(bool) {
			ret = "1"
		} else {
			ret = "0"
		}

	case VARIABLE:
		ret = fmt.Sprintf("[%s]", token.Value.(string))

	case NUMERIC:
		ret = fmt.Sprintf("%g", token.Value.(float64))

	case COMPARATOR:
		switch COMPARATOR_SYMBOLS[token.Value.(string)] {

		case EQ:
			ret = "="
		case NEQ:
			ret = "<>"
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}

	case PREFIX:
		switch PREFIX_SYMBOLS[token.Value.(string)] {

		case INVERT:
			ret = fmt.Sprintf("NOT")
		default:

			right, err := this.findNextSQLString(stream, transactions)
			if(err != nil) {
				return "", err
			}
			
			ret = fmt.Sprintf("%s%s", token.Value.(string), right)
		}
	case MODIFIER:

		switch MODIFIER_SYMBOLS[token.Value.(string)] {

		case EXPONENT:

			left := transactions.rollback()
			right, err := this.findNextSQLString(stream, transactions)
			if(err != nil) {
				return "", err
			}

			ret = fmt.Sprintf("POW(%s, %s)", left, right)
		case MODULUS:

			left := transactions.rollback()
			right, err := this.findNextSQLString(stream, transactions)
			if(err != nil) {
				return "", err
			}

			ret = fmt.Sprintf("MOD(%s, %s)", left, right)
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}
	case CLAUSE:
		ret = "("
	case CLAUSE_CLOSE:
		ret = ")"

	default:
		errorMsg := fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind)
		return "", errors.New(errorMsg)
	}

	return ret, nil
}
