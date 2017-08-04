package main

import "github.com/akwick/gotcha/tests/exampleCode/h"

func main() {
	l := h.NewList("Hello World")
	tainted := h.Source()
	l.AddData(tainted)
	l.AddData("Gophers are welcome")
	l.AddData(tainted)

	s := l.GetData(0)
	// @expectedflow: true
	h.Sink(s)
	s := l.GetData(2)
	// @expectedflow: true
	h.Sink(s)
	s := l.GetData(1)
	// @expectedflow: false
	h.Sink(s)
	s := l.GetData(3)
	// @expectedflow: false
	h.Sink(s)

}
