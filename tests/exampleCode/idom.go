// Very simple example program which import nothing so that the test is simpler and faster.
package main

func main() {
	s := "Hello"
	if s == "Hello" {
		s += " World"
	} else {
		s = " not executed"
	}
	s += " "
	if s == "Hello W" {
		s += "orld"
	} else {
		s += "and goodbye"
	}
}
