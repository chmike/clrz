package clrcore

import (
	"fmt"
	"sort"
)

// A LexemeType identifies a type of lexeme.
type LexemeType struct {
	Parent   *LexemeType   // nil if LexemType is a lexeme class
	Children []*LexemeType // nil if has no children
	Name     string        // full name (e.g. "MultiLine.Comment.Lexeme")
}

func (t *LexemeType) String() string {
	return t.Name
}

var lexemeClassTypes = make([]*LexemeType, 0, 8)

// LexemeClassTypes return a list of all lexeme type classes.
// A lexeme class type is a lexeme type with a nil parent.
func LexemeClassTypes() []*LexemeType {
	return append([]*LexemeType(nil), lexemeClassTypes...)
}

var lexemeTypes = make([]*LexemeType, 0, 32)

// LexemeTypes return a list of all lexeme types.
func LexemeTypes() []*LexemeType {
	return append([]*LexemeType(nil), lexemeTypes...)
}

var lexemesTypesByName = make(map[string]*LexemeType)

// LexemeTypeByName return the lexeme type with the given name or nil if none is found.
func LexemeTypeByName(name string) *LexemeType {
	return lexemesTypesByName[name]
}

// NewLexemeType instantiate and register a new lexeme type.
// The intended use is to declare lexeme types at the global scope:
// var CommentMultiline = NewLexemeType(Comment, "Comment.Multiline").
// Duplicate name panics.
func NewLexemeType(parent *LexemeType, name string) *LexemeType {
	if name == "" {
		panic("lexeme type name is empty string")
	}
	if _, ok := lexemesTypesByName[name]; ok {
		panic(fmt.Sprintf("lexeme name %q already defined", name))
	}
	t := &LexemeType{Parent: parent, Name: name}
	lexemesTypesByName[name] = t
	s := &lexemeClassTypes
	if parent != nil {
		s = &parent.Children
		if *s == nil {
			*s = make([]*LexemeType, 0, 4)
		}
	}
	*s = append(*s, t)
	if len(*s) > 1 {
		sort.Slice(*s, func(i, j int) bool { return (*s)[i].Name < (*s)[i].Name })
	}
	lexemeTypes = append(lexemeTypes, t)
	sort.Slice(lexemeTypes, func(i, j int) bool { return lexemeTypes[i].Name < lexemeTypes[i].Name })
	return t
}

// IsA return true if l2 is the same or a parent of l.
func (l *LexemeType) IsA(l2 *LexemeType) bool {
	if l == l2 {
		return true
	}
	if l2 != nil {
		for l1 := l.Parent; l1 != nil; l1 = l1.Parent {
			if l1 == l2 {
				return true
			}
		}
	}
	return false
}

// Class return the class lexeme type of l.
func (l *LexemeType) Class() *LexemeType {
	for l.Parent != nil {
		l = l.Parent
	}
	return l
}

