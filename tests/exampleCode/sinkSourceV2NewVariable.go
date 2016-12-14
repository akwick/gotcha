package main

func main() {
	var s string
	s = source()
	// should be reported
	sink(s)
	t := "Hello World"
	// shouldn't reported
	sink(t)
}
func sink(s string) {
}

func source() string {
	return "secret"
}
