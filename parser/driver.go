//go:generate goyacc -o parser.go -v "" parser.y

package parser

// Parse is a wrapper around goyacc's yyParse.
func Parse(l *Lexer) int {
	return yyParse(l)
}
