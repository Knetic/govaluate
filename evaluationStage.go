package govaluate

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	logicalErrorFormat    string = "Value '%v' cannot be used with the logical operator '%v', it is not a bool"
	modifierErrorFormat   string = "Value '%v' cannot be used with the modifier '%v', it is not a number"
	comparatorErrorFormat string = "Value '%v' cannot be used with the comparator '%v', it is not a number"
	ternaryErrorFormat    string = "Value '%v' cannot be used with the ternary operator '%v', it is not a bool"
	prefixErrorFormat     string = "Value '%v' cannot be used with the prefix '%v'"
)

type evaluationOperator func(left interface{}, right interface{}, parameters Parameters) (interface{}, error)
type stageTypeCheck func(value interface{}) bool
type stageCombinedTypeCheck func(left interface{}, right interface{}) bool

type evaluationStage struct {
	symbol OperatorSymbol

	leftStage, rightStage *evaluationStage

	// the operation that will be used to evaluate this stage (such as adding [left] to [right] and return the result)
	operator evaluationOperator

	// ensures that both left and right values are appropriate for this stage. Returns an error if they aren't operable.
	leftTypeCheck  stageTypeCheck
	rightTypeCheck stageTypeCheck

	// if specified, will override whatever is used in "leftTypeCheck" and "rightTypeCheck".
	// primarily used for specific operators that don't care which side a given type is on, but still requires one side to be of a given type
	// (like string concat)
	typeCheck stageCombinedTypeCheck

	// regardless of which type check is used, this string format will be used as the error message for type errors
	typeErrorFormat string
}

var (
	_true  = interface{}(true)
	_false = interface{}(false)
)

func (this *evaluationStage) swapWith(other *evaluationStage) {

	temp := *other
	other.setToNonStage(*this)
	this.setToNonStage(temp)
}

func (this *evaluationStage) setToNonStage(other evaluationStage) {

	this.symbol = other.symbol
	this.operator = other.operator
	this.leftTypeCheck = other.leftTypeCheck
	this.rightTypeCheck = other.rightTypeCheck
	this.typeCheck = other.typeCheck
	this.typeErrorFormat = other.typeErrorFormat
}

func (this *evaluationStage) isShortCircuitable() bool {

	switch this.symbol {
	case AND:
		fallthrough
	case OR:
		fallthrough
	case TERNARY_TRUE:
		fallthrough
	case TERNARY_FALSE:
		fallthrough
	case COALESCE:
		return true
	}

	return false
}

func noopStageRight(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	return right, nil
}

func addStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {

	// string concat if either are strings
	if isString(left) || isString(right) {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(leftNum+rightNum, 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(leftNum+rightNum, 10)), nil
}
func subtractStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(leftNum-rightNum, 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(leftNum-rightNum, 10)), nil
}
func multiplyStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(leftNum*rightNum, 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(leftNum*rightNum, 10)), nil
}
func divideStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(leftNum/rightNum, 'f', 10, 64)), nil
	}
	leftNum, err := left.(json.Number).Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := right.(json.Number).Int64()
	if err != nil {
		return json.Number("0"), err
	}
	res := float64(leftNum) / float64(rightNum)
	if float64(int64(res)) == res {
		return json.Number(strconv.FormatInt(int64(res), 10)), nil
	}
	return json.Number(strconv.FormatFloat(res, 'f', 10, 64)), nil
}
func exponentStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(math.Pow(leftNum, rightNum), 'f', 10, 64)), nil
	}
	leftNum, err := left.(json.Number).Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := right.(json.Number).Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(int64(math.Pow(float64(leftNum), float64(rightNum))), 10)), nil
}

func modulusStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(math.Mod(leftNum, rightNum), 'f', 10, 64)), nil
	}
	leftNum, err := left.(json.Number).Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := right.(json.Number).Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(int64(math.Mod(float64(leftNum), float64(rightNum))), 10)), nil
}

func gteStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) >= right.(string)), nil
	}
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		return boolIface(leftNum > rightNum), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	return boolIface(leftNum >= rightNum), nil
}
func gtStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) > right.(string)), nil
	}
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		return boolIface(leftNum > rightNum), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	return boolIface(leftNum > rightNum), nil
}
func lteStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) <= right.(string)), nil
	}
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		return boolIface(leftNum <= rightNum), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	return boolIface(leftNum <= rightNum), nil
}
func ltStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) < right.(string)), nil
	}
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return boolIface(false), err
		}
		return boolIface(leftNum < rightNum), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return boolIface(false), err
	}
	return boolIface(leftNum < rightNum), nil
}
func equalStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	return boolIface(reflect.DeepEqual(left, right)), nil
}
func notEqualStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	return boolIface(!reflect.DeepEqual(left, right)), nil
}
func andStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	return boolIface(left.(bool) && right.(bool)), nil
}
func orStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	return boolIface(left.(bool) || right.(bool)), nil
}
func negateStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	rightNumStr := right.(json.Number).String()
	if len(rightNumStr) <= 0 {
		return json.Number("0"), errors.New("empty string")
	} else if rightNumStr[0] == '-' {
		return json.Number(rightNumStr[1:]), nil
	}
	return json.Number(fmt.Sprintf("-%s", rightNumStr)), nil
}
func invertStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	return boolIface(!right.(bool)), nil
}
func bitwiseNotStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	rightNumStr := right.(json.Number)
	if strings.Contains(rightNumStr.String(), ".") {
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(float64(^int64(rightNum)), 'f', 10, 64)), nil
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(^rightNum, 10)), nil
}
func ternaryIfStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	if left.(bool) {
		return right, nil
	}
	return nil, nil
}
func ternaryElseStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	if left != nil {
		return left, nil
	}
	return right, nil
}

func regexStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {

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

func notRegexStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {

	ret, err := regexStage(left, right, parameters)
	if err != nil {
		return nil, err
	}

	return !(ret.(bool)), nil
}

func bitwiseOrStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(float64(int64(leftNum)|int64(rightNum)), 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(leftNum|rightNum, 10)), nil
}
func bitwiseAndStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(float64(int64(leftNum)&int64(rightNum)), 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(leftNum&rightNum, 10)), nil
}
func bitwiseXORStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(float64(int64(leftNum)^int64(rightNum)), 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(leftNum^rightNum, 10)), nil
}
func leftShiftStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(float64(uint64(leftNum)<<uint64(rightNum)), 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(int64(uint64(leftNum)<<uint64(rightNum)), 10)), nil
}
func rightShiftStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
	leftNumStr, rightNumStr := left.(json.Number), right.(json.Number)
	if strings.Contains(leftNumStr.String(), ".") || strings.Contains(rightNumStr.String(), ".") {
		leftNum, err := leftNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		rightNum, err := rightNumStr.Float64()
		if err != nil {
			return json.Number("0"), err
		}
		return json.Number(strconv.FormatFloat(float64(uint64(leftNum)>>uint64(rightNum)), 'f', 10, 64)), nil
	}
	leftNum, err := leftNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	rightNum, err := rightNumStr.Int64()
	if err != nil {
		return json.Number("0"), err
	}
	return json.Number(strconv.FormatInt(int64(uint64(leftNum)>>uint64(rightNum)), 10)), nil
}

func makeParameterStage(parameterName string) evaluationOperator {

	return func(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
		value, err := parameters.Get(parameterName)
		if err != nil {
			return nil, err
		}

		return value, nil
	}
}

func makeLiteralStage(literal interface{}) evaluationOperator {
	return func(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {
		return literal, nil
	}
}

func makeFunctionStage(function ExpressionFunction) evaluationOperator {

	return func(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {

		if right == nil {
			val, err := function()
			return castToNumber(val), err
		}
		var (
			val interface{}
			err error
		)

		switch right.(type) {
		case []interface{}:
			val, err = function(right.([]interface{})...)
		default:
			val, err = function(right)
		}
		return castToNumber(val), err
	}
}

func typeConvertParam(p reflect.Value, t reflect.Type) (ret reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("Argument type conversion failed: failed to convert '%s' to '%s'", p.Kind().String(), t.Kind().String())
			err = errors.New(errorMsg)
			ret = p
		}
	}()

	return p.Convert(t), nil
}

func typeConvertParams(method reflect.Value, params []reflect.Value) ([]reflect.Value, error) {

	methodType := method.Type()
	numIn := methodType.NumIn()
	numParams := len(params)

	if numIn != numParams {
		if numIn > numParams {
			return nil, fmt.Errorf("Too few arguments to parameter call: got %d arguments, expected %d", len(params), numIn)
		}
		return nil, fmt.Errorf("Too many arguments to parameter call: got %d arguments, expected %d", len(params), numIn)
	}

	for i := 0; i < numIn; i++ {
		t := methodType.In(i)
		p := params[i]
		pt := p.Type()

		if t.Kind() != pt.Kind() {

			np, err := typeConvertParam(p, t)
			if err != nil {
				return nil, err
			}
			params[i] = np
		}
	}

	return params, nil
}

