package tests

import "testing"

func TestSinkTS(t *testing.T) {
	t1 := &testDataStruct{[]string{"./exampleCode/structTest.go"}, 1}
	td := []*testDataStruct{t1}
	check(td, t)
}

func TestFieldSensitiv(t *testing.T) {
	t1 := &testDataStruct{[]string{"./exampleCode/structTestV2.go"}, 0}
	td := []*testDataStruct{t1}
	check(td, t)
}

func TestAssignTSToNewVariable(t *testing.T) {
	t1 := &testDataStruct{[]string{"./exampleCode/structTestV4.go"}, 2}
	td := []*testDataStruct{t1}
	check(td, t)
}

func TestPointerAsParameterSignature(t *testing.T) {
	t1 := &testDataStruct{[]string{"./exampleCode/structTestV3Ref.go"}, 2}
	t2 := &testDataStruct{[]string{"./exampleCode/structTestV3RefSimple.go"}, 1}
	td := []*testDataStruct{t1, t2}
	check(td, t)
}

func TestPointerAsParameterCall(t *testing.T) {
	t1 := &testDataStruct{[]string{"./exampleCode/structTestV3Val.go"}, 1}
	td := []*testDataStruct{t1}
	check(td, t)
}

func TestFunctionAsParameter(t *testing.T) {
	// TODO implement feature such that the tool can handle functions as a parameter tool
	//	t1 := &testDataStruct{[]string{"./exampleCode/pointerToFunc.go"}, 1}
	t1 := &testDataStruct{[]string{"./exampleCode/pointerToFunc.go"}, 0}
	testData := []*testDataStruct{t1}
	check(testData, t)
}

func TestAliasingThroughFields(t *testing.T) {
	t1 := &testDataStruct{[]string{"./exampleCode/aliasingThroughFields.go"}, 2}
	t2 := &testDataStruct{[]string{"./exampleCode/aliasingThroughFields2.go"}, 2}
	testData := []*testDataStruct{t1, t2}
	check(testData, t)
}
