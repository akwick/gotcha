package main

func main() {
	var s string
	s = source()
	// @expectedflow: true
	sink(s)
	s = "Hello World"
	// @expectedflow: false
	sink(s)
}
func sink(s string) {
}

func source() string {
	return "secret"
}
