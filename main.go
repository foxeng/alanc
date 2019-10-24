package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/foxeng/alanc/parser"
)

func main() {
	// TODO: Use proper command line parsing packages.
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source file>\n", os.Args[0])
		os.Exit(1)
	}
	fin, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "open %q: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	defer fin.Close()

	l := parser.NewLexer(bufio.NewReader(fin))
	_, err = parser.Parse(&l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		os.Exit(1)
	}
}
