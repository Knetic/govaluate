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
)
