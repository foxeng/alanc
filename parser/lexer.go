// Package parser defines the lexer and the parser for Alan.
package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

// EOF is the token number for the endmarker
const EOF = 0

var operators = []byte{'=', '+', '-', '*', '/', '%', '!', '&', '|', '<', '>'}
var separators = []byte{'(', ')', '[', ']', '{', '}', ',', ':', ';'}

// Lexer is the lexer for Alan.
type Lexer struct {
	pbs *posByteScanner
}

// NewLexer returns a new Lexer.
func NewLexer(bs io.ByteScanner) Lexer {
	return Lexer{
		pbs: newPosByteScanner(bs),
	}
}

func (l Lexer) printError(err error) (int, error) {
	return fmt.Fprintf(os.Stderr, "lexer: %v (%v)\n", err, l.pbs)
}

// Lex returns the next token identfier and places the relevant token information on lval.
func (l Lexer) Lex(lval *yySymType) int {
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
	var handler func(byte, io.ByteScanner, *yySymType) (int, error)
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
	fmt.Fprintf(os.Stderr, "parse: %s (around %v)\n", e, l.pbs)
}
