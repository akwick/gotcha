// +build !junittest

package worklist

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelloWorldConvey(t *testing.T) {
	source := []string{"../tests/exampleCode/hello.go"}
	wlInit(path, source, taintfile, false, "", true)

	lattices := make(map[string]string)
	lattices["t0"] = "t0 : Untainted"
	lattices["world"] = "\" World\":string : Untainted"
	lattices["hello"] = "\"Hello\":string : Untainted"
	should := []td{
		{name: "\"Hello\":string + \" World\":string",
			updatedin: []string{},
			flowout:   []string{lattices["t0"], lattices["world"], lattices["hello"]},
			wlLen:     1,
			wlLenRet:  1},
		{name: "return",
			updatedin: []string{lattices["t0"], lattices["world"], lattices["hello"]},
			flowout:   []string{lattices["t0"], lattices["world"], lattices["hello"]},
			wlLen:     0,
			wlLenRet:  0},
	}

	iterateWl(should, t)
}

func TestSimpleLinkedListConvey(t *testing.T) {
	source := []string{"../tests/exampleCode/simpleLinkedListOnlyTaint.go"}
	wlInit(path, source, taintfile, false, "", true)

	lattices := make(map[string]string)
	lattices["t0"] = "t0 : Tainted"
	lattices["t1"] = "t1 : Tainted"
	lattices["t1.0"] = "t1 : Uninitialized"
	lattices["t2"] = "t2 : Tainted"
	lattices["t2.0"] = "t2 : Uninitialized"
	lattices["t3"] = "t3 : Tainted"
	lattices["t3.0"] = "t3 : Uninitialized"
	lattices["0"] = "0:int : Untainted"
	len := 5
	len2 := 28

	should := []td{
		// Context 0
		/*
			func main():
			0:                                                                entry P:0 S:0
				t0 = source()                                                    string
				t1 = NewList(t0)                                            *LinkedList
				t2 = (*LinkedList).GetData(t1, 0:int)                            string
				t3 = sink(t2)                                                        ()
				return
		*/
		{name: "source()",
			updatedin: []string{},
			flowout:   []string{lattices["t0"]},
			wlLen:     len, // no node is added because source() is a source
			wlLenRet:  len},
		{name: "NewList(t0)",
			updatedin: []string{lattices["t0"]},
			flowout:   []string{lattices["t0"], lattices["t1.0"]},
			wlLen:     len - 1 + 9,  // NewList has 9 nodes
			wlLenRet:  len - 1 + 9}, // 14
		{name: "(*LinkedList).GetData(t1, 0:int)",
			updatedin: []string{lattices["t0"], lattices["t1.0"]},
			flowout:   []string{lattices["t0"], lattices["t1.0"], lattices["0"], lattices["t2.0"]},
			wlLen:     len - 1 + 9 - 1 + 19, // GetData has 19 nodes
			wlLenRet:  len - 1 + 9 - 1 + 19},
		{name: "sink(t2)",
			updatedin: []string{lattices["t0"], lattices["t1.0"], lattices["0"], lattices["t2.0"]},
			flowout:   []string{lattices["t0"], lattices["t1.0"], lattices["0"], lattices["t2.0"], lattices["t3.0"]},
			wlLen:     len - 1 + 9 - 1 + 19 - 1, // no additional node because sink() is a sink
			wlLenRet:  len - 1 + 9 - 1 + 19 - 1},
		{name: "return",
			updatedin: []string{lattices["t0"], lattices["t1.0"], lattices["0"], lattices["t2.0"], lattices["t3.0"]},
			flowout:   []string{lattices["t0"], lattices["t1.0"], lattices["0"], lattices["t2.0"], lattices["t3.0"]},
			wlLen:     len - 1 + 9 - 1 + 19 - 1 - 1,  // no call
			wlLenRet:  len - 1 + 9 - 1 + 19 - 1 - 1}, // no context transfer to context 0 -> add no nodes
		// Context 1 - start with node 5
		/*
					func NewList(s string) *LinkedList:
			0:                                                                entry P:0 S:0
				t0 = new llNode (complit)                                       *llNode
				t1 = &t0.data [#0]                                              *string
				t2 = &t0.next [#1]                                             **llNode
				*t1 = s
				*t2 = nil:*llNode
				t3 = new LinkedList (complit)                               *LinkedList
				t4 = &t3.head [#0]                                             **llNode
				*t4 = t0
				return t3
		*/
		{name: "NewList(t0)",
			updatedin: []string{lattices["t0"]},
			flowout:   []string{lattices["t0"]},
			wlLen:     len2,  // no call
			wlLenRet:  len2}, // Context 0 node 0 -> Context 1
	}
	// TODO: Continue this test
	iterateWl(should, t)
}

func iterateWl(should []td, t *testing.T) {
	Convey("Run worklist as long as elements are in the worklist", t, FailureContinues, func() {
		for i := 0; i < len(should); i++ {
			Convey("Remove element "+strconv.Itoa(i)+" from the worklist", func() {
				n := worklist.RemoveFirst()
				So(n.String(), ShouldContainSubstring, should[i].name)
				Convey("Update the entry context of the element "+strconv.Itoa(i), func() {
					updateEntryContext(n)
					for _, s := range should[i].updatedin {
						So(n.GetIn().String(), ShouldContainSubstring, s)
					}

					Convey("Flow the node of element "+strconv.Itoa(i), func() {
						flow(n)
						for _, s := range should[i].flowout {
							So(n.GetOut().String(), ShouldContainSubstring, s)
						}
						Convey("Check whether the flow has updated the number of elements in the worklist"+strconv.Itoa(i)+". If yes: Handle the change", func() {
							n.CheckAndHandleChange()
							So(worklist.Len(), ShouldEqual, should[i].wlLen)
							Convey("Check whether the node ( "+strconv.Itoa(i)+")is a return statement. If yes: Add all contexts which have a transition to the node to the worklist", func() {
								checkAndHandleReturn(n)
								So(worklist.Len(), ShouldEqual, should[i].wlLenRet)
							})
						})
					})
				})

			})
		}
	})
}

type td struct {
	name      string
	updatedin []string
	flowout   []string
	wlLen     int
	wlLenRet  int
}

var (
	taintfile = "../sourcesAndSinks.txt"
	path      = "github.com/akwick/gotcha"
)
