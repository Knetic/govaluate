package govaluate

import "math"

const (
	EFMax   = "Max"
	EFMin   = "Min"
	EFRound = "Round"
)

func BuiltinFunctions() map[string]ExpressionFunction {
	return map[string]ExpressionFunction{
		EFMax:   ExpFuncMax,
		EFMin:   ExpFuncMin,
		EFRound: ExpFuncRound,
	}
}

func ExpFuncMax(parameters Parameters, args ...interface{}) (interface{}, error) {
	val := 0.0
	for _, a := range args {
		if a.(float64) > val {
			val = a.(float64)
		}
	}
	return val, nil
}

func ExpFuncMin(parameters Parameters, args ...interface{}) (interface{}, error) {
	val := math.MaxFloat64
	for _, a := range args {
		if a.(float64) < val {
			val = a.(float64)
		}
	}
	return val, nil
}

func ExpFuncRound(parameters Parameters, args ...interface{}) (interface{}, error) {
	return math.Round(args[0].(float64)), nil
}
