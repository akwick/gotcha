// Struct 1 element
// Report 2 flows[17,28]
// 17: Element of struct (t.s) as parameter to function.
// 28: Assign element of struct to a new variable and pass the new variable to function.
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

	v := t.s
	sink(v) // sink, leak

	w := u.s
	sink(w) // sink, no leak
}

func sink(s string) {
}

func source() string {
	return "secret"
}
