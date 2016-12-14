package main

func main() {
	x := "Hello World"
	ch := make(chan string)
	go f(ch)
	sink(x)
	x = source()
	ch <- x
	go add(x, ch)
	g(ch)

	// Channel defined in an other method than the main method
	c := newChannel()
	go f(c)
}

func f(ch chan string) {
	y := <-ch
	// A flow should be reported
	sink(y)
}
func g(ch chan string) {
	go func() {
		y := <-ch
		// A flow should be reported
		sink(y)
	}()
}

func add(x string, ch chan string) {
	ch <- x
}

func newChannel() chan string {
	c := make(chan string)
	return c
}

func sink(s string) {
}

func source() string {
	return "secret"
}
