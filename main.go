package main

//Bibliotecas importadas
import (
	"fmt"
	"github.com/bbuck/go-lexer"
	"strings"
	"unicode"
	"os"
	"io/ioutil"
)

// Tokens
const (
	_ = iota
	IdentifierToken
	KeywordToken
	OpAndPuncToken
	IntLiteralToken
	FloatLiteralToken
	ImagLiteralToken
	RuneLiteralToken
	StringLiteralToken
)

//Palavras Chave da linguagem
var Keywords = [25]string{
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}

//Operadores e pontuações da linguagem
var OpAndPuncs = [47]string{
	"+", "&", "+=", "&=", "&&", "==", "!=", "(", ")",
	"-", "|", "-=", "|=", "||", "<", "<=", "[", "]",
	"*", "^", "*=", "^=", "<-", ">", ">=", "{", "}",
	"/", "<<", "/=", "<<=", "++", "=", ":=", ",", ";",
	"%", ">>", "%=", ">>=", "--", "!", "...", ".", ":",
	"&^", "&^=",
}

//Associação do token com um nome legível, em string
var Tokens = map[int]string{
	IdentifierToken: "Identificador",
	KeywordToken: "Palavra Chave",
	OpAndPuncToken: "Operador/Pontuação",
	IntLiteralToken: "Literal Int",
	FloatLiteralToken: "Literal Float",
	ImagLiteralToken: "Literal Imaginário",
	RuneLiteralToken: "Literal Rune",
	StringLiteralToken: "Literal String",
}

//Função para identificar um dígito hexadecimal
func isHex(r rune) bool {
	return strings.ContainsRune("0123456789abcdfABCDF", r)
}

//Função para identificar um dígito octal
func isOctal(r rune) bool {
	return strings.ContainsRune("01234567", r)
}

//Função para identificar se o caractér passado está dentro dos operadores e pontuações da lnguagem
func isOpAndPunc(r rune) bool {
	return strings.ContainsRune("+-*/%&|^<>=!:.([{,)]};", r)
}

func main() {
	code := os.Args[1]
	f, _ := os.Open(code)
	txt, _ := ioutil.ReadAll(f)
	//fmt.Print(string(txt))

	l := lexer.New(string(txt), Root)

	l.Start()

	//Loop para retirar tokens da cadeia de tokens
	for tok, done := l.NextToken(); !done; tok, done = l.NextToken() {
		fmt.Printf("%q - %q\n", tok.Value, Tokens[int(tok.Type)])
	}
}

//Função raiz, para análise do código
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
	if r == '/' {
		if l.Peek() == '/' {
			return Comment(l)
		} else {
			l.Error(fmt.Sprintf("Simbolo inesperado: %q", r))
			return nil
		}

	} else if r == '_' {
		return Identifier(l)
	} else if r == '\'' {
		return Runes(l)
	} else if r == '"' {
		return Strings(l)
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
	}  else if r != lexer.EOFRune {
		l.Error(fmt.Sprintf("Símbolo desconhecido %q", r))
	}

	return nil
}

//Função para checar se a palavra analisada é uma palavra-chave
func CheckKeyword(l *lexer.L) lexer.StateFunc {
	//fmt.Println("Entrou CK")
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

//Função para checar se a palavra analisada é um identificador
func Identifier(l *lexer.L) lexer.StateFunc {
	//fmt.Println("Entrou Id")
	r := l.Next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		//fmt.Println("ID: " + string(r))
		r = l.Next()
	}
	l.Rewind()

	l.Emit(IdentifierToken)
	return Root(l)
}

//Função para checar se a palavra analisada é um operador ou pontuação
func OpAndPunc(l *lexer.L) lexer.StateFunc {
	var counter int
	for isOpAndPunc(l.Next()) {
		counter++
	}
	l.Rewind()
	for ; counter > 0; counter-- {
		for _, smb := range OpAndPuncs {
			if l.Current() == smb {
				l.Emit(OpAndPuncToken)
				return Root(l)
			}
		}
		l.Rewind()
	}
	l.Next()
	l.Error(fmt.Sprintf("Operador/Pontuação desconhecido: %q", l.Current()))
	return nil
}

//Função para checar se a palavra analisada é um literal numérico
func Numbers(l *lexer.L) lexer.StateFunc {
	//fmt.Println("Entrou Numb")
	l.Take("0123456789")

	r := l.Peek()
	if !unicode.IsLetter(r) || r == lexer.EOFRune {
		l.Emit(IntLiteralToken)
		return Root(l)
	} else if r == '.' {
		return Float(l)
	} else if r == 'i' {
		l.Next()
		r = l.Peek()
		if !unicode.IsLetter(r) || r == lexer.EOFRune {
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

//Função para checar se a palavra analisada é um literal numérico flutuante
func Float(l *lexer.L) lexer.StateFunc {
	//fmt.Println("Entrou float")
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

//Função para checar se a palavra analisada é um literal rune (ou caractér)
func Runes(l *lexer.L) lexer.StateFunc {
	//fmt.Println("Entrou Runes")
	l.Next()
	r := l.Next()
	if r == '\\' {
		r = l.Next()
		if strings.ContainsRune("abcfnrtv\\'", r) {
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
			for ok && r != '\'' {
				ok = isHex(r)
				r = l.Next()
				cont++
			}
			if ok && cont == 8 {
				l.Emit(RuneLiteralToken)
				return Root(l)
			}
		} else {
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

//Função para checar se a palavra analisada é um literal string válido
func Strings(l *lexer.L) lexer.StateFunc {
	//fmt.Println("Entrou string")
	r := l.Next()
	if r == '`' {
		for ok := true; ok; r = l.Next() {
			if r == '`' {
				ok = false
			} else if r == lexer.EOFRune {
				l.Error(fmt.Sprintf("` Esperado."))
				return nil
			}
		}
		l.Emit(StringLiteralToken)
		return Root(l)
	} else {
		r = l.Next()
		for ; r != '"' && r != lexer.EOFRune; r = l.Next() {
			for unicode.IsSpace(r) {
				if r == '\n' {
					l.Error(fmt.Sprintf("Quebra de linha inesperada: %q", l.Current()))
					return nil
				}
				r = l.Next()
			}
			if r == '\\' {
				r = l.Next()
				if strings.ContainsRune("abcfnrtv\\\"", r) {
					continue
				} else if isOctal(r) {
					if isOctal(l.Next()) && isOctal(l.Next()) {
						continue
					}
				} else if r == 'x' {
					r = l.Next()
					if isHex(r) && isHex(l.Next()) {
						continue
					}
				} else if r == 'u' {
					r = l.Next()
					if isHex(r) && isHex(l.Next()) && isHex(l.Next()) {
						continue
					}
				} else if r == 'U' {
					r = l.Next()
					var cont int
					ok := true
					for ok && !unicode.IsSpace(r) && r != '"' {
						ok = isHex(r)
						r = l.Next()
						cont++
					}
					if ok && cont == 8 {
						continue
					}
				} else {
					l.Error(fmt.Sprintf("Literal string inválido: %q", l.Current()))
					return nil
				}
			}
		}
		if r == lexer.EOFRune {
			l.Error(fmt.Sprintf("\" Esperado."))
			return nil
		}
	}
	l.Emit(StringLiteralToken)
	return Root(l)
}

//Função para ignorar comentários
func Comment(l *lexer.L) lexer.StateFunc {
	l.Rewind()
	for l.Next() != '\n' {
		l.Ignore()
	}
	l.Ignore()
	return Root(l)
}