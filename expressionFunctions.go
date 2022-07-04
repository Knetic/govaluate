package govaluate

/*
	Represents a function that can be called from within an expression.
	This method must return an error if, for any reason, it is unable to produce exactly one unambiguous result.
	An error returned will halt execution of the expression.
*/

type Callable interface {
	Call(ctx interface{}, arguments ...interface{}) (interface{}, error)
}

type callable struct {
	fn ExpressionFunction
}

func (c callable) Call(ctx interface{}, arguments ...interface{}) (interface{}, error) {
	return c.fn(arguments...)
}

func NewCallable(fn ExpressionFunction) Callable {
	return &callable{fn: fn}
}

type ExpressionFunction func(arguments ...interface{}) (interface{}, error)
