package worklist

import (
	"github.com/akwick/gotcha/lattice"
	"github.com/akwick/gotcha/lattice/taint"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/tools/go/ssa"
)

func TestPredecessors(t *testing.T) {
	src = []string{"../tests/exampleCode/chanPaper0.go"}
	allpkgs = false
	_, _ = initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	var testInstr *ContextCallSite
	testInstr = worklist.order[2]
	pred := predeccessors(testInstr.Node())
	assert.Equal(t, 2, len(pred))
	expected := []string{
		"new string (x)",
		"*t0 = \"hello, world\":string",
	}
	for instr := range pred {
		_, i0 := instr.(*ssa.Alloc)
		if i0 {
			assert.Equal(t, expected[0], instr.String())
		} else {
			assert.Equal(t, expected[1], instr.String())
		}
	}
}

func BenchmarkInitContext(b *testing.B) {
	mainFunc, _ := initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	for n := 0; n < b.N; n++ {
		initContext(mainFunc)
	}
}

func BenchmarkDoAnalysis(b *testing.B) {
	// Tests without modification in the lattice -> is only for the performance
	for n := 0; n < b.N; n++ {
		DoAnalysis("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", false, "", true)
	}
}

func TestUpdatePredecessors(t *testing.T) {
	allpkgs = false
	mainFunc, _ := initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	//make interface{} <- string ("Hello, World!":string) = testinstr
	var instr ssa.Value
	testInstr := worklist.order[2]
	// Generate a lattice where testInstr is tainted.
	taintedLattice := make(taint.Lattice)
	var ok bool
	instr, ok = testInstr.Node().(ssa.Value)
	if !ok {
		t.Error("Type conversion error")
	}
	taintedLattice[instr] = 1

	// Set the exit lattice of testInstr to the tainted Lattice (= instr 2 marked as tainted).
	worklist.order[0].out = taintedLattice.DeepCopy()
	worklist.order[1].out = taintedLattice.DeepCopy()
	worklist.order[2].in = taintedLattice.DeepCopy()
	ccsPool[0].out = taintedLattice.DeepCopy()
	ccsPool[1].out = taintedLattice.DeepCopy()
	ccsPool[2].in = taintedLattice.DeepCopy()

	var wl1, wl2 *ContextCallSite
	wl1 = worklist.order[1]
	wl2 = worklist.order[2]
	for _, ccs := range ccsPool {
		if ccs.Node() == worklist.getFirstCCS().Node() || ccs.Node() == wl1.Node() {
			ccs.Context().exitValue = taintedLattice.DeepCopy()
		}
	}

	n := wl2

	updateEntryContext(n)
	for _, ccs := range ccsPool {
		if ccs.Node() == worklist.getFirstCCS().Node() || ccs.Node() == wl1.Node() {
			ccs.Context().exitValue = taintedLattice.DeepCopy()
		}
	}

	var emptyLattice, filledLattice lattice.Latticer
	if isPointer {
		emptyLattice = taint.NewLatticePointer(0, valToPtr)
	} else {
		emptyLattice = taint.NewLattice(0)
	}
	filledLattice = taintedLattice
	filledLattice.SetVal(instr, taint.Tainted)
	for _, ccs := range ccsPool {
		// Case first node:
		// expect as in lattice: empty lattice.
		// expect as out lattice: filledLattice = instr is tainted.
		if ccs.Node() == worklist.getFirstCCS().Node() || ccs.Node() == wl1.Node() {
			ok, err := ccs.GetIn().Equal(emptyLattice)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("Expected the empty Lattice, but get: %v", ccs.GetIn())
			}
			ok, err = ccs.GetOut().Equal(filledLattice)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("Expected Lattice %v as exitValue of %v, but get: %v", filledLattice, ccs.node, ccs.GetOut())
			}
		} else {
			// Case: second node:
			// expect as in lattie: filled Lattice (instr is untainted (=out lattice from node 1))
			// expect as out lattice:  filled lattice
			if ccs.Node() == wl2.Node() {
				ok, err := ccs.GetIn().Equal(filledLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Logf("ccs node: %v | entryValue: %v | exitValue: %v", ccs.Node(), ccs.Context().EntryValue(), ccs.Context().ExitValue())
					t.Errorf("Expected the Lattice %v as in of %v, but get: %v", filledLattice, ccs.node, ccs.GetIn())
				}
				ok, err = ccs.GetOut().Equal(emptyLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Errorf("Expected the empty Lattice, but get: %v", ccs.GetOut())
				}
			} else {
				// all other Lattices:
				// should be empty because testing the predecessors and not successors.
				ok, err := ccs.GetIn().Equal(emptyLattice)
				if !ok {
					t.Errorf("Expected for the context (%v) the empty lattice:, but get: %v", ccs.ValueContext, ccs.GetIn())
				}
				ok, err = ccs.GetOut().Equal(emptyLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Errorf("Expected the empty Lattice, but get: %v", ccs.GetOut())
				}
			}
		}
	}
}

