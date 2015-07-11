package govaluate

/*
  Tests to make sure evaluation fails in the expected ways.
*/
import (
  "testing"
)

type DebugStruct struct {
  x int
}

func TestComplexParameter(test *testing.T) {

  expression, _ := NewEvaluableExpression("1")
  parameters := map[string]interface{} {
    "foo": 1i,
  }

  _, err := expression.Evaluate(parameters)
  if(err == nil) {
    test.Logf("Expected to  fail when giving a complex value, did not")
    test.Fail()
  }
}

func TestStructParameter(test *testing.T) {

  expression, _ := NewEvaluableExpression("1")
  parameters := map[string]interface{} {
    "foo": DebugStruct{},
  }

  _, err := expression.Evaluate(parameters)
  if(err == nil) {
    test.Logf("Expected to  fail when giving a struct value, did not")
    test.Fail()
  }
}
