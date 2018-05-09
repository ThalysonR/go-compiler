package main

//Bibliotecas importadas
import (
	"os"
	"io/ioutil"
	"github.com/thalysonr/go-compiler/Lex"
	"github.com/thalysonr/go-compiler/Synth"
	"fmt"
)

//Associação do token com um nome legível, em string
var Tokens = map[int]string{
	Lex.IdentifierToken: "Identificador",
	Lex.KeywordToken: "Palavra Chave",
	Lex.OpAndPuncToken: "Operador/Pontuação",
	Lex.IntLiteralToken: "Literal Int",
	Lex.FloatLiteralToken: "Literal Float",
	Lex.ImagLiteralToken: "Literal Imaginário",
	Lex.RuneLiteralToken: "Literal Rune",
	Lex.StringLiteralToken: "Literal String",
}

var Declarations = map[int]string {
	Synth.PackageClause: "Package Clause",
	Synth.ImportDeclaration: "Import Declaration",
}

func main() {
	code := os.Args[1]
	f, _ := os.Open(code)
	txt, _ := ioutil.ReadAll(f)
	l := Lex.StartSync(string(txt))
	//Loop para retirar tokens da cadeia de tokens
	//for tok, done := l.NextToken(); !done; tok, done = l.NextToken() {
	//	fmt.Printf("%q - %q; Linha: %d - Coluna: %d\n", tok.Value, Tokens[int(tok.Type)], tok.Linha, tok.Coluna)
	//}
	s := Synth.Start(l)
	for dec, done := s.NextDeclaration(); !done; dec, done = s.NextDeclaration() {
		fmt.Printf("%q - %q\n", dec.Value, Declarations[int(dec.Type)])
	}

}