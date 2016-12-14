package tests

import "testing"

func TestTaintForGoInMain(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/chanPaper1.go"}, 1}
	td := []*testDataStruct{t0}
	check(td, t)
}

func TestTaintAfterGoInMain(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/chanPaper3.go"}, 1}
	t1 := &testDataStruct{[]string{"./exampleCode/chanPaper0.go"}, 1}
	td := []*testDataStruct{t0, t1}
	check(td, t)
}

func TestTaintAfterGoInFAsClosure(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/chanPaper2.go"}, 1}
	td := []*testDataStruct{t0}
	check(td, t)
}

func TestChangeChannelValueInClosure(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/chanClosureChangeValue.go"}, 1}
	t1 := &testDataStruct{[]string{"./exampleCode/chanClosureChangeValueBeforeClosure.go"}, 1}
	td := []*testDataStruct{t0, t1}
	check(td, t)
}

func TestChannelWithPointer(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/pointerChan.go"}, 1}
	td := []*testDataStruct{t0}
	check(td, t)
}
