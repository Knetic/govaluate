package govaluate

import (
	"encoding/json"
	"strconv"
)

// sanitizedParameters is a wrapper for Parameters that does sanitization as
// parameters are accessed.
type sanitizedParameters struct {
	orig Parameters
}

func (p sanitizedParameters) Get(key string) (interface{}, error) {
	value, err := p.orig.Get(key)
	if err != nil {
		return nil, err
	}

	return castToNumber(value), nil
}

func castToNumber(value interface{}) interface{} {
	switch value.(type) {
	case uint8:
		return json.Number(strconv.FormatInt(int64(value.(uint8)), 10))
	case uint16:
		return json.Number(strconv.FormatInt(int64(value.(uint16)), 10))
	case uint32:
		return json.Number(strconv.FormatInt(int64(value.(uint32)), 10))
	case uint64:
		return json.Number(strconv.FormatInt(int64(value.(uint64)), 10))
	case uint:
		return json.Number(strconv.FormatInt(int64(value.(uint)), 10))
	case int8:
		return json.Number(strconv.FormatInt(int64(value.(int8)), 10))
	case int16:
		return json.Number(strconv.FormatInt(int64(value.(int16)), 10))
	case int32:
		return json.Number(strconv.FormatInt(int64(value.(int32)), 10))
	case int64:
		return json.Number(strconv.FormatInt(value.(int64), 10))
	case int:
		return json.Number(strconv.FormatInt(int64(value.(int)), 10))
	case float32:
		return json.Number(strconv.FormatFloat(float64(value.(float32)), 'f', 10, 64))
	case float64:
		return json.Number(strconv.FormatFloat(value.(float64), 'f', 10, 64))
	}

	return value
}
