package govaluate

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var findMap map[TokenKind]func(EvaluableExpression, parsableInput) (string, error)

type TokenData struct {
	Expression string            `json:"name"`
	Tokens     []ExpressionToken `json:"tokens"`
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
	var transactions *expressionOutputStream
	var transaction string
	var err error

	stream = newTokenStream(this.tokens)
	transactions = new(expressionOutputStream)

	for stream.hasNext() {

		transaction, err = this.findNextSQLString(stream, transactions)
		if err != nil {
			return "", err
		}

		transactions.add(transaction)
	}

	return transactions.createString(" "), nil
}

type parsableInput struct {
	token        ExpressionToken
	stream       *tokenStream
	transactions *expressionOutputStream
}

var findNextSQLSubstringMap = map[TokenKind]func(EvaluableExpression, parsableInput) (string, error){
	STRING:       String,
	PATTERN:      Pattern,
	TIME:         Time,
	LOGICALOP:    LogicalOP,
	BOOLEAN:      Boolean,
	VARIABLE:     Variable,
	NUMERIC:      Numeric,
	COMPARATOR:   Comparator,
	TERNARY:      Ternary,
	PREFIX:       Prefix,
	MODIFIER:     Modifier,
	CLAUSE:       Clause,
	CLAUSE_CLOSE: Clause_Close,
	SEPARATOR:    Separator,
}

func init() {
	findMap = findNextSQLSubstringMap
}

func (this EvaluableExpression) findNextSQLString(stream *tokenStream, transactions *expressionOutputStream) (string, error) {

	token := stream.next()
	sqlSubstringFunc, found := findMap[token.Kind]
	if found {
		return sqlSubstringFunc(this, parsableInput{token: token, stream: stream, transactions: transactions})
	}
	return "", errors.New(fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind))
}

func String(this EvaluableExpression, parsable parsableInput) (string, error) {
	return fmt.Sprintf("'%v'", parsable.token.Value), nil
}
func Pattern(this EvaluableExpression, parsable parsableInput) (string, error) {
	return fmt.Sprintf("'%s'", parsable.token.Value.(*regexp.Regexp).String()), nil
}
func Time(this EvaluableExpression, parsable parsableInput) (string, error) {
	return fmt.Sprintf("'%s'", parsable.token.Value.(time.Time).Format(this.QueryDateFormat)), nil
}
func LogicalOP(this EvaluableExpression, parsable parsableInput) (string, error) {
	if logicalSymbols[parsable.token.Value.(string)] == AND {
		return "AND", nil
	}
	return "OR", nil
}
func Boolean(this EvaluableExpression, parsable parsableInput) (string, error) {
	if parsable.token.Value.(bool) {
		return "1", nil
	}
	return "0", nil
}
func Variable(this EvaluableExpression, parsable parsableInput) (string, error) {
	return fmt.Sprintf("[%s]", parsable.token.Value.(string)), nil
}
func Numeric(g EvaluableExpression, parsable parsableInput) (string, error) {
	return fmt.Sprintf("%g", parsable.token.Value.(float64)), nil
}
func Comparator(this EvaluableExpression, parsable parsableInput) (string, error) {
	comparatorSymbol := comparatorSymbols[parsable.token.Value.(string)]
	symbol, found := comparatorSymbolsReverse[comparatorSymbol]
	if found {
		return symbol, nil
	}
	return fmt.Sprintf("%s", parsable.token.Value.(string)), nil
}
func Ternary(this EvaluableExpression, parsable parsableInput) (string, error) {

	ternarySymbol := ternarySymbols[parsable.token.Value.(string)]
	if ternarySymbol == COALESCE {

		left := parsable.transactions.rollback()
		right, err := this.findNextSQLString(parsable.stream, parsable.transactions)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("COALESCE(%v, %v)", left, right), nil
	}
	return "", errors.New("Ternary operators are unsupported in SQL output")
}
func Prefix(this EvaluableExpression, parsable parsableInput) (string, error) {
	prefixSymbol := prefixSymbols[parsable.token.Value.(string)]
	if prefixSymbol == INVERT {
		return fmt.Sprintf("NOT"), nil
	}

	right, err := this.findNextSQLString(parsable.stream, parsable.transactions)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", parsable.token.Value.(string), right), nil
}
func Modifier(this EvaluableExpression, parsable parsableInput) (string, error) {
	modifierSymbol := modifierSymbols[parsable.token.Value.(string)]
	if modifierSymbol == EXPONENT {
		left := parsable.transactions.rollback()
		right, err := this.findNextSQLString(parsable.stream, parsable.transactions)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("POW(%s, %s)", left, right), nil
	}

	if modifierSymbol == MODULUS {
		left := parsable.transactions.rollback()
		right, err := this.findNextSQLString(parsable.stream, parsable.transactions)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("MOD(%s, %s)", left, right), nil
	}
	return fmt.Sprintf("%s", parsable.token.Value.(string)), nil
}
func Clause(this EvaluableExpression, parsable parsableInput) (string, error) {
	return "(", nil
}
func Clause_Close(this EvaluableExpression, parsable parsableInput) (string, error) {
	return ")", nil
}
func Separator(this EvaluableExpression, parsable parsableInput) (string, error) {
	return ",", nil
}

var comparatorSymbolsReverse = map[OperatorSymbol]string{
	EQ:   "=",
	NEQ:  "<>",
	REQ:  "RLIKE",
	NREQ: "NOT RLIKE",
}

func (this EvaluableExpression) MarshalTokens() ([]byte, error) {

	var index int
	tokens := this.Tokens()

	for position, token := range tokens {

		if token.Kind == FUNCTION {
			token.Value = this.functions[index]
			tokens[position] = token
			index++
		}
	}

	data := TokenData{
		Tokens:     tokens,
		Expression: this.String(),
	}

	buffer, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func UnmarshalTokens(bytes []byte) (TokenData, error) {

	data := TokenData{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return TokenData{}, err
	}

	return data, nil
}
