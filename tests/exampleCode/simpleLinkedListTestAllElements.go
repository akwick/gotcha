package main

import "github.com/akwick/gotcha/tests/exampleCode/h"

func main() {
	l := h.NewList("Hello World")
	tainted := h.Source()
	l.AddData(tainted)
	l.AddData("Gophers are welcome")
	l.AddData(tainted)

	for i := 0; i < 4; i++ {
		s := l.GetData(i)
		h.Sink(s) // sink, no leak if i%2==0, leak if i%2 == 1
	}
}
