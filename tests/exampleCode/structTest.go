package main

// T is a struct with one element
// The elements in a struct are accesed via pointers
type T struct {
	s string
}

func main() {
	t := new(T)
	t.s = source()
	// This statement decomposes into two instructions in SSA:
	// y = &t.0 (FieldAddr)
	// *y = source()
	sink(t.s) // sink, leak
}

func sink(s string) {
}

func source() string {
	return "secret"
}
