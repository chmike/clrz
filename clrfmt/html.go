package clrfmt

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/chmike/clrz/clrcore"
)

// HTML writes the text formatted in HTML, and return the number of bytes written.
// When the lexer stops, whatever the reason, the remaining text is written out unformatted.
// TODO 1. add formatting options: susbstitute chars (e.g. \n -> <br>), table with line numbers,
// add highlighted section of text. Use io.Writer equivalent ?
// TODO 2. escape < and > characters.
func HTML(w io.Writer, info *clrcore.LexerInfo, text string) (score int, n int, err error) {
	if info == nil {
		return 0, 0, errors.New("LexerInfo is nil")
	}
	lexer, err := info.NewLexer(text)
	if err != nil {
		return 0, 0, err
	}
	var buf bytes.Buffer
	var bytesWritten int
	for {
		lexeme := lexer.NextLexeme()
		if lexeme.Type == nil {
			return 0, bytesWritten, errors.New("lexeme with undefined LexemeType")
		}
		if lexeme.IsA(clrcore.Stop) {
			n, err := w.Write([]byte(lexer.RemainingText()))
			return lexer.Score(), bytesWritten + n, err
		}
		className := typeClassNameMap[lexeme.Type]
		if className == "" {
			return 0, bytesWritten, fmt.Errorf("unknown LexemeType %s", lexeme.Type)
		}
		buf.Reset()
		_, err := fmt.Fprintf(&buf, "<scan class=\"%s\">", className)
		if err != nil {
			return 0, bytesWritten, err
		}
		n, err = w.Write(buf.Bytes())
		bytesWritten += n
		if err != nil {
			return 0, bytesWritten, err
		}
		n, err = w.Write([]byte(lexeme.Str))
		bytesWritten += n
		if err != nil {
			return 0, bytesWritten, err
		}
		n, err = w.Write([]byte("</span>"))
		bytesWritten += n
		if err != nil {
			return 0, bytesWritten, err
		}
	}
}

type classEntry struct {
	className, typeName string
}
type styleEntry struct {
	style   string
	classes []classEntry
}

// CSS return the CSS encoded style.
func CSS(w io.Writer, style clrcore.Style, prefix ...string) (int, error) {
	// build the css style selector format strings
	var selectorFmt, lastSelectorFmt string
	if len(prefix) == 0 {
		selectorFmt = ".%-2[1]s, /* %[3]*[2]s */\n"
		lastSelectorFmt = ".%-2[1]s  /* %[3]*[2]s */ {%[4]s}\n"
	} else {
		selectorFmt = strings.Join(prefix, " .%-2[1]s, ") + " .%-2[1]s, /* %[3]*[2]s */\n"
		lastSelectorFmt = strings.Join(prefix, " .%-2[1]s, ") + " .%-2[1]s  /* %[3]*[2]s */ {%[4]s} \n"
	}
	// build the list of styles and classe names
	const maxStyleStrLen = 20 + 19 + 27 + 33
	var maxTypeNameLength int
	buf := bytes.NewBuffer(make([]byte, 0, maxStyleStrLen))
	styles := make([]styleEntry, len(style))
	styleIdx := 0
	for typeStyle, lexemeTypes := range style {
		var styleStr string
		if typeStyle != 0 {
			buf.Reset()
			if typeStyle.Italic() {
				buf.WriteString("font-style: italic; ")
			}
			if typeStyle.Bold() {
				buf.WriteString("font-weight: bold; ")
			}
			if typeStyle.HasTextColor() {
				r, g, b := typeStyle.TextColor()
				fmt.Fprintf(buf, "font-color: #%02x%02x%02x; ", r, g, b)
			}
			if typeStyle.HasBackColor() {
				r, g, b := typeStyle.BackColor()
				fmt.Fprintf(buf, "background-color: #%02x%02x%02x; ", r, g, b)
			}
			styleStr = buf.String()[:buf.Len()-2]
		}
		styles[styleIdx].style = styleStr
		classes := make([]classEntry, len(lexemeTypes))
		classIdx := 0
		for _, t := range lexemeTypes {
			if className, ok := typeClassNameMap[t]; ok {
				if len(t.Name) > maxTypeNameLength {
					maxTypeNameLength = len(t.Name)
				}
				classes[classIdx].className = className
				classes[classIdx].typeName = t.Name
				classIdx++
			}
		}
		classes = classes[:classIdx]
		sort.Slice(classes, func(i, j int) bool {
			return classes[i].className < classes[j].className
		})
		styles[styleIdx].classes = classes
		styleIdx++
	}
	styles = styles[:styleIdx]
	sort.Slice(styles, func(i, j int) bool { return styles[i].style < styles[j].style })
	// generate the style
	var bytesWritten int
	for _, s := range styles {
		last := len(s.classes) - 1
		for i, c := range s.classes {
			if i == last {
				n, err := fmt.Fprintf(w, lastSelectorFmt, c.className, c.typeName, maxTypeNameLength, s.style)
				bytesWritten += n
				if err != nil {
					return bytesWritten, err
				}
			} else {
				n, err := fmt.Fprintf(w, selectorFmt, c.className, c.typeName, maxTypeNameLength)
				bytesWritten += n
				if err != nil {
					return bytesWritten, err
				}
			}
		}
	}
	return bytesWritten, nil
}

var typeClassNameMap = map[*clrcore.LexemeType]string{}
var classNameToLexemeTypeMap = map[string]*clrcore.LexemeType{}

// generate class names: a - z,aa - az,ba - bz,ca - cz, ...
func classNameFromIdx(idx int) string {
	buf := make([]byte, 8)
	strIdx := len(buf) - 1
	buf[strIdx] = byte('a' + (idx % 26))
	for idx /= 26; idx > 0; idx /= 26 {
		idx--
		strIdx--
		buf[strIdx] = byte('a' + (idx % 26))
	}
	return string(buf[strIdx:])
}

func init() {
	for i, t := range clrcore.LexemeTypes() {
		if t.Class() == clrcore.Stop {
			continue
		}
		className := classNameFromIdx(i)
		typeClassNameMap[t] = className
		classNameToLexemeTypeMap[className] = t
	}
}
