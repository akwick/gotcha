package main

func main() {
	// @expectedflow: true
	sink(source())
}
func sink(s string) {
}

func source() string {
	return "secret"
}
