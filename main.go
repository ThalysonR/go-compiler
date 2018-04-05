package main

import (
	"github.com/bbuck/go-lexer"
	"fmt"
	"unicode"
	"strings"
)

const (
	_  = iota
	IdentifierToken
	KeywordToken
	OpAndPuncToken
	IntLiteralToken
	FloatLiteralToken
	ImagLiteralToken
	RuneLiteralToken
	StringLiteralToken

)
var Keywords  = [25]string {
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}
var OpAndPuncs = [47]string {
	"+",    "&",     "+=",    "&=",     "&&",    "==",    "!=",    "(",    ")",
	"-",    "|",     "-=",    "|=",     "||",    "<",     "<=",    "[",    "]",
	"*",    "^",     "*=",    "^=",     "<-",    ">",     ">=",    "{",    "}",
	"/",    "<<",    "/=",    "<<=",    "++",    "=",     ":=",    ",",    ";",
    "%",    ">>",    "%=",    ">>=",    "--",    "!",     "...",   ".",    ":",
	"&^",          "&^=",
}

func isHex(r rune) bool {
	return strings.ContainsRune("0123456789abcdfABCDF", r)
}

func isOctal(r rune) bool {
	return strings.ContainsRune("01234567", r)
}

func main()  {
	teste_1 := "    case && switcheroo .10254512i '\U99999999'"

	l := lexer.New(teste_1, Root)

	l.Start()

	for tok, done := l.NextToken(); !done; tok,done = l.NextToken() {
		fmt.Printf("%q - %q\n", tok.Value, tok.Type)
	}
	//for r := l.Next(); r != lexer.EOFRune ; {
	//	fmt.Println(string(r))
	//	r = l.Next()
	//}

}

func Root(l *lexer.L) lexer.StateFunc {
	for unicode.IsSpace(l.Next()) {
		l.Ignore()
	}
	l.Rewind()

	r := l.Peek()
	//fmt.Println("r: " + string(r) + "\n")
	if unicode.IsLetter(r) {
		return CheckKeyword(l)
	}


	if r == '_' {
		return Identifier(l)
		} else if r == '\'' {
		return Runes(l)
		} else if unicode.IsSymbol(r) || unicode.IsPunct(r) {
			if r == '.' {
				l.Next()
				if unicode.IsDigit(l.Peek()) {
					l.Rewind()
					return Float(l)
				}
				l.Rewind()
			}
			return OpAndPunc(l)
		} else if unicode.IsDigit(r) {
			return Numbers(l)
		} else if r == '.' {
			return Float(l)
		//} else if r == '"' {
		//	return Strings(l)
	} else if r != lexer.EOFRune {
		l.Error(fmt.Sprintf("Símbolo desconhecido %q", r))
	}

	return nil
}

func CheckKeyword(l *lexer.L) lexer.StateFunc {
	r := l.Next()
	for unicode.IsLetter(r) {
		r = l.Next()
	}
	l.Rewind()
	for _, word := range Keywords {
		if l.Current() == word {
			l.Emit(KeywordToken)
			return Root(l)
		}
	}


	return Identifier(l)
}

func Identifier(l *lexer.L) lexer.StateFunc {
	r := l.Next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		fmt.Println("ID: " + string(r))
		r = l.Next()
	}
	l.Rewind()

	l.Emit(IdentifierToken)
	return Root(l)
}

func OpAndPunc(l *lexer.L) lexer.StateFunc {
	l.Take("+-*/%&|^<>=!:.([{,)]};")
	for _, smb := range OpAndPuncs {
		if l.Current() == smb {
			l.Emit(OpAndPuncToken)
			return Root(l)
		}
	}
	l.Error(fmt.Sprintf("Operador/Pontuação desconhecido: %q", l.Current()))
	return nil
}

func Numbers(l *lexer.L) lexer.StateFunc {
	l.Take("0123456789")

	r := l.Peek()
	if unicode.IsSpace(r) || r == lexer.EOFRune {
		l.Emit(IntLiteralToken)
		return Root(l)
	} else if r == '.' {
		return Float(l)
	} else if r == 'i' {
		l.Next()
		r = l.Peek()
		if unicode.IsSpace(r) || r == lexer.EOFRune {
			l.Emit(ImagLiteralToken)
			return Root(l)
		} else {
			l.Error(fmt.Sprintf("Caractér inesperado: %q", r))
		}
	} else {
		l.Error(fmt.Sprintf("Caractér inesperado: %q", r))
	}
	return nil
}

func Float(l *lexer.L) lexer.StateFunc {
	l.Next()
	l.Take("0123456789")

	r := l.Peek()
	if unicode.IsSpace(r) || r == lexer.EOFRune {
		l.Emit(FloatLiteralToken)
		return Root(l)
	} else if r == 'e' || r == 'E' { // Bloco Exponencial
		l.Next()
		r = l.Peek()
		if r == '+' || r == '-' {
			l.Next()
		}
		l.Take("0123456789")
		r = l.Peek()
		if unicode.IsSpace(r) || r == lexer.EOFRune {
			l.Emit(FloatLiteralToken)
			return Root(l)
		}
	}

	if r == 'i' { // Bloco Imaginário
		l.Next()
		r = l.Peek()
		if unicode.IsSpace(r) || r == lexer.EOFRune {
			l.Emit(ImagLiteralToken)
			return Root(l)
		}
	}
	l.Next()
	l.Error(fmt.Sprintf("Literal numérico inválido: %q", l.Current()))
	return nil
}

func Runes(l *lexer.L) lexer.StateFunc {
	l.Next()
	r := l.Next()
	if r == '\\' {
		r = l.Next()
		if r == ('a'|'b'|'f'|'n'|'r'|'t'|'v'|'\\'|'\'') {
			r = l.Next()
			if r == '\'' {
				l.Emit(RuneLiteralToken)
				return Root(l)
			}
		} else if r == 'x' {
			r = l.Next()
			if isHex(r) && isHex(l.Next()) && l.Next() == '\'' {
				l.Emit(RuneLiteralToken)
				return Root(l)
			}
		} else if r == 'u' {
			r = l.Next()
			if isHex(r) && isHex(l.Next()) && isHex(l.Next()) && l.Next() == '\'' {
				l.Emit(RuneLiteralToken)
				return Root(l)
			}
		} else if r == 'U' {
			r = l.Next()
			var cont int
			ok := true
			for ; ok && r != '\'';  {
				ok = isHex(r)
				r = l.Next()
				cont++
			}
			if ok && cont == 8 {
				l.Emit(RuneLiteralToken)
				return Root(l)
			}
		}  else {
			l.Error(fmt.Sprintf("Literal rune inválido: %q", l.Current()))
		}
	} else {
		r = l.Next()
		if r == '\'' {
			l.Emit(RuneLiteralToken)
			return Root(l)
		} else {
			l.Error(fmt.Sprintf("Literal rune inválido: %q", l.Current()))
		}
	}
	return nil
}

func String(l *lexer.L) lexer.StateFunc {
	r := l.Next()
	if r == '`' {
		for ok := true; ok ; r = l.Next() {
			if r == '`' {
				ok = false
			}
		}
		l.Emit(StringLiteralToken)
		return Root(l)
	}
}