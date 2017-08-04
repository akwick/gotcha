package main

func main() {
	var s string
	s = source()
	// @expectedflow: true
	sink(s)
	t := "Hello World"
	// @expectedflow: false
	sink(t)
}
func sink(s string) {
}

func source() string {
	return "secret"
}
