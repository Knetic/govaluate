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

// Convenience array that describes all symbols that count as "additive", which is a subset of modifiers that is evaluated last if a sequence of modifiers are used.
var ADDITIVE_MODIFIERS = []OperatorSymbol {
	PLUS, MINUS,
}

// Convenience array that describes all symbols that count as "additive", which is a subset of modifiers that is evaluated second if a sequence of modifiers are used.
var MULTIPLICATIVE_MODIFIERS = []OperatorSymbol {
	MULTIPLY, DIVIDE, MODULUS,
}

// Convenience array that describes all symbols that count as "additive", which is a subset of modifiers that is evaluated first if a sequence of modifiers are used.
var EXPONENTIAL_MODIFIERS = []OperatorSymbol {
	EXPONENT,
}

var PREFIX_MODIFIERS = []OperatorSymbol {
	NEGATE, INVERT,
}

/*
	Returns true if this operator is contained by the given array of candidate symbols.
	False otherwise.
*/
func (this OperatorSymbol) IsModifierType(candidate []OperatorSymbol) bool {

	for _, symbolType := range candidate {
		if(this == symbolType) {
			return true
		}
	}

	return false
}
