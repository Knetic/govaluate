package govaluate

import (
	"errors"
)

/*
	Struct used to test "parameter calls".
*/
type dummyParameter struct {
	String    string
	Int       int
	BoolFalse bool
	Nil       interface{}
	Nested    dummyNestedParameter
}

func (this dummyParameter) Func() string {
	return "funk"
}

func (this dummyParameter) Func2() (string, error) {
	return "frink", nil
}

func (this dummyParameter) FuncArgStr(arg1 string) string {
	return arg1
}

func (this dummyParameter) AlwaysFail() (interface{}, error) {
	return nil, errors.New("function should always fail")
}

type dummyNestedParameter struct {
	Funk string
}

func (this dummyNestedParameter) Dunk(arg1 string) string {
	return arg1 + "dunk"
}

var dummyParameterInstance = dummyParameter{
	String:    "string!",
	Int:       101,
	BoolFalse: false,
	Nil:       nil,
	Nested: dummyNestedParameter{
		Funk: "funkalicious",
	},
}

var fooParameter = EvaluationParameter{
	Name:  "foo",
	Value: dummyParameterInstance,
}

var fooPtrParameter = EvaluationParameter{
	Name:  "fooptr",
	Value: &dummyParameterInstance,
}

var fooFailureParameters = map[string]interface{}{
	"foo":    fooParameter.Value,
	"fooptr": &fooPtrParameter.Value,
}
