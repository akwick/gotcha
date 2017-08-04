// Struct 1 element
// 1 Flow [23]
// 23: pointer to struct as parameter
package main

// T is a simple test struct
type T struct {
	s string
}

func main() {
	t := new(T)
	t.s = source()

	u := new(T)
	u.s = "Hello World"
	// @expectedflow: false
	sink(u.s) // sink, no leak

	a(*u) // u is untainted
	a(*t) // t is tainted
}

func a(t T) {
	sink(t.s) // sink, leak if t is Tainted
}

func sink(s string) {
}

func source() string {
	return "secret"
}
