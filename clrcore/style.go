package clrcore

import (
	"bytes"
	"fmt"
	"strings"
)

// TypeStyle is an encoded style for a lexeme type.
type TypeStyle uint64

const (
	// TextColorFlag is set if a text color is defined.
	textColorFlag TypeStyle = 1 << iota
	// BackColorFlag is set if a background color is defined.
	backColorFlag
	// ItalicFlag is the text is italic.
	italicFlag
	// BoldFlag is the text is bold.
	boldFlag
)

// MakeStyle return a Style set with the given value: "italic" sets the italic
// flag, "bold" sets the Bold flag, "text#RRGGBB" sets the text color,
// and "back#RRGGBB" sets the bacground color. RR, GG, BB are the color values
// for red (RR), green (GG) and blue (BB) in hexadecimal.
// Return an error if the style attribute is not rcognized.
//
//   MakeStyle() // default style
//   MakeStyle("bold", "italic") // style is bold and italic with default color
//   MakeStyle("text#123456") // style is text in color #123456
func MakeStyle(style ...string) (s TypeStyle, err error) {
	for _, str := range style {
		switch str {
		case "bold":
			s |= boldFlag
		case "italic":
			s |= italicFlag
		default:
			var red, green, blue TypeStyle
			if n, err := fmt.Sscanf(str, "text#%2x%2x%2x", &red, &green, &blue); err == nil && n == 3 {
				s &= ^TypeStyle(0xFFFFFF0000)
				s |= textColorFlag | red<<16 | green<<24 | blue<<32
			} else if n, err := fmt.Sscanf(str, "back#%2x%2x%2x", &red, &green, &blue); err == nil && n == 3 {
				s &= ^TypeStyle(0xFFFFFF0000000000)
				s |= backColorFlag | red<<40 | green<<48 | blue<<56
			} else {
				return 0, fmt.Errorf("invalid style attribute '%s'", str)
			}
		}
	}
	return
}

// Bold return true if the type style is bold.
func (s TypeStyle) Bold() bool {
	return s&boldFlag != 0
}

// Italic return true if the type style is italic.
func (s TypeStyle) Italic() bool {
	return s&italicFlag != 0
}

// HasTextColor return true if the type style has a text color defined.
func (s TypeStyle) HasTextColor() bool {
	return s&textColorFlag != 0
}

// TextColor return the text color.
func (s TypeStyle) TextColor() (red byte, green byte, blue byte) {
	return byte(s >> 16), byte(s >> 24), byte(s >> 32)
}

// HasBackColor return true if the type style has a back color defined.
func (s TypeStyle) HasBackColor() bool {
	return s&backColorFlag != 0
}

// BackColor return true if the type style has a back color defined.
func (s TypeStyle) BackColor() (red byte, green byte, blue byte) {
	return byte(s >> 40), byte(s >> 48), byte(s >> 56)
}

func (s TypeStyle) String() string {
	if uint16(s) == 0 {
		return ""
	}
	var buf bytes.Buffer
	if s.Italic() {
		buf.WriteString("italic ")
	}
	if s.Bold() {
		buf.WriteString("bold ")
	}
	if s.HasTextColor() {
		r, g, b := s.TextColor()
		fmt.Fprintf(&buf, "text#%02X%02X%02X ", r, g, b)
	}
	if s.HasBackColor() {
		r, g, b := s.BackColor()
		fmt.Fprintf(&buf, "back#%02X%02X%02X ", r, g, b)
	}
	return buf.String()[:buf.Len()-1]
}

// Style defines a formatting style
type Style map[TypeStyle][]*LexemeType

type forAllTypesFunc func(*LexemeType, TypeStyle, forAllTypesFunc)

// NewStyle return a Style initialized with the text style definition.
// Each non empty line specify the style for a lexeme type and all it's subtypes.
// and are followed by the style specification which will apply by default
// to all the  type's children.
// The Stop lexeme types can't have a style defined.
func NewStyle(text string) (Style, error) {
	index := make(map[*LexemeType]TypeStyle)
	for i, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		t := LexemeTypeByName(fields[0])
		if t == nil {
			return nil, fmt.Errorf("unknown LexemeType '%s'", fields[0])
		}
		if t.Class() == Stop {
			return nil, fmt.Errorf("LexemeType '%s' can't have style", t.Name)
		}
		s, err := MakeStyle(fields[1:]...)
		if err != nil {
			return nil, fmt.Errorf("%s in line %d", err, i+1)
		}
		index[t] = s
	}
	// invert type style index
	style := make(map[TypeStyle][]*LexemeType)
	forAllTypesInClass := func(t *LexemeType, s TypeStyle, f forAllTypesFunc) {
		if tmp, ok := index[t]; ok {
			s = tmp
		}
		style[s] = append(style[s], t)
		for _, c := range t.Children {
			f(c, s, f)
		}
	}
	for _, t := range LexemeClassTypes() {
		if t == Stop {
			continue
		}
		forAllTypesInClass(t, 0, forAllTypesInClass)
	}
	return style, nil
}
