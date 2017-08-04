package main

func main() {
	x := "hello, world"
	ch := make(chan string)
	go f(ch)
	// @expectedflow: false
	sink(&x) // no leak (flow-insensitive analysis for pointers)
	x = source()
	ch <- x
}
func f(ch_1 chan string) {
	y := <-ch_1
	// @expectedflow: true
	sink(&y) // leak
}
func sink(s *string) {}
func source() string {
	return "secret"
}
