// Struct: 2 Elements
// Flow: 0
// Field sensitive
package main

// T is a struct with two elements
type T struct {
	s string
	t string
}

func main() {
	t := new(T)
	t.s = source()
	t.t = "Hello World"
	sink(t.t) // sink, no leak
}

func sink(s string) {
}

func source() string {
	return "secret"
}
