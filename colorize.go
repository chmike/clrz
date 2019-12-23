package clrz

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/chmike/clrz/clrcore"
	"github.com/chmike/clrz/clrfmt"
)

// FormatCSS writes into w a list of CSS classes with styles definition for HTML formatted text.
// The textStyle is a multiline string specifying the style of different lexeme types.
// A line starts with a lexeme type name (e.g. Code.String), and is followed by the style
// specification. The style specification is a combination of "italic", "bold", "text#RRGGBB",
// "back#RRGGBB" where RRGGBB is the red, green and blue intensity in hexadecimal.
//
//	    Code.Identifier text#FF0000
//      Code.Comment text#a0a0a0
//      Code.Identifier.Variable italic
//      Code.Identifier.Keyword bold
//      Code.String text#fc03df
//
// The prefixes are tages les "<code>" or "<code><pre>" to limit the scope of the defined classes.
func FormatCSS(w io.Writer, textStyle string, prefix ...string) (int, error) {
	style, err := clrcore.NewStyle(textStyle)
	if err != nil {
		return 0, err
	}
	return clrfmt.CSS(w, style, prefix...)
}

// FormatHTML return an HTML encoded string using the CSS style classes.
// When more than one or no language are specified, the language with highest score is picked.
func FormatHTML(w io.Writer, text string, lang ...string) (int, error) {
	if len(lang) == 1 {
		_, n, err := clrfmt.HTML(w, clrcore.LexerByName(lang[0]), text)
		return n, err
	}
	var bestScore int
	var bestLexerInfo *clrcore.LexerInfo
	var lexerInfos []*clrcore.LexerInfo
	if len(lang) == 0 {
		lexerInfos = clrcore.Lexers()
	} else {
		unique := make(map[*clrcore.LexerInfo]struct{})
		lexerInfos = make([]*clrcore.LexerInfo, 0, len(lang))
		for _, name := range lang {
			lexerInfo := clrcore.LexerByName(name)
			if lexerInfo == nil {
				continue
			}
			if _, ok := unique[lexerInfo]; ok {
				continue
			}
			unique[lexerInfo] = struct{}{}
			lexerInfos = append(lexerInfos, lexerInfo)
		}
	}
	for _, lexerInfo := range lexerInfos {
		score, _, err := clrfmt.HTML(ioutil.Discard, lexerInfo, text)
		if err != nil {
			if len(lang) == 1 {
				return 0, err
			}
			continue
		}
		if score > bestScore {
			bestLexerInfo = lexerInfo
			bestScore = score
		}
	}
	if bestLexerInfo == nil {
		return 0, fmt.Errorf("failed HTML formatting: no matching lexer found")
	}
	_, n, _ := clrfmt.HTML(w, bestLexerInfo, text)
	return n, nil
}
