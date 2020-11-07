.PHONY: all parser alanc clean distclean

all: alanc

parser: parser/parser.y
	go generate github.com/foxeng/alanc/parser

parser/parser.go: parser/parser.y
	go generate github.com/foxeng/alanc/parser

alanc: parser/parser.go
	go build

clean:
	rm -f parser/parser.go

distclean: clean
	rm -f alanc
