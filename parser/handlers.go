package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/foxeng/alanc/semantic"
)

var hexDigits = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	'a', 'b', 'c', 'd', 'e', 'f'}

// eofToUnexpectedEOF returns io.UnexpectedEOF if err is io.EOF. It returns err otherwise.
func eofToUnexpectedEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

// consumeSpace consumes as many whitespace characters from bs.
func consumeSpace(bs io.ByteScanner) error {
	for {
		b, err := bs.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if !unicode.IsSpace(rune(b)) {
			if err = bs.UnreadByte(); err != nil {
				panic("no byte to unread")
			}
			break
		}
	}
	return nil
}

// consumeLineComment consumes the body of a line comment from bs (including the newline in the
// end). It assumes the leading "--" has already been consumed.
func consumeLineComment(bs io.ByteScanner) error {
	for {
		b, err := bs.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if b == '\n' {
			break
		}
	}
	return nil
}

// consumeBlockComment consumes the body of a block comment from bs (up to and including the closing
// "*)"). It assumes the leading "(*" has already been consumed.
func consumeBlockComment(bs io.ByteScanner) error {
	for {
		b0, err := bs.ReadByte()
		if err != nil {
			return eofToUnexpectedEOF(err)
		}
		if b0 == '(' || b0 == '*' {
			b1, err := bs.ReadByte()
			if err != nil {
				return eofToUnexpectedEOF(err)
			}
			if b0 == '(' && b1 == '*' {
				if err = consumeBlockComment(bs); err != nil {
					return eofToUnexpectedEOF(err)
				}
			} else if b0 == '*' && b1 == ')' {
				return nil
			}
		}
	}
}

// handleKwdOrIdent returns an keyword or identifier from bs starting with b0. It assumes b0 is a
// letter (the last byte **read** from bs).
func handleKwdOrIdent(b0 byte, bs io.ByteScanner, lval *yySymType) (int, error) {
	// Read as many letters, digits and underscores
	var buf strings.Builder
	buf.WriteByte(b0)
	for {
		b, err := bs.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1, err
		}
		if unicode.IsLetter(rune(b)) || unicode.IsDigit(rune(b)) || b == '_' {
			buf.WriteByte(b)
		} else {
			if err = bs.UnreadByte(); err != nil {
				panic("no byte to unread")
			}
			break
		}
	}

	switch word := buf.String(); word {
	case "byte":
		return BYTE, nil
	case "else":
		return ELSE, nil
	case "false":
		return FALSE, nil
	case "if":
		return IF, nil
	case "int":
		return INT, nil
	case "proc":
		return PROC, nil
	case "reference":
		return REFERENCE, nil
	case "return":
		return RETURN, nil
	case "while":
		return WHILE, nil
	case "true":
		return TRUE, nil
	default:
		lval.id = semantic.ID(word)
		return IDENT, nil
	}
}

// handleIntConst returns the an integer constant from bs starting with b0. It assumes b0 is a
// decimal digit (the last byte **read** from bs).
func handleIntConst(b0 byte, bs io.ByteScanner, lval *yySymType) (int, error) {
	// Read as many decimal digits
	var buf strings.Builder
	buf.WriteByte(b0)
	for {
		b, err := bs.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1, err
		}
		if unicode.IsDigit(rune(b)) {
			buf.WriteByte(b)
		} else {
			if err = bs.UnreadByte(); err != nil {
				panic("no byte to unread")
			}
			break
		}
	}

	i, err := strconv.Atoi(buf.String())
	if err != nil {
		if err == strconv.ErrSyntax {
			// This shouldn't happen because the digits are checked above
			panic(err)
		}
		return -1, err
	}
	lval.iconst = semantic.IntConstExpr{
		Val: i,
	}
	return INT_CONST, nil
}

// nextChar returns the next character from bs, interpreting escape sequences. It returns quotation
// marks ('\'' and '"') like all normal characters, leaving their interpretation to the caller. If
// it reaches EOF, it reports it to the caller.
func nextChar(bs io.ByteScanner) (byte, error) {
	b0, err := bs.ReadByte()
	if err != nil {
		return 0, err
	}
	if !unicode.IsPrint(rune(b0)) {
		return 0, fmt.Errorf("non-printable character %q", b0)
	}
	if b0 == '\\' {
		b1, err := bs.ReadByte()
		if err != nil {
			return 0, err
		}
		switch b1 {
		case 'n':
			return '\n', nil
		case 't':
			return '\t', nil
		case 'r':
			return '\r', nil
		case '0':
			return '\x00', nil
		case '\\':
			return '\\', nil
		case '\'':
			return '\'', nil
		case '"':
			return '"', nil
		case 'x':
			var d [2]byte
			for i := range d {
				d[i], err = bs.ReadByte()
				if err != nil {
					return 0, err
				}
				if !bytes.ContainsRune(hexDigits, rune(d[i])) {
					return 0, fmt.Errorf("not a hex digit in hex escape sequence: %q", d[i])
				}
			}
			return 16*d[0] + d[1], nil
		default:
			return 0, fmt.Errorf("invalid first character of escape sequence: %q", b1)
		}
	} else {
		return b0, nil
	}
}

// handleCharLit returns a character literal from bs. It assumes the first argument is the starting
// '\'' (the last byte **read** from bs).
func handleCharLit(_ byte, bs io.ByteScanner, lval *yySymType) (int, error) {
	c, err := nextChar(bs)
	if err != nil {
		return -1, eofToUnexpectedEOF(err)
	}
	if c == '\'' {
		return -1, fmt.Errorf("empty character literal")
	}

	// Check for the closing '\''
	b, err := bs.ReadByte()
	if err != nil {
		return -1, eofToUnexpectedEOF(err)
	}
	if b != '\'' {
		return -1, fmt.Errorf("too many characters in character literal")
	}

	lval.cconst = semantic.CharConstExpr{
		Val: rune(c),
	}
	return CHAR_LIT, nil
}

// handleStrLit returns a string literal from bs. It assumes the first argument is '"' (the last
// byte **read** from bs).
func handleStrLit(_ byte, bs io.ByteScanner, lval *yySymType) (int, error) {
	var buf strings.Builder
	for {
		c, err := nextChar(bs)
		if err != nil {
			return -1, eofToUnexpectedEOF(err)
		}
		if c == '"' {
			break
		}
		buf.WriteByte(c)
	}
	lval.strlit = semantic.StrLitExpr{
		Val: buf.String(),
	}
	return STR_LIT, nil
}

// handleOp returns an operator from bs. It assumes b0 is in operators (the last byte **read** from
// bs).
func handleOp(b0 byte, bs io.ByteScanner, lval *yySymType) (int, error) {
	if b0 != '!' && b0 != '=' && b0 != '<' && b0 != '>' {
		return int(b0), nil
	}
	b1, err := bs.ReadByte()
	if err != nil {
		if err == io.EOF {
			return int(b0), nil
		}
		return -1, err
	}
	if b1 == '=' {
		switch b0 {
		case '=':
			return EQ, nil
		case '!':
			return NE, nil
		case '<':
			return LE, nil
		case '>':
			return GE, nil
		default:
			// This shouldn't happen, it's checked above
			panic(fmt.Sprintf("invalid operator %q", fmt.Sprintf("%c%c", b0, b1)))
		}
	} else {
		if err = bs.UnreadByte(); err != nil {
			panic("no byte to unread")
		}
		return int(b0), err
	}
}

// handleSep returns a separator from bs. It assumes b0 is in separators (the last byte **read**
// from bs).
func handleSep(b0 byte, _ io.ByteScanner, _ *yySymType) (int, error) {
	return int(b0), nil
}
