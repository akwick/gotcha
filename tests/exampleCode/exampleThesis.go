package main

import "fmt"

func main() {
	ch := make(chan string)
	go f(ch)
	x := source()
	ch <- x
}

func f(c chan string) {
	y := <-c
	sink(y)
	a := "Hello Gopher"
	b := g(a)
	fmt.Printf(" %s\n", b)
}

func g(s string) string {
	t := s + " 1"
	return t
}

func sink(s string) {
	fmt.Printf("A gopher reaches a sink: %s \n", s)
}

func source() string {
	return "I am an evil gopher"
}
