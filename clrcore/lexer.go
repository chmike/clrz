package clrcore

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LexerInfo groups information of a lexer and a function to instatiate a lexer.
type LexerInfo struct {
	// Names is a list or names identifying the language.
	// The first wiil be considered the canonical name.
	Names []string
	// MimeTypes identifying text in this language.
	MimeTypes []string
	// File name patterns used to identify language by file names.
	FileNames []string
	// NewLexer is a function instantiating a lexer ready to parse the given text.
	// The lexer will return a StopLexer when one of the stopMarkers is found in the text.
	NewLexer func(text string, stopMarkers ...string) (Lexer, error)
}

// A Lexer has a lexeme iterator method NextLexeme() iterator methods.
type Lexer interface {
	// NextLexeme return the next lexeme extracted from the input text until a stop
	// lexeme is output. The stop lexeme is then returned forever.
	NextLexeme() Lexeme

	// RemainingText return the remaining text to parse.
	RemainingText() string

	// Score returns a matching score for the parsed language. Its value only
	// make sense once a Stop lexeme has been reached. When trying multiple lexer
	// on a piece of text, the lexer returning the highest score should be picked.
	Score() int
}

var (
	lexersList       = make([]*LexerInfo, 0, 10)     // list of all lexers
	lexersByName     = make(map[string]*LexerInfo)   // index lexers by name
	lexersByMimeType = make(map[string][]*LexerInfo) // index of lexers by mimeTypes
	lexersByFileName = make(map[string][]*LexerInfo) // index of lexers by file name pattern
)

// RegisterLexer register a Lexer. Panics il the lexer is invalid.
func RegisterLexer(l *LexerInfo) {
	if l == nil {
		return
	}
	if len(l.Names) == 0 {
		panic("lexer has no names defined")
	}
	if l.NewLexer == nil {
		panic("lexer has no NewLexer function defined")
	}
	// check for name duplicates
	for _, name := range l.Names {
		if _, ok := lexersByName[name]; ok {
			panic(fmt.Sprintf("lexer name %q already registered", name))
		}
	}
	// check that fileNames are valid
	for _, fileName := range l.FileNames {
		if _, err := filepath.Match(fileName, "     "); err != nil {
			panic(fmt.Sprintf("lexer %q has invalid file name pattern %q", l.Names[0], fileName))
		}
	}
	for _, name := range l.Names {
		lexersByName[name] = l
	}
	lexersList = append(lexersList, l)
	for _, mimeType := range l.MimeTypes {
		lexersByMimeType[mimeType] = append(lexersByMimeType[mimeType], l)
	}
	for _, fileName := range l.FileNames {
		lexersByFileName[fileName] = append(lexersByFileName[fileName], l)
	}
}

// Lexers return a copy of list of all the LexerInfo.
func Lexers() []*LexerInfo {
	return append([]*LexerInfo(nil), lexersList...)
}

// LexerByName return the lexer associated to the given name.
func LexerByName(name string) *LexerInfo {
	return lexersByName[name]
}

// LexersByMimeType return a copy of the list of lexers associated to the given mime type.
func LexersByMimeType(mimeType string) []*LexerInfo {
	q := lexersByMimeType[mimeType]
	if q == nil {
		return nil
	}
	return append([]*LexerInfo(nil), q...)
}

// LexersByFileName return a list of lexers whose file name pattern match the given file name.
func LexersByFileName(fileName string) []*LexerInfo {
	if fileName == "" || fileName[len(fileName)-1] == os.PathSeparator {
		return nil
	}
	for i := len(fileName) - 2; i >= 0; i-- {
		if fileName[i] == os.PathSeparator {
			fileName = fileName[i+1:]
			break
		}
	}
	matchedLexers := make(map[*LexerInfo]struct{})
	for k, v := range lexersByFileName {
		if ok, _ := filepath.Match(k, fileName); ok {
			for _, l := range v {
				matchedLexers[l] = struct{}{}
			}
		}
	}
	res := make([]*LexerInfo, 0, len(matchedLexers))
	for l := range matchedLexers {
		res = append(res, l)
	}
	return res
}

// LexerByScore select the lexer from the lexers slice with the highest score.
// When two or more lexers have the same highest score, the first one is picked.
// It is thus advised to order the lexers by decreasing preference order.
// If lexers is nil, the search is performed on all lexers.
func LexerByScore(text string, lexers []*LexerInfo, stopMarkers ...string) *LexerInfo {
	if lexers == nil {
		lexers = lexersList
	}
	var bestScore int
	var bestLexer *LexerInfo
	for _, lexer := range lexers {
		l, err := lexer.NewLexer(text, stopMarkers...)
		if err != nil || l == nil {
			continue
		}
		for !l.NextLexeme().Type.IsA(Stop) {
		}
		if l.Score() > bestScore {
			bestScore = l.Score()
			bestLexer = lexer
		}
	}
	return bestLexer
}

// Normalize replaces all \r, \r\n, and \n\r by \n, replaces tabs by tabLen spaces when tabSize is >= 0,
// and skip all other characters below ' ' and DEL.
func Normalize(str string, tabLen int) string {
	var buf bytes.Buffer
	var tabStr = "\t"
	if tabLen >= 0 {
		buf.Grow(len(str) + tabLen*len(str)/10)
		tabStr = strings.Repeat(" ", tabLen)
	} else {
		buf.Grow(len(str))
	}
	var beg, end int
	const bitMask uint64 = (1 << byte('\n')) | (1 << byte('\r')) | (1 << byte('\t'))
	for end < len(str) {
		c := str[end]
		if c >= 32 && c != 127 {
			end++
			continue
		}
		// skip non printable chars except if \r, \n or \t
		if (bitMask & (1 << c)) == 0 {
			if beg != end {
				buf.WriteString(str[beg:end])
			}
			end++
			beg = end
			continue
		}
		switch c {
		case '\r':
			if beg != end {
				buf.WriteString(str[beg:end])
			}
			buf.WriteByte('\n')
			end++
			if end < len(str) && str[end] == '\n' {
				end++
			}
			beg = end
		case '\n':
			if beg != end {
				buf.WriteString(str[beg:end])
			}
			buf.WriteByte('\n')
			end++
			if end < len(str) && str[end] == '\r' {
				end++
			}
			beg = end
		case '\t':
			if beg != end {
				buf.WriteString(str[beg:end])
			}
			buf.WriteString(tabStr)
			end++
			beg = end
		}
	}
	if beg == 0 && end == len(str) {
		return str
	}
	if beg != end {
		buf.WriteString(str[beg:end])
	}
	return buf.String()
}
