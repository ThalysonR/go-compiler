package Funcoes

import (
	. "github.com/thalysonr/go-compiler/Synth"
	"github.com/thalysonr/go-compiler/Lex"
)

func Type(s *S) bool {
	d := s.Peek()

	if d.Value == "(" {
		s.Next()
		if Type(s) {
			if s.Next().Value == ")" {
				return true
			} else {
				s.Rewind()
				s.Rewind()
			}
		} else {
			s.Rewind()
		}
	} else if d.Type == Lex.IdentifierToken {
		s.Next()
		if s.Peek().Value == "." {
			s.Next()
			if s.Next().Type == Lex.IdentifierToken {
				return true
			} else {
				s.Rewind()
				s.Rewind()
				s.Rewind()
			}
		} else {
			return true
		}
	} else {
		return false
	}
	return false
}

func TypeLit(s *S) bool {
	switch s.Peek().Value {
	case "[":
		return ArraySliceType(s)
	case "struct":
		return StructType(s)
	case "*":
		return PointerType(s)
	case "func":
		return FunctionType(s)
	case "interface":
		return InterfaceType(s)
	case "map":
		return MapType(s)
	case "chan":
		return ChannelType(s)
	default:
		return false
	}
}

func ArraySliceType(s *S) bool {
	s.Next()
	if d := s.Next(); d.Value == "]" {
		return Type(s)
	}
}