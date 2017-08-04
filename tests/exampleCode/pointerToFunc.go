package main

func main() {
	// The tool currently does not handle function arguments
	// @expectedflow: false
	sink(source)
}

func sink(f func() string) {
	f()
}

func source() string {
	return "secret"
}
