package govaluate

import (
	"errors"
	"fmt"
	"reflect"
)

// Parameters is a collection of named parameters that are accessible via the
// Get method.
type Parameters interface {
	// Get gets the parameter of the given name
	Get(name string) (interface{}, error)
}

// MapParameters is an implementation of the Parameters interface using a map.
type MapParameters map[string]interface{}

// Get implemetns the method from Parameters
func (p MapParameters) Get(name string) (interface{}, error) {
	return p[name], nil
}

type SanitizedParameters struct {
	orig Parameters
}

func (p SanitizedParameters) Get(key string) (interface{}, error) {
	value, err := p.orig.Get(key)
	if err != nil {
		return nil, err
	}
	// make sure that the parameter is a valid type.
	err = checkValidType(key, value)
	if err != nil {
		return nil, err
	}

	// should be converted to fixed point?
	if isFixedPoint(value) {
		return castFixedPoint(value), nil
	}

	return value, nil
}

func checkValidType(key string, value interface{}) error {

	switch value.(type) {
	case complex64:
		errorMsg := fmt.Sprintf("Parameter '%s' is a complex64 integer, which is not evaluable", key)
		return errors.New(errorMsg)
	case complex128:
		errorMsg := fmt.Sprintf("Parameter '%s' is a complex128 integer, which is not evaluable", key)
		return errors.New(errorMsg)
	}

	if reflect.ValueOf(value).Kind() == reflect.Struct {
		errorMsg := fmt.Sprintf("Parameter '%s' is a struct, which is not evaluable", key)
		return errors.New(errorMsg)
	}

	return nil
}

func isFixedPoint(value interface{}) bool {

	switch value.(type) {
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case int:
		return true
	}
	return false
}

func castFixedPoint(value interface{}) float64 {
	switch value.(type) {
	case uint8:
		return float64(value.(uint8))
	case uint16:
		return float64(value.(uint16))
	case uint32:
		return float64(value.(uint32))
	case uint64:
		return float64(value.(uint64))
	case int8:
		return float64(value.(int8))
	case int16:
		return float64(value.(int16))
	case int32:
		return float64(value.(int32))
	case int64:
		return float64(value.(int64))
	case int:
		return float64(value.(int))
	}

	return 0.0
}

func isString(value interface{}) bool {

	switch value.(type) {
	case string:
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
