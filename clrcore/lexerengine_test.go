package clrcore

import (
	"fmt"
	"testing"
)

func TestLexerEngineQueue(t *testing.T) {
	lexer := &LexerEngine{
		outBuf: make([]Lexeme, 0, 4),
	}
	if !lexer.QueueEmpty() {
		t.Error("unexpected non-empty queue")
	}
	lexeme := lexer.UnqueueLexeme()
	if !lexeme.IsA(StopError) {
		t.Errorf("unexpected lexeme: %s", lexeme)
	}
	lexer.err = nil
	lexer.QueueLexeme(Lexeme{Text, "..."})
	if lexer.QueueEmpty() {
		t.Error("unexpected empty queue")
	}
	lexer.QueueLexeme(Lexeme{Text, "..."})
	lexer.QueueLexeme(Lexeme{Text, "..."})
	if len(lexer.outBuf) != 3 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 3)
	}
	if cap(lexer.outBuf) != 4 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 4)
	}
	lexer.QueueLexeme(Lexeme{Text, "..."})
	if len(lexer.outBuf) != 4 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 4)
	}
	if cap(lexer.outBuf) != 4 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 4)
	}
	lexer.QueueLexeme(Lexeme{Text, "..."})
	if len(lexer.outBuf)-lexer.outIdx != 5 {
		t.Errorf("got queue len %d, expected %d", len(lexer.outBuf)-lexer.outIdx, 5)
	}
	if len(lexer.outBuf) != 5 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 5)
	}
	if cap(lexer.outBuf) != 8 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 8)
	}

	lexer.QueueLexeme(Lexeme{Text, "..."})
	lexer.UnqueueLexeme()
	lexer.UnqueueLexeme()
	lexer.UnqueueLexeme()
	if len(lexer.outBuf)-lexer.outIdx != 3 {
		t.Errorf("got queue len %d, expected %d", len(lexer.outBuf)-lexer.outIdx, 3)
	}
	if len(lexer.outBuf) != 6 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 6)
	}
	if cap(lexer.outBuf) != 8 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 8)
	}
	lexer.UnqueueLexeme()
	lexer.UnqueueLexeme()
	if len(lexer.outBuf)-lexer.outIdx != 1 {
		t.Errorf("got queue len %d, expected %d", len(lexer.outBuf)-lexer.outIdx, 1)
	}
	if len(lexer.outBuf) != 6 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 6)
	}
	if cap(lexer.outBuf) != 8 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 8)
	}

	lexer.QueueLexeme(Lexeme{Text, "..."})
	lexer.QueueLexeme(Lexeme{Text, "..."})
	lexer.QueueLexeme(Lexeme{Text, "..."})
	if len(lexer.outBuf)-lexer.outIdx != 4 {
		t.Errorf("got queue len %d, expected %d", len(lexer.outBuf)-lexer.outIdx, 4)
	}
	if len(lexer.outBuf) != 4 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 4)
	}
	if cap(lexer.outBuf) != 8 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 8)
	}

	lexer.QueueLexeme(Lexeme{Text, "..."})
	if len(lexer.outBuf)-lexer.outIdx != 5 {
		t.Errorf("got queue len %d, expected %d", len(lexer.outBuf)-lexer.outIdx, 5)
	}
	if len(lexer.outBuf) != 5 {
		t.Errorf("got outBuf len %d, expected %d", len(lexer.outBuf), 5)
	}
	if cap(lexer.outBuf) != 8 {
		t.Errorf("got outBuf cap %d, expected %d", cap(lexer.outBuf), 8)
	}
	if lexer.err != nil {
		t.Errorf("unexpected non-nil error: %s", lexer.err)
	}
}

func TestLexerEngineStopMarker(t *testing.T) {
	def := &LexerDef{
		Name: "TestLexerDefPushPop",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&RegexDefRule{Re: "[A-Z]+", Do: PopMatch(CodeIdentifier)},
				}},
			}
		},
	}
	lexer, err := NewLexerEngine(def, "AB stop CD", []string{"stop"}, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes := []Lexeme{
		{CodeIdentifier, "AB"},
		{TextWhiteSpace, " "},
		{StopLexer, "stop"},
		{StopLexer, "stop"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}
	if lexer.RemainingText() != " CD" {
		t.Errorf("got %q, expected %q", lexer.RemainingText(), " CD")
	}
	if lexer.Score() != 0 {
		t.Errorf("got score %d, expected %d", lexer.Score(), 0)
	}
}

func TestLexerEnginePushPop(t *testing.T) {
	def := &LexerDef{
		Name: "TestLexerEnginePushPop",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&RegexDefRule{Re: "[A-Z]+", Do: All(PopMatch(CodeIdentifier), ScoreAdd(5))},
					&RegexDefRule{Re: `\(`, Do: All(PopMatch(CodeDelimiter), PushMode("param"))},
					&RegexDefRule{Re: `{`, Do: All(PopMatch(CodeDelimiter), PushMode("xxx"))},
					&RegexDefRule{Re: `\)`, Do: All(PopMatch(CodeDelimiter), PopMode())},
				}},
				{Name: "param", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&RegexDefRule{Re: "[a-z]+", Do: PopMatch(CodeIdentifierType)},
					&RegexDefRule{Re: `\(`, Do: All(PopMatch(CodeDelimiter), PushMode("param"))},
					&RegexDefRule{Re: `\)`, Do: All(PopMatch(CodeDelimiter), PopMode())},
				}},
			}
		},
	}

	lexer, err := NewLexerEngine(def, "AB{a b(c d)) a", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes := []Lexeme{
		{CodeIdentifier, "AB"},
		{CodeDelimiter, "{"},
		{StopError, "LexerDef 'TestLexerEnginePushPop' has no mode 'xxx'"},
		{StopError, "LexerDef 'TestLexerEnginePushPop' has no mode 'xxx'"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}
	if lexer.RemainingText() != "a b(c d)) a" {
		t.Errorf("got %q, expected %q", lexer.RemainingText(), "a b(c d)) a")
	}
	if lexer.Score() != 5 {
		t.Errorf("got score %d, expected %d", lexer.Score(), 5)
	}

	fmt.Println("----------------------------")
	lexer, err = NewLexerEngine(def, "AB(a)) a", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{CodeIdentifier, "AB"},
		{CodeDelimiter, "("},
		{CodeIdentifierType, "a"},
		{CodeDelimiter, ")"},
		{CodeDelimiter, ")"},
		{StopError, "pop empty mode stack"},
		{StopError, "pop empty mode stack"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}
	if lexer.RemainingText() != " a" {
		t.Errorf("got %q, expected %q", lexer.RemainingText(), " a")
	}
	if lexer.Score() != 5 {
		t.Errorf("got score %d, expected %d", lexer.Score(), 5)
	}
}
