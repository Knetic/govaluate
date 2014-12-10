package govaluate

type EvaluableExpression struct {

	Tokens []ExpressionToken
	inputExpression string
}


func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	var ret *EvaluableExpression;
	var err error

	ret = new(EvaluableExpression)
	ret.inputExpression = expression;
	ret.Tokens, err = parseTokens(expression)

	if(err != nil) {
		return nil, err
	}
	return ret, nil
}

func (this EvaluableExpression) Evaluate(parameters map[string]interface{}) interface{} {

	return false
}

func (this EvaluableExpression) String() string {

	return this.inputExpression;
}
