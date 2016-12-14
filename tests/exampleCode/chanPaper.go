package main

import "fmt"

func main() {
	x := "Hello World"
	ch := make(chan string)
	go f(ch)
	sink(x) // sink, no leak
	x = secret()
	ch <- x
}

func f(ch chan string) {
	y := <-ch
	sink(y) // sink, leak
}

func sink(s string) {
	fmt.Println(s)
}

func secret() string {
	return "secret"
}
