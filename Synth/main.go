package Synth

import (
	"github.com/thalysonr/go-compiler/Lex"
	"fmt"
)

const (
	_ = iota
	PackageClause
	ImportDeclaration
	ConstDeclaration
)

func Start(l *Lex.L) *S {
	s := New(l, Package)
	s.Start()
	return s
}

func StartSync(l *Lex.L) *S {
	s := New(l, Package)
	s.StartSync()
	return s
}

func Package(s *S) StateFunc {
	d := s.Next()
	if n := s.Next(); d.Value == "package" && d.Type == Lex.KeywordToken && n.Type == Lex.IdentifierToken && n.Value != "_" {
		s.Emit(PackageClause)
		return Root(s)
	}
	return nil
}

func Root(s *S) StateFunc {
	d := s.Peek()
	if d.Type == Lex.KeywordToken {
		switch d.Value {
		case "import":
			return ImportDecl(s)
		default:
			return nil
		}
	}
	return nil
}

func ImportDecl(s *S) StateFunc {
	d := s.Next()
	if d.Type == Lex.KeywordToken && d.Value == "import" {
		d = s.Peek()
		if d.Value == "(" {
			d = s.Next()
			for d.Value != ")" {
				if !ImportSpec(s) {
					s.Error(fmt.Sprintf("Import Spec esperado. Linha %d; Valor - %q", d.Linha, d.Value))
					return nil
				}
				d = s.Peek()
			}
			s.Next()
		} else {
			if !ImportSpec(s) {
				s.Error(fmt.Sprintf("Import Spec esperado. Linha %d", d.Linha))
				return nil
			}
		}
		s.Emit(ImportDeclaration)
	}
	return Root(s)
}

func ImportSpec(s *S) bool {
	d := s.Next()
	if d.Value == "." || d.Type == Lex.IdentifierToken {
		d = s.Next()
	}
	if d.Type == Lex.StringLiteralToken {
		return true
	} else {
		return false
	}
}

func TopLevelDecl(s *S) StateFunc {
	d := s.Peek()
	if d.Value == "func" {
		s.Next()
		d = s.Peek()
		s.Rewind()
		if d.Type == Lex.IdentifierToken {
			return FunctionDecl(s)
		} else {
			return MethodDecl(s)
		}
	} else {
		return Decl(s)
	}
}

func Decl(s *S) StateFunc {
	d := s.Next()
	if d.Type == Lex.KeywordToken {
		switch d.Value {
		case "const":
			return ConstDecl(s)
		case "type":
			return TypeDecl(s)
		case "var":
			return VarDecl(s)
		default:
			s.Error(fmt.Sprintf("Declaração esperada. Linha %d", d.Linha))
			return nil
		}
	} else {
		s.Error(fmt.Sprintf("Declaração esperada. Linha %d", d.Linha))
	}
}

func FunctionDecl(s *S) StateFunc {
	return nil
}

func MethodDecl(s *S) StateFunc {
	return nil
}'
'

func ConstDecl(s *S) StateFunc {
	d := s.Peek()
	if d.Value == "(" {
		d = s.Next()
		for d.Value != ")" {
			if !ConstSpec(s) {
				s.Error(fmt.Sprintf("Const Spec esperado. Linha %d; Valor - %q", d.Linha, d.Value))
				return nil
			}
			d = s.Peek()
		}
		s.Next()
	} else {
		if !ConstSpec(s) {
			s.Error(fmt.Sprintf("Const Spec esperado. Linha %d", d.Linha))
			return nil
		}
	}
	s.Emit(ConstDeclaration)
	return Root(s)
}

func ConstSpec(s *S) bool {
	if IdentifierList(s) {

	}
}

func IdentifierList(s *S) bool {
	d := s.Next()
	ok := true

	if d.Type == Lex.IdentifierToken {
		d = s.Peek()
		for ; d.Value == ","; d = s.Peek() {
			s.Next()
			d = s.Next()
			if d.Type != Lex.IdentifierToken {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	return ok
}