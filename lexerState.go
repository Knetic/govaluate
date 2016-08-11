package govaluate

type lexerState struct {
	isEOF          bool
	kind           TokenKind
	validNextKinds []TokenKind
}

// lexer states.
// Constant for all purposes except compiler.
var validLexerStates = []lexerState{

	lexerState{

		kind:  CLAUSE,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},

	lexerState{

		kind:  CLAUSE_CLOSE,
		isEOF: true,
		validNextKinds: []TokenKind{

			COMPARATOR,
			MODIFIER,
			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
			LOGICALOP,
			TERNARY,
			SEPARATOR,
		},
	},

	lexerState{

		kind:  NUMERIC,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:  BOOLEAN,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:  STRING,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:  TIME,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			SEPARATOR,
		},
	},
	lexerState{

		kind:  VARIABLE,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
			TERNARY,
			SEPARATOR,
		},
	},
	lexerState{

		kind:  MODIFIER,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			VARIABLE,
			FUNCTION,
			STRING,
			BOOLEAN,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState{

		kind:  COMPARATOR,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState{

		kind:  LOGICALOP,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState{

		kind:  PREFIX,
		isEOF: false,
		validNextKinds: []TokenKind{

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			FUNCTION,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},

	lexerState{

		kind:  TERNARY,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			STRING,
			TIME,
			VARIABLE,
			FUNCTION,
			CLAUSE,
			SEPARATOR,
		},
	},
	lexerState{

		kind:  FUNCTION,
		isEOF: false,
		validNextKinds: []TokenKind{
			CLAUSE,
		},
	},
	lexerState{

		kind:  SEPARATOR,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			STRING,
			TIME,
			VARIABLE,
			FUNCTION,
			CLAUSE,
		},
	},
}

func (this lexerState) canTransitionTo(kind TokenKind) bool {

	for _, validKind := range this.validNextKinds {

		if validKind == kind {
			return true
		}
	}

	return false
}
