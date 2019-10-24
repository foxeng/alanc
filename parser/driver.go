//go:generate goyacc -o parser.go -v "" parser.y

package parser

import (
	"errors"

	"github.com/foxeng/alanc/ast"
)

// Parse is a wrapper around goyacc's yyParse. When yyParse accepts, this returns the AST produced.
func Parse(l *Lexer) (*ast.Ast, error) {
	r := yyParse(l)
	if r != 0 {
		return nil, errors.New("parser rejected")
	}
	return _ast, nil
}
