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
var Keywords  = [25]string {
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}

func main()  {
	teste := "case switch"

	l := lexer.New(teste, Root)

	l.Start()

	for tok, done := l.NextToken(); !done; tok,done = l.NextToken() {
		fmt.Printf("%q - %q", tok.Value, tok.Type)
	}
	//for r := l.Next(); r != lexer.EOFRune ; {
	//	fmt.Println(string(r))
	//	r = l.Next()
	//}

}

func Root(l *lexer.L) lexer.StateFunc {
	//r := l.Next()
	fmt.Print("Passa 0")
	for unicode.IsSpace(l.Next()) {
		l.Ignore()
	}
	l.Rewind()

	r := l.Peek()
	fmt.Print("Passa 1")
	fmt.Println("r : " + string(r))
	if unicode.IsLetter(l.Peek()) {
		return CheckKeyword(l)
	}
	fmt.Print("Passa 2")


	if r == '_' {
		return Identifier(l)
		//} else if unicode.IsSymbol(r) || unicode.IsPunct(r) {
		//	return OpAndPunc(l)
		//} else if unicode.IsDigit(r) || r == '.'  {
		//	return Numbers(l)
		//} else if r == '\'' {
		//	return Runes(l)
		//} else if r == '"' {
		//	return Strings(l)
	} else if r != lexer.EOFRune {
		l.Error(fmt.Sprintf("Símbolo desconhecido %q", r))
	}

	return nil
}

func CheckKeyword(l *lexer.L) lexer.StateFunc {
	r := l.Next()
	fmt.Println("r: " + string(r))
	//kw := Keywords
	//fmt.Print(Keywords)
	for _, word := range Keywords {
		counter := 0
		for _, lt := range word {
			//fmt.Printf("%q = %q\n", string(r), string(lt))
			if lt != r {
				break
			} else {
				counter++
				r = l.Next()
			}
		}
		if unicode.IsSpace(r) || r == lexer.EOFRune {
			l.Emit(KeywordToken)
			fmt.Println("Match: " + word + "\n")
			return Root(l)
		}
		if counter > 0 {
			for  ;counter >= 0; counter-- {
				l.Rewind()
			}
			r = l.Next()
		}
	}

	return Identifier(l)
}

func Identifier(l *lexer.L) lexer.StateFunc {
	r := l.Next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		fmt.Println("ID: " + string(r))
		r = l.Next()
	}
	l.Rewind()

	l.Emit(IdentifierToken)
	return Root(l)
}
