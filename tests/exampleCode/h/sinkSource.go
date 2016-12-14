package h

// Sink is a sink
func Sink(s string) {
}

// Source returns a tainted string
func Source() string {
	return "secret"
}
