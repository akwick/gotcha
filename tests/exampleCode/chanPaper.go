package main

import "fmt"

func main() {
	x := "Hello World"
	ch := make(chan string)
	go f(ch)
	// @expectedflow: false
	sink(x) // sink, no leak
	x = secret()
	ch <- x
}

func f(ch chan string) {
	y := <-ch
	// @expectedflow: true
	sink(y) // sink, leak
}

func sink(s string) {
	fmt.Println(s)
}

func secret() string {
	return "secret"
}