func makeAccessorStage(pair []string) evaluationOperator {

	reconstructed := strings.Join(pair, ".")

	return func(left interface{}, right interface{}, parameters Parameters) (ret interface{}, err error) {

		var params []reflect.Value

		value, err := parameters.Get(pair[0])
		if err != nil {
			return nil, err
		}

		// while this library generally tries to handle panic-inducing cases on its own,
		// accessors are a sticky case which have a lot of possible ways to fail.
		// therefore every call to an accessor sets up a defer that tries to recover from panics, converting them to errors.
		defer func() {
			if r := recover(); r != nil {
				errorMsg := fmt.Sprintf("Failed to access '%s': %v", reconstructed, r.(string))
				err = errors.New(errorMsg)
				ret = nil
			}
		}()

		for i := 1; i < len(pair); i++ {

			coreValue := reflect.ValueOf(value)

			var corePtrVal reflect.Value

			// if this is a pointer, resolve it.
			if coreValue.Kind() == reflect.Ptr {
				corePtrVal = coreValue
				coreValue = coreValue.Elem()
			}

			if coreValue.Kind() != reflect.Struct {
				return nil, errors.New("Unable to access '" + pair[i] + "', '" + pair[i-1] + "' is not a struct")
			}

			field := coreValue.FieldByName(pair[i])
			if field != (reflect.Value{}) {
				value = field.Interface()
				continue
			}

			method := coreValue.MethodByName(pair[i])
			if method == (reflect.Value{}) {
				if corePtrVal.IsValid() {
					method = corePtrVal.MethodByName(pair[i])
				}
				if method == (reflect.Value{}) {
					return nil, errors.New("No method or field '" + pair[i] + "' present on parameter '" + pair[i-1] + "'")
				}
			}

			switch right.(type) {
			case []interface{}:

				givenParams := right.([]interface{})
				params = make([]reflect.Value, len(givenParams))
				for idx, p := range givenParams {
					val, ok := p.(json.Number)
					if !ok {
						params[idx] = reflect.ValueOf(givenParams[idx])
						continue
					}
					valStr := val.String()
					if strings.Contains(valStr, ".") {
						valFloat, numErr := val.Float64()
						if numErr != nil {
							return nil, numErr
						}
						params[idx] = reflect.ValueOf(valFloat)
					} else {
						valInt, numErr := val.Int64()
						if numErr != nil {
							return nil, numErr
						}
						params[idx] = reflect.ValueOf(valInt)
					}
				}

			default:

				if right == nil {
					params = []reflect.Value{}
					break
				}

				params = []reflect.Value{reflect.ValueOf(right.(interface{}))}
			}

			params, err = typeConvertParams(method, params)

			if err != nil {
				return nil, errors.New("Method call failed - '" + pair[0] + "." + pair[1] + "': " + err.Error())
			}

			returned := method.Call(params)
			retLength := len(returned)

			if retLength == 0 {
				return nil, errors.New("Method call '" + pair[i-1] + "." + pair[i] + "' did not return any values.")
			}

			if retLength == 1 {

				value = returned[0].Interface()
				continue
			}

			if retLength == 2 {

				errIface := returned[1].Interface()
				err, validType := errIface.(error)

				if validType && errIface != nil {
					return returned[0].Interface(), err
				}

				value = returned[0].Interface()
				continue
			}

			return nil, errors.New("Method call '" + pair[0] + "." + pair[1] + "' did not return either one value, or a value and an error. Cannot interpret meaning.")
		}

		value = castToNumber(value)
		return value, nil
	}
}

func separatorStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {

	var ret []interface{}

	switch left.(type) {
	case []interface{}:
		ret = append(left.([]interface{}), right)
	default:
		ret = []interface{}{left, right}
	}

	return ret, nil
}

func inStage(left interface{}, right interface{}, parameters Parameters) (interface{}, error) {

	for _, value := range right.([]interface{}) {
		if left == value {
			return true, nil
		}
	}
	return false, nil
}

//

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

func isNumber(value interface{}) bool {
	switch value.(type) {
	case json.Number:
		return true
	}
	return false
}

/*
	Addition usually means between numbers, but can also mean string concat.
	String concat needs one (or both) of the sides to be a string.
*/
func additionTypeCheck(left interface{}, right interface{}) bool {

	if isNumber(left) && isNumber(right) {
		return true
	}
	if !isString(left) && !isString(right) {
		return false
	}
	return true
}

/*
	Comparison can either be between numbers, or lexicographic between two strings,
	but never between the two.
*/
func comparatorTypeCheck(left interface{}, right interface{}) bool {

	if isNumber(left) && isNumber(right) {
		return true
	}
	if isString(left) && isString(right) {
		return true
	}
	return false
}

func isArray(value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		return true
	}
	return false
}

/*
	Converting a boolean to an interface{} requires an allocation.
	We can use interned bools to avoid this cost.
*/
func boolIface(b bool) interface{} {
	if b {
		return _true
	}
	return _false
}