func TestFlow(t *testing.T) {
	mainFunc, _ := initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	//make interface{} <- string ("Hello, World!":string)
	var instr ssa.Value
	var wl1, wl2, wl3 *ContextCallSite
	wl1 = worklist.order[1]
	wl2 = worklist.order[2]
	wl3 = worklist.order[3]

	testInstr := wl2
	// Generate a lattice where testInstr is tainted.
	taintedLattice := make(taint.Lattice)
	var ok bool
	instr, ok = testInstr.Node().(ssa.Value)
	if !ok {
		t.Error("Type conversion error")
	}
	taintedLattice[instr] = 2
	t.Logf("taintedLattice %v\n", taintedLattice)

	// Set the exit lattice of testInstr to the tainted Lattice.
	worklist.getFirstCCS().out = taintedLattice.DeepCopy()
	wl1.out = taintedLattice.DeepCopy()
	ccsPool[0].out = taintedLattice.DeepCopy()
	ccsPool[1].out = taintedLattice.DeepCopy()

	t.Logf("entry value of worklist[0]: %v", worklist.getFirstCCS().Context().EntryValue())
	for _, ccs := range ccsPool {
		if ccs.Node() == worklist.getFirstCCS().Node() || ccs.Node() == wl1.Node() {
			ccs.Context().exitValue = taintedLattice.DeepCopy()
			t.Logf("node: %v\n  entryValue: %v\n     exitValue: %v\n", ccs.Node(), ccs.Context().EntryValue(), ccs.Context().ExitValue())
		}
	}

	t.Logf("worklist 2 node: %v | context %v\n | entryValue %v | exitValue %v\n", wl2.Node(), wl2.Context(), wl2.Context().EntryValue(), wl2.Context().exitValue)
	t.Logf("worklist 3 node: %v | context %v\n | entryValue %v | exitValue %v\n", wl3.Node(), wl3.Context(), wl3.Context().EntryValue(), wl3.Context().exitValue)
	n := wl2
	updateEntryContext(n)
	t.Logf("worklist 2 node: %v | context %v\n | entryValue %v | exitValue %v\n", wl2.Node(), wl2.Context(), wl2.Context().EntryValue(), wl2.Context().ExitValue())
	for _, ccs := range ccsPool {
		if ccs.Node() == worklist.getFirstCCS().Node() || ccs.Node() == wl1.Node() {
			ccs.Context().exitValue = taintedLattice.DeepCopy()
			t.Logf("node: %v\n  entryValue: %v\n     exitValue: %v\n", ccs.Node(), ccs.GetIn(), ccs.GetOut())
		}
	}
	err := flow(wl2)
	if err != nil {
		t.Error(err.Error())
	}

	emptyLattice := make(taint.Lattice)
	filledLattice := make(taint.Lattice)
	filledLattice[instr] = 2
	for _, ccs := range ccsPool {
		if ccs.Node() == worklist.getFirstCCS().Node() || ccs.Node() == wl1.Node() {
			ok, err := ccs.GetIn().Equal(emptyLattice)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("Expected the empty Lattice, but get: %v", ccs.GetIn())
			}
			ok, err = ccs.GetOut().Equal(filledLattice)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("Expected Lattice %v as exitValue of %v, but get: %v", filledLattice, ccs.node, ccs.GetOut())
			}
		} else {
			if ccs.Node() == wl2.Node() {
				ok, err := ccs.GetIn().Equal(filledLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Logf("ccs node: %v | entryValue: %v | exitValue: %v", ccs.Node(), ccs.Context().EntryValue(), ccs.Context().ExitValue())
					t.Errorf("Expected the Lattice %v as in of %v, but get: %v", filledLattice, ccs.node, ccs.GetIn())
				}
				ok, err = ccs.GetOut().Equal(filledLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Errorf("Expected the Lattice %v, but get: %v", filledLattice, ccs.GetOut())
				}
			} else {
				ok, err := ccs.GetIn().Equal(emptyLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Errorf("Expected for the context (%v) the lattice: %v, but get: %v", ccs.Context(), filledLattice, ccs.GetIn())
				}
				ok, err = ccs.GetOut().Equal(emptyLattice)
				if err != nil {
					t.Error(err)
				}
				if !ok {
					t.Errorf("Expected the empty Lattice, but get: %v", ccs.GetOut())
				}
			}
		}
	}
}
