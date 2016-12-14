package main

func main() {
	var s string
	s = source()
	// should be reported
	sink(s)
	s = "Hello World"
	// shouldn't reported
	sink(s)
}
func sink(s string) {
}

func source() string {
	return "secret"
}
