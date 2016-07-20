package govaluate

import (
	"fmt"
	"errors"
	"regexp"
	"math"
)

const (
	TYPEERROR_LOGICAL 		string = "Value '%v' cannot be used with the logical operator '%v', it is not a bool"
	TYPEERROR_MODIFIER 		string = "Value '%v' cannot be used with the modifier '%v', it is not a number"
	TYPEERROR_COMPARATOR 	string = "Value '%v' cannot be used with the comparator '%v', it is not a number"
	TYPEERROR_TERNARY 		string = "Value '%v' cannot be used with the ternary operator '%v', it is not a bool"
)

type evaluationOperator func(left interface{}, right interface{}) (interface{}, error)
type stageTypeCheck func(value interface{}) bool

type evaluationStage struct {

	rightStage *evaluationStage

	// the operation that will be used to evaluate this stage (such as adding [left] to [right] and return the result)
	operator evaluationOperator

	// ensures that both left and right values are appropriate for this stage. Returns an error if they aren't operable.
	leftTypeCheck stageTypeCheck
	rightTypeCheck stageTypeCheck
	typeErrorFormat string
}

func addStage(left interface{}, right interface{}) (interface{}, error) {

	// string concat if either are strings
	if isString(left) || isString(right) {
		return fmt.Sprintf("%v%v", left, right), nil
	}

	return left.(float64) + right.(float64), nil
}
func subtractStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) - right.(float64), nil
}
func multiplyStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) * right.(float64), nil
}
func divideStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) / right.(float64), nil
}
func exponentStage(left interface{}, right interface{}) (interface{}, error) {
	return math.Pow(left.(float64), right.(float64)), nil
}
func modulusStage(left interface{}, right interface{}) (interface{}, error) {
	return math.Mod(left.(float64), right.(float64)), nil
}
func gteStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) >= right.(float64), nil
}
func gtStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) > right.(float64), nil
}
func lteStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) >= right.(float64), nil
}
func ltStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(float64) > right.(float64), nil
}
func andStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(bool) && right.(bool), nil
}
func orStage(left interface{}, right interface{}) (interface{}, error) {
	return left.(bool) || right.(bool), nil
}
func negateStage(left interface{}, right interface{}) (interface{}, error) {
	return -left.(float64), nil
}
func invertStage(left interface{}, right interface{}) (interface{}, error) {
	return !left.(bool), nil
}
func ternaryIfStage(left interface{}, right interface{}) (interface{}, error) {
	if(left.(bool)) {
		return right, nil
	}
	return nil, nil
}
func ternaryElseStage(left interface{}, right interface{}) (interface{}, error) {
	if(left == nil) {
		return right, nil
	}
	return left, nil
}

func regexStage(left interface{}, right interface{}) (interface{}, error) {

	var pattern *regexp.Regexp
	var err error

	switch right.(type) {
	case string:
		pattern, err = regexp.Compile(right.(string))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to compile regexp pattern '%v': %v", right, err))
		}
	case *regexp.Regexp:
		pattern = right.(*regexp.Regexp)
	}

	return pattern.Match([]byte(left.(string))), nil
}

func regexNotStage(left interface{}, right interface{}) (interface{}, error) {

	ret, err := regexStage(left, right)
	if(err != nil) {
		return nil, err
	}

	return !(ret.(bool)), nil
}

func valueStage(left interface{}, right interface{}) (interface{}, error) {
	return left, nil
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
