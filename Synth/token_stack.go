package Synth

import (
	"github.com/thalysonr/go-compiler/Lex"
	"errors"
)

type tokenNode struct {
	t    Lex.Token
	next *tokenNode
}

type tokenStack struct {
	start *tokenNode
}

func newTokenStack() tokenStack {
	return tokenStack{}
}

func (s *tokenStack) push(t Lex.Token) {
	node := &tokenNode{t: t}
	if s.start == nil {
		s.start = node
	} else {
		node.next = s.start
		s.start = node
	}
}

func (s *tokenStack) pop() (Lex.Token, error) {
	if s.start == nil {
		return Lex.Token{}, errors.New("Sem tokens")
	} else {
		n := s.start
		s.start = n.next
		return n.t, nil
	}
}

func (s *tokenStack) clear() {
	s.start = nil
}
