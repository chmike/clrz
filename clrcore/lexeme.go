package clrcore

import "fmt"

// A Lexeme is a unit of text used or produced by the lexer.
type Lexeme struct {
	Type *LexemeType
	Str  string
}

func (l Lexeme) String() string {
	return fmt.Sprintf("[%s] %q", l.Type, l.Str)
}

// IsEmpty return true if the lexeme has an empty string value.
func (l Lexeme) IsEmpty() bool {
	return l.Str == ""
}

// IsNil return true if the lexeme has a nil type value.
func (l Lexeme) IsNil() bool {
	return l.Type == nil
}

// IsA return true if the lexeme has the same or a parent type of t.
func (l Lexeme) IsA(t *LexemeType) bool {
	if l.Type == nil {
		return t == nil
	}
	return l.Type.IsA(t)
}
