package govaluate

import (
	"unicode/utf8"
)

type lexerStream struct {
	source   string
	position int
	length   int
}

func newLexerStream(source string) *lexerStream {

	var ret *lexerStream

	ret = new(lexerStream)
	ret.source = source
	ret.length = len(source)
	return ret
}

func (this *lexerStream) readCharacter() rune {

	var character rune

	character, _ = utf8.DecodeRuneInString(this.source[this.position:])
	this.position += 1
	return character
}

func (this *lexerStream) rewind(amount int) {
	this.position -= amount
}

func (this lexerStream) canRead() bool {
	return this.position < this.length
}
