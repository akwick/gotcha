package main

import "fmt"

func main() {
	a := "Hello Gophers"
	f(a)
	a = source()
	f(a)
}

func f(s string) {
	sink(&s)
}

func sink(s *string) {
	fmt.Printf("A gopher reaches a sink: %s \n", *s)
}

// copied from 10_FuncAndVarDecl
// no guaranty for equality
func source() string {
	return "I am an evil gopher"
}