var (
	// Stop is lexeme type signaling the stop of the lexer. When not empty, Str provides the reason.
	Stop = NewLexemeType(nil, "Stop")
	// StopError is a lexeme signaling an error. Str is the error message.
	StopError = NewLexemeType(Stop, "Stop.Error")
	// StopEndOfString is a lexeme returned when the end of string is reached. Str is "".
	StopEndOfString = NewLexemeType(Stop, "Stop.EndofString")
	// StopLexer is a lexeme returned when the lexer stops.
	// Str contains the stopMarker when this stopped the lexer.
	// If Str is empty, the lexer was stopped because no lexeme could be extracted.
	StopLexer = NewLexemeType(Stop, "Stop.Lexer")

	// Text is a type of lexeme containing a sequence of unicode characters. It may contain white
	// spaces or new lines.
	Text = NewLexemeType(nil, "Text")
	// TextWhiteSpace is a type of lexeme containing white space characters (" ", \t, \a, \b, \f, \v).
	TextWhiteSpace = NewLexemeType(Text, "Text.WhiteSpace")
	// TextNewLine is a type of lexeme containing new line characters (\n,\r,\r\n,\n\r).
	TextNewLine = NewLexemeType(Text, "Text.NewLine")
	// TextInvalid is a sequence of unicade characters that don't constitute a valid lexeme.
	TextInvalid = NewLexemeType(Text, "Text.Invalid")
	// TextPunctuation is a type of lexeme containing text ponctuations.
	TextPunctuation = NewLexemeType(Text, "Text.Punctuation")
	// TextPunctuationSeparator is a type of lexeme containing text separation ponctuations (., ,, ;, :...).
	TextPunctuationSeparator = NewLexemeType(TextPunctuation, "Text.Punctuation.Separator")
	// TextPunctuationDelimiter is a type of lexeme containing text delimiter ponctuations ( {, }, [, ], ...).
	TextPunctuationDelimiter = NewLexemeType(TextPunctuation, "Text.Punctuation.Delimiter")
	// TextOperator is a type of lexeme containing operators ( +, -, *, /, !, <, >, =, ...).
	TextOperator = NewLexemeType(Text, "Text.Operator")
	// TextWord is a type of lexeme containing only unicode letters.
	TextWord = NewLexemeType(Text, "Text.Word")
	// TextNumber is a type of lexeme containing only the digits 0 to 9.
	TextNumber = NewLexemeType(Text, "Text.Number")
	// TextOther is a text lexeme that is not a WhiteSpace, NewLine, Invalid, Punctuation, Operator,
	// Word or Number.
	TextOther = NewLexemeType(Text, "Text.Other")

	// Code is a class of lexemes for code text
	Code = NewLexemeType(nil, "Code")
	// CodeIdentifier is an identifier
	CodeIdentifier = NewLexemeType(Code, "Code.Identifier")
	// CodeIdentifierVariable is a variable name
	CodeIdentifierVariable = NewLexemeType(CodeIdentifier, "Code.Identifier.Variable")
	// CodeIdentifierFunction is a function name
	CodeIdentifierFunction = NewLexemeType(CodeIdentifier, "Code.Identifier.Function")
	// CodeIdentifierMethod is a method name
	CodeIdentifierMethod = NewLexemeType(CodeIdentifier, "Code.Identifier.Method")
	// CodeIdentifierType is a type name
	CodeIdentifierType = NewLexemeType(CodeIdentifier, "Code.Identifier.Type")
	// CodeIdentifierClass is a class name
	CodeIdentifierClass = NewLexemeType(CodeIdentifier, "Code.Identifier.Class")
	// CodeIdentifierNamespace is a name space name
	CodeIdentifierNamespace = NewLexemeType(CodeIdentifier, "Code.Identifier.Namespace")
	// CodeIdentifierKeyword is a language keyword
	CodeIdentifierKeyword = NewLexemeType(CodeIdentifier, "Code.Identifier.Keyword")
	// CodeIdentifierLiteral is a language keyword
	CodeIdentifierLiteral = NewLexemeType(CodeIdentifier, "Code.Identifier.Literal")
	// CodeIdentifierOperator is an operator identifier (e.g. not, or, and)
	CodeIdentifierOperator = NewLexemeType(CodeIdentifier, "Code.Identifier.Operator")
	// CodeString is a string
	CodeString = NewLexemeType(Code, "Code.String")
	// CodeStringSingle is a single quoted string
	CodeStringSingle = NewLexemeType(CodeString, "Code.String.Single")
	// CodeStringDouble is a double quoted string
	CodeStringDouble = NewLexemeType(CodeString, "Code.String.Double")
	// CodeStringRaw is a raw string
	CodeStringRaw = NewLexemeType(CodeString, "Code.String.Raw")
	// CodeStringUnicode is a unicode string
	CodeStringUnicode = NewLexemeType(CodeString, "Code.String.Unicode")
	// CodeStringMultiline is a multiline string
	CodeStringMultiline = NewLexemeType(CodeString, "Code.String.Multiline")
	// CodeNumber is a number
	CodeNumber = NewLexemeType(Code, "Code.Number")
	// CodeNumberInteger is an integer number
	CodeNumberInteger = NewLexemeType(CodeNumber, "Code.Number.Integer")
	// CodeNumberHexadecimal is an hexadecimal integer number
	CodeNumberHexadecimal = NewLexemeType(CodeNumber, "Code.Number.Hexadecimal")
	// CodeNumberOctal is an octal integer number
	CodeNumberOctal = NewLexemeType(CodeNumber, "Code.Number.Octal")
	// CodeNumberBinary is a binary integer number
	CodeNumberBinary = NewLexemeType(CodeNumber, "Code.Number.Binary")
	// CodeNumberDecimal is an decimal number
	CodeNumberDecimal = NewLexemeType(CodeNumber, "Code.Number.Decimal")
	// CodeComment is a comment
	CodeComment = NewLexemeType(Code, "Code.Comment")
	// CodeOperator is an unary or binary operator
	CodeOperator = NewLexemeType(Code, "Code.Operator")
	// CodeOperatorAssignment is an assignment operator
	CodeOperatorAssignment = NewLexemeType(CodeOperator, "Code.Operator.Assignment")
	// CodeOperatorArithmetic is an arithmetic operator
	CodeOperatorArithmetic = NewLexemeType(CodeOperator, "Code.Operator.Arithmetic")
	// CodeOperatorLogical is an logical operator
	CodeOperatorLogical = NewLexemeType(CodeOperator, "Code.Operator.Logical")
	// CodeOperatorBinary is an binary operator
	CodeOperatorBinary = NewLexemeType(CodeOperator, "Code.Operator.Binary")
	// CodePunctuation is a punctuation (e.g. . , or ; )
	CodePunctuation = NewLexemeType(Code, "Code.Punctuation")
	// CodeDelimiter is a delimiter (e.g. (), {}, [])
	CodeDelimiter = NewLexemeType(Code, "Code.Delimiter")
)
