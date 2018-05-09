package Lex

//Bibliotecas importadas
import (
	"fmt"
	"strings"
	"unicode"
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

var Linha = 0
var Coluna = 0

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

func Start(txt string) *L {
	l := New(txt, Root)
	l.Start()
	return l
}

func StartSync(txt string) *L {
	l := New(txt, Root)
	l.StartSync()
	return l
}

//Função raiz, para análise do código
func Root(l *L) StateFunc {
	r := l.Next()
	for unicode.IsSpace(r) {
		l.Ignore()
		Coluna++
		if r == '\n' {
			Linha++
			Coluna = 0
		}
		r = l.Next()
	}
	l.Rewind()
	//r = l.Peek()
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
	}  else if r != EOFRune {
		l.Error(fmt.Sprintf("Símbolo desconhecido %q", r))
	}

	return nil
}

//Função para checar se a palavra analisada é uma palavra-chave
func CheckKeyword(l *L) StateFunc {
	//fmt.Println("Entrou CK")
	r := l.Next()
	Coluna++
	for unicode.IsLetter(r) {
		r = l.Next()
		Coluna++
	}
	l.Rewind()
	Coluna--
	for _, word := range Keywords {
		if l.Current() == word {
			l.Emit(KeywordToken, Linha, Coluna)
			return Root(l)
		}
	}

	return Identifier(l)
}

//Função para checar se a palavra analisada é um identificador
func Identifier(l *L) StateFunc {
	//fmt.Println("Entrou Id")
	r := l.Next()
	Coluna++

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		//fmt.Println("ID: " + string(r))
		r = l.Next()
		Coluna++
	}
	l.Rewind()
	Coluna--

	l.Emit(IdentifierToken, Linha, Coluna)
	return Root(l)
}

//Função para checar se a palavra analisada é um operador ou pontuação
func OpAndPunc(l *L) StateFunc {
	var counter int
	for isOpAndPunc(l.Next()) {
		counter++
		Coluna++
	}
	l.Rewind()

	for ; counter > 0; counter-- {
		for _, smb := range OpAndPuncs {
			if l.Current() == smb {
				l.Emit(OpAndPuncToken, Linha, Coluna)
				return Root(l)
			}
		}
		l.Rewind()
		Coluna--
	}
	l.Next()
	Coluna++
	l.Error(fmt.Sprintf("Operador/Pontuação desconhecido: %q", l.Current()))
	return nil
}

//Função para checar se a palavra analisada é um literal numérico
func Numbers(l *L) StateFunc {
	//fmt.Println("Entrou Numb")
	//l.Take("0123456789")
	r := l.Next()
	Coluna++
	for strings.ContainsRune("0123456789", r) {
		r = l.Next()
		Coluna++
	}
	//r = l.Peek()
	l.Rewind()
	if !unicode.IsLetter(r) || r == EOFRune {
		l.Emit(IntLiteralToken, Linha, Coluna)
		return Root(l)
	} else if r == '.' {
		return Float(l)
	} else if r == 'i' {
		l.Next()
		Coluna++
		r = l.Peek()
		if !unicode.IsLetter(r) || r == EOFRune {
			l.Emit(ImagLiteralToken, Linha, Coluna)
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
func Float(l *L) StateFunc {
	//fmt.Println("Entrou float")
	l.Next()
	//l.Take("0123456789")
	r := l.Next()
	Coluna++
	for strings.ContainsRune("0123456789", r) {
		r = l.Next()
		Coluna++
	}
	l.Rewind()
	//r = l.Peek()
	if unicode.IsSpace(r) || r == EOFRune {
		l.Emit(FloatLiteralToken, Linha, Coluna)
		return Root(l)
	} else if r == 'e' || r == 'E' { // Bloco Exponencial
		l.Next()
		Coluna++
		r = l.Peek()
		if r == '+' || r == '-' {
			l.Next()
			Coluna++
		}
		//l.Take("0123456789")
		//r = l.Peek()
		r := l.Next()
		Coluna++
		for strings.ContainsRune("0123456789", r) {
			r = l.Next()
			Coluna++
		}
		l.Rewind()
		if unicode.IsSpace(r) || r == EOFRune {
			l.Emit(FloatLiteralToken, Linha, Coluna)
			return Root(l)
		}
	}

	if r == 'i' { // Bloco Imaginário
		l.Next()
		Coluna++
		r = l.Peek()
		if unicode.IsSpace(r) || r == EOFRune {
			l.Emit(ImagLiteralToken, Linha, Coluna)
			return Root(l)
		}
	}
	l.Next()
	l.Error(fmt.Sprintf("Literal numérico inválido: %q", l.Current()))
	return nil
}

//Função para checar se a palavra analisada é um literal rune (ou caractér)
func Runes(l *L) StateFunc {
	//fmt.Println("Entrou Runes")
	l.Next()
	Coluna++
	r := l.Next()
	Coluna++
	if r == '\\' {
		r = l.Next()
		Coluna++
		if strings.ContainsRune("abcfnrtv\\'", r) {
			r = l.Next()
			Coluna++
			if r == '\'' {
				l.Emit(RuneLiteralToken, Linha, Coluna)
				return Root(l)
			}
		} else if r == 'x' {
			r = l.Next()
			Coluna++
			if isHex(r) && isHex(l.Next()) && l.Next() == '\'' {
				l.Emit(RuneLiteralToken, Linha, Coluna)
				return Root(l)
			}
		} else if r == 'u' {
			r = l.Next()
			Coluna++
			if isHex(r) && isHex(l.Next()) && isHex(l.Next()) && l.Next() == '\'' {
				l.Emit(RuneLiteralToken, Linha, Coluna)
				return Root(l)
			}
		} else if r == 'U' {
			r = l.Next()
			Coluna++
			var cont int
			ok := true
			for ok && r != '\'' {
				ok = isHex(r)
				r = l.Next()
				Coluna++
				cont++
			}
			if ok && cont == 8 {
				l.Emit(RuneLiteralToken, Linha, Coluna)
				return Root(l)
			}
		} else {
			l.Error(fmt.Sprintf("Literal rune inválido: %q", l.Current()))
		}
	} else {
		r = l.Next()
		Coluna++
		if r == '\'' {
			l.Emit(RuneLiteralToken, Linha, Coluna)
			return Root(l)
		} else {
			l.Error(fmt.Sprintf("Literal rune inválido: %q", l.Current()))
		}
	}
	return nil
}

//Função para checar se a palavra analisada é um literal string válido
func Strings(l *L) StateFunc {
	//fmt.Println("Entrou string")
	r := l.Next()
	Coluna++
	if r == '`' {
		for ok := true; ok; r = l.Next() {
			if r == '`' {
				ok = false
			} else if r == EOFRune {
				l.Error(fmt.Sprintf("` Esperado."))
				return nil
			}
		}
		l.Emit(StringLiteralToken, Linha, Coluna)
		return Root(l)
	} else {
		r = l.Next()
		Coluna++
		for ; r != '"' && r != EOFRune; r = l.Next() {
			Coluna++
			for unicode.IsSpace(r) {
				if r == '\n' {
					l.Error(fmt.Sprintf("Quebra de linha inesperada: %q", l.Current()))
					return nil
				}
				r = l.Next()
				Coluna++
			}
			if r == '\\' {
				r = l.Next()
				Coluna++
				if strings.ContainsRune("abcfnrtv\\\"", r) {
					continue
				} else if isOctal(r) {
					if isOctal(l.Next()) && isOctal(l.Next()) {
						Coluna += 2
						continue
					}
				} else if r == 'x' {
					r = l.Next()
					Coluna++
					if isHex(r) && isHex(l.Next()) {
						Coluna++
						continue
					}
				} else if r == 'u' {
					r = l.Next()
					Coluna++
					if isHex(r) && isHex(l.Next()) && isHex(l.Next()) {
						Coluna += 2
						continue
					}
				} else if r == 'U' {
					r = l.Next()
					Coluna++
					var cont int
					ok := true
					for ok && !unicode.IsSpace(r) && r != '"' {
						ok = isHex(r)
						r = l.Next()
						Coluna++
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
		if r == EOFRune {
			l.Error(fmt.Sprintf("\" Esperado."))
			return nil
		}
	}
	l.Emit(StringLiteralToken, Linha, Coluna)
	return Root(l)
}

//Função para ignorar comentários
func Comment(l *L) StateFunc {
	l.Rewind()
	for l.Next() != '\n' {
		l.Ignore()
	}
	//l.Ignore()
	return Root(l)
}