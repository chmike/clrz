package clrcore

import (
	"errors"
	"fmt"
	"strings"
)

// LexerEngine is an engine to decompose an input text into lexemes based on
// the given definition and optional extend information.
type LexerEngine struct {
	def         *LexerDef       // The LexerDef currently used.
	stopMarkers []string        // String markers stopping the lexer.
	score       int             // The score of the parsed text with the given definition.
	str         string          // Text remaining to be parsed (slice of text.Str).
	outBuf      []Lexeme        // Buffer for output lexemes.
	outIdx      int             // Index of next lexeme in outBuf to output.
	mode        *LexerDefMode   // Current mode in RegexLexerDef.
	ruleIdx     int             // Current rule index of the current mode.
	modeStack   []*LexerDefMode // Stack of modes.
	err         error           // Last error.
	stopLexeme  Lexeme          // First stop lexeme issued or nil lexeme.
	extend      interface{}     // Language specific additionnal information.
}

// NewLexerEngine returns a LexerEngine that will use the LexerDef to
// parse the input text into lexemes. The parsing will stop when an end marker is found
// in the text or the end of the text is reached.
func NewLexerEngine(def *LexerDef, text string, stopMarkers []string, extend interface{}) (*LexerEngine, error) {
	if err := def.Init(); err != nil {
		return nil, err
	}
	return &LexerEngine{
		def:         def,
		stopMarkers: stopMarkers,
		str:         text,
		outBuf:      make([]Lexeme, 0, 4),
		mode:        def.Modes[0],
		extend:      extend,
	}, nil
}

// Score returns a matching score for the parsed language. Its value only
// make sense once a Stop lexeme has been reached. When trying multiple lexer
// on a piece of text, the lexer with the highest score will be picked.
func (l *LexerEngine) Score() int {
	return l.score
}

// RemainingText return a Text lexeme containing the remaining text to parse.
func (l *LexerEngine) RemainingText() string {
	return l.str
}

// NextLexeme return the next lexeme extracted from the input text until a stop
// lexeme is returned. The stop lexeme is then returned on every call.
func (l *LexerEngine) NextLexeme() (lexeme Lexeme) {
	if !l.stopLexeme.IsNil() {
		return l.stopLexeme
	}
	if l.QueueEmpty() {
		l.GetLexemes()
	}
	lexeme = l.UnqueueLexeme()
	if lexeme.IsA(Stop) {
		l.stopLexeme = lexeme
	}
	return
}

// QueueLexeme appends lexeme to the back of outBuf, growing it when required.
func (l *LexerEngine) QueueLexeme(lexeme Lexeme) {
	if l.outIdx == len(l.outBuf) { // queue is empty
		l.outBuf = l.outBuf[:1]
		l.outIdx = 0
	} else if len(l.outBuf) < cap(l.outBuf) { // room at end of queue
		l.outBuf = l.outBuf[:len(l.outBuf)+1]
	} else if l.outIdx > 4 { // room in front of queue
		l.outBuf = l.outBuf[:copy(l.outBuf, l.outBuf[l.outIdx:])+1]
		l.outIdx = 0
	} else { // queue buffer must grow
		tmp := make([]Lexeme, 2*cap(l.outBuf))
		l.outBuf = tmp[:copy(tmp, l.outBuf[l.outIdx:])+1]
		l.outIdx = 0
	}
	l.outBuf[len(l.outBuf)-1] = lexeme
}

// QueueEmpty return true if the lexeme queue is empty.
func (l *LexerEngine) QueueEmpty() bool {
	return l.outIdx == len(l.outBuf)
}

// UnqueueLexeme extract and return the front most lexeme from the queue,
// or a StopError lexeme if the queue is empty.
func (l *LexerEngine) UnqueueLexeme() (lexeme Lexeme) {
	if l.outIdx == len(l.outBuf) {
		return Lexeme{StopError, "can't unqueue lexeme from an empty queue"}
	}
	lexeme = l.outBuf[l.outIdx]
	l.outIdx++
	return
}

// GetLexemes extracts lexemes from the text and queue them in outBuf.
func (l *LexerEngine) GetLexemes() {
	if len(l.str) == 0 {
		l.QueueLexeme(Lexeme{Type: StopEndOfString})
		return
	}
	for l.ruleIdx = 0; l.ruleIdx < len(l.mode.Rules); l.ruleIdx++ {
		for _, stopMarker := range l.stopMarkers {
			if strings.HasPrefix(l.str, stopMarker) {
				l.QueueLexeme(Lexeme{StopLexer, stopMarker})
				l.str = l.str[len(stopMarker):]
				return
			}
		}
		done := l.mode.Rules[l.ruleIdx].Exec(l)
		if l.err != nil {
			l.QueueLexeme(Lexeme{Type: StopError, Str: l.err.Error()})
			return
		}
		if done {
			return
		}
	}
	l.QueueLexeme(Lexeme{Type: StopLexer})
}

// PushMode set the current mode to the named mode.
func (l *LexerEngine) PushMode(name string) {
	for _, m := range l.def.Modes {
		if m.Name == name {
			l.modeStack = append(l.modeStack, l.mode)
			l.mode = m
			return
		}
	}
	l.err = fmt.Errorf("LexerDef %q has no mode %q", l.def.Name, name)
}

// PopMode set the current mode to the stacked mode.
// If PopMode is called on an empty mode stack, mode is left unmodified and
// l.err is set to an error.
func (l *LexerEngine) PopMode() {
	if len(l.modeStack) == 0 {
		l.err = errors.New("pop empty mode stack")
		return
	}
	last := len(l.modeStack) - 1
	l.mode = l.modeStack[last]
	l.modeStack = l.modeStack[:last]
}
