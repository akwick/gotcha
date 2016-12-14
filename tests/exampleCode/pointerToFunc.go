package main

func main() {
	sink(source)
}

func sink(f func() string) {
	f()
}

func source() string {
	return "secret"
}
