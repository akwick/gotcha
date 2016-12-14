package main

func main() {
	sink(source())
}
func sink(s string) {
}

func source() string {
	return "secret"
}
