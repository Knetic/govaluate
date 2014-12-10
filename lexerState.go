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
		},
	},

	lexerState {

		kind: NUMERIC,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: BOOLEAN,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: STRING,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: VARIABLE,
		isEOF: true,
		validNextKinds: []TokenKind {

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
		},
	},
	lexerState {

		kind: MODIFIER,
		isEOF: false,
		validNextKinds: []TokenKind {

			NUMERIC,
			VARIABLE,
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

