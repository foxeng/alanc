package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/foxeng/alanc/lexer"
)

func main() {
	fin, err := os.Open(os.Args[1])
	if err != nil {
		return
	}

	l := lexer.New(bufio.NewReader(fin))
	var lval lexer.YySymType
	for t := l.Lex(&lval); t != lexer.EOF; t = l.Lex(&lval) {
		fmt.Printf("%d: %#v\n", t, lval)
	}
}
