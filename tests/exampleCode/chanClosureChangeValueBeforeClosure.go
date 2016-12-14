// f is not called via a go function, instead the go function is inside the body of f.
package main

func main() {
	ch := make(chan string)
	f(ch)
}

func f(ch chan string) {
	x := "Hello World"
	x = source()
	ch <- x
	// *ssa.MakeClosure
	go func() {
		y := <-ch
		// A flow should be reported
		sink(y)
	}()
	sink(x) // sink, leak
}

func sink(s string) {
}

func source() string {
	return "secret"
}
