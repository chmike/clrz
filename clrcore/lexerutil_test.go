package clrcore

import (
	"testing"
)

var nilLexeme Lexeme

func TestPopSingleQuotedString(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: `'abc'`, out: Lexeme{CodeStringSingle, `'abc'`}, rem: ``},
		{in: `'ab\'c'`, out: Lexeme{CodeStringSingle, `'ab\'c'`}, rem: ``},
		{in: `'abc' `, out: Lexeme{CodeStringSingle, `'abc'`}, rem: ` `},
		{in: `'ab\'c' `, out: Lexeme{CodeStringSingle, `'ab\'c'`}, rem: ` `},
		{in: "'ab\nc' ", out: Lexeme{CodeStringSingle, `'ab`}, rem: "\nc' "},
		{in: `'ab`, out: Lexeme{CodeStringSingle, `'ab`}, rem: ``},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopSingleQuotedString()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopDoubleQuotedString(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: `"abc"`, out: Lexeme{CodeStringDouble, `"abc"`}, rem: ``},
		{in: `"ab\"c"`, out: Lexeme{CodeStringDouble, `"ab\"c"`}, rem: ``},
		{in: `"abc" `, out: Lexeme{CodeStringDouble, `"abc"`}, rem: ` `},
		{in: `"ab\"c" `, out: Lexeme{CodeStringDouble, `"ab\"c"`}, rem: ` `},
		{in: "\"ab\nc\" ", out: Lexeme{CodeStringDouble, `"ab`}, rem: "\nc\" "},
		{in: `"ab`, out: Lexeme{CodeStringDouble, `"ab`}, rem: ``},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopDoubleQuotedString()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopBackTickString(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "`abc`", out: Lexeme{CodeStringMultiline, "`abc`"}, rem: ``},
		{in: "`ab\nc`", out: Lexeme{CodeStringMultiline, "`ab\nc`"}, rem: ``},
		{in: "`` ", out: Lexeme{CodeStringMultiline, "``"}, rem: " "},
		{in: "`abc \n ", out: Lexeme{CodeStringMultiline, "`abc \n "}, rem: ""},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopBackTickString()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestTripleDoubleQuotedString(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: `"""abc"""`, out: Lexeme{CodeStringMultiline, `"""abc"""`}},
		{in: `"""ab " cd "" """`, out: Lexeme{CodeStringMultiline, `"""ab " cd "" """`}},
		{in: `"""ab " cd "" """"`, out: Lexeme{CodeStringMultiline, `"""ab " cd "" """`}, rem: `"`},
		{in: `"""""a`, out: Lexeme{CodeStringMultiline, `"""""a`}},
		{in: `""""""`, out: Lexeme{CodeStringMultiline, `""""""`}},
		{in: `"""a"""`, out: Lexeme{CodeStringMultiline, `"""a"""`}},
		{in: `"""ab"""`, out: Lexeme{CodeStringMultiline, `"""ab"""`}},
		{in: `"""abc"""`, out: Lexeme{CodeStringMultiline, `"""abc"""`}},
		{in: `"""abcd"""`, out: Lexeme{CodeStringMultiline, `"""abcd"""`}},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopTripleDoubleQuotedString()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopOneCharLineComment(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "#abc", out: Lexeme{CodeComment, "#abc"}, rem: ``},
		{in: "#ab\nc", out: Lexeme{CodeComment, "#ab"}, rem: "\nc"},
		{in: "'\n ", out: Lexeme{CodeComment, "'"}, rem: "\n "},
		{in: ";abc \n ", out: Lexeme{CodeComment, ";abc "}, rem: "\n "},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopOneCharLineComment()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopTwoCharLineComment(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "//abc", out: Lexeme{CodeComment, "//abc"}, rem: ``},
		{in: "//ab\nc", out: Lexeme{CodeComment, "//ab"}, rem: "\nc"},
		{in: "//\n ", out: Lexeme{CodeComment, "//"}, rem: "\n "},
		{in: "///", out: Lexeme{CodeComment, "///"}, rem: ""},
		{in: "//", out: Lexeme{CodeComment, "//"}, rem: ""},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopTwoCharLineComment()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopSlashStarComment(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "/*", out: Lexeme{CodeComment, "/*"}, rem: ""},
		{in: "/*abc", out: Lexeme{CodeComment, "/*abc"}, rem: ""},
		{in: "/*ab\nc", out: Lexeme{CodeComment, "/*ab\nc"}, rem: ""},
		{in: "/*\n */", out: Lexeme{CodeComment, "/*\n */"}, rem: ""},
		{in: "/*/a\n*/b", out: Lexeme{CodeComment, "/*/a\n*/"}, rem: "b"},
		{in: "/**/a", out: Lexeme{CodeComment, "/**/"}, rem: "a"},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopSlashStarComment()
		if out != test.out {
			t.Errorf("got lexeme %q, expected %q for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain %q, expected %q for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopDecimalNumber(t *testing.T) {
	tests := []struct {
		in  string
		end int
		out Lexeme
		rem string
	}{
		{in: "0.", end: 2, out: Lexeme{CodeNumberDecimal, "0."}, rem: ""},
		{in: "0.012q", end: 2, out: Lexeme{CodeNumberDecimal, "0.012"}, rem: "q"},
		{in: "0.e", end: 2, out: Lexeme{CodeNumberDecimal, "0."}, rem: "e"},
		{in: "0.e ", end: 2, out: Lexeme{CodeNumberDecimal, "0."}, rem: "e "},
		{in: "0.e+", end: 2, out: Lexeme{CodeNumberDecimal, "0."}, rem: "e+"},
		{in: "0.e+ ", end: 2, out: Lexeme{CodeNumberDecimal, "0."}, rem: "e+ "},
		{in: "0.e-1", end: 2, out: Lexeme{CodeNumberDecimal, "0.e-1"}, rem: ""},
		{in: "0.e-1 ", end: 2, out: Lexeme{CodeNumberDecimal, "0.e-1"}, rem: " "},
		{in: "0.e123izz", end: 2, out: Lexeme{CodeNumberDecimal, "0.e123"}, rem: "izz"},
		{in: "0.e1izz", end: 2, out: Lexeme{CodeNumberDecimal, "0.e1"}, rem: "izz"},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopDecimalNumber(test.end)
		if out != test.out {
			t.Errorf("got lexeme {%s}, expected {%s} for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain '%s', expected '%s' for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopNumber(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "", out: nilLexeme, rem: ""},
		{in: "a", out: nilLexeme, rem: "a"},
		{in: "a  ", out: nilLexeme, rem: "a  "},
		{in: "0.012q", out: Lexeme{CodeNumberDecimal, "0.012"}, rem: "q"},
		{in: "0.e", out: Lexeme{CodeNumberDecimal, "0."}, rem: "e"},
		{in: "0.e ", out: Lexeme{CodeNumberDecimal, "0."}, rem: "e "},
		{in: "0.e+", out: Lexeme{CodeNumberDecimal, "0."}, rem: "e+"},
		{in: "0.e+ ", out: Lexeme{CodeNumberDecimal, "0."}, rem: "e+ "},
		{in: "0.e-1", out: Lexeme{CodeNumberDecimal, "0.e-1"}, rem: ""},
		{in: "0.e-1 ", out: Lexeme{CodeNumberDecimal, "0.e-1"}, rem: " "},
		{in: "0.e123izz", out: Lexeme{CodeNumberDecimal, "0.e123"}, rem: "izz"},
		{in: "0.e1izz", out: Lexeme{CodeNumberDecimal, "0.e1"}, rem: "izz"},
		{in: ".1", out: Lexeme{CodeNumberDecimal, ".1"}, rem: ""},
		{in: ".1e", out: Lexeme{CodeNumberDecimal, ".1"}, rem: "e"},
		{in: ".e", out: nilLexeme, rem: ".e"},
		{in: ".e1", out: nilLexeme, rem: ".e1"},
		{in: ".1e1", out: Lexeme{CodeNumberDecimal, ".1e1"}, rem: ""},
		{in: "1.1e1", out: Lexeme{CodeNumberDecimal, "1.1e1"}, rem: ""},
		{in: "0.1e1", out: Lexeme{CodeNumberDecimal, "0.1e1"}, rem: ""},
		{in: "1.1", out: Lexeme{CodeNumberDecimal, "1.1"}, rem: ""},
		{in: "0.1", out: Lexeme{CodeNumberDecimal, "0.1"}, rem: ""},
		{in: "1.1 ", out: Lexeme{CodeNumberDecimal, "1.1"}, rem: " "},
		{in: "0.1 ", out: Lexeme{CodeNumberDecimal, "0.1"}, rem: " "},
		{in: ".1E1", out: Lexeme{CodeNumberDecimal, ".1E1"}, rem: ""},
		{in: "0", out: Lexeme{CodeNumberInteger, "0"}, rem: ""},
		{in: "1", out: Lexeme{CodeNumberInteger, "1"}, rem: ""},
		{in: "-0", out: Lexeme{CodeNumberInteger, "-0"}, rem: ""},
		{in: "0 ", out: Lexeme{CodeNumberInteger, "0"}, rem: " "},
		{in: "1 ", out: Lexeme{CodeNumberInteger, "1"}, rem: " "},
		{in: "-0 ", out: Lexeme{CodeNumberInteger, "-0"}, rem: " "},
		{in: "+", out: nilLexeme, rem: "+"},
		{in: ".", out: nilLexeme, rem: "."},
		{in: "00", out: Lexeme{CodeNumberOctal, "00"}, rem: ""},
		{in: "00 ", out: Lexeme{CodeNumberOctal, "00"}, rem: " "},
		{in: "08", out: nilLexeme, rem: "08"},
		{in: "08 ", out: nilLexeme, rem: "08 "},
		{in: "0x0", out: Lexeme{CodeNumberHexadecimal, "0x0"}, rem: ""},
		{in: "0x0 ", out: Lexeme{CodeNumberHexadecimal, "0x0"}, rem: " "},
		{in: "0X0 ", out: Lexeme{CodeNumberHexadecimal, "0X0"}, rem: " "},
		{in: "0x", out: nilLexeme, rem: "0x"},
		{in: "0x ", out: nilLexeme, rem: "0x "},
		{in: "0XAF09af", out: Lexeme{CodeNumberHexadecimal, "0XAF09af"}, rem: ""},
		{in: "0XAF09af ", out: Lexeme{CodeNumberHexadecimal, "0XAF09af"}, rem: " "},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopNumber()
		if out != test.out {
			t.Errorf("got lexeme {%s}, expected {%s} for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain '%s', expected '%s' for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopWhiteSpaces(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: " ", out: Lexeme{TextWhiteSpace, " "}, rem: ""},
		{in: "  ", out: Lexeme{TextWhiteSpace, "  "}, rem: ""},
		{in: "  \n", out: Lexeme{TextWhiteSpace, "  "}, rem: "\n"},
		{in: " \t", out: Lexeme{TextWhiteSpace, " \t"}, rem: ""},
		{in: " \t\n", out: Lexeme{TextWhiteSpace, " \t"}, rem: "\n"},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopWhiteSpaces()
		if out != test.out {
			t.Errorf("got lexeme {%s}, expected {%s} for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain '%s', expected '%s' for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopNewLines(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "\n", out: Lexeme{TextNewLine, "\n"}, rem: ""},
		{in: "\r\n", out: Lexeme{TextNewLine, "\r\n"}, rem: ""},
		{in: "\n\n\n", out: Lexeme{TextNewLine, "\n\n\n"}, rem: ""},
		{in: "\n\n\n ", out: Lexeme{TextNewLine, "\n\n\n"}, rem: " "},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopNewLines()
		if out != test.out {
			t.Errorf("got lexeme {%s}, expected {%s} for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain '%s', expected '%s' for %+v", l.Str, test.rem, test)
		}
	}
}

func TestPopASCIIIdentifier(t *testing.T) {
	tests := []struct {
		in  string
		out Lexeme
		rem string
	}{
		{in: "Aa", out: Lexeme{CodeIdentifier, "Aa"}, rem: ""},
		{in: "Aa ", out: Lexeme{CodeIdentifier, "Aa"}, rem: " "},
		{in: "_A ", out: Lexeme{CodeIdentifier, "_A"}, rem: " "},
		{in: "a9A ", out: Lexeme{CodeIdentifier, "a9A"}, rem: " "},
		{in: "aé ", out: Lexeme{CodeIdentifier, "a"}, rem: "é "},
	}
	for _, test := range tests {
		l := Lexeme{Str: test.in}
		out := l.PopASCIIIdentifier()
		if out != test.out {
			t.Errorf("got lexeme {%s}, expected {%s} for %+v", out, test.out, test)
		}
		if l.Str != test.rem {
			t.Errorf("got remain '%s', expected '%s' for %+v", l.Str, test.rem, test)
		}
	}
}
