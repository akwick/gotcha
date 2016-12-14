package tests

import "testing"

func TestSimpleLinkList(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/simpleLinkedList.go"}, 1}
	t1 := &testDataStruct{[]string{"./exampleCode/simpleLinkedListTestAllElements.go"}, 2}
	//t2 := &testDataStruct{[]string{"./exampleCode/simpleLinkedList0.go"}, 1}
	td := []*testDataStruct{t0, t1}
	check(td, t)
}
