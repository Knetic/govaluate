package govaluate


/*
	Represents all valid types of tokens that a token can be.
*/
type TokenKind int


const (

	UNKNOWN TokenKind = iota

	NUMERIC 
	BOOLEAN
	STRING
	VARIABLE

	COMPARATOR
	LOGICALOP
	MODIFIER

	CLAUSE
	CLAUSE_CLOSE
)

/*
	GetTokenKindString returns a string that describes the given TokenKind.
	e.g., when passed the NUMERIC TokenKind, this returns the string "NUMERIC".
*/
func GetTokenKindString(kind TokenKind) string {

	switch(kind) {

		case	NUMERIC		:	return "NUMERIC";
		case	BOOLEAN		:	return "BOOLEAN";
		case	STRING		:	return "STRING";
		case	VARIABLE	:	return "VARIABLE";
		case	COMPARATOR	:	return "COMPARATOR";
		case	LOGICALOP	:	return "LOGICALOP";
		case	MODIFIER	:	return "MODIFIER";
		case	CLAUSE		:	return "CLAUSE";
		case	CLAUSE_CLOSE	:	return "CLAUSE_CLOSE";
	}

	return "UNKNOWN";
}
