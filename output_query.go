package govaluate

import (

	"fmt"
	"bytes"
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
