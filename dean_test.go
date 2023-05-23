package govaluate

import (
	"fmt"
	"log"
	"testing"
)

var counter = 0

func decorate(stage *evaluationStage) {
}

func TestDecoratorCollector(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	exp, err := NewEvaluableExpression("foo > bar")
	if err != nil {
		t.Error("couldn't parse expression: ", err)
	}

	fmt.Println(exp.evaluationStages.symbol)
	decorate(exp.evaluationStages)

	params := make(MapParameters)
	params["foo"] = 1
	params["bar"] = 2
	res, err := exp.Eval(params)

	if err != nil {
		t.Error(err)
	}
	log.Println(exp, "==", res)
}
