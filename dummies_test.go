package govaluate

import (
	"errors"
	"fmt"
	"net"
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
	IP1       net.IP
	IP2       net.IP
	CIDR1     net.IPNet
	CIDR2     net.IPNet
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
	IP1:   net.ParseIP("127.0.0.1"),
	IP2:   net.ParseIP("127.0.0.3"),
	CIDR1: mustParseCIDR("127.0.0.4/22"),
	CIDR2: mustParseCIDR("27.0.0.0/12"),
}

var fooParameter = EvaluationParameter{
	Name:  "foo",
	Value: dummyParameterInstance,
}

var fooParameterEmptyIP = EvaluationParameter{
	Name: "foo",
	Value: dummyParameter{
		IP1:   nil,
		IP2:   nil,
		CIDR1: mustParseCIDR("127.0.0.4/22"),
		CIDR2: mustParseCIDR("27.0.0.0/12"),
	},
}

var fooPtrParameter = EvaluationParameter{
	Name:  "fooptr",
	Value: &dummyParameterInstance,
}

var fooFailureParameters = map[string]interface{}{
	"foo":    fooParameter.Value,
	"fooptr": &fooPtrParameter.Value,
}

func mustParseCIDR(cidr string) net.IPNet {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return *ipNet

}

var CIDRTestFunction = map[string]ExpressionFunction{
	"InNetwork": func(args ...interface{}) (interface{}, error) {
		ip, ok1 := args[0].(net.IP)
		ipNet, ok2 := args[1].(net.IPNet)

		if !ok1 {
			return nil, fmt.Errorf("variable %s not ip", args[0])
		}
		if !ok2 {
			return nil, fmt.Errorf("variable %s not IPnet its %T", args[1], args[1])
		}

		return ipNet.Contains(ip), nil
	},
}
