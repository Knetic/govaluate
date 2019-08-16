package govaluate

// Reduce does a partial parameter evaluation and returns simplified expression
func (expr EvaluableExpression) Reduce(parameters map[string]interface{}) (EvaluableExpression, error) {

	if expr.evaluationStages == nil {
		return expr, nil
	}

	var p Parameters = DUMMY_PARAMETERS
	if parameters != nil {
		p = &sanitizedParameters{MapParameters(parameters)}
	}

	// simplify expression
	reducedStage, err := reduceStage(expr.evaluationStages, p)
	if err != nil {
		return expr, err
	}

	// build new token list that corresponds to the reduced expression
	tokens, err := buildTokens(reducedStage)
	if err != nil {
		return expr, err
	}

	// unwrap top-level parentheses
	for len(tokens) > 0 && tokens[0].Kind == CLAUSE {
		tokens = tokens[1 : len(tokens)-1]
	}

	expr.evaluationStages = reducedStage
	expr.tokens = tokens
	return expr, nil
}

func reduceStage(stage *evaluationStage, parameters Parameters) (*evaluationStage, error) {
	if stage == nil || stage.symbol == LITERAL {
		// can't reduce nil or literal
		return stage, nil
	}

	// process variable
	if stage.symbol == VALUE {
		if value, err := stage.operator(nil, nil, parameters); err == nil {
			// variable is known, replace it with value literal
			return newLiteral(value), nil
		}
		// variable is unknown, return as is
		return stage, nil
	}

	// reduce left operand
	leftStage, err := reduceStage(stage.leftStage, parameters)
	if err != nil {
		return nil, err
	}
	left, leftIsLiteral := getLiteralValue(leftStage)

	if leftIsLiteral && stage.isShortCircuitable() {
		switch stage.symbol {
		case AND:
			if left == false {
				// false && ... -> false
				return newLiteral(false), nil
			}
		case OR:
			if left == true {
				// true || ... -> true
				return newLiteral(true), nil
			}
		case COALESCE:
			if left != nil {
				// x ?: y, where x != nil -> x
				return leftStage, nil
			}

		case TERNARY_TRUE:
			if left == false {
				// false ? x : y -> nil : y -> y
				// return nil here so the parent TERNARY_FALSE will evaluate to the right operand
				return nil, nil
			}
		}
	}

	// reduce right operand
	rightStage, err := reduceStage(stage.rightStage, parameters)
	if err != nil {
		return nil, err
	}
	right, rightIsLiteral := getLiteralValue(rightStage)

	if (leftStage == nil || leftIsLiteral) && (rightStage == nil || rightIsLiteral) {
		// both operands are known, perform the operation
		value, err := stage.operator(left, right, parameters)
		if err != nil {
			return nil, err
		}
		return newLiteral(value), nil
	}

	// optimizations
	switch stage.symbol {
	case AND:
		if leftIsLiteral && left == true {
			// true && x -> x
			return rightStage, nil
		}
		if rightIsLiteral && right == true {
			// x && true -> x
			return leftStage, nil
		}
		if rightIsLiteral && right == false {
			// x && false -> false
			return newLiteral(false), nil
		}

	case OR:
		if leftIsLiteral && left == false {
			// false || x -> x
			return rightStage, nil
		}
		if rightIsLiteral && right == false {
			// x || false -> x
			return leftStage, nil
		}
		if rightIsLiteral && right == true {
			// x || true -> true
			return newLiteral(true), nil
		}

	case PLUS:
		if leftIsLiteral && left == 0.0 {
			// 0 + x -> x
			return rightStage, nil
		}
		if rightIsLiteral && right == 0.0 {
			// x + 0 -> x
			return leftStage, nil
		}

	case MINUS:
		if leftIsLiteral && left == 0.0 {
			// 0 - x -> -x
			return &evaluationStage{
				symbol:     NEGATE,
				rightStage: rightStage,
				operator:   negateStage,
			}, nil
		}
		if rightIsLiteral && right == 0.0 {
			// x - 0 -> x
			return leftStage, nil
		}

	case MULTIPLY:
		if leftIsLiteral && left == 0.0 || rightIsLiteral && right == 0.0 {
			// 0 * x -> 0, x * 0 -> 0
			return newLiteral(0.0), nil
		}
		if leftIsLiteral && left == 1.0 {
			// 1 * x -> x
			return rightStage, nil
		}
		if rightIsLiteral && right == 1.0 {
			// x * 1 -> x
			return leftStage, nil
		}

	case DIVIDE:
		if rightIsLiteral && right == 1.0 {
			// x / 1 -> x
			return leftStage, nil
		}

	case COALESCE:
		if rightIsLiteral && right == nil {
			// x ?: nil -> x
			return leftStage, nil
		}

	case TERNARY_TRUE:
		if leftIsLiteral && left == true {
			// true ? x : y -> x
			return rightStage, nil
		}

	case TERNARY_FALSE:
		if leftStage == nil {
			// false ? x : y -> y
			return rightStage, nil
		}
		if leftStage.symbol != TERNARY_TRUE {
			// true ? x : y -> x
			return leftStage, nil
		}
	}

	return &evaluationStage{
		symbol:          stage.symbol,
		leftStage:       leftStage,
		rightStage:      rightStage,
		operator:        stage.operator,
		leftTypeCheck:   stage.leftTypeCheck,
		rightTypeCheck:  stage.rightTypeCheck,
		typeCheck:       stage.typeCheck,
		typeErrorFormat: stage.typeErrorFormat,
	}, nil
}

