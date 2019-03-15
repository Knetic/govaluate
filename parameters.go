package govaluate

import (
	"errors"
	"strings"
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

	parts := strings.Split(name, ".")
	var value interface{}
	current := p
	for i, p := range parts {
		var found bool
		value, found = current[p]

		if !found {
			errorMessage := "No parameter '" + name + "' found."
			return nil, errors.New(errorMessage)
		}

		if i != len(parts) - 1 {
			var ok bool
			current, ok = value.(map[string]interface{})
			if !ok {
				errorMessage := "No parameter '" + name + "' found."
				return nil, errors.New(errorMessage)
			}

		}
	}

	return value, nil
}
