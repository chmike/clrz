package clrcore

import (
	"reflect"
	"strings"
	"testing"
)

func TestMakeStyle(t *testing.T) {
	tests := []struct {
		in  []string
		out string
		err string
	}{
		{in: []string{}, out: ""},
		{in: []string{"bold"}, out: "bold"},
		{in: []string{"bold", "italic"}, out: "italic bold"},
		{in: []string{"test", "bold"}, err: "invalid style attribute 'test'"},
		{in: []string{"text#123456"}, out: "text#123456"},
		{in: []string{"back#789ABC", "bold", "italic"}, out: "italic bold back#789ABC"},
		{in: []string{"#123"}, err: "invalid style attribute '#123'"},
		{in: []string{"italic", "bold", "#1234"}, err: "invalid style attribute '#1234'"},
	}
	for _, test := range tests {
		style, err := MakeStyle(test.in...)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("got error '%s', expected '%s' for '%s'", err, test.err, strings.Join(test.in, ", "))
			}
		} else if out := style.String(); out != test.out {
			t.Errorf("got '%s', expected '%s' for '%s'", out, test.out, strings.Join(test.in, ", "))
		}
	}
}

func TestTypeStyle(t *testing.T) {
	var style TypeStyle
	if style.Italic() {
		t.Error("unexpected italic style")
	}
	if style.Bold() {
		t.Error("unexpected bold style")
	}
	if style.HasTextColor() {
		t.Error("unexpected text color style")
	}
	r, g, b := style.TextColor()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("got text color #%02X%02X%02X, expected #000000", r, g, b)
	}
	if style.HasBackColor() {
		t.Error("unexpected back color style")
	}
	r, g, b = style.BackColor()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("got back color #%02X%02X%02X, expected #000000", r, g, b)
	}

	style, err := MakeStyle("bold", "italic", "text#123456", "back#556677", "back#789ABC")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if !style.Italic() {
			t.Error("unexpected non-italic style")
		}
		if !style.Bold() {
			t.Error("unexpected non-bold style")
		}
		if !style.HasTextColor() {
			t.Error("unexpected text uncolored style")
		}
		r, g, b := style.TextColor()
		if r != 0x12 || g != 0x34 || b != 0x56 {
			t.Errorf("got text color #%02X%02X%02X, expected #123456", r, g, b)
		}
		if !style.HasBackColor() {
			t.Error("unexpected back uncolored style")
		}
		r, g, b = style.BackColor()
		if r != 0x78 || g != 0x9A || b != 0xBC {
			t.Errorf("got back color #%02X%02X%02X, expected #789ABC", r, g, b)
		}
	}
}

func TestNewStyleFromText(t *testing.T) {
	colorStyle, _ := MakeStyle("text#123456")
	italicStyle, _ := MakeStyle("italic")
	boldStyle, _ := MakeStyle("bold")
	dfltStyle, _ := MakeStyle()

	textStyle := `
	Code.Identifier italic
	Code.Identifier.Variable bold
	Code.Identifier.Namespace text#123456
	`
	style, err := NewStyle(textStyle)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		gotTypes := style[boldStyle]
		expectTypes := []*LexemeType{CodeIdentifierVariable}
		if !reflect.DeepEqual(gotTypes, expectTypes) {
			t.Errorf("got bold types %v, expect %v", gotTypes, expectTypes)
		}
		gotTypes = style[colorStyle]
		expectTypes = []*LexemeType{CodeIdentifierNamespace}
		if !reflect.DeepEqual(gotTypes, expectTypes) {
			t.Errorf("got color types %v, expect %v", gotTypes, expectTypes)
		}
		gotTypes = style[italicStyle]
		expectTypes = []*LexemeType{CodeIdentifier, CodeIdentifierFunction, CodeIdentifierMethod,
			CodeIdentifierType, CodeIdentifierClass, CodeIdentifierKeyword, CodeIdentifierLiteral,
			CodeIdentifierOperator}
		if !reflect.DeepEqual(gotTypes, expectTypes) {
			t.Errorf("got italic types %v, expect %+v", gotTypes, expectTypes)
		}
		gotTypes = style[dfltStyle]
		expectTypes = []*LexemeType{Text, TextWhiteSpace, TextNewLine, TextInvalid, TextPunctuation,
			TextPunctuationSeparator, TextPunctuationDelimiter, TextOperator, TextWord, TextNumber,
			TextOther, Code, CodeString, CodeStringSingle, CodeStringDouble, CodeStringRaw,
			CodeStringUnicode, CodeStringMultiline, CodeNumber, CodeNumberInteger,
			CodeNumberHexadecimal, CodeNumberOctal, CodeNumberBinary, CodeNumberDecimal,
			CodeComment, CodeOperator, CodeOperatorAssignment, CodeOperatorArithmetic,
			CodeOperatorLogical, CodeOperatorBinary, CodePunctuation, CodeDelimiter}
		if !reflect.DeepEqual(gotTypes, expectTypes) {
			t.Errorf("got default types %v, expect %+v", gotTypes, expectTypes)
		}
	}

	textStyle = "Stop.Error bold"
	style, err = NewStyle(textStyle)
	if err == nil {
		t.Errorf("unexpected nil error")
	}

	textStyle = "toto badam"
	style, err = NewStyle(textStyle)
	if err == nil {
		t.Errorf("unexpected nil error")
	}
	textStyle = "Code badam"
	style, err = NewStyle(textStyle)
	if err == nil {
		t.Errorf("unexpected nil error")
	}

}
