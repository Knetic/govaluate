package govaluate

import (
  "testing"
)

/*
  The most common use case, a single variable, modified slightly, compared to a constant.
  This is the "expected" use case of govaluate.
*/
func BenchmarkSimpleParse(bench *testing.B) {

  for i := 0; i < bench.N; i++ {
    NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90")
  }
}

/*
  Benchmarks all syntax possibilities in one expression.
*/
func BenchmarkFullParse(bench *testing.B) {

  var expression string

  // represents all the major syntax possibilities.
  expression = "2 > 1 &&" +
               "'something != nothing || " +
               "'2014-01-20' < 'Wed Jul  8 23:07:35 MDT 2015' &&" +
               "[escapedVariable name with spaces] <= unescaped\\-variableName &&" +
               "modifierTest + 1000 / 2 > (80 * 100 % 2)"

  for i := 0; i < bench.N; i++ {
    NewEvaluableExpression(expression)
  }
}

/*
  Benchmarks evaluation times of literals (no variables, no modifiers)
*/
func BenchmarkEvaluationNumericLiteral(bench *testing.B) {

  expression, _ := NewEvaluableExpression("2 > 1")

  for i := 0; i < bench.N; i++ {
    expression.Evaluate(nil)
  }
}

/*
  Benchmarks evaluation times of literals with modifiers
*/
func BenchmarkEvaluationLiteralModifiers(bench *testing.B) {

  expression, _ := NewEvaluableExpression("2 + 2 == 4")

  for i := 0; i < bench.N; i++ {
    expression.Evaluate(nil)
  }
}

/*
  Benchmarks evaluation times of parameters + literals with modifiers
*/
func BenchmarkEvaluationParametersModifiers(bench *testing.B) {

  expression, _ := NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90")
  parameters := map[string]interface{} {
    "requests_made": 99,
    "requests_succeeded": 90,
  }

  for i := 0; i < bench.N; i++ {
      expression.Evaluate(parameters)
  }
}
