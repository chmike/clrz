package clrcore

import (
	"bytes"
	"fmt"
	"regexp"
	"sync"
)

// LexerDef is a list of LexerDefMode instances initialized once.
type LexerDef struct {
	Name     string
	Modes    []*LexerDefMode  // List of lexer modes, first is entry point.
	InitFunc LexerDefInitFunc // Lexer definition initialization function.
	once     sync.Once        // Ensure the definition is initialized only once at first use.
}

// LexerDefInitFunc is a function that initialize a LexerDef.
// It is called once at first LexerDef use.
type LexerDefInitFunc func(d *LexerDef)

// LexerDefMode is a named group of LexerRules initialized once.
// Rules are executed in sequence until one return true, in which case
// execution restart from the first rule.
// It's an error if no rule return true.
type LexerDefMode struct {
	Name  string         // LexerDefMode name used for referencing (reserved: "root")
	Rules []LexerDefRule // Rules executed in sequence until one return true.
}

// A LexerDefRule defines a condition and action to perform.
type LexerDefRule interface {
	Init() error            // Initialize the rule only once.
	Exec(*LexerEngine) bool // Return true if must exec from first mode rule.
}

// Init initializes the LexerDef and return an error if any.
func (d *LexerDef) Init() (err error) {
	if d.InitFunc == nil {
		return fmt.Errorf("LexerDef '%s' has undefined InitFunc", d.Name)
	}
	d.once.Do(func() {
		d.InitFunc(d)
		if len(d.Modes) == 0 {
			err = fmt.Errorf("LexerDef '%s' has no modes", d.Name)
			return
		}
		for _, m := range d.Modes {
			if len(m.Rules) == 0 {
				err = fmt.Errorf("mode '%s' in LexerDef '%s' has no rules", m.Name, d.Name)
				return
			}
			for i, r := range m.Rules {
				if err = r.Init(); err != nil {
					err = fmt.Errorf("%s (LexerDef='%s', Mode='%s', Rule=%d)", err, d.Name, m.Name, i)
					return
				}
			}
		}
	})
	return
}

// FuncDefRule is for an Exec function defined in the LexerDef.
type FuncDefRule struct {
	ExecFunc FuncDefRuleExec
}

// FuncDefRuleExec executes the condition test and action.
// Return true when execution of rules must restart from the
// first rule of the mode.
type FuncDefRuleExec func(*LexerEngine) bool

// Init initilizes the FuncDefRule.
func (r *FuncDefRule) Init() error {
	return nil
}

// Exec executes the defined function.
func (r *FuncDefRule) Exec(l *LexerEngine) bool {
	return r.ExecFunc(l)
}

// RegexDefRule is a regex rule with an associated function.
type RegexDefRule struct {
	Re   string           // regex pattern to trigger RegexDefRule
	Do   RegexDefRuleFunc // action to perform when this RegexDefRule is triggered
	cp   *regexp.Regexp   // compiled regex RegexDefRule wrapped in "^(?<re>)"
	once sync.Once        // to ensure it's compiled only once
}

// RegexDefRuleFunc is a function called when the associated regex is triggered.
// The returned value is returned as Exec return value.
type RegexDefRuleFunc func(l *LexerEngine, match []int) bool

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Init initialize the regex def rule only once.
func (r *RegexDefRule) Init() (err error) {
	r.once.Do(func() {
		buf := bufPool.Get().(*bytes.Buffer)
		defer bufPool.Put(buf)
		buf.Reset()
		fmt.Fprintf(buf, `\A(?:%s)`, r.Re)
		if r.cp, err = regexp.Compile(buf.String()); err != nil {
			return
		}
	})
	return
}

// Exec executes the regex def rule. Return true if the
func (r *RegexDefRule) Exec(l *LexerEngine) bool {
	match := r.cp.FindStringSubmatchIndex(l.str)
	if match == nil {
		return false
	}
	//fmt.Printf("rule[%d]: %q matched: %q\n", l.ruleIdx, r.Re, l.str[match[0]:match[1]])
	return r.Do(l, match)
}

// All execute list of RegexLexerActionFunc in sequence, abort when l.err is not nil.
// The last action return value yields the return value.
func All(actions ...RegexDefRuleFunc) RegexDefRuleFunc {
	return func(l *LexerEngine, match []int) (res bool) {
		for _, a := range actions {
			res = a(l, match)
			if l.err != nil {
				return
			}
		}
		return
	}
}

