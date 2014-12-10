package govaluate

type OperatorSymbol int

const (
	EQ	OperatorSymbol = 1
	NEQ	OperatorSymbol = iota
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
)

// map of all valid symbols
var COMPARATOR_SYMBOLS = map[string]OperatorSymbol {

	"==": EQ,
	"!=": NEQ,
	">": GT,
	">=": GTE,
	"<": LT,
	"<=": LTE,
};

var LOGICAL_SYMBOLS = map[string]OperatorSymbol {

	"&&": AND,
	"||": OR,
};

var MODIFIER_SYMBOLS = map[string]OperatorSymbol {
	
	"+": PLUS,
	"-": MINUS,
	"*": MULTIPLY,
	"/": DIVIDE,
};
