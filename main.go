package main

import (
	"github.com/bbuck/go-lexer"
	"log"
	"fmt"
	"unicode"
	"unicode/utf8"
)

const (
	_  = iota
	StringToken
	IntegerToken
)

func main()  {
	teste := `123456789`

	l := lexer.New(teste, NumberState)
	l.Start()

	tok, done := l.NextToken()
	if done {
		log.Print("Falhou")
	}
	fmt.Println(tok.Type)
}

func Program(L *lexer.L) lexer.StateFunc {
	return ProgramHeading(L)
}

func ProgramHeading(L *lexer.L) lexer.StateFunc {
	if unicode.IsLetter(L.Peek()) {
		return Identifier(L)
	} else {
		return nil
	}
}

func NumberState(L *lexer.L) lexer.StateFunc {
	L.Take("0123456789")
	L.Emit(IntegerToken)
	return nil
}

func Identifier(L *lexer.L) lexer.StateFunc {

}