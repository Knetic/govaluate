govaluate
====

[![Build Status](https://travis-ci.org/Knetic/govaluate.svg?branch=master)](https://travis-ci.org/Knetic/govaluate)
[![Godoc](https://godoc.org/github.com/Knetic/govaluate?status.png)](https://godoc.org/github.com/Knetic/govaluate)


Provides support for evaluating arbitrary C-like artithmetic/string expressions.

Why can't you just write these expressions in code?
--

Sometimes, you can't know ahead-of-time what an expression will look like, or you want those expressions to be configurable.
Perhaps you've got a set of data running through your application, and you want to allow your users to specify some validations to run on it before committing it to a database. Or maybe you've written a monitoring framework which is capable of gathering a bunch of metrics, then evaluating a few expressions to see if any metrics should be alerted upon, but the conditions for alerting are different for each monitor.

A lot of people wind up writing their own half-baked style of evaluation language that fits their needs, but isn't complete. Or they wind up baking the expression into the actual executable, even if they know it's subject to change. These strategies may work, but they take time to implement, time for users to learn, and induce technical debt as requirements change. This library is meant to cover all the normal C-like expressions, so that you don't have to reinvent one of the oldest wheels on a computer.

How do I use it?
--

You create a new EvaluableExpression, then call "Evaluate" on it.

```go
	expression, err := govaluate.NewEvaluableExpression("10 > 0");
	result, err := expression.Evaluate(nil);
	// result is now set to "true", the bool value.
```

Cool, but how about with parameters?

```go
	expression, err := govaluate.NewEvaluableExpression("foo > 0");

	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = -1;

	result, err := expression.Evaluate(parameters);
	// result is now set to "false", the bool value.
```

That's cool, but we can almost certainly have done all that in code. What about a complex use case that involves some math?

```go
	expression, err := govaluate.NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90");

	parameters := make(map[string]interface{}, 8)
	parameters["requests_made"] = 100;
	parameters["requests_succeeded"] = 80;

	result, err := expression.Evaluate(parameters);
	// result is now set to "false", the bool value.
```

Or maybe you want to check the status of an alive check ("smoketest") page, which will be a string?

```go
	expression, err := govaluate.NewEvaluableExpression("http_response_body == 'service is ok'");

	parameters := make(map[string]interface{}, 8)
	parameters["http_response_body"] = "service is ok";

	result, err := expression.Evaluate(parameters);
	// result is now set to "true", the bool value.
```

These examples have all returned boolean values, but it's equally possible to return numeric ones.

```go
	expression, err := govaluate.NewEvaluableExpression("(mem_used / total_mem) * 100");

	parameters := make(map[string]interface{}, 8)
	parameters["total_mem"] = 1024;
	parameters["mem_used"] = 512;

	result, err := expression.Evaluate(parameters);
	// result is now set to "50.0", the float64 value.
```

You can also do date parsing, though the formats are somewhat limited. Stick to RF3339, ISO8061, unix date, or ruby date formats. If you're having trouble getting a date string to parse, check the list of formats actually used: [parsing.go:248](https://github.com/Knetic/govaluate/blob/0580e9b47a69125afa0e4ebd1cf93c49eb5a43ec/parsing.go#L258).

```go
	expression, err := govaluate.NewEvaluableExpression("'2014-01-02' > '2014-01-01 23:59:59'");
	result, err := expression.Evaluate(nil);

	// result is now set to true
```

Expressions are parsed once, and can be re-used multiple times. Parsing is the compute-intensive phase of the process, so if you intend to use the same expression with different parameters, just parse it once. Like so;

```go
	expression, err := govaluate.NewEvaluableExpression("response_time <= 100");
	parameters := make(map[string]interface{}, 8)

	for {
		parameters["response_time"] = pingSomething();
		result, err := expression.Evaluate(parameters)
	}
```

The normal C-standard order of operators is respected. When writing an expression, be sure that you either order the operators correctly, or use parenthesis to clarify which portions of an expression should be run first. 

Escaping characters
--

Sometimes you'll have parameters that have spaces, slashes, pluses, ampersands or some other character
that this library interprets as something special. For example, the following expression will not
act as one might expect:

	"response-time < 100"

As written, the library will parse it as "[response] minus [time] is less than 100". In reality,
"response-time" is meant to be one variable that just happens to have a dash in it.

There are two ways to work around this. First, you can escape the entire parameter name:

 	"[response-time] < 100"

Or you can use backslashes to escape only the minus sign.

	"response\\-time < 100"

Backslashes can be used anywhere in an expression to escape the very next character. Square bracketed parameter names can be used instead of plain parameter names at any time.


What operators and types does this support?
--

* Modifiers: `+` `-` `/` `*` `^` `%`
* Comparators: `>` `>=` `<` `<=` `==` `!=` `=~` `!~`
* Logical ops: `||` `&&`
* Numeric constants, as 64-bit floating point (`12345.678`)
* String constants (single quotes: `'foobar'`)
* Date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
* Boolean constants: `true` `false`
* Parenthesis to control order of evaluation `(` `)`
* Prefixes: `!` `-`
* Ternary conditional `?` `:`

Note: for those not familiar, `=~` is "regex-equals" and `!~` is "regex-not-equals".

If a ternary operator resolves to false, it returns nil. So `false ? 10` will return `nil`, whereas `true ? 10` will return `10.0`.

Types
--

Some operators don't make sense when used with some types. For instance, what does it mean to get the modulo of a string? What happens if you check to see if two numbers are logically AND'ed together?

Everyone has a different intuition about the answers to these questions. To prevent confusion, this library will _refuse to operate_ upon types for which there is not an unambiguous meaning for the operation. The table is listed below.

Any time you attempt to use an operator on a type which doesn't explicitly support it (indicated by a bold "X" in the table below), the expression will fail to evaluate, and return an error indicating the problem.

Note that this table shows what each type supports - if you use an operator then _both_ types need to support the operator, otherwise an error will be returned.

|                            	| Number/Date           	| String          	| Boolean         	|
|----------------------------	|-----------------------	|-----------------	|-----------------	|
| +                          	| Adds                  	| Concatenates    	| **X**           	|
| -                          	| Subtracts             	| **X**           	| **X**           	|
| /                          	| Divides               	| **X**           	| **X**           	|
| *                          	| Multiplies            	| **X**           	| **X**           	|
| ^                          	| Takes to the power of 	| **X**           	| **X**           	|
| %                          	| Modulo                	| **X**           	| **X**           	|
| Greater/Lesser (> >= < <=) 	| Valid                 	| **X**           	| **X**           	|
| Equality (== !=)           	| Checks by value       	| Checks by value 	| Checks by value 	|
| Ternary (? :)                 | **X**                     | **X**             | Checks by value   |
| Regex (=~ !~)                 | **X**                     | Regex             | **X**             |
| !                          	| **X**                 	| **X**           	| Inverts         	|
| Negate (-)                 	| Multiplies by -1        	| **X**           	| **X**           	|

It may, at first, not make sense why a Date supports all the same things as a number. In this library, dates are treated as the unix time. That is, the number of seconds since epoch. In practice this means that sub-second precision with this library is impossible (drop an issue in Github if this is a deal-breaker for you). It also, by association, means that you can do operations that you may not expect, like taking a date to the power of two. The author sees no harm in this. Your date probably appreciates it.

Complex types, arrays, and structs are not supported as literals nor parameters. All numeric constants and variables are converted to float64 for evaluation.

Benchmarks
--

If you're concerned about the overhead of this library, a good range of benchmarks are built into this repo. You can run them with `go test -bench=.`. The library is built with an eye towards being quick, but has not been aggressively profiled and optimized. For most applications, though, it is completely fine.

For a very rough idea of performance, here are the results output from a benchmark run on my 3rd-gen Macbook Pro (Linux Mint 17.1).

```
BenchmarkSingleParse-12                          2000000               768 ns/op
BenchmarkSimpleParse-12                           200000              6842 ns/op
BenchmarkFullParse-12                             200000             12791 ns/op
BenchmarkEvaluationSingle-12                    10000000               142 ns/op
BenchmarkEvaluationNumericLiteral-12             3000000               577 ns/op
BenchmarkEvaluationLiteralModifiers-12           2000000               675 ns/op
BenchmarkEvaluationParameters-12                 2000000               883 ns/op
BenchmarkEvaluationParametersModifiers-12        1000000              1305 ns/op
BenchmarkComplexExpression-12                    1000000              1308 ns/op
BenchmarkRegexExpression-12                       100000             22751 ns/op
BenchmarkConstantRegexExpression-12               500000              2599 ns/op
ok
```

Branching
--

I use green masters, and heavily develop with private feature branches. Full releases are pinned and unchangeable, representing the best available version with the best documentation and test coverage. Master branch, however, should always have all tests pass and implementations considered "working", even if it's just a first pass. Master should never panic.

License
--

This project is licensed under the MIT general use license. You're free to integrate, fork, and play with this code as you feel fit without consulting the author, as long as you provide proper credit to the author in your works.


Activity
--

If this repository hasn't been updated in a while, it's probably because I don't have any outstanding issues to work on - it's not because I've abandoned the project. If you have questions, issues, or patches; I'm completely open to pull requests, issues opened on github, or emails from out of the blue.
