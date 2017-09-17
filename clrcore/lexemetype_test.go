package clrcore

import (
	"testing"
)

func TestLexemeTypeString(t *testing.T) {
	tests := []struct {
		in  *LexemeType
		out string
	}{
		{Text, "Text"},
		{TextInvalid, "Text.Invalid"},
		{TextPunctuationDelimiter, "Text.Punctuation.Delimiter"},
	}
	for _, test := range tests {
		out := test.in.String()
		if out != test.out {
			t.Errorf("got %q, expected %q", out, test.out)
		}
	}
}

func TestLexemeClassTypes(t *testing.T) {
	for i, l := range LexemeClassTypes() {
		if lexemeClassTypes[i] != l {
			t.Errorf("got %q, expected %q", l, lexemeClassTypes[i])
		}
	}
}

func TestLexemeByName(t *testing.T) {
	for _, l1 := range LexemeClassTypes() {
		for _, l2 := range l1.Children {
			l3 := LexemeTypeByName(l2.Name)
			if l3 != l2 {
				t.Errorf("got %q, expected %q", l3, l2)
			}
		}
	}
	l := LexemeTypeByName("%*&@#")
	if l != nil {
		t.Errorf("got %q, expected nil", l)
	}
}

func TestLexemeClass(t *testing.T) {
	var class *LexemeType
	for _, l := range LexemeTypes() {
		// This works because lexeme types are sorted by name
		if l.Parent == nil {
			class = l
		}
		c := l.Class()
		if c != class {
			t.Errorf("got %q, expected %q", c, class)
		}
	}
}

func TestLexemeIsA(t *testing.T) {
	var class *LexemeType
	for _, l := range LexemeTypes() {
		// This works because lexeme types are sorted by name
		if l.Parent == nil {
			class = l
		}
		res := l.IsA(class)
		if !res {
			t.Errorf("%q should be a %q LexemeType", l, class)
		}
	}
	if TextPunctuation.IsA(Stop) {
		t.Errorf("expected %q not to be a %q LexemeType", TextPunctuation, Stop)
	}
	if TextPunctuation.IsA(nil) {
		t.Errorf("expected %q not to be a nil LexemeType", TextPunctuation)
	}
}

func TestNewLexemeTypePanicNoName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewLexemeType(nil, "")
}

func TestNewLexemeTypePanicDuplicateName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewLexemeType(nil, LexemeClassTypes()[0].Name)
}
