package clrcore

import (
	"testing"
)

func TestLexeme(t *testing.T) {
	l := Lexeme{Text, "blabla"}
	if l.String() != "[Text] \"blabla\"" {
		t.Errorf("got %q, expected %q from %q", l.String(), "[Text] \"blabla\"", "blabla")
	}
	if l.IsEmpty() {
		t.Errorf("expected %q empty", l)
	}
	if !l.IsA(Text) {
		t.Error("expected lexeme to be of type Text")
	}
	if l.IsA(Stop) {
		t.Error("expected lexeme to not be of type Stop")
	}
	l = Lexeme{Text, ""}
	if !l.IsEmpty() {
		t.Errorf("expected %q not empty", l)
	}
	l = Lexeme{}
	if l.String() != `[<nil>] ""` {
		t.Errorf("expected %q, got %q", `[<nil>] ""`, l.String())
	}
	if !l.IsNil() {
		t.Error("got not nil lexeme, expected nil")
	}
	if !l.IsA(nil) {
		t.Error("expected lexeme type to be nil")
	}
	l.Type = Text
	if l.IsNil() {
		t.Error("got nil lexeme, expected not nil")
	}
	if l.IsA(nil) {
		t.Error("expected lexeme type not to be nil")
	}
}
