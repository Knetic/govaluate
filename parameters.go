/*
The MIT License (MIT)

Copyright (c) 2014-2016 George Lester

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package govaluate

import (
	"errors"
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

	return value, nil
}
