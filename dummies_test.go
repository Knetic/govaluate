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
	"fmt"
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

func (this *dummyParameter) Func3() string {
	return "fronk"
}

func (this dummyParameter) FuncArgStr(arg1 string) string {
	return arg1
}

func (this dummyParameter) TestArgs(str string, ui uint, ui8 uint8, ui16 uint16, ui32 uint32, ui64 uint64, i int, i8 int8, i16 int16, i32 int32, i64 int64, f32 float32, f64 float64, b bool) string {
	
	var sum float64
	
	sum = float64(ui) + float64(ui8) + float64(ui16) + float64(ui32) + float64(ui64)
	sum += float64(i) + float64(i8) + float64(i16) + float64(i32) + float64(i64)
	sum += float64(f32)

	if b {
		sum += f64
	}

	return fmt.Sprintf("%v: %v", str, sum)
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
