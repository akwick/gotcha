package main

import (
	"goretech/analysis/tests/exampleCode/h"
)

func main() {
	l := h.NewList("Hello World")
	tainted := h.Source()
	l.AddData(tainted)
	l.AddData("Gophers are welcome")
	l.AddData(tainted)

	s0 := l.GetData(0) // untainted
	s1 := l.GetData(1) // tainted

	h.Sink(s0) // sink, no leak
	h.Sink(s1) // sink, leak
}