// PopMatch returns a function that extracts matching lexemes from the input string,
// and assign the given lexeme type to each group.
// The number of lexeme types must match the number of groups, or it must be one
// if the regexp contains no group.
// Text outside of groups is output as lexeme of type Text.
// Groups containing no chars are not output and their lexeme type is skipped.
func PopMatch(lexemeTypes ...*LexemeType) RegexDefRuleFunc {
	return func(l *LexerEngine, match []int) bool {
		// If the regexp had no groups
		if len(match) == 2 {
			if len(lexemeTypes) != 1 {
				l.err = fmt.Errorf("invalid number of lexemeTypes (LexerDef='%s', Mode='%s', Rule=%d)", l.def.Name, l.mode.Name, l.ruleIdx)
				return true
			}
			// Never occurs because the regex always starts with \A
			// if match[0] != 0 {
			// 	l.QueueLexeme(Lexeme{Type: Text, Str: l.str[:match[0]]})
			// }
			l.QueueLexeme(Lexeme{Type: lexemeTypes[0], Str: l.str[match[0]:match[1]]})
			l.str = l.str[match[1]:]
			return true
		}
		// The regexp has one or more groups
		n := len(match)/2 - 1
		if len(lexemeTypes) != n {
			l.err = fmt.Errorf("invalid number of lexemeTypes (LexerDef='%s', Mode='%s', Rule=%d)", l.def.Name, l.mode.Name, l.ruleIdx)
			return true
		}
		maxEnd := 0
		prevEnd := 0
		for i := 1; i < n; i++ {
			beg, end := match[i*2], match[i*2+1]
			if end <= maxEnd {
				l.err = fmt.Errorf("overlapping regex groups (LexerDef='%s', Mode='%s', Rule=%d)", l.def.Name, l.mode.Name, l.ruleIdx)
				return true
			}
			if prevEnd < beg {
				l.QueueLexeme(Lexeme{Type: Text, Str: l.str[prevEnd:beg]})
			}
			if beg != end {
				l.QueueLexeme(Lexeme{Type: lexemeTypes[i], Str: l.str[beg:end]})
			}
			if maxEnd < end {
				maxEnd = end
			}
			prevEnd = end
		}
		beg, end := match[len(match)-1], match[1]
		if beg < end {
			if end <= maxEnd {
				l.err = fmt.Errorf("overlapping regex groups (LexerDef='%s', Mode='%s', Rule=%d)", l.def.Name, l.mode.Name, l.ruleIdx)
				return true
			}
			l.QueueLexeme(Lexeme{Type: Text, Str: l.str[beg:end]})
		}
		l.str = l.str[end:]
		return true
	}
}

// ScoreAdd add val to the current score
func ScoreAdd(val int) RegexDefRuleFunc {
	return func(l *LexerEngine, match []int) bool {
		l.score += val
		return true
	}
}

// PushMode add current mode to stack and set current mode to name.
func PushMode(name string) RegexDefRuleFunc {
	return func(l *LexerEngine, match []int) bool {
		l.PushMode(name)
		return true
	}
}

// PopMode set current mode to the one removed from the top of the mode stack.
func PopMode() RegexDefRuleFunc {
	return func(l *LexerEngine, match []int) bool {
		l.PopMode()
		return true
	}
}

// Predefined RegexLexerRules
var (
	WhiteSpaceRule        = &RegexDefRule{Re: `[ \t\f\v]+`, Do: PopMatch(TextWhiteSpace)}
	NewLineRule           = &RegexDefRule{Re: `[\n\r]+`, Do: PopMatch(TextNewLine)}
	ASCIIIdentifierRule   = &RegexDefRule{Re: `[a-zA-Z][a-zA-Z0-9]*`, Do: PopMatch(CodeIdentifier)}
	AssignmentRule        = &RegexDefRule{Re: `=|\+=|-=|\*=|/=|!=|%=|\|=|&=`, Do: PopMatch(CodeOperatorAssignment)}
	MathOpRule            = &RegexDefRule{Re: `[\+-\*/%]`, Do: PopMatch(CodeOperatorArithmetic)}
	LogicalOpRule         = &RegexDefRule{Re: `&&|\|\||!`, Do: PopMatch(CodeOperatorLogical)}
	BinaryOpRule          = &RegexDefRule{Re: `[\|&^]`, Do: PopMatch(CodeOperatorBinary)}
	OperatorsRule         = &RegexDefRule{Re: `\+=|-=|\*=|/=|!=|\|\|&&|[<>%!\+\-\*/=]`, Do: PopMatch(CodeOperator)}
	PunctuationRule       = &RegexDefRule{Re: `[\.,;]`, Do: PopMatch(CodePunctuation)}
	DelimiterRule         = &RegexDefRule{Re: `[{}\[\]\(\)]`, Do: PopMatch(CodeDelimiter)}
	SlashSlashCommentRule = &RegexDefRule{Re: `//.*[^\n]`, Do: All(PopMatch(CodeComment), ScoreAdd(1))}
	SlashStarCommentRule  = &RegexDefRule{Re: `/\*(?:.|\r|\n)*?\*/`, Do: All(PopMatch(CodeComment), ScoreAdd(1))}
	HexadecimalNumberRule = &RegexDefRule{Re: `0[xX][A-Faf0-9]+`, Do: PopMatch(CodeNumberHexadecimal)}
	OctalNumberRule       = &RegexDefRule{Re: `0[0-7]+`, Do: PopMatch(CodeNumberOctal)}
	IntegerNumberRule     = &RegexDefRule{Re: `[-+]?(?:0|(?:[1-9][0-9]*))[lL]?`, Do: PopMatch(CodeNumberInteger)}
	DecimalNumberRule     = &RegexDefRule{Re: `[+-]?(?:(?:(?:[1-9][0-9]*[eE][+-]?[1-9][0-9]*))|(?:(?:0|(?:[1-9][0-9]*))\.(?:(?:[1-9]?|[0-9]+))?)|(?:\.(?:[1-9]|[0-9]+))(?:[eE][+-]?[1-9][0-9]*)?)[fFlL]?`,
		Do: PopMatch(CodeNumberDecimal)}
)
