package govaluate

/*
	Represents the valid symbols for operators.

*/
type OperatorSymbol int

const (
	EQ OperatorSymbol = iota
	NEQ
	GT
	LT
	GTE
	LTE

	AND
	OR

	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULUS
	EXPONENT

	NEGATE
	INVERT
)

/*
	Map of all valid comparators, and their string equivalents.
	Used during parsing of expressions to determine if a symbol is, in fact, a comparator.
	Also used during evaluation to determine exactly which comparator is being used.
*/
var COMPARATOR_SYMBOLS = map[string]OperatorSymbol{

	"==": EQ,
	"!=": NEQ,
	">":  GT,
	">=": GTE,
	"<":  LT,
	"<=": LTE,
}

/*
	Map of all valid logical operators, and their string equivalents.
	Used during parsing of expressions to determine if a symbol is, in fact, a logical operator.
	Also used during evaluation to determine exactly which logical operator is being used.
*/
var LOGICAL_SYMBOLS = map[string]OperatorSymbol{

	"&&": AND,
	"||": OR,
}

/*
	Map of all valid modifiers, and their string equivalents.
	Used during parsing of expressions to determine if a symbol is, in fact, a modifier.
	Also used during evaluation to determine exactly which modifier is being used.
*/
var MODIFIER_SYMBOLS = map[string]OperatorSymbol{

	"+": PLUS,
	"-": MINUS,
	"*": MULTIPLY,
	"/": DIVIDE,
	"%": MODULUS,
	"^": EXPONENT,
}

var PREFIX_SYMBOLS = map[string]OperatorSymbol{

	"-": NEGATE,
	"!": INVERT,
}
