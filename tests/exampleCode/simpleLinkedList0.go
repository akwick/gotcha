package main

func main() {
	l := NewList("Hello World")
	tainted := source()
	l.AddData(tainted)
	l.AddData("Gophers are welcome")

	s0 := l.GetData(0) // untainted
	s1 := l.GetData(1) // tainted

	// @expectedflow: false
	sink(s0) // sink, no leak
	// @expectedflow: true
	sink(s1) // sink, leak
}

/*
Simple implementation of a linked List
*/

// LinkedList is a simple type which provides some operations upon a linked list
type LinkedList struct {
	head *llNode
}

// NewList creates a new linked list with s as first element.
func NewList(s string) *LinkedList {
	h := &llNode{data: s, next: nil}
	return &LinkedList{head: h}
}

// GetData returns the string in position i or an empty string if i does not exist
func (l *LinkedList) GetData(i int) string {
	j := 0
	node := l.head
	for node.next != nil {
		if j == i {
			return node.data
		}
		i++
		node = node.next
	}
	return ""
}

// AddData adds a new element to the end of l.
func (l *LinkedList) AddData(s string) {
	node := l.head
	for node.next != nil {
		node = node.next
	}
	newnode := &llNode{data: s, next: nil}
	node.setNext(newnode)
}

type llNode struct {
	data string
	next *llNode
}

func (l *llNode) getNext() *llNode {
	return l.next
}

func (l *llNode) setNext(n *llNode) {
	l.next = n
}

func (l *llNode) getData() string {
	return l.data
}

func (l *llNode) add(s string) {
	newl := &llNode{next: nil, data: s}
	l.next = newl
}

// Sink is a sink
func sink(s string) {
}

// Source returns a tainted string
func source() string {
	return "secret"
}
