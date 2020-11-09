//go:generate goyacc -o parser.go -v "" parser.y

package parser

import (
	"errors"
	"fmt"

	"github.com/foxeng/alanc/semantic"
)

// Parse is a wrapper around goyacc's yyParse. When yyParse accepts, this returns the AST produced.
func Parse(l *Lexer) (*semantic.Ast, error) {
	r := yyParse(l)
	if lexErr != nil {
		return nil, fmt.Errorf("lexer: %v", lexErr)
	}
	if r != 0 {
		return nil, errors.New("parser rejected")
	}
	return ast, nil
}
