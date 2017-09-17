package clrcore

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexerDefInit(t *testing.T) {
	def := &LexerDef{Name: "TestLexerDefInit"}
	if _, err := NewLexerEngine(def, "", nil, nil); err == nil {
		t.Error("unexpected nil error")
	}

	def = &LexerDef{
		Name:     "TestLexerDefInit",
		InitFunc: func(d *LexerDef) {}, // leaves Modes to nil
	}
	if _, err := NewLexerEngine(def, "", nil, nil); err == nil {
		t.Error("unexpected nil error")
	}

	def = &LexerDef{
		Name: "TestLexerDefInit",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root"}, // define a mode with no rules
			}
		},
	}
	if _, err := NewLexerEngine(def, "", nil, nil); err == nil {
		t.Error("unexpected nil error")
	}

	def = &LexerDef{
		Name: "TestLexerDefInit",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&FuncDefRule{},            // Init succeeds
					&RegexDefRule{Re: "[a-b"}, // Bogus regex rule, init should fail
				}},
			}
		},
	}
	if _, err := NewLexerEngine(def, "", nil, nil); err == nil {
		t.Error("unexpected nil error")
	}

	def = &LexerDef{
		Name: "TestLexerDefInit",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&FuncDefRule{},               // Init succeeds
					&RegexDefRule{Re: "[a-z ]*"}, // Valid regex rule
				}},
			}
		},
	}
	if _, err := NewLexerEngine(def, "", nil, nil); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestLexerDefExec(t *testing.T) {
	def := &LexerDef{
		Name: "TestLexerDefExec",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&FuncDefRule{ExecFunc: func(l *LexerEngine) bool {
						word := "test"
						if !strings.HasPrefix(l.str, word) {
							return false
						}
						l.QueueLexeme(Lexeme{CodeIdentifierClass, word})
						l.str = l.str[len(word):]
						l.score++
						return true
					}}, // Init succeeds
					&RegexDefRule{Re: "[a-z]+", Do: PopMatch(CodeIdentifierVariable)},
					&RegexDefRule{Re: "[A-Z]+", Do: All(PopMatch(CodeIdentifierFunction), ScoreAdd(5))},
					&RegexDefRule{Re: "10([a-zA-Z0-9]*) +(test|toto) *\n",
						Do: PopMatch(CodeIdentifier, CodeIdentifierKeyword)},
				}},
			}
		},
	}
	lexer, err := NewLexerEngine(def, "  abc\n  test", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes := []Lexeme{
		{TextWhiteSpace, "  "},
		{CodeIdentifierVariable, "abc"},
		{TextNewLine, "\n"},
		{TextWhiteSpace, "  "},
		{CodeIdentifierClass, "test"},
		{StopEndOfString, ""},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}

	lexer, err = NewLexerEngine(def, "10Test2 toto \n  ABC", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{Text, "10"},
		{CodeIdentifierKeyword, "Test2"},
		{Text, " \n"},
		{TextWhiteSpace, "  "},
		{CodeIdentifierFunction, "ABC"},
		{StopEndOfString, ""},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}

	def = &LexerDef{
		Name: "TestLexerDefExec",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&RegexDefRule{Re: "10([a-zA-Z0-9]*)",
						Do: PopMatch(CodeIdentifier, CodeIdentifierKeyword)},
					&RegexDefRule{Re: "11([a-zA-Z0-9]*)",
						Do: All(func(l *LexerEngine, match []int) bool {
							l.err = fmt.Errorf("dummy test error")
							return true
						})},
					&RegexDefRule{Re: "12(?:[a-zA-Z0-9]*)",
						Do: PopMatch(CodeIdentifier, CodeIdentifierKeyword)},
					&RegexDefRule{Re: "13((Test|Toto)[a-zA-Z0-9]*)",
						Do: PopMatch(CodeIdentifier, CodeIdentifierKeyword)},
					&RegexDefRule{Re: "14((Test|Toto)[a-zA-Z0-9]*)(toto|test)",
						Do: PopMatch(CodeIdentifier, CodeIdentifierKeyword, CodeIdentifierKeyword)},
				}},
			}
		},
	}
	lexer, err = NewLexerEngine(def, "10Test", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{StopError, "invalid number of lexemeTypes (LexerDef='TestLexerDefExec', Mode='root', Rule=2)"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}

	lexer, err = NewLexerEngine(def, "11Test", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{StopError, "dummy test error"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}

	lexer, err = NewLexerEngine(def, "12Test", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{StopError, "invalid number of lexemeTypes (LexerDef='TestLexerDefExec', Mode='root', Rule=4)"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}

	lexer, err = NewLexerEngine(def, "13Test123", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{Text, "13"},
		{CodeIdentifierKeyword, "Test123"},
		{StopError, "overlapping regex groups (LexerDef='TestLexerDefExec', Mode='root', Rule=5)"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}
	lexer, err = NewLexerEngine(def, "14Test123toto", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes = []Lexeme{
		{Text, "14"},
		{CodeIdentifierKeyword, "Test123"},
		{StopError, "overlapping regex groups (LexerDef='TestLexerDefExec', Mode='root', Rule=6)"},
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}
}

func TestLexerDefPushPop(t *testing.T) {
	def := &LexerDef{
		Name: "TestLexerDefPushPop",
		InitFunc: func(d *LexerDef) {
			d.Modes = []*LexerDefMode{
				{Name: "root", Rules: []LexerDefRule{
					WhiteSpaceRule, NewLineRule,
					&RegexDefRule{Re: "[A-Z]+", Do: All(PopMatch(CodeIdentifier), ScoreAdd(5))},
					&RegexDefRule{Re: `\(`, Do: All(PopMatch(CodeDelimiter), PushMode("param"))},
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
	lexer, err := NewLexerEngine(def, "AB(a b(c d)) a", nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	lexemes := []Lexeme{
		{CodeIdentifier, "AB"},
		{CodeDelimiter, "("},
		{CodeIdentifierType, "a"},
		{TextWhiteSpace, " "},
		{CodeIdentifierType, "b"},
		{CodeDelimiter, "("},
		{CodeIdentifierType, "c"},
		{TextWhiteSpace, " "},
		{CodeIdentifierType, "d"},
		{CodeDelimiter, ")"},
		{CodeDelimiter, ")"},
		{TextWhiteSpace, " "},
		{StopLexer, ""}, // No matching rule in root mode
		{StopLexer, ""}, // No matching rule in root mode
	}
	for _, expect := range lexemes {
		lexeme := lexer.NextLexeme()
		if lexeme != expect {
			t.Errorf("got %s, expected %s", lexeme, expect)
		}
	}
	if lexer.RemainingText() != "a" {
		t.Errorf("got %q, expected %q", lexer.RemainingText(), "a")
	}
	if lexer.Score() != 5 {
		t.Errorf("got score %d, expected %d", lexer.Score(), 5)
	}
}
