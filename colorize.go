package clrz

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/chmike/clrz/clrcore"
	"github.com/chmike/clrz/clrfmt"
)

// FormatCSS return a CSS encode style for HTML formatted text.
// ... need a detailed explanation of textStyle specification.
func FormatCSS(w io.Writer, textStyle string, prefix ...string) (int, error) {
	style, err := clrcore.NewStyle(textStyle)
	if err != nil {
		return 0, err
	}
	return clrfmt.CSS(w, style, prefix...)
}

// FormatHTML return an HTML encoded string using the CSS style classes.
func FormatHTML(w io.Writer, text string, lang ...string) (int, error) {
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
