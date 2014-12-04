govaluation
====

[![Build Status](https://travis-ci.org/Knetic/govaluation.svg?branch=master)](https://travis-ci.org/Knetic/govaluation)
[![Godoc](https://godoc.org/github.com/Knetic/govaluation?status.png)](https://godoc.org/github.com/Knetic/govaluation)


Provides support for evaluating arbitrary artithmetic/string expressions. 

Why can't you just write these expressions in code?
--

If you're writing software that needs to accept an evaluable expression from configuration values, you can't just write code for it - you need to parse that expression (a string), figure out what parameters it wants, and then evaluate it. That's what this library does - take a string representing the expression, and a set of parameters to use during evaluation, and allows you to evaluate the expression.

License
--

This implementation of Go named parameter queries is licensed under the MIT general use license. You're free to integrate, fork, and play with this code as you feel fit without consulting the author, as long as you provide proper credit to the author in your works. If you have questions, issues, or patches, I'm completely open to pull requests, issues opened on github, or emails from out of the blue.
