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

// lexErr communicates lexer errors to the parser driver.
// TODO: Figure out a way to do this via the parser, not bypassing it.
var lexErr error

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

// Lex returns the next token identifier and places the relevant token information on lval.
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
			lexErr = fmt.Errorf("starting new token: %v (%v)", err, l.pbs)
			return -1
		}
		if unicode.IsSpace(rune(b0)) {
			if err = consumeSpace(l.pbs); err != nil {
				if err == io.EOF {
					return EOF
				}
				lexErr = fmt.Errorf("consuming white space: %v (%v)", err, l.pbs)
				return -1
			}
		} else if b0 == '-' || b0 == '(' {
			b1, err := l.pbs.ReadByte()
			if err != nil {
				if err == io.EOF {
					break // return b0 and let the parser handle the EOF in the next call
				}
				lexErr = fmt.Errorf("checking for comment: %v (%v)", err, l.pbs)
				return -1
			}
			if b0 == '-' && b1 == '-' {
				if err = consumeLineComment(l.pbs); err != nil {
					lexErr = fmt.Errorf("consuming line comment: %v (%v)", err, l.pbs)
					return -1
				}
			} else if b0 == '(' && b1 == '*' {
				if err = consumeBlockComment(l.pbs); err != nil {
					lexErr = fmt.Errorf("consuming block comment: %v (%v)", err, l.pbs)
					return -1
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
		lexErr = fmt.Errorf("unexpected character: %c (code point %d) (%v)", b0, b0, l.pbs)
		return -1
	}

	tok, err := handler(b0, l.pbs, lval)
	if err != nil {
		// TODO OPT: Report what token was being scanned (thus, specialize for each switch case
		// above)
		lexErr = fmt.Errorf("%v (%v)", err, l.pbs)
		return -1
	}
	return tok
}

// Error reports a parser error, e.
func (l Lexer) Error(e string) {
	fmt.Fprintf(os.Stderr, "parse: %s (around %v)\n", e, l.pbs)
}
