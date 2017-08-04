package main

func main() {
	x := "hello, world"
	ch := make(chan *string)
	go f(ch)
	// @expectedflow: false
	sink(&x) // sink, no leak
	x = source()
	// @expectedflow: true
	ch <- &x // sink, leak
}
func f(ch_1 chan *string) {
	y := <-ch_1
	sink(y)
}
func sink(s *string) {}
func source() string {
	return "secret"
}
