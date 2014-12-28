govaluation
====

[![Build Status](https://travis-ci.org/Knetic/govaluate.svg?branch=master)](https://travis-ci.org/Knetic/govaluate)
[![Godoc](https://godoc.org/github.com/Knetic/govaluate?status.png)](https://godoc.org/github.com/Knetic/govaluate)


Provides support for evaluating arbitrary artithmetic/string expressions. 

How do I use it?
--

You create a new EvaluableExpression, then call "Evaluate" on it.

	expression, err := govaluate.NewEvaluableExpression("10 > 0");
	result := expression.Evaluate(nil);

	// result is now set to "true", the bool value.

Cool, but how about with parameters?

	expression, err := govaluate.NewEvaluableExpression("foo > 0");

	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = -1;

	result := expression.Evaluate(parameters);
	// result is now set to "false", the bool value.

That's cool, but we can almost certainly have done all that in code. What about a complex use case that involves some math?

	expression, err := govaluate.NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90");

	parameters := make(map[string]interface{}, 8)
	parameters["requests_made"] = 100;
	parameters["requests_succeeded"] = 80;

	result := expression.Evaluate(parameters);
	// result is now set to "false", the bool value.

Or maybe you want to check the status of an alive check ("smoketest") page, which will be a string?

	expression, err := govaluate.NewEvaluableExpression("http_response_body == 'service is ok'");

	parameters := make(map[string]interface{}, 8)
	parameters["http_response_body"] = "service is ok";

	result := expression.Evaluate(parameters);
	// result is now set to "true", the bool value.

These examples have all returned boolean values, but it's equally possible to return numeric ones. 

	expression, err := govaluate.NewEvaluableExpression("total_mem * mem_used / 100");

	parameters := make(map[string]interface{}, 8)
	parameters["total_mem"] = 1024;
	parameters["mem_used"] = 512;

	result := expression.Evaluate(parameters);
	// result is now set to "50.0", the float64 value.

You can also do date parsing, though the formats are somewhat limited. Stick to RF3339, ISO8061, unix date, or ruby date formats. If you're having trouble getting a date string to parse, check the list of formats actually used: [parsing.go:248](https://github.com/Knetic/govaluate/blob/master/parsing.go#L248).

	expression, err := govaluate.NewEvaluableExpression("'2014-01-02' > '2014-01-01 23:59:59'");
	result := expression.Evaluate(nil);

	// result is now set to true


Why can't you just write these expressions in code?
--

Sometimes, you can't know ahead-of-time what an expression looks like. Commonly, you'll have written a monitoring framework which is capable of gathering a bunch of metrics, then evaluating a few expressions to see if any metrics should be alerted upon. Or perhaps you've got a set of data running through your application, and you want to allow your DBA's to run some validations on it before committing it to a database, but neither of you can predict what those validations will be.

A lot of people (myself included, for a long time) wind up writing their own half-baked style of evaluation language that fits their needs, but isn't complete. Or they wind up baking their monitor logic into the actual monitor executable. This library is meant to cover all the normal ALGOL and C-like expressions, so that you don't have to reinvent one of the oldest wheels on a computer.

What operators and types does this support?
--

Modifiers: + - / *

Comparators: > >= < <= == !=

Logical ops: || &&

Numeric constants, including 64-bit floating point (12345)

String constants (single quotes: 'foobar')

Date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date)

Boolean constants: true false

Parenthesis to control order of evaluation

Branching
--

I use green masters, and heavily develop with private feature branches. Full releases are pinned and unchangeable, representing the best available version with the best documentation and test coverage. Master branch, however, should always have all tests pass and implementations considered "working", even if it's just a first pass. Master should never panic.

License
--

This project is licensed under the MIT general use license. You're free to integrate, fork, and play with this code as you feel fit without consulting the author, as long as you provide proper credit to the author in your works. 


Activity
--

If this repository hasn't been updated in a while, it's probably because I don't have any outstanding issues to work on - it's not because I've abandoned the project. If you have questions, issues, or patches; I'm completely open to pull requests, issues opened on github, or emails from out of the blue.
