package clrfmt

import (
	"bytes"
	"colorize/clrcore"
	"errors"
	"io"
	"testing"
)

func TestHTML(t *testing.T) {
	def := &clrcore.LexerDef{
		Name: "TestLexerEnginePushPop",
		InitFunc: func(d *clrcore.LexerDef) {
			d.Modes = []*clrcore.LexerDefMode{
				{Name: "root", Rules: []clrcore.LexerDefRule{
					clrcore.WhiteSpaceRule, clrcore.NewLineRule,
					&clrcore.RegexDefRule{Re: "[A-Z]+", Do: clrcore.All(clrcore.PopMatch(clrcore.CodeIdentifier), clrcore.ScoreAdd(5))},
					&clrcore.RegexDefRule{Re: `\(`, Do: clrcore.All(clrcore.PopMatch(clrcore.CodeDelimiter), clrcore.PushMode("param"))},
					&clrcore.RegexDefRule{Re: `{`, Do: clrcore.All(clrcore.PopMatch(clrcore.CodeDelimiter), clrcore.PushMode("root"))},
					&clrcore.RegexDefRule{Re: `}`, Do: clrcore.All(clrcore.PopMatch(clrcore.CodeDelimiter), clrcore.PopMode())},
				}},
				{Name: "param", Rules: []clrcore.LexerDefRule{
					clrcore.WhiteSpaceRule, clrcore.NewLineRule,
					&clrcore.RegexDefRule{Re: "[a-z]+", Do: clrcore.PopMatch(clrcore.CodeIdentifierType)},
					&clrcore.RegexDefRule{Re: `\(`, Do: clrcore.All(clrcore.PopMatch(clrcore.CodeDelimiter), clrcore.PushMode("param"))},
					&clrcore.RegexDefRule{Re: `\)`, Do: clrcore.All(clrcore.PopMatch(clrcore.CodeDelimiter), clrcore.PopMode())},
				}},
			}
		},
	}

	lexerInfo := &clrcore.LexerInfo{
		Names: []string{"test"},
		NewLexer: func(text string, stopMarkers ...string) (clrcore.Lexer, error) {
			return clrcore.NewLexerEngine(def, text, stopMarkers, nil)
		},
	}
	var buf bytes.Buffer
	score, n, err := HTML(&buf, lexerInfo, "AB{ (a b(c d)) } CD")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if score != 10 {
		t.Errorf("got score %d, expected %d", score, 10)
	}
	if n != len(buf.String()) {
		t.Errorf("get n %d, expected %d", n, len(buf.String()))
	}
	expect := `<scan class="q">AB</span><scan class="at">{</span><scan class="f"> </span><scan class="at">(</span><scan class="u">a</span><scan class="f"> </span><scan class="u">b</span><scan class="at">(</span><scan class="u">c</span><scan class="f"> </span><scan class="u">d</span><scan class="at">)</span><scan class="at">)</span><scan class="f"> </span><scan class="at">}</span><scan class="f"> </span><scan class="q">CD</span>`
	if buf.String() != expect {
		t.Errorf("got:\n%s\n, expect:\n%s", buf.String(), expect)
	}
}

func TestCSS1(t *testing.T) {
	style, err := clrcore.NewStyle(`
	Code.Identifier text#FF0000
	Code.Comment text#00FF00
	Code.Identifier.Variable italic
	Code.Identifier.Keyword bold
	`)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, 3000))
	_, err = CSS(buf, style, "<code>", "<code><pre>")
	if err != nil {
		t.Errorf("unexpected error:%s", err)
	} else {
		expect := `<code> .aa, <code><pre> .aa, /*                Code.String */
<code> .ab, <code><pre> .ab, /*         Code.String.Single */
<code> .ac, <code><pre> .ac, /*         Code.String.Double */
<code> .ad, <code><pre> .ad, /*            Code.String.Raw */
<code> .ae, <code><pre> .ae, /*        Code.String.Unicode */
<code> .af, <code><pre> .af, /*      Code.String.Multiline */
<code> .ag, <code><pre> .ag, /*                Code.Number */
<code> .ah, <code><pre> .ah, /*        Code.Number.Integer */
<code> .ai, <code><pre> .ai, /*    Code.Number.Hexadecimal */
<code> .aj, <code><pre> .aj, /*          Code.Number.Octal */
<code> .ak, <code><pre> .ak, /*         Code.Number.Binary */
<code> .al, <code><pre> .al, /*        Code.Number.Decimal */
<code> .an, <code><pre> .an, /*              Code.Operator */
<code> .ao, <code><pre> .ao, /*   Code.Operator.Assignment */
<code> .ap, <code><pre> .ap, /*   Code.Operator.Arithmetic */
<code> .aq, <code><pre> .aq, /*      Code.Operator.Logical */
<code> .ar, <code><pre> .ar, /*       Code.Operator.Binary */
<code> .as, <code><pre> .as, /*           Code.Punctuation */
<code> .at, <code><pre> .at, /*             Code.Delimiter */
<code> .e , <code><pre> .e , /*                       Text */
<code> .f , <code><pre> .f , /*            Text.WhiteSpace */
<code> .g , <code><pre> .g , /*               Text.NewLine */
<code> .h , <code><pre> .h , /*               Text.Invalid */
<code> .i , <code><pre> .i , /*           Text.Punctuation */
<code> .j , <code><pre> .j , /* Text.Punctuation.Separator */
<code> .k , <code><pre> .k , /* Text.Punctuation.Delimiter */
<code> .l , <code><pre> .l , /*              Text.Operator */
<code> .m , <code><pre> .m , /*                  Text.Word */
<code> .n , <code><pre> .n , /*                Text.Number */
<code> .o , <code><pre> .o , /*                 Text.Other */
<code> .p , <code><pre> .p   /*                       Code */ {} 
<code> .am, <code><pre> .am  /*               Code.Comment */ {font-color: #00ff00} 
<code> .q , <code><pre> .q , /*            Code.Identifier */
<code> .s , <code><pre> .s , /*   Code.Identifier.Function */
<code> .t , <code><pre> .t , /*     Code.Identifier.Method */
<code> .u , <code><pre> .u , /*       Code.Identifier.Type */
<code> .v , <code><pre> .v , /*      Code.Identifier.Class */
<code> .w , <code><pre> .w , /*  Code.Identifier.Namespace */
<code> .y , <code><pre> .y , /*    Code.Identifier.Literal */
<code> .z , <code><pre> .z   /*   Code.Identifier.Operator */ {font-color: #ff0000} 
<code> .r , <code><pre> .r   /*   Code.Identifier.Variable */ {font-style: italic} 
<code> .x , <code><pre> .x   /*    Code.Identifier.Keyword */ {font-weight: bold} 
`
		styleStr := buf.String()
		if styleStr != expect {
			t.Errorf("got:\n%s-----\nexpected:\n%s-----", styleStr, expect)
		}
	}
}

