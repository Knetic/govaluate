package govaluate

const STAGE_SYMBOL_MAP = map[OperatorSymbol]evaluationOperator = {
	EQ: equalsStage,
	NEQ: notEqualsStage,
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

/*
	Creates a `evaluationStageList` object which represents an execution plan (or tree)
	which is used to completely evaluate a set of tokens at evaluation-time.
	The three stages of evaluation can be thought of as parsing strings to tokens, then tokens to a stage list, then evaluation with parameters.
*/
func planStages(tokens []ExpressionToken) (*evaluationStage, rror) {

	stream := newTokenStream(tokens)
	stages := newEvaluationStageList()

	stage, err := evaluateTokens(stream)
	if(err != nil) {
		return nil, err
	}

	return stage, nil
}

func evaluateTokens(stream *tokenStream) (*evaluationStage, error) {

	if !stream.hasNext() {
		return nil, nil
	}

	value := &precedencePlanner {
		validSymbols: ,
		leftTypeCheck: isFloat,
		rightTypeCheck: isFloat,
	}
	exponential := &precedencePlanner {
		validSymbols: ,
		leftTypeCheck: isFloat,
		rightTypeCheck: isFloat,
		next: value,
	}
	multiplicative := &precedencePlanner {
		validSymbols: ,
		leftTypeCheck: isFloat,
		rightTypeCheck: isFloat,
		next: exponential,
	}
	additive := &precedencePlanner {
		validSymbols: ,
		leftTypeCheck: isFloat,
		rightTypeCheck: isFloat,
		next: multiplicative,
	}

	// TODO
	comparator := &precedencePlanner {
		validSymbols: COMPARATOR_SYMBOLS,
		leftTypeCheck: //TODO
		rightTypeCheck: // TODO
		next: ,
	}

	logical := &precedencePlanner {
		validSymbols: LOGICAL_SYMBOLS,
		leftTypeCheck: isBool,
		rightTypeCheck: isBool,
		next: comparator,
	}
	ternary := &precedencePlanner {
		validSymbols: TERNARY_SYMBOLS,
		leftTypeCheck: isBool,
		rightTypeCheck: isBool,
		next: logical,
	}
}

type precedencePlanner struct {

	validSymbols map[string]OperatorSymbol

	leftTypeCheck stageTypeCheck
	rightTypeCheck stageTypeCheck

	next *precedencePlanner
)

func evaluatePrecedenceLevel(
	stream *tokenStream,
	leftTypeCheck stageTypeCheck,
	rightTypeCheck stageTypeCheck,
	validSymbols map[string]OperatorSymbol,
	next precedencePlanner
	) (*evaluationStage, error) {

	var token ExpressionToken
	var symbol OperatorSymbol
	var leftStage, rightStage *evaluationStage
	var err error
	var keyFound bool

	leftStage, err = next(stream)
	if err != nil {
		return nil, err
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

		rightStage, err = next(stream)
		if err != nil {
			return nil, err
		}

		stage := &evaluationStage {
			operator: operator,
			leftTypeCheck: leftTypeCheck,
			rightTypeCheck: rightTypeCheck,
		}
	}

	stream.rewind()
	return value, nil
}


func evaluateComparator(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var pattern *regexp.Regexp
	var err error
	var keyFound bool

	value, err = evaluateAdditiveModifier(stream)

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

		rightValue, err = evaluateAdditiveModifier(stream)
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

func evaluateAdditiveModifier(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateMultiplicativeModifier(stream)

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

		rightValue, err = evaluateAdditiveModifier(stream)
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

func evaluateMultiplicativeModifier(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateExponentialModifier(stream)

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

		rightValue, err = evaluateMultiplicativeModifier(stream)
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

func evaluateExponentialModifier(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var value, rightValue interface{}
	var symbol OperatorSymbol
	var err error
	var keyFound bool

	value, err = evaluateValue(stream)

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

		rightValue, err = evaluateExponentialModifier(stream)
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

func evaluatePrefix(stream *tokenStream) (*evaluationStage, error) {

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

		value, err = evaluateValue(stream)
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

func evaluateValue(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var value interface{}
	var errorMessage, variableName string
	var err error

	token = stream.next()

	switch token.Kind {

	case CLAUSE:
		value, err = evaluateTokens(stream)
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

		value, err = evaluatePrefix(stream)
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
