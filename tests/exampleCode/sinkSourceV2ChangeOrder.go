// Simple test to test flow-sensitive
package main

func main() {
	var s string
	s = "Hello World"
	// @expectedflow: false
	sink(s) // sink, no leak
	s = source()
	// @expectedflow: true
	sink(s) // sink, leak
}
func sink(s string) {
}

func source() string {
	return "secret"
}
