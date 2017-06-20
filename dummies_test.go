package govaluate

/*
	Struct used to test "parameter calls".
*/
type dummyParameter struct {
	String string
	Int int
	Nested dummyNestedParameter
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

type dummyNestedParameter struct {
	Funk string
}

func (this dummyNestedParameter) Dunk(arg1 string) string {
	return arg1 + "dunk"
}

var fooParameter = EvaluationParameter {
	Name: "foo",
	Value: dummyParameter {
		String: "string!",
		Int: 101,
		Nested: dummyNestedParameter {
			Funk: "funkalicious",
		},
	},
}