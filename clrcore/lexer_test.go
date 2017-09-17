package clrcore

import (
	"path"
	"testing"
)

func TestRegisterLexerPanicNoNames(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	RegisterLexer(&LexerInfo{})
}

func TestRegisterLexerPanicNoNew(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	RegisterLexer(&LexerInfo{Names: []string{"test"}})
}

func TestRegisterLexerPanicDuplicateNames(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	newFunc := func(text string, stopMarkers ...string) (Lexer, error) {
		return nil, nil
	}
	RegisterLexer(&LexerInfo{Names: []string{"test"}, NewLexer: newFunc})
	RegisterLexer(&LexerInfo{Names: []string{"test2", "test"}, NewLexer: newFunc})
}

func TestRegisterLexerPanicInvalidFileName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	newFunc := func(text string, stopMarkers ...string) (Lexer, error) {
		return nil, nil
	}
	RegisterLexer(&LexerInfo{
		Names:     []string{"test3"},
		FileNames: []string{"*[].go"},
		NewLexer:  newFunc,
	})
}

func isLexerInSlice(lexer *LexerInfo, lexers []*LexerInfo) bool {
	for _, l := range lexers {
		if l == lexer {
			return true
		}
	}
	return false
}

func TestRegisterLexer(t *testing.T) {
	RegisterLexer(nil)
	newFunc := func(text string, stopMarkers ...string) (Lexer, error) {
		return nil, nil
	}
	mimeType := "text/x-gosrc"
	l := &LexerInfo{
		Names:     []string{"test4"},
		MimeTypes: []string{mimeType},
		FileNames: []string{"*.go"},
		NewLexer:  newFunc,
	}
	RegisterLexer(l)
	if !isLexerInSlice(l, Lexers()) {
		t.Errorf("lexer %q not registered", l.Names[0])
	}
	if LexerByName(l.Names[0]) != l {
		t.Errorf("failed retrieving lexer %q by name", l.Names[0])
	}
	if !isLexerInSlice(l, LexersByMimeType(mimeType)) {
		t.Errorf("failed retrieving lexer by mime type %q", mimeType)
	}
	mimeType = "text/x-csrc"
	if isLexerInSlice(l, LexersByMimeType(mimeType)) {
		t.Errorf("unexpected lexer %q retrieved by mime type %q", l.Names[0], mimeType)
	}
	fileName := "toto.go"
	if !isLexerInSlice(l, LexersByFileName(fileName)) {
		t.Errorf("failed retrieving lexer by file name %q", fileName)
	}
	fileName = "tutu.c"
	if isLexerInSlice(l, LexersByFileName(fileName)) {
		t.Errorf("unexpected lexer %q retrieved by file name %q", l.Names[0], fileName)
	}
	fileName = ""
	if isLexerInSlice(l, LexersByFileName(fileName)) {
		t.Errorf("unexpected lexer %q retrieved by file name %q", l.Names[0], fileName)
	}
	fileName = path.Join("xxx", "tutu.c")
	if isLexerInSlice(l, LexersByFileName(fileName)) {
		t.Errorf("unexpected lexer %q retrieved by file name %q", l.Names[0], fileName)
	}
}

type DummyLexer struct {
	score int
	stop  bool
}

func (l *DummyLexer) NextLexeme() Lexeme {
	if l.stop {
		return Lexeme{Type: Stop, Str: ""}
	}
	l.stop = true
	return Lexeme{Type: Text, Str: ""}
}

func (l *DummyLexer) RemainingText() string {
	return ""
}

func (l *DummyLexer) Score() int {
	return l.score
}

func TestLexerByScore(t *testing.T) {
	l1 := &LexerInfo{
		Names: []string{"l1"},
		NewLexer: func(text string, stopMarkers ...string) (Lexer, error) {
			return &DummyLexer{score: 10}, nil
		},
	}
	RegisterLexer(l1)

	l2 := &LexerInfo{
		Names: []string{"l2"},
		NewLexer: func(text string, stopMarkers ...string) (Lexer, error) {
			return &DummyLexer{score: 15}, nil
		},
	}
	RegisterLexer(l2)

	l3 := &LexerInfo{
		Names: []string{"l3"},
		NewLexer: func(text string, stopMarkers ...string) (Lexer, error) {
			return &DummyLexer{score: 15}, nil
		},
	}
	RegisterLexer(l3)

	l := LexerByScore("tadaaa", nil)
	if l == nil {
		t.Error("unexpected nil lexer")
	} else if l != l2 {
		t.Errorf("got lexer %q, expected %q", l.Names[0], l2.Names[0])
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		inStr    string
		inTabLen int
		outStr   string
	}{
		{inStr: "a test", inTabLen: -1, outStr: "a test"},
		{inStr: "a \rtest", inTabLen: -1, outStr: "a \ntest"},
		{inStr: "a \ntest", inTabLen: -1, outStr: "a \ntest"},
		{inStr: "a \n\rtest", inTabLen: -1, outStr: "a \ntest"},
		{inStr: "a \r\ntest", inTabLen: -1, outStr: "a \ntest"},
		{inStr: "a \r\n\rtest", inTabLen: -1, outStr: "a \n\ntest"},
		{inStr: "a test\r", inTabLen: -1, outStr: "a test\n"},
		{inStr: "a test\n", inTabLen: -1, outStr: "a test\n"},
		{inStr: "a test\r\n", inTabLen: -1, outStr: "a test\n"},
		{inStr: "a test\n\r", inTabLen: -1, outStr: "a test\n"},
		{inStr: "a test\f", inTabLen: -1, outStr: "a test"},
		{inStr: "\ta test", inTabLen: 0, outStr: "a test"},
		{inStr: "\ta test", inTabLen: 2, outStr: "  a test"},
		{inStr: "\ta test\nx\tother", inTabLen: 2, outStr: "  a test\nx  other"},
	}
	for _, test := range tests {
		outStr := Normalize(test.inStr, test.inTabLen)
		if outStr != test.outStr {
			t.Errorf("got %q, expected %q for %q", outStr, test.outStr, test.inStr)
		}
	}
}
