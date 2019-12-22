package clrcore

// PopLexeme extract lexeme of type Type up to end from target.
func (l *Lexeme) PopLexeme(Type *LexemeType, end int) (lexeme Lexeme) {
	lexeme = Lexeme{Type: Type, Str: l.Str[:end]}
	l.Str = l.Str[end:]
	return lexeme
}

// PopSingleQuotedString extracts the single quoted string from the front of the
// target lexeme, and return it.
// A single quoted string may contain \' and is terminated by ', \n (not included) or
// the end of the target string.
// It requires that the target string is at least 1 char long and start with '.
func (l *Lexeme) PopSingleQuotedString() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 1; i < len(l.Str); i++ {
		c := l.Str[i]
		if c == '\'' && l.Str[i-1] != '\\' {
			end = i + 1
			break
		}
		if c == '\n' {
			end = i
			break
		}
	}
	return l.PopLexeme(CodeStringSingle, end)
}

// PopDoubleQuotedString extracts the double quoted string from the front of the
// target lexeme, and return it.
// A double quoted string may contain \" and is terminated by ", \n (not included) or
// the end of the target string.
// It requires that the target string is at least 1 char long and start with ".
func (l *Lexeme) PopDoubleQuotedString() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 1; i < len(l.Str); i++ {
		c := l.Str[i]
		if c == '"' && l.Str[i-1] != '\\' {
			end = i + 1
			break
		}
		if c == '\n' {
			end = i
			break
		}
	}
	return l.PopLexeme(CodeStringDouble, end)
}

// PopBackTickString extracts the back tick quoted string from the front of the
// target lexeme, and return it.
// A back tick quoted string is terminated by ` or the end of the target string.
// It requires that the target string is at least 1 char long and start with `.
func (l *Lexeme) PopBackTickString() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 1; i < len(l.Str); i++ {
		if l.Str[i] == '`' {
			end = i + 1
			break
		}
	}
	return l.PopLexeme(CodeStringMultiline, end)
}

// PopTripleDoubleQuotedString extracts the triple double quoted string from the
// front of the target lexeme, and return it.
// A triple double quoted string is terminated by """ or the end of the target string.
// It requires that the target string is at least 3 char long and start with """.
func (l *Lexeme) PopTripleDoubleQuotedString() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 5; i < len(l.Str); i += 3 {
		if l.Str[i] == '"' {
			if l.Str[i-1] == '"' && l.Str[i-2] == '"' {
				end = i + 1
				break
			}
			i -= 2
		}
	}
	return l.PopLexeme(CodeStringMultiline, end)
}

// PopOneCharLineComment extracts a comment ending at \n (not included) from the
// front of the target lexeme, and return it.
// It requires that the target string is at least 1 char long.
// If no \n is encountered, the full target string is extracted and returned as
// lexeme. Typical use is with line comments starting whith character #, ', ! or ;.
func (l *Lexeme) PopOneCharLineComment() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 1; i < len(l.Str); i++ {
		if l.Str[i] == '\n' {
			end = i
			break
		}
	}
	return l.PopLexeme(CodeComment, end)
}

// PopTwoCharLineComment extracts a comment starting with two chars (e.g. //) and
// ending at \n (not included) from the front of the target lexeme, and return it.
// It requires that the target string is at least 2 char long.
// If no \n is encountered, the full target string is extracted and returned as
// lexeme.
func (l *Lexeme) PopTwoCharLineComment() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 2; i < len(l.Str); i++ {
		if l.Str[i] == '\n' {
			end = i
			break
		}
	}
	return l.PopLexeme(CodeComment, end)
}

// PopSlashStarComment extracts a comment starting with /* and ending with
// */ from the front of the target lexeme, and return it.
// Requires that the target string is 2 char long and starts with /*.
// If no */ is encountered, the full target string is extracted and returned as
// lexeme.
func (l *Lexeme) PopSlashStarComment() (lexeme Lexeme) {
	end := len(l.Str)
	for i := 3; i < len(l.Str); i += 2 {
		if l.Str[i] == '*' {
			i--
			continue
		}
		if l.Str[i] == '/' && l.Str[i-1] == '*' {
			end = i + 1
			break
		}
	}
	return l.PopLexeme(CodeComment, end)
}

