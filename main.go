package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/foxeng/alanc/parser"
)

func main() {
	fin, err := os.Open(os.Args[1])
	if err != nil {
		return
	}

	l := parser.NewLexer(bufio.NewReader(fin))
	fmt.Println(parser.Parse(&l) == 0)
}
