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


func main()  {
	teste := .25

}

func Root(l *lexer.L) lexer.StateFunc {
	for unicode.IsSpace(l.Peek()) {
		l.Ignore()
	}
	if unicode.IsLetter(l.Peek()) {
		return CheckKeyword(l)
	}

	r := l.Next()
	 if r == '_' {
		return Identifier(l)
	} else if unicode.IsSymbol(r) || unicode.IsPunct(r) {
		return OpAndPunc(l)
	} else if unicode.IsDigit(r) || r == '.'  {
		return Numbers(l)
	} else if r == '\'' {
		return Runes(l)
	} else if r == '"' {
		return Strings(l)
	} else {
		l.Error(fmt.Sprintf("Símboolo desconhecido %q", r))
	}

	return nil
}

func CheckKeyword(l *lexer.L) lexer.StateFunc {
	var Keywords  = [25]string {
		"break", "case", "chan", "const", "continue",
		"default", "defer", "else", "fallthrough", "for",
		"func", "go", "goto", "if", "import",
		"interface", "map", "package", "range", "return",
		"select", "struct", "switch", "type", "var",
	}
	r := l.Next()
	for i, word := range Keywords {
		for j, c := range word {
			if r != c {
				break
			} else {
				r = l.Next()
			}
			if l.Current() == word && unicode.IsSpace(l.Peek()) {
				l.Emit(KeywordToken)
				return Root(l)
			} else {

			}
		}
	}
}
