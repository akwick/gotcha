/*
One test case can handle only a certain amount of classes ~3 simple because of the time constraints.
*/
package tests

import "testing"

// TODO teardown which deletes files if no error occurs
func TestReportedFlows_sinkSource(t *testing.T) {
	t0 := &testDataStruct{[]string{"./exampleCode/hello.go"}, 0}
	t1 := &testDataStruct{[]string{"./exampleCode/sinkSourceSimple.go"}, 1}
	t2 := &testDataStruct{[]string{"./exampleCode/sinkSourceV2.go"}, 1}
	t3 := &testDataStruct{[]string{"./exampleCode/sinkSourceV3.go"}, 2}
	t4 := &testDataStruct{[]string{"./exampleCode/sinkSourceV2NewVariable.go"}, 1}
	t5 := &testDataStruct{[]string{"./exampleCode/sinkSourceV2ChangeOrder.go"}, 1}
	testData := []*testDataStruct{t0, t1, t2, t3, t4, t5}
	t.Log(testData[0])
	check(testData, t)
}
