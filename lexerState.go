package govaluate

type lexerState struct {

	isEOF bool
	kind TokenKind
	validNextKinds []TokenKind
}

// lexer states.
// Constant for all purposes except compiler.
var VALID_LEXER_STATES = []lexerState {

	lexerState {

		kind: CLAUSE,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},

	lexerState {

		kind: CLAUSE_CLOSE,
		isEOF: true,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			CLAUSE,			
			CLAUSE_CLOSE,
		},
	},

	lexerState {

		kind: NUMERIC,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	lexerState {

		kind: BOOLEAN,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	lexerState {

		kind: STRING,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	lexerState {

		kind: VARIABLE,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	lexerState {

		kind: MODIFIER,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			VARIABLE,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState {

		kind: COMPARATOR,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	lexerState {

		kind: LOGICALOP,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
}

func (this lexerState) canTransitionTo(kind TokenKind) bool {

	for _, validKind := range this.validNextKinds {

		if(validKind == kind) {
			return true
		}
	}

	return false
}

