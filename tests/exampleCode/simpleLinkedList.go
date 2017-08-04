package main

import (
	"github.com/akwick/gotcha/tests/exampleCode/h"
)

func main() {
	l := h.NewList("Hello World")
	tainted := h.Source()
	l.AddData(tainted)
	l.AddData("Gophers are welcome")
	l.AddData(tainted)

	s0 := l.GetData(0) // untainted
	s1 := l.GetData(1) // tainted

	// @expectedflow: false
	h.Sink(s0) // sink, no leak
	// @expectedflow: true
	h.Sink(s1) // sink, leak
}
