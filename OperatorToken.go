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
	REQ
	NREQ

	AND
	OR

	PLUS
	MINUS
	BITWISE_AND
	BITWISE_OR
	BITWISE_XOR
	MULTIPLY
	DIVIDE
	MODULUS
	EXPONENT

	NEGATE
	INVERT
	BITWISE_NOT

	TERNARY_TRUE
	TERNARY_FALSE
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
	"=~": REQ,
	"!~": NREQ,
}

var LOGICAL_SYMBOLS = map[string]OperatorSymbol{

	"&&": AND,
	"||": OR,
}

var MODIFIER_SYMBOLS = map[string]OperatorSymbol{

	"+":  PLUS,
	"-":  MINUS,
	"&":  BITWISE_AND,
	"|":  BITWISE_OR,
	"^":  BITWISE_XOR,
	"*":  MULTIPLY,
	"/":  DIVIDE,
	"%":  MODULUS,
	"**": EXPONENT,
}

var PREFIX_SYMBOLS = map[string]OperatorSymbol{

	"-": NEGATE,
	"!": INVERT,
	"~": BITWISE_NOT,
}

var TERNARY_SYMBOLS = map[string]OperatorSymbol{
	"?": TERNARY_TRUE,
	":": TERNARY_FALSE,
}

var ADDITIVE_MODIFIERS = []OperatorSymbol{
	PLUS, MINUS,
}

var BITWISE_MODIFIERS = []OperatorSymbol{
	BITWISE_AND, BITWISE_OR, BITWISE_XOR,
}

var MULTIPLICATIVE_MODIFIERS = []OperatorSymbol{
	MULTIPLY, DIVIDE, MODULUS,
}

var EXPONENTIAL_MODIFIERS = []OperatorSymbol{
	EXPONENT,
}

var PREFIX_MODIFIERS = []OperatorSymbol{
	NEGATE, INVERT, BITWISE_NOT,
}

var NUMERIC_COMPARATORS = []OperatorSymbol{
	GT, GTE, LT, LTE,
}

var STRING_COMPARATORS = []OperatorSymbol{
	REQ, NREQ,
}

/*
	Returns true if this operator is contained by the given array of candidate symbols.
	False otherwise.
*/
func (this OperatorSymbol) IsModifierType(candidate []OperatorSymbol) bool {

	for _, symbolType := range candidate {
		if this == symbolType {
			return true
		}
	}

	return false
}
