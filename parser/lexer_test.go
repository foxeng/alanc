package parser

import (
	"fmt"
	"strings"
	"testing"
)

var tokenNames = map[int]string{
	EOF:       "EOF",
	BYTE:      "byte",
	ELSE:      "else",
	FALSE:     "false",
	IF:        "if",
	INT:       "int",
	PROC:      "proc",
	REFERENCE: "reference",
	RETURN:    "return",
	WHILE:     "while",
	TRUE:      "true",
	IDENT:     "<id>",
	INT_CONST: "<int>",
	CHAR_LIT:  "<char>",
	STR_LIT:   "<string>",
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
	}`: {IDENT, int('('), int(')'), int(':'), PROC,
		int('{'),
		IDENT, int('('), STR_LIT, int(')'), int(';'),
		int('}')},
	`writeString(".\n");`: {IDENT, int('('), STR_LIT, int(')'), int(';')},
	` { -- hanoi`:         {int('{')},
	`if (rings >= 1) {`:   {IF, int('('), IDENT, GE, INT_CONST, int(')'), int('{')},
	`	else if ( n % 2 == 0) return 0;`: {ELSE, IF, int('('), IDENT, int('%'), INT_CONST, EQ,
		INT_CONST, int(')'), RETURN, INT_CONST, int(';')},
	`while (i <= n / 2) {`: {WHILE, int('('), IDENT, LE, IDENT, int('/'), INT_CONST, int(')'), int('{')},
	`i = i + 2;`:           {IDENT, int('='), IDENT, int('+'), INT_CONST, int(';')},
	`limit: int;`:          {IDENT, int(':'), INT, int(';')},
}

func TestLexer(t *testing.T) {
	var lval yySymType
	for test, wants := range tests {
		wants = append(wants, EOF)
		l := NewLexer(strings.NewReader(test))
		for i, want := range wants {
			if got := l.Lex(&lval); got != want {
				t.Errorf("Lex(%q) [%d] = %q, want %q", test, i, tokToName(got), tokToName(want))
			}
		}
	}
}
