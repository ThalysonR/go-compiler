package Synth

import (
	"errors"
	"github.com/thalysonr/go-compiler/Lex"
)

type StateFunc func(*S) StateFunc

type DeclarationType int

type Declaration struct {
	Type  DeclarationType
	Value []Lex.Token
}

type S struct {
	source          *Lex.L
	start, position int
	startState      StateFunc
	Err             error
	declarations    chan Declaration
	ErrorHandler    func(e string)
	rewind          tokenStack
	tokens 		    []Lex.Token
}

// New creates a returns a lexer ready to parse the given source code.
func New(src *Lex.L, start StateFunc) *S {
	return &S{
		source:     src,
		startState: start,
		start:      0,
		position:   0,
		rewind:     newTokenStack(),
	}
}

// Start begins executing the Lexer in an asynchronous manner (using a goroutine).
func (s *S) Start() {
	// Take half the string length as a buffer size.
	//buffSize := len(l.source) / 2
	//if buffSize <= 0 {
	//	buffSize = 1
	//}
	s.declarations = make(chan Declaration)
	go s.run()
}

func (s *S) StartSync() {
	// Take half the string length as a buffer size.
	//buffSize := len(l.source) / 2
	//if buffSize <= 0 {
	//	buffSize = 1
	//}
	s.declarations = make(chan Declaration)
	s.run()
}

// Current returns the value being being analyzed at this moment.
func (s *S) Current() []Lex.Token {
	return s.tokens[s.start:s.position]
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (s *S) Emit(d DeclarationType) {
	dec := Declaration{
		Type:  d,
		Value: s.Current(),
	}
	s.declarations <- dec
	s.start = s.position
	s.rewind.clear()
}

// Peek performs a Next operation immediately followed by a Rewind returning the
// peeked rune.
func (s *S) Peek() *Lex.Token {
	t := s.Next()
	s.Rewind()

	return t
}

// Rewind will take the last rune read (if any) and rewind back. Rewinds can
// occur more than once per call to Next but you can never rewind past the
// last point a token was emitted.
func (s *S) Rewind() {
	_, err := s.rewind.pop()
	if err == nil {
		s.position -= 1
		if s.position < s.start {
			s.position = s.start
		}
	}
}

// Next pulls the next rune from the Lexer and returns it, moving the position
// forward in the source.
func (s *S) Next() *Lex.Token {
	var (
		t *Lex.Token
		i int
	)
	//str := l.source[l.position:]
	if len(s.tokens) == s.position {
		tok, done := s.source.NextToken()
		if done {
			t, i =  nil, 0
		} else {
			t, i = tok, 1
			s.tokens = append(s.tokens, *tok)
		}
	} else {
		t = &s.tokens[s.position]
		i = 1
	}
	s.position += i
	s.rewind.push(*t)

	return t
}

// NextToken returns the next token from the lexer and a value to denote whether
// or not the token is finished.
func (s *S) NextDeclaration() (*Declaration, bool) {
	if dec, ok := <-s.declarations; ok {
		return &dec, false
	} else {
		return nil, true
	}
}

// Partial yyLexer implementation

func (s *S) Error(e string) {
	if s.ErrorHandler != nil {
		s.Err = errors.New(e)
		s.ErrorHandler(e)
	} else {
		panic(e)
	}
}

// Private methods

func (s *S) run() {
	state := s.startState
	for state != nil {
		state = state(s)
	}
	close(s.declarations)
}

