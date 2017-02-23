package main

import (
	"fmt"

	"github.com/knetic/govaluate"
)

type mystruct struct {
	category string
	values   []interface{}
}

func main() {
	mys := mystruct{}
	mys.category = "cat1"
	mys.values = append(mys.values, "25")
	mys.values = append(mys.values, "val1")
	mys.values = append(mys.values, "val2")
	mys.values = append(mys.values, "val3")

	mymap := make(map[string]interface{})
	mymap[mys.category] = mys.values
	expression, err := govaluate.NewEvaluableExpression("'25' in cat1 && !('val4' in cat1)  && 'al' contains cat1")
	result, err := expression.Evaluate(mymap)
	fmt.Printf("result: %v\nerror: %v\n", result, err)
	// should return
	// result: true
    // error: <nil>
}

