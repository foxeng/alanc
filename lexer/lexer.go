// Package lexer defines the lexer for Alan.
package lexer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

const (
	// EOF is the end of file marker
	EOF = iota
	// Byte is the keyword "byte"
	Byte
	// Else is the keyword "else"
	Else
	// False is the keyword "false"
	False
	// If is the keyword "if"
	If
	// Int is the keyword "int"
	Int
	// Proc is the keyword "proc"
	Proc
	// Reference is the keyword "reference"
	Reference
	// Return is the keyword "return"
	Return
	// While is the keyword "while"
	While
	// True is the keyword "true"
	True
	// Ident is an identifier
	Ident
	// IntConst is an integer constant
	IntConst
	// CharLit is a character literal
	CharLit
	// StrLit is a string literal
	StrLit
	// EQ is the separator '=='
	EQ
	// NE is the separator '!='
	NE
	// LE is the separator '<='
	LE
	// GE is the separator '>='
	GE
)

var operators = []byte{'=', '+', '-', '*', '/', '%', '!', '&', '|', '<', '>'}
var separators = []byte{'(', ')', '[', ']', '{', '}', ',', ':', ';'}

// Lexer is the lexer for Alan.
type Lexer struct {
	pbs *posByteScanner
}

// New returns a new Lexer.
func New(bs io.ByteScanner) Lexer {
	return Lexer{
		pbs: newPosByteScanner(bs),
	}
}

type YySymType struct {
	id string
	i  int
	c  byte
	s  string
}

func (l Lexer) printError(err error) (int, error) {
	return fmt.Fprintf(os.Stderr, "lexer: %v (%v)\n", err, l.pbs)
}

// Lex returns the next token identfier and places the relevant token information on lval.
func (l Lexer) Lex(lval *YySymType) int {
	// Consume whitespace and comments
	var b0 byte
	var err error
	for {
		b0, err = l.pbs.ReadByte()
		if err != nil {
			if err == io.EOF {
				return EOF
			}
			l.printError(fmt.Errorf("starting new token: %v", err))
		}
		if unicode.IsSpace(rune(b0)) {
			if err = consumeSpace(l.pbs); err != nil {
				if err == io.EOF {
					return EOF
				}
				l.printError(fmt.Errorf("consuming white space: %v", err))
			}
		} else if b0 == '-' || b0 == '(' {
			b1, err := l.pbs.ReadByte()
			if err != nil {
				if err == io.EOF {
					break // return b0 and let the parser handle the EOF in the next call
				}
				l.printError(fmt.Errorf("checking for comment: %v", err))
			}
			if b0 == '-' && b1 == '-' {
				if err = consumeLineComment(l.pbs); err != nil {
					l.printError(fmt.Errorf("consuming line comment: %v", err))
				}
			} else if b0 == '(' && b1 == '*' {
				if err = consumeBlockComment(l.pbs); err != nil {
					l.printError(fmt.Errorf("consuming block comment: %v", err))
				}
			} else {
				if err = l.pbs.UnreadByte(); err != nil {
					panic("no byte to unread")
				}
				break
			}
		} else {
			break
		}
	}

	// Return token starting with b0
	var handler func(byte, io.ByteScanner, *YySymType) (int, error)
	switch {
	case unicode.IsLetter(rune(b0)):
		handler = handleKwdOrIdent
	case unicode.IsDigit(rune(b0)):
		handler = handleIntConst
	case b0 == '\'':
		handler = handleCharLit
	case b0 == '"':
		handler = handleStrLit
	case bytes.ContainsRune(operators, rune(b0)):
		handler = handleOp
	case bytes.ContainsRune(separators, rune(b0)):
		handler = handleSep
	default:
		l.printError(fmt.Errorf("unexpected character: %c (code point %d)", b0, b0))
	}

	tok, err := handler(b0, l.pbs, lval)
	if err != nil {
		// TODO OPT: Report what token was being scanned (thus, specialize for each switch case
		// above)
		l.printError(err)
	}
	return tok
}

// Error reports a parser error, e.
func (l Lexer) Error(e string) {
	fmt.Fprintf(os.Stderr, "parser: %s (around %v)\n", e, l.pbs)
}
