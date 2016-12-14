package worklist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var src = []string{"../tests/exampleCode/hello.go"}

func TestNumberElementsWorklist1(t *testing.T) {
	mainFunc, _ := initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	assert.Equal(t, 2, worklist.Len())
}

func TestNumberValueContexts(t *testing.T) {
	mainFunc, _ := initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	assert.Equal(t, 1, vcs.len())
}

// expect: 0 elements in transitions
func TestNumberTransitions(t *testing.T) {
	_, _ = initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	assert.Equal(t, 0, len(transitions))
}

//// expect: correct ssa instructions
func TestElementsWorklist(t *testing.T) {
	mainFunc, _ = initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt", "")
	initContext(mainFunc)
	should := [...]string{
		"\"Hello\":string + \" World\":string",
		"return",
	}
	//assumption: contexts contains only the correct value context.
	if assert.Equal(t, len(should), worklist.Len()) {
		for i := 0; i < len(worklist.wlMap); i++ {
			// TODO correct assertion to an map (it is not necessary to know the order of alle contexts)
			/*getContext := worklist.order[i].Context()
			shouldContext := vcs.ctx[0]
			assert.True(t, shouldContext.Equal(getContext)) */

			getInstr := worklist.order[i].Node()
			shouldInstr := should[i]
			assert.Equal(t, shouldInstr, getInstr.String())
			i++
		}
	}
}

// expect: value context{function: mainfunction, init: empty, exit: empty}
func TestElementValueContext(t *testing.T) {
	t.Skip()
	// TODO Test must be changed due to the change to a map [acessing elements in order is not any longer possible.]
	/*	mainFunc, _ := initSSAandPTA("github.com/akwick/gotcha", src, "../sourcesAndSinks.txt")
		initContext(mainFunc)
		var shouldLattice lattice.Latticer
		if isPointer {
			shouldLattice = taint.NewLatticePointer(0, valToPtr)
		} else {
			shouldLattice = taint.NewLattice(0)
		}
		shouldFunction := "main"
		if assert.Equal(t, vcs.len(), 1, "contexts should contain 1 element") {
			c := vcs.ctx[0]
			getFunctionName := c.method.Name()
			getEntryLattice := c.entryValue
			getExitLattice := c.exitValue
			if shouldFunction != getFunctionName {
				t.Errorf("Returned value should: %v, but is: %v", shouldFunction, getFunctionName)
			}
			ok, err := shouldLattice.Equal(getEntryLattice)
			if err != nil {
				t.Errorf("%s", err)
			}
			if !ok {
				t.Errorf("Returned value should: %v, but is: %v", shouldLattice, getEntryLattice)
			}
			ok, err = shouldLattice.Equal(getExitLattice)
			if err != nil {
				t.Errorf("%s", err)
			}
			if !ok {
				t.Errorf("Returned value should: %v, but is: %v", shouldLattice, getExitLattice)
			}
		}*/
}