func newLiteral(value interface{}) *evaluationStage {
	return &evaluationStage{
		symbol:   LITERAL,
		operator: makeLiteralStage(value),
	}
}

func getLiteralValue(stage *evaluationStage) (interface{}, bool) {
	if stage != nil && stage.symbol == LITERAL {
		value, err := stage.operator(nil, nil, nil)
		if err == nil {
			return value, true
		}
	}
	return nil, false
}

func buildTokens(stage *evaluationStage) ([]ExpressionToken, error) {
	if stage == nil {
		return []ExpressionToken{}, nil
	}

	switch stage.symbol {
	case VALUE:
		// get variable name: invoke the operator and look what parameter it has queried
		p := paramCaptor{}
		if _, err := stage.operator(nil, nil, &p); err != nil {
			return nil, err
		}
		variableName := p.lastParameterName
		return []ExpressionToken{
			newToken(VARIABLE, variableName),
		}, nil

	case LITERAL:
		// get literal value
		value, err := stage.operator(nil, nil, nil)
		if err != nil {
			return nil, err
		}
		if value == nil {
			// nil -> null(), because there is no null literal
			return []ExpressionToken{
				newToken(FUNCTION, "null"),
				newToken(CLAUSE, nil),
				newToken(CLAUSE_CLOSE, nil),
			}, nil
		}
		kind := STRING
		switch value.(type) {
		case bool:
			kind = BOOLEAN
		case float64:
			kind = NUMERIC
		}
		return []ExpressionToken{
			newToken(kind, value),
		}, nil

	case NOOP:
		// parentheses
		rightTokens, err := buildTokens(stage.rightStage)
		if err != nil {
			return nil, err
		}
		tokens := []ExpressionToken{
			newToken(CLAUSE, nil),
		}
		tokens = append(tokens, rightTokens...)
		tokens = append(tokens, newToken(CLAUSE_CLOSE, nil))
		return tokens, nil
	}

	leftTokens, err := buildTokens(stage.leftStage)
	if err != nil {
		return nil, err
	}

	rightTokens, err := buildTokens(stage.rightStage)
	if err != nil {
		return nil, err
	}

	tokenString := stage.symbol.String()
	if stage.symbol == EQ {
		tokenString = "==" // for some reason EQ.String() is '='
	}
	kind := UNKNOWN
	if _, found := prefixSymbols[tokenString]; found && stage.leftStage == nil {
		kind = PREFIX
	} else if _, found := modifierSymbols[tokenString]; found {
		kind = MODIFIER
	} else if _, found := logicalSymbols[tokenString]; found {
		kind = LOGICALOP
	} else if _, found := comparatorSymbols[tokenString]; found {
		kind = COMPARATOR
	} else if _, found := ternarySymbols[tokenString]; found {
		kind = TERNARY
	}

	tokens := append(leftTokens, newToken(kind, tokenString))
	tokens = append(tokens, rightTokens...)
	return tokens, nil
}

func newToken(kind TokenKind, value interface{}) ExpressionToken {
	return ExpressionToken{
		Kind:  kind,
		Value: value,
	}
}

type paramCaptor struct {
	lastParameterName string
}

func (p *paramCaptor) Get(key string) (interface{}, error) {
	p.lastParameterName = key
	return 0.0, nil
}
