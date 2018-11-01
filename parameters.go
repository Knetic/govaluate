package govaluate

import (
	"errors"
	"reflect"
)

/*
	Parameters is a collection of named parameters that can be used by an EvaluableExpression to retrieve parameters
	when an expression tries to use them.
*/
type Parameters interface {

	/*
		Get gets the parameter of the given name, or an error if the parameter is unavailable.
		Failure to find the given parameter should be indicated by returning an error.
	*/
	Get(name string) (interface{}, error)
}

type MapParameters map[string]interface{}

func (p MapParameters) Get(name string) (interface{}, error) {

	value, found := p[name]

	if !found {
		errorMessage := "No parameter '" + name + "' found."
		return nil, errors.New(errorMessage)
	}
	
	if value != nil {
	    s := reflect.ValueOf(value)
	    if s.Kind() == reflect.Slice {
		ret := []interface{}{}
		for i := 0; i < s.Len(); i++ {
			ret = append(ret, s.Index(i).Interface())
		}
		return ret, nil
	    }
	}

	return value, nil
}
