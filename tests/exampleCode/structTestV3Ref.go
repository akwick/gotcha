// Struct with 1 element
// Report 2 Flows[17;32]
// 17: Call sink with element of struct
// 32: Call sink with T passed as pointer to a function
package main

// T is a simple test struct
type T struct {
	s string
}

func main() {
	t := new(T)
	t.s = source()
	sink(t.s) // sink, leak

	u := new(T)
	u.s = "Hello World"
	sink(u.s) // sink, no leak

	a(u) // u is untainted
	a(t) // t is tainted
}

// Test a pointer as parameter
func a(t *T) {
	sink(t.s) // sink, leak if t is tainted
}

func sink(s string) {
}

func source() string {
	return "secret"
}
