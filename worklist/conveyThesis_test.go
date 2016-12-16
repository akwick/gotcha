package worklist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelExample(t *testing.T) {
	source := []string{"../tests/exampleCode/exampleThesis.go"}
	wlInit(path, source, taintfile, false, "", true)

	// Variable with the should values of the lattice
	l := make(map[string]string)
	l["t0"] = "t0 : Untainted"
	l["t0.1"] = "t0 : Both"
	l["t1"] = "t1 : Tainted"
	l["t0_f"] = "t0 : Untainted"
	l["t0_f.1"] = "t0 : Both"
	l["t1_f"] = "t1 : Uninitialized"
	l["t2_f"] = "t2 : Uninitialized"
	l["t3_f"] = "t3 : Untainted"
	l["t3_f.1"] = "t3 : Uninitialized"
	l["t4_f"] = "t4 : Untainted"
	l["t4_f.1"] = "t4 : Uninitialized"
	l["t5_f"] = "t5 : Uninitialized"
	l["t6_f"] = "t6 : Uninitialized"
	l["t7_f"] = "t7 : Uninitialized"
	l["c_f"] = "c : Untainted"
	l["c_f.1"] = "c : Both"
	l["t0_g"] = "s : Untainted"
	len0 := 5          // length after initialization
	len1 := 11         // length after first iteration of first context
	len2 := 5 - 1 + 10 // length after first iteration of second context
	len3 := 12
	len4 := 11
	len5 := 3

	// Should Data
	should := []td{
		/* Context 0 - starts with 5 nodes
		                                                                 entry P:0 S:0
		   	t0 = make chan string 0:int                                 chan string
		   	go f(t0)
		   	t1 = source()                                                    string
		   	send t0 <- t1
		   	return
		*/
		{name: "make chan string 0:int",
			updatedin: []string{},
			flowout:   []string{l["t0"]},
			wlLen:     len0,  // statement is no call}
			wlLenRet:  len0}, // 5
		{name: "go f(t0)",
			updatedin: []string{l["t0"]},
			flowout:   []string{l["t0"]},
			wlLen:     len0 - 1 + 10,  // call to f and f has 10 statements
			wlLenRet:  len0 - 1 + 10}, // 7
		{name: "source()",
			updatedin: []string{l["t0"]},
			flowout:   []string{l["t0"], l["t1"]},
			wlLen:     len0 - 1 + 10 - 1, // call to source, but source is a source -> no new value context
			wlLenRet:  len0 - 1 + 10 - 1},
		{name: "send t0 <- t1",
			updatedin: []string{l["t0"], l["t1"]},
			flowout:   []string{l["t0.1"], l["t1"]},
			wlLen:     len0 - 1 + 10 - 1 - 1 + 1, // send statement add statement 2
			wlLenRet:  len0 - 1 + 10 - 1 - 1 + 1},
		{name: "return",
			updatedin: []string{l["t0.1"], l["t1"]},
			flowout:   []string{l["t0.1"], l["t1"]},
			wlLen:     len0 - 1 + 10 - 1 - 1 + 1 - 1,
			wlLenRet:  len0 - 1 + 10 - 1 - 1 + 1 - 1}, // no transition to this value context (beginning context)
		// analyzing the context for the go node
		{name: "go f(t0)",
			updatedin: []string{l["t0"]},
			flowout:   []string{l["t0"]},
			wlLen:     len0 - 1 + 10 - 1 - 1 + 1 - 1 - 1,
			wlLenRet:  len0 - 1 + 10 - 1 - 1 + 1 - 1 - 1}, // 11

		/*func f(c chan string):
		0:                                                                entry P:0 S:0
			t0 = <-c                                                         string
			t1 = sink(t0)                                                        ()
			t2 = g("Hello Gopher":string)                                    string
			t3 = new [1]interface{} (varargs)                       *[1]interface{}
			t4 = &t3[0:int]                                            *interface{}
			t5 = make interface{} <- string (t2)                        interface{}
			*t4 = t5
			t6 = slice t3[:]                                          []interface{}
			t7 = fmt.Printf(" %s\n":string, t6...)               (n int, err error)
			return
		*/
		// Element 6 of the worklist
		{name: "<-c",
			updatedin: []string{l["c_f"]},
			flowout:   []string{l["c_f"], l["t0_f"]},
			wlLen:     len1,
			wlLenRet:  len1},
		{name: "sink(t0)",
			updatedin: []string{l["c_f"], l["t0_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"]},
			wlLen:     len1 - 1, // call to sink but sink is sink  -> no new value context
			wlLenRet:  len1 - 1},
		{name: "g(\"Hello Gopher\":string)",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"]},
			wlLen:     len1 - 1 - 1 + 2, // call to g and g has 2 elements
			wlLenRet:  len1 - 1 - 1 + 2},
		{name: "new [1]interface{} (varargs)",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 + 2 - 1},
		{name: "&t3[0:int]",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 + 2 - 1 - 1},
		{name: "make interface{} <- string (t2)",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"], l["t5_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 + 2 - 1 - 1 - 1},
		{name: "t5",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"], l["t5_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1},
		{name: "slice t3[:]",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1 - 1},
		{name: ":string, t6...)",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"], l["t7_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1 - 1 - 1, // call but not analyzed as not part of the package
			wlLenRet:  len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1 - 1 - 1},
		{name: "return",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"], l["t7_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"], l["t7_f"]},
			wlLen:     len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1 - 1 - 1 - 1,
			wlLenRet:  len1 - 1 - 1 + 2 - 1 - 1 - 1 - 1 - 1 - 1 - 1 + 1}, // transition to line 7
		// Element 16 of the worklist (go f(t0) with t0=Both - additional flow through channel)
		{name: "go f(t0)",
			updatedin: []string{l["t0_f.1"]},
			flowout:   []string{l["t0_f.1"]},
			wlLen:     len2, // return one element, add 10 for new context
			wlLenRet:  len2},
		{name: "sink(t0)",
			updatedin: []string{l["c_f"], l["t0_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"]},
			wlLen:     len2 - 1,
			wlLenRet:  len2 - 1},
		// Element 18 of the worklist - first of context 2 method g
		/*	func g(s string) string:
			 	0:                                                                entry P:0 S:0
					t0 = s + " 1":string                                             string
					return t0 */
		{name: "s + \" 1\":string",
			updatedin: []string{l["t0_g"]},
			flowout:   []string{l["t0_g"]},
			wlLen:     len3,
			wlLenRet:  len3},
		{name: "return t0",
			updatedin: []string{l["t0_g"]},
			flowout:   []string{l["t0_g"]},
			wlLen:     len3 - 1,
			wlLenRet:  len3 - 1 + 1},
		// Element 20 of the worklist go f(t0) again through channel communication
		{name: "go f(t0)",
			updatedin: []string{l["t0_f"]},
			flowout:   []string{l["t0_f"]},
			wlLen:     len4,
			wlLenRet:  len4},
		// Element 21 of the worklist iterate context 4 (f with both)
		{name: "<-c",
			updatedin: []string{l["c_f.1"]},
			flowout:   []string{l["c_f.1"]},
			wlLen:     len4,
			wlLenRet:  len4},
		{name: "sink(t0)",
			updatedin: []string{l["c_f.1"], l["t0_f.1"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f.1"]},
			wlLen:     len1 - 1, // call to sink but sink is sink  -> no new value context
			wlLenRet:  len1 - 1},
		{name: "g(\"Hello Gopher\":string)",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f.1"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f.1"], l["t2_f"]},
			wlLen:     len1 - 1 - 1, // call to g and g has 2 elements
			wlLenRet:  len1 - 1 - 1},
		{name: "new [1]interface{} (varargs)",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f"]},
			wlLen:     len1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 - 1},
		{name: "&t3[0:int]",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"]},
			wlLen:     len1 - 1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 - 1 - 1},
		{name: "make interface{} <- string (t2)",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"], l["t5_f"]},
			wlLen:     len1 - 1 - 1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 - 1 - 1 - 1},
		{name: "t5",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f"], l["t4_f"], l["t5_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"]},
			wlLen:     len1 - 1 - 1 - 1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 - 1 - 1 - 1 - 1},
		{name: "slice t3[:]",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"]},
			wlLen:     len1 - 1 - 1 - 1 - 1 - 1 - 1 - 1, // no call - no new elements
			wlLenRet:  len1 - 1 - 1 - 1 - 1 - 1 - 1 - 1},
		{name: ":string, t6...)",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"], l["t7_f"]},
			wlLen:     len1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 - 1, // call but not analyzed as not part of the package
			wlLenRet:  len1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 - 1},
		{name: "return",
			updatedin: []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"], l["t7_f"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f"], l["t2_f"], l["t3_f.1"], l["t4_f.1"], l["t5_f"], l["t6_f"], l["t7_f"]},
			wlLen:     len1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 - 1,
			wlLenRet:  len1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 - 1 + 2}, // two transitions
		// element 31 of the worklist
		// nodes caused by the transitions
		{name: "g(\"Hello Gopher\":string)",
			updatedin: []string{l["c_f"], l["t0_f"], l["t1_f"]},
			flowout:   []string{l["c_f"], l["t0_f"], l["t1_f"], l["t2_f"]},
			wlLen:     len5,
			wlLenRet:  len5}, // no change
		{name: "sink(t0)",
			updatedin: []string{l["c_f.1"], l["t0_f.1"]},
			flowout:   []string{l["c_f.1"], l["t0_f.1"], l["t1_f.1"]},
			wlLen:     len5 - 1,
			wlLenRet:  len5 - 1},
		{name: "go f(t0)",
			updatedin: []string{l["t0_f.1"]},
			flowout:   []string{l["t0_f.1"]},
			wlLen:     len5 - 1 - 1,
			wlLenRet:  len5 - 1 - 1},
		{name: "go f(t0)",
			updatedin: []string{l["t0_f"]},
			flowout:   []string{l["t0_f"]},
			wlLen:     len5 - 1 - 1 - 1,
			wlLenRet:  len5 - 1 - 1 - 1},
	}

	iterateWl(should, t)
	assert.Equal(t, worklist.Len(), 0)
	logging()
}