type DummyWriter struct {
	w      io.Writer
	n, max int
}

func (w *DummyWriter) Write(d []byte) (int, error) {
	if w.n+len(d) < w.max {
		n, err := w.w.Write(d)
		w.n += n
		return n, err
	}
	n, err := w.w.Write(d[:w.max-w.n])
	w.n += n
	if err == nil {
		err = errors.New("dummy error")
	}
	return n, err
}

func TestCSS2(t *testing.T) {
	style, err := clrcore.NewStyle(`
	Code.Identifier text#FF0000 back#000055
	Code.Comment text#00FF00
	Code.Identifier.Variable italic
	Code.Identifier.Keyword bold
	`)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, 3000))
	n, err := CSS(buf, style)
	if err != nil {
		t.Errorf("unexpected error:%s", err)
	} else {
		expect := `.aa, /*                Code.String */
.ab, /*         Code.String.Single */
.ac, /*         Code.String.Double */
.ad, /*            Code.String.Raw */
.ae, /*        Code.String.Unicode */
.af, /*      Code.String.Multiline */
.ag, /*                Code.Number */
.ah, /*        Code.Number.Integer */
.ai, /*    Code.Number.Hexadecimal */
.aj, /*          Code.Number.Octal */
.ak, /*         Code.Number.Binary */
.al, /*        Code.Number.Decimal */
.an, /*              Code.Operator */
.ao, /*   Code.Operator.Assignment */
.ap, /*   Code.Operator.Arithmetic */
.aq, /*      Code.Operator.Logical */
.ar, /*       Code.Operator.Binary */
.as, /*           Code.Punctuation */
.at, /*             Code.Delimiter */
.e , /*                       Text */
.f , /*            Text.WhiteSpace */
.g , /*               Text.NewLine */
.h , /*               Text.Invalid */
.i , /*           Text.Punctuation */
.j , /* Text.Punctuation.Separator */
.k , /* Text.Punctuation.Delimiter */
.l , /*              Text.Operator */
.m , /*                  Text.Word */
.n , /*                Text.Number */
.o , /*                 Text.Other */
.p   /*                       Code */ {}
.am  /*               Code.Comment */ {font-color: #00ff00}
.q , /*            Code.Identifier */
.s , /*   Code.Identifier.Function */
.t , /*     Code.Identifier.Method */
.u , /*       Code.Identifier.Type */
.v , /*      Code.Identifier.Class */
.w , /*  Code.Identifier.Namespace */
.y , /*    Code.Identifier.Literal */
.z   /*   Code.Identifier.Operator */ {font-color: #ff0000; background-color: #000055}
.r   /*   Code.Identifier.Variable */ {font-style: italic}
.x   /*    Code.Identifier.Keyword */ {font-weight: bold}
`
		styleStr := buf.String()
		if styleStr != expect {
			t.Errorf("got:\n%s-----\nexpected:\n%s-----", styleStr, expect)
		}
		if n != buf.Len() {
			t.Errorf("got n %d, expected %d", n, buf.Len())
		}
	}
	buf.Reset()
	dw := &DummyWriter{w: buf, max: 75}
	n, err = CSS(dw, style)
	if err == nil {
		t.Errorf("unexpected nil error")
	}
	buf.Reset()
	dw = &DummyWriter{w: buf, max: 1168}
	n, err = CSS(dw, style)
	if err == nil {
		t.Errorf("unexpected nil error")
	}

}
