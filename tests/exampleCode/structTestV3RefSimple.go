// Simplified Version of structTestV3Ref
package main

// T is a simple test struct
type T struct {
	s string
}

func main() {
	t := new(T)
	t.s = source()
	a(t) // t is tainted
}

func a(t *T) {
	sink(t.s) // sink, leak if t is tainted
}

func sink(s string) {
}

func source() string {
	return "secret"
}
