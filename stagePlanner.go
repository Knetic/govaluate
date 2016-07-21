package govaluate

import (
	"errors"
	"time"
)

var stageSymbolMap = map[OperatorSymbol]evaluationOperator {
	EQ: equalStage,
	NEQ: notEqualStage,
	GT: gtStage,
	LT: ltStage,
	GTE: gteStage,
	LTE: lteStage,
	REQ: regexStage,
	NREQ: notRegexStage,
	AND: andStage,
	OR: orStage,
	PLUS: addStage,
	MINUS: subtractStage,
	MULTIPLY: multiplyStage,
	DIVIDE: divideStage,
	MODULUS: modulusStage,
	EXPONENT: exponentStage,
	NEGATE: negateStage,
	INVERT: invertStage,
	TERNARY_TRUE: ternaryIfStage,
	TERNARY_FALSE: ternaryElseStage,
}

type precedent func(stream *tokenStream) (*evaluationStage, error)

type precedencePlanner struct {

	validSymbols map[string]OperatorSymbol

	leftTypeCheck stageTypeCheck
	rightTypeCheck stageTypeCheck

	next precedent
}

func makePrecedentFromPlanner(planner *precedencePlanner) precedent {

	return func(stream *tokenStream) (*evaluationStage, error) {
		return planPrecedenceLevel(
			stream,
			planner.leftTypeCheck,
			planner.rightTypeCheck,
			planner.validSymbols,
			planner.next,
		)
	}
}

var planPrefix precedent
var planExponential precedent
var planMultiplicative precedent
var planAdditive precedent
var planLogical precedent
var planTernary precedent

func init() {

	planPrefix = makePrecedentFromPlanner(&precedencePlanner {
		validSymbols: PREFIX_SYMBOLS,
	})
	planExponential = makePrecedentFromPlanner(&precedencePlanner {
		validSymbols: EXPONENTIAL_SYMBOLS,
		leftTypeCheck: isFloat64,
		rightTypeCheck: isFloat64,
		next: planValue,
	})
	planMultiplicative = makePrecedentFromPlanner(&precedencePlanner {
		validSymbols: MULTIPLICATIVE_SYMBOLS,
		leftTypeCheck: isFloat64,
		rightTypeCheck: isFloat64,
		next: planExponential,
	})
	planAdditive = makePrecedentFromPlanner(&precedencePlanner {
		validSymbols: ADDITIVE_SYMBOLS,
		leftTypeCheck: isFloat64,
		rightTypeCheck: isFloat64,
		next: planMultiplicative,
	})
	planLogical = makePrecedentFromPlanner(&precedencePlanner {
		validSymbols: LOGICAL_SYMBOLS,
		leftTypeCheck: isBool,
		rightTypeCheck: isBool,
		next: planComparator,
	})
	planTernary = makePrecedentFromPlanner(&precedencePlanner {
		validSymbols: TERNARY_SYMBOLS,
		leftTypeCheck: isBool,
		rightTypeCheck: isBool,
		next: planLogical,
	})
}

/*
	Creates a `evaluationStageList` object which represents an execution plan (or tree)
	which is used to completely evaluate a set of tokens at evaluation-time.
	The three stages of evaluation can be thought of as parsing strings to tokens, then tokens to a stage list, then evaluation with parameters.
*/
func planStages(tokens []ExpressionToken) (*evaluationStage, error) {

	stream := newTokenStream(tokens)

	stage, err := planTokens(stream)
	if(err != nil) {
		return nil, err
	}

	return stage, nil
}

func planTokens(stream *tokenStream) (*evaluationStage, error) {

	if !stream.hasNext() {
		return nil, nil
	}

	return planTernary(stream)
}

func planPrecedenceLevel(
	stream *tokenStream,
	leftTypeCheck stageTypeCheck,
	rightTypeCheck stageTypeCheck,
	validSymbols map[string]OperatorSymbol,
	next precedent) (*evaluationStage, error) {

	var token ExpressionToken
	var symbol OperatorSymbol
	var leftStage, rightStage *evaluationStage
	var err error
	var keyFound bool

	if(next != nil) {
		leftStage, err = next(stream)
		if err != nil {
			return nil, err
		}
	}

	for stream.hasNext() {

		token = stream.next()
		if !isString(token.Value) {
			break
		}

		symbol, keyFound = validSymbols[token.Value.(string)]
		if !keyFound {
			break
		}

		if(next != nil) {
			rightStage, err = next(stream)
			if err != nil {
				return nil, err
			}
		}

		return &evaluationStage {

			leftStage: leftStage,
			rightStage: rightStage,
			operator: stageSymbolMap[symbol],

			leftTypeCheck: leftTypeCheck,
			rightTypeCheck: rightTypeCheck,
		}, nil
	}

	stream.rewind()
	return leftStage, nil
}


func planComparator(stream *tokenStream) (*evaluationStage, error) {

	// comparators can operate on a bunch of different types.
	// this is mostly a copy of `planPredecenceLevel`, except with multiple possible type checks based on the comparator.
	var token ExpressionToken
	var leftStage, rightStage *evaluationStage
	var symbol OperatorSymbol
	var leftTypeCheck, rightTypeCheck stageTypeCheck
	var err error
	var keyFound bool

	leftStage, err = planAdditive(stream)

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

		rightStage, err = planAdditive(stream)
		if err != nil {
			return nil, err
		}

		// make sure that we're only operating on the appropriate types
		if symbol.IsModifierType(NUMERIC_COMPARATORS) {
			leftTypeCheck = isFloat64
			rightTypeCheck = isFloat64
		}

		if symbol.IsModifierType(STRING_COMPARATORS) {
			leftTypeCheck = isString
			rightTypeCheck = isRegexOrString
		}

		return &evaluationStage {

			operator: stageSymbolMap[symbol],
			leftStage: leftStage,
			rightStage: rightStage,

			leftTypeCheck: leftTypeCheck,
			rightTypeCheck: rightTypeCheck,

			typeErrorFormat: TYPEERROR_COMPARATOR,
		}, nil
	}

	stream.rewind()
	return leftStage, nil
}

func planValue(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var ret *evaluationStage
	var operator evaluationOperator
	var err error

	token = stream.next()

	switch token.Kind {

	case CLAUSE:

		ret, err = planTokens(stream)
		if(err != nil) {
			return nil, err
		}

		// advance past the CLAUSE_CLOSE token. We know that it's a CLAUSE_CLOSE, because at parse-time we check for unbalanced parens.
		stream.next()
		return ret, nil

	case VARIABLE:
		operator = makeParameterStage(token.Value.(string))

	case NUMERIC:
		fallthrough
	case STRING:
		fallthrough
	case PATTERN:
		fallthrough
	case BOOLEAN:
		operator = makeLiteralStage(token.Value)
	case TIME:
		operator = makeLiteralStage(float64(token.Value.(time.Time).Unix()))

	case PREFIX:
		stream.rewind()
		return planPrefix(stream)
	}

	if(operator == nil) {
		return nil, errors.New("Unable to plan token kind: " + GetTokenKindString(token.Kind))
	}

	return &evaluationStage {
		operator: operator,
	}, nil
}
