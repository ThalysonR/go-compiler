package main

import (
	"github.com/bbuck/go-lexer"
	"fmt"
	//"unicode"
	"unicode"
)

const (
	_  = iota
	IdentifierToken
	KeywordToken
	OpAndPuncToken
	LiteralToken

)
var Keywords  = [...]string {
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}

func main()  {
	teste := .25

}

func Root(l *lexer.L) lexer.StateFunc {
	for unicode.IsSpace(l.Peek()) {
		l.Ignore()
	}
	r := l.Next()

	if unicode.IsLetter(r) {
		CheckKeyword(l)
	} else if r == '_' {
		Identifier(l)
	} else if unicode.IsSymbol(r) || unicode.IsPunct(r) {
		OpAndPunc(l)
	} else if unicode.IsDigit(r) || r == '.'  {
		Numbers(l)
	} else if r == '\'' {
		Runes(l)
	} else if r == '"' {
		Strings(l)
	} else {
		l.Error(fmt.Sprintf("Símboolo desconhecido %q", r))
	}

	return nil
}


