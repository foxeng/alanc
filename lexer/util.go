package lexer

import (
	"fmt"
	"io"
)

// posByteScanner adds position handling to io.ByteScanner.
type posByteScanner struct {
	bs          io.ByteScanner
	line, col   int
	prevCol     int  // the last column of the previous line
	lastNewLine bool // whether the last byte read was a newline
}

func newPosByteScanner(bs io.ByteScanner) *posByteScanner {
	return &posByteScanner{
		bs:          bs,
		line:        1,
		col:         1,
		lastNewLine: false,
	}
}

// ReadByte adds position handling to io.ByteScanner.ReadByte.
func (pbs *posByteScanner) ReadByte() (byte, error) {
	b, err := pbs.bs.ReadByte()
	if err == nil {
		if b == '\n' {
			pbs.line++
			pbs.prevCol = pbs.col
			pbs.col = 1
			pbs.lastNewLine = true
		} else {
			pbs.col++
			pbs.lastNewLine = false
		}
	}
	return b, err
}

// UnreadByte adds position handling to io.ByteScanner.UnreadByte.
func (pbs *posByteScanner) UnreadByte() error {
	if pbs.lastNewLine {
		pbs.line--
		pbs.col = pbs.prevCol
	} else {
		pbs.col--
	}
	return pbs.bs.UnreadByte()
}

func (pbs posByteScanner) String() string {
	return fmt.Sprintf("line %d, column %d", pbs.line, pbs.col)
}
