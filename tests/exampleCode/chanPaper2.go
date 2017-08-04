// f is not called via a go function, instead the go function is inside the body of f.
package main

func main() {
	x := "Hello World"
	ch := make(chan string)
	f(ch)
	// @expectedflow: false
	sink(x)
	x = source()
	ch <- x
}

func f(ch chan string) {
	// *ssa.MakeClosure
	go func() {
		y := <-ch
		// @expectedflow: true
		sink(y)
	}()
}

func sink(s string) {
}

func source() string {
	return "secret"
}
