# alanc

This is the compiler for the [Alan](http://courses.softlab.ntua.gr/compilers/2018a/alan2018.pdf)
programming language, as part of the [compilers](http://courses.softlab.ntua.gr/compilers/2018a/)
course at ECE NTUA (spring 2018).

## Alan

For the specification of the language (in Greek), see `alan2018.pdf`.

Various example programs in Alan can be found in the `examples` directory.

## Build guide

To build the compiler you need the Go [toolchain](https://golang.org/dl/).

There is a single dependency, on [goyacc](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc). To
install the latest version:

```
go get -u golang.org/x/tools/cmd/goyacc
```

Finally, to build alanc:

```
make
```
