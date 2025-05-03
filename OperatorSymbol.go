package govaluate

/*
	Represents the valid symbols for operators.

*/
type OperatorSymbol int

const (
	VALUE OperatorSymbol = iota
	LITERAL
	NOOP
	EQ
	NEQ
	GT
	LT
	GTE
	LTE
	REQ
	NREQ
	IN

	AND
	OR

	PLUS
	MINUS
	BITWISE_AND
	BITWISE_OR
	BITWISE_XOR
	BITWISE_LSHIFT
	BITWISE_RSHIFT
	MULTIPLY
	DIVIDE
	MODULUS
	EXPONENT

	NEGATE
	INVERT
	BITWISE_NOT

	TERNARY_TRUE
	TERNARY_FALSE
	COALESCE

	FUNCTIONAL
	ACCESS
	SEPARATE
)

type operatorPrecedence int

const (
	noopPrecedence operatorPrecedence = iota
	valuePrecedence
	functionalPrecedence
	prefixPrecedence
	exponentialPrecedence
	additivePrecedence
	bitwisePrecedence
	bitwiseShiftPrecedence
	multiplicativePrecedence
	comparatorPrecedence
	ternaryPrecedence
	logicalAndPrecedence
	logicalOrPrecedence
	separatePrecedence
)

var precendence = map[OperatorSymbol]operatorPrecedence{
	NOOP:           noopPrecedence,
	VALUE:          valuePrecedence,
	EQ:             comparatorPrecedence,
	NEQ:            comparatorPrecedence,
	GT:             comparatorPrecedence,
	LT:             comparatorPrecedence,
	GTE:            comparatorPrecedence,
	LTE:            comparatorPrecedence,
	REQ:            comparatorPrecedence,
	NREQ:           comparatorPrecedence,
	IN:             comparatorPrecedence,
	AND:            logicalAndPrecedence,
	OR:             logicalOrPrecedence,
	BITWISE_AND:    bitwisePrecedence,
	BITWISE_OR:     bitwisePrecedence,
	BITWISE_XOR:    bitwisePrecedence,
	BITWISE_LSHIFT: bitwiseShiftPrecedence,
	BITWISE_RSHIFT: bitwiseShiftPrecedence,
	PLUS:           additivePrecedence,
	MINUS:          additivePrecedence,
	MULTIPLY:       multiplicativePrecedence,
	DIVIDE:         multiplicativePrecedence,
	MODULUS:        multiplicativePrecedence,
	EXPONENT:       exponentialPrecedence,
	BITWISE_NOT:    prefixPrecedence,
	INVERT:         prefixPrecedence,
	NEGATE:         prefixPrecedence,
	COALESCE:       ternaryPrecedence,
	TERNARY_TRUE:   ternaryPrecedence,
	TERNARY_FALSE:  ternaryPrecedence,
	ACCESS:         functionalPrecedence,
	FUNCTIONAL:     functionalPrecedence,
	SEPARATE:       separatePrecedence,
}

func findOperatorPrecedenceForSymbol(symbol OperatorSymbol) operatorPrecedence {

	precendenceValue, found := precendence[symbol]
	if found {
		return precendenceValue
	}

	return valuePrecedence
}

/*
	Map of all valid comparators, and their string equivalents.
	Used during parsing of expressions to determine if a symbol is, in fact, a comparator.
	Also used during evaluation to determine exactly which comparator is being used.
*/
var comparatorSymbols = map[string]OperatorSymbol{
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	">=": GTE,
	"<":  LT,
	"<=": LTE,
	"=~": REQ,
	"!~": NREQ,
	"in": IN,
}

var logicalSymbols = map[string]OperatorSymbol{
	"&&": AND,
	"||": OR,
}

var bitwiseSymbols = map[string]OperatorSymbol{
	"^": BITWISE_XOR,
	"&": BITWISE_AND,
	"|": BITWISE_OR,
}

var bitwiseShiftSymbols = map[string]OperatorSymbol{
	">>": BITWISE_RSHIFT,
	"<<": BITWISE_LSHIFT,
}

var additiveSymbols = map[string]OperatorSymbol{
	"+": PLUS,
	"-": MINUS,
}

var multiplicativeSymbols = map[string]OperatorSymbol{
	"*": MULTIPLY,
	"/": DIVIDE,
	"%": MODULUS,
}

var exponentialSymbolsS = map[string]OperatorSymbol{
	"**": EXPONENT,
}

var prefixSymbols = map[string]OperatorSymbol{
	"-": NEGATE,
	"!": INVERT,
	"~": BITWISE_NOT,
}

var ternarySymbols = map[string]OperatorSymbol{
	"?":  TERNARY_TRUE,
	":":  TERNARY_FALSE,
	"??": COALESCE,
}

// this is defined separately from additiveSymbols et al because it's needed for parsing, not stage planning.
var modifierSymbols = map[string]OperatorSymbol{
	"+":  PLUS,
	"-":  MINUS,
	"*":  MULTIPLY,
	"/":  DIVIDE,
	"%":  MODULUS,
	"**": EXPONENT,
	"&":  BITWISE_AND,
	"|":  BITWISE_OR,
	"^":  BITWISE_XOR,
	">>": BITWISE_RSHIFT,
	"<<": BITWISE_LSHIFT,
}

var separatorSymbols = map[string]OperatorSymbol{
	",": SEPARATE,
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

var symbolToStringMap = map[OperatorSymbol]string{
	NOOP:           "NOOP",
	VALUE:          "VALUE",
	EQ:             "=",
	NEQ:            "!=",
	GT:             ">",
	LT:             "<",
	GTE:            ">=",
	LTE:            "<=",
	REQ:            "=~",
	NREQ:           "!~",
	AND:            "&&",
	OR:             "||",
	IN:             "in",
	BITWISE_AND:    "&",
	BITWISE_OR:     "|",
	BITWISE_XOR:    "^",
	BITWISE_LSHIFT: "<<",
	BITWISE_RSHIFT: ">>",
	PLUS:           "+",
	MINUS:          "-",
	MULTIPLY:       "*",
	DIVIDE:         "/",
	MODULUS:        "%",
	EXPONENT:       "**",
	NEGATE:         "-",
	INVERT:         "!",
	BITWISE_NOT:    "~",
	TERNARY_TRUE:   "?",
	TERNARY_FALSE:  ":",
	COALESCE:       "??",
}

/*
	Generally used when formatting type check errors.
	We could store the stringified symbol somewhere else and not require a duplicated codeblock to translate
	OperatorSymbol to string, but that would require more memory, and another field somewhere.
	Adding operators is rare enough that we just stringify it here instead.
*/
func (this OperatorSymbol) String() string {

	symbolString, found := symbolToStringMap[this]
	if found {
		return symbolString
	}

	return ""
}