// PopDecimalNumber extracts a decimal number from the front of the
// target lexeme, and return it.
// It requires that the char . is before end.
func (l *Lexeme) PopDecimalNumber(end int) (lexeme Lexeme) {
	// scann decimal digits
	for ; end < len(l.Str); end++ {
		c := l.Str[end]
		// if reached end of decimal digits
		if c < '0' || c > '9' {
			// if not followed with an exponent, done
			if c|0x20 != 'e' {
				break
			}
			exp := end
			if exp++; exp == len(l.Str) {
				break
			}
			c = l.Str[exp]
			// scan sign if any
			if c == '+' || c == '-' {
				exp++
				if len(l.Str) == exp {
					break
				}
				c = l.Str[exp]
			}
			// if e or exponent sign not followed with a digit, no exponent
			if c < '0' || c > '9' {
				break
			}
			// scan exponent digits
			for end = exp + 1; end < len(l.Str); end++ {
				c = l.Str[end]
				if c < '0' || c > '9' {
					break
				}
			}
			break
		}
	}
	return l.PopLexeme(CodeNumberDecimal, end)
}

// PopNumber extracts the number from the front of the
// target lexeme, and return it, or return a nil lexeme if not a number.
func (l *Lexeme) PopNumber() (lexeme Lexeme) {
	if len(l.Str) == 0 {
		return
	}
	var end int
	c := l.Str[0]
	// scan sign if any
	if c == '-' || c == '+' {
		end++
		if len(l.Str) == end {
			return
		}
		c = l.Str[1]
	}
	// if . and followed by a digit, it's a decimal number, otherwise it's not a number.
	if c == '.' {
		end++
		if len(l.Str) == end {
			return
		}
		c = l.Str[end]
		if c < '0' || c > '9' {
			return
		}
		return l.PopDecimalNumber(end + 1)
	}
	// if not a digit, it's not a number
	if c < '0' || c > '9' {
		return
	}
	// if first number digit is 0
	if c == '0' {
		if end++; end == len(l.Str) {
			return l.PopLexeme(CodeNumberInteger, end)
		}
		c = l.Str[end]
		// if it's an octal digit, it's an octal number
		if c >= '0' && c <= '7' {
			for end++; end < len(l.Str); end++ {
				c := l.Str[end]
				if c < '0' || c > '7' {
					break
				}
			}
			return l.PopLexeme(CodeNumberOctal, end)
		}
		// if it's an X followed by a hexadecial digit, it's a hexadecimal number
		if (c | 0x20) == 'x' {
			if end++; end == len(l.Str) {
				return
			}
			c = l.Str[end] | 0x20
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
				return
			}
			for end++; end < len(l.Str); end++ {
				c := l.Str[end]
				if c >= '0' && c <= '9' {
					continue
				}
				c |= 0x20
				if c < 'a' || c > 'f' {
					break
				}
			}
			return l.PopLexeme(CodeNumberHexadecimal, end)
		}
		// if it's a decimal number
		if c == '.' {
			return l.PopDecimalNumber(end + 1)
		}
		// single 0 digit integer
		if c < '0' || c > '9' {
			return l.PopLexeme(CodeNumberInteger, end)
		}
		// not a number
		return
	}
	for end++; end < len(l.Str); end++ {
		c = l.Str[end]
		if c < '0' || c > '9' {
			if c == '.' {
				return l.PopDecimalNumber(end + 1)
			}
			break
		}
	}
	return l.PopLexeme(CodeNumberInteger, end)
}

// PopWhiteSpaces extracts white spaces from the front of the
// target lexeme, and return it.
// A white space character is smaller or equal than ' ' and not '\n', '\r' or DEL.
// It requires the target string is at least 1 char long and that the first char
// is a white space.
func (l *Lexeme) PopWhiteSpaces() (lexeme Lexeme) {
	for end := 1; end < len(l.Str); end++ {
		c := l.Str[end]
		if c != 127 && (c > ' ' || c == '\n' || c == '\r') {
			return l.PopLexeme(TextWhiteSpace, end)
		}
	}
	return l.PopLexeme(TextWhiteSpace, len(l.Str))
}

// PopNewLines extracts newlines from the front of the
// target lexeme, and return it. A new line character is equal to '\n' or '\r.
// It requires the target string is at least 1 char long and that the first char
// is a new line character.
func (l *Lexeme) PopNewLines() (lexeme Lexeme) {
	for end := 1; end < len(l.Str); end++ {
		c := l.Str[end]
		if c != '\n' && c != '\r' {
			return l.PopLexeme(TextNewLine, end)
		}
	}
	return l.PopLexeme(TextNewLine, len(l.Str))
}

// PopASCIIIdentifier extracts an ASCII identifier from the front of the
// target lexeme, and return it.
// It requires the target string is at least 1 char long and that the first char
// is an ASCII letter or underscore.
func (l *Lexeme) PopASCIIIdentifier() (lexeme Lexeme) {
	for end := 1; end < len(l.Str); end++ {
		c := l.Str[end]
		if c == '_' || (c >= '0' && c <= '9') {
			continue
		}
		c |= 0x20
		if c < 'a' || c > 'z' {
			return l.PopLexeme(CodeIdentifier, end)
		}
	}
	return l.PopLexeme(CodeIdentifier, len(l.Str))
}
