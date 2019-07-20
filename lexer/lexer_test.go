package lexer

import (
	"fmt"
	"strings"
	"testing"
)

var tokenNames = map[int]string{
	EOF:       "EOF",
	Byte:      "byte",
	Else:      "else",
	False:     "false",
	If:        "if",
	Int:       "int",
	Proc:      "proc",
	Reference: "reference",
	Return:    "return",
	While:     "while",
	True:      "true",
	Ident:     "<id>",
	IntConst:  "<int>",
	CharLit:   "<char>",
	StrLit:    "<string>",
	EQ:        "==",
	NE:        "!=",
	LE:        "<=",
	GE:        ">=",
	int('='):  "=",
	int('+'):  "+",
	int('-'):  "-",
	int('*'):  "*",
	int('/'):  "/",
	int('%'):  "%",
	int('!'):  "!",
	int('&'):  "&",
	int('|'):  "|",
	int('<'):  "<",
	int('>'):  ">",
	int('('):  "(",
	int(')'):  ")",
	int('['):  "[",
	int(']'):  "]",
	int('{'):  "{",
	int('}'):  "}",
	int(','):  ",",
	int(':'):  ":",
	int(';'):  ";",
}

func tokToName(token int) string {
	name, ok := tokenNames[token]
	if !ok {
		return fmt.Sprintf("Unknown (%d)", token)
	}
	return name
}

var tests = map[string][]int{
	"":              {},
	" ":             {},
	"\n":            {},
	"\t":            {},
	`--`:            {},
	`--comment`:     {},
	`(* comment *)`: {},
	`(* nesting (*
		comments *) *)`: {},
	`hello() : proc
	{
		writeString("Hello world!\n");
	}`: {Ident, int('('), int(')'), int(':'), Proc,
		int('{'),
		Ident, int('('), StrLit, int(')'), int(';'),
		int('}')},
	`writeString(".\n");`: {Ident, int('('), StrLit, int(')'), int(';')},
	` { -- hanoi`:         {int('{')},
	`if (rings >= 1) {`:   {If, int('('), Ident, GE, IntConst, int(')'), int('{')},
	`	else if ( n % 2 == 0) return 0;`: {Else, If, int('('), Ident, int('%'), IntConst, EQ,
		IntConst, int(')'), Return, IntConst, int(';')},
	`while (i <= n / 2) {`: {While, int('('), Ident, LE, Ident, int('/'), IntConst, int(')'), int('{')},
	`i = i + 2;`:           {Ident, int('='), Ident, int('+'), IntConst, int(';')},
	`limit: int;`:          {Ident, int(':'), Int, int(';')},
}

func TestLexer(t *testing.T) {
	var lval YySymType
	for test, wants := range tests {
		wants = append(wants, EOF)
		l := New(strings.NewReader(test))
		for i, want := range wants {
			if got := l.Lex(&lval); got != want {
				t.Errorf("Lex(%q) [%d] = %q, want %q", test, i, tokToName(got), tokToName(want))
			}
		}
	}
}
