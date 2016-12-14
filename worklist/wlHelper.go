package worklist

import (
	"github.com/akwick/gotcha/lattice"
	"github.com/akwick/gotcha/lattice/taint"
	"log"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/ssa"
)

// ErrInFlows holds several ErrInFlow errors.
type ErrInFlows struct {
	errs []*taint.ErrLeak
}

func newErrInFlows() *ErrInFlows {
	e := &ErrInFlows{errs: make([]*taint.ErrLeak, 0)}
	return e
}
func (e *ErrInFlows) add(err *taint.ErrLeak) {
	for _, errs := range e.errs {
		if errs.Error() == err.Error() {
			return
		}
	}
	e.errs = append(e.errs, err)
}

// Error returns a string of all flows beeing in e.
func (e *ErrInFlows) Error() (s string) {
	for _, err := range e.errs {
		if err != nil {
			s += err.Error()
		}
	}
	return
}

// NumberOfFlows returns the number of taint.ErrLeaks in ErrInFlows
func (e *ErrInFlows) NumberOfFlows() int {
	return len(e.errs)
}

// Returns a slice with the predeccessors of instruction n
func predeccessors(n ssa.Instruction) map[ssa.Instruction]bool {
	pred := make(map[ssa.Instruction]bool)
	// Add all predecessors within the basic block
	basicBlock := n.Block()
	for _, instr := range basicBlock.Instrs {
		if instr == n {
			break
		}
		_, alreadyIn := pred[instr]
		if !alreadyIn {
			pred[instr] = true
		}
	}

	// Add instructions of the basic blocks which are before instruction n
	predBB := basicBlock.Preds
	////log.Printf("instruction: %s | pred basicBlocks: %v", n.String(), predBB)
	var predRec func(b *ssa.BasicBlock, p map[ssa.Instruction]bool) map[ssa.Instruction]bool
	predRec = func(b *ssa.BasicBlock, p map[ssa.Instruction]bool) map[ssa.Instruction]bool {
		bbNew := true
		for _, instr := range b.Instrs {
			_, alreadyIn := p[instr]
			if !alreadyIn {
				p[instr] = true
				bbNew = false
			}
		}
		if bbNew {
			for _, bb := range b.Preds {
				return predRec(bb, p)
			}
		}
		return p
	}
	for _, bb := range predBB {
		//		//log.Printf("Before predRec(%s), pred= %v", bb.String(), pred)
		pred = predRec(bb, pred)
		//		//log.Printf("After predRec(%s), pred= %v", bb.String(), pred)
	}
	////log.Printf("instruction: %s | Predecessors: %v", n.String(), pred)
	return pred
}

// succs returns a slice of dominees of instruction n
// Use map to avoid a duplication of instructions
func succs(n ssa.Instruction) map[ssa.Instruction]bool {
	succ := make(map[ssa.Instruction]bool)
	// All successors within the basic block
	basicBlock := n.Block()
	meetInstr := false
	for _, instr := range basicBlock.Instrs {
		if meetInstr {
			_, alreadyIn := succ[instr]
			if !alreadyIn {
				succ[instr] = true
			}
		}
		if instr == n {
			meetInstr = true
		}
	}

	// Add instructions of the basic blocks which are dominees of n's basic block.
	bDominees := basicBlock.Dominees()
	for _, dominee := range bDominees {
		if dominee != nil {
			for _, instr := range dominee.Instrs {
				_, alreadyIn := succ[instr]
				if !alreadyIn {
					succ[instr] = true
				}
			}
		}
		for _, bb := range dominee.Dominees() {
			bDominees = append(bDominees, bb)
		}
	}

	return succ
}

// updateEntryContext updates the in lattice of n with the idoms of n.node
func updateEntryContext(n *ContextCallSite) error {
	// Update the nodes =/= entry node
	entryNode := isEntryNode(n.node)
	if !entryNode {
		var lupLattice lattice.Latticer
		if isPointer {
			lupLatticep := taint.NewLatticePointer(0, valToPtr)
			// Add all pointers to the node
			lupLatticep.SetPtrs(valToPtr)
			lupLattice = lupLatticep
		} else {
			lupLattice = taint.NewLattice(0)
		}
		// Bug #751
		// Predecessors = all direct! Predecessors
		node := n.Node()
		block := node.Block()
		// 1) node is first element in Block (=/= entry node (first block and first element within the first block):
		//    Get last instruction of predecessor (idom)
		// 2) node is withing the Block:
		//    Get the lattice of the instruction before node
		if block.Instrs[0] == node {
			preds := block.Preds
			for i, b := range preds {
				lasti := b.Instrs[len(b.Instrs)-1]
			LoopCcsPool:
				for _, ccs := range ccsPool {
					if ccs.Node() == lasti {
						if ccs.Context().SameID(n.Context()) {
							if i == 0 {
								lupLattice = ccs.GetOut()
								// we have multiple predecessors
							} else {
								var err error
								lupLattice, err = lupLattice.LeastUpperBound(ccs.GetOut())
								if err != nil {
									return errors.Wrap(err, "failed lup of predecessors")
								}
							}
							continue LoopCcsPool
						}
					}
				}
			}
		} else {
		LoopInstrs:
			for i, instr := range block.Instrs {

				if instr == node {
					nbefore := block.Instrs[i-1]
				LoopCcsPool2:
					for _, ccs := range ccsPool {
						if ccs.Node() == nbefore {
							if ccs.Context().SameID(n.Context()) {
								lupLattice = ccs.GetOut()

								break LoopCcsPool2
							}
						}
					}
					break LoopInstrs
				}
			}
		}
		n.SetIn(lupLattice)
	} else {
		context := n.Context()
		n.SetIn(context.EntryValue())
	}
	return nil
}

// Returns true if the node is the entry node of the function
func isEntryNode(n ssa.Instruction) bool {
	parentFunc := n.Parent()
	if parentFunc.Blocks[0].Instrs[0] == n {
		return true
	}
	return false
}

// getIdomsBlocks returns a slice of basic blocks which idoms(dominates) b.
// The parameter b should be the first idom of the block the caller is interested in.
func getIdomsBlocks(b *ssa.BasicBlock) []*ssa.BasicBlock {
	var idoms []*ssa.BasicBlock
	// b will be nil for the entry node and a recover node (see API)
	if b != nil {
		idoms = append(idoms, b)
		// Iterate over all idiom of a basic block.
		// Idiom() only returns one basic block which directly dominates b
		idom := b.Idom()
		for idom != nil {
			idoms = append(idoms, idom)
			idomNew := idom.Idom()
			// Stop condition: If idom == idom.Idom()
			if idomNew != nil && idomNew.Index != idom.Index {
				idom = idomNew
			} else {
				idom = nil
			}
		}
	}
	return idoms
}

// for all edges <X', c> -> X return c in a slice
func getTransToX(x *ContextCallSite) []*ContextCallSite {
	var ccss []*ContextCallSite
	if x != nil && x.ValueContext != nil {
		for _, t := range transitions {
			log.Printf("t: %s   | len transition %d\n", t.String(), len(transitions))
			if t.targetContext.Equal(x.ValueContext) && t.Context() != t.targetContext {
				for _, ccs := range ccsPool {
					// A (context) -> B (target context)
					// Context out of the pool (ccs.Context) should be equal to the non target context (t.context).
					// Each transition is activated through a node  (t.node) and only this node of ccs should be added.
					//					if ccs.Context().Equal(t.context) && x.Node() == t.node {
					if ccs.Context().Equal(t.Context()) && ccs.Node() == t.Node() {
						ccss = append(ccss, ccs)
					}
				}
			}
		}
	}
	return ccss
}

// build lup of l1 and l2 and returns a taint.Lattice
func buildLUP(l1 lattice.Latticer, l2 lattice.Latticer) (lattice.Latticer, error) {
	lupl, err := l2.LeastUpperBound(l1)
	if err != nil {
		return nil, err
	}
	return lupl, nil
}

// Match statement tries to match the return values of the calle to the returned values of the caller
func matchRetStatements(calleeL lattice.Latticer, callerL lattice.Latticer, calleeI *ssa.Function, callerI ssa.Instruction) lattice.Latticer {
	// Build the LUP of the return statements of the callee
	blocks := calleeI.Blocks
	var lupValer []lattice.Valuer
	for _, bb := range blocks {
		// only the last element of a basic block can be a return statement
		retVal := bb.Instrs[len(bb.Instrs)-1]
		retValRes, ok := retVal.(*ssa.Return)
		if ok {
			res := retValRes.Results
			for _, r := range res {
				lupValer = append(lupValer, calleeL.GetVal(r))
			}
		}
	}
	// Unprecise: Match lupValer to one value
	var lupval lattice.Valuer
	lupval = taint.Uninitialized
	if len(lupValer) > 0 {
		lupval = lupValer[0]
		for _, val := range lupValer[1:] {
			var err error
			lupval, err = lupval.LeastUpperBound(val)
			if err != nil {
				errors.Wrap(err, "")
			}
		}
	}

	// Unprecise set lupVal to callee node (unprecise because of possible extract)
	var callValue ssa.Value
	switch c := callerI.(type) {
	case *ssa.Call:
		callValue = c
	case *ssa.Defer:
		callValue = c.Common().StaticCallee()
	case *ssa.Go:
		callValue = c.Common().StaticCallee()
	}
	retLat := callerL
	//retLat := callerL.DeepCopy()
	retLat.SetVal(callValue, lupval)
	return retLat
	//	callerL.SetVal(callValue, lupval)
	//return callerL
}

// getSuccessors returns the directly successors of i.
// When i is not the last instruction of a basic block b only the next instruction of b is returned.
// In the other case the first instruction of all Dominees of b will be returned.
func getSuccessors(i ssa.Instruction) []ssa.Instruction {
	var succs []ssa.Instruction
	b := i.Block()
	if i == b.Instrs[len(b.Instrs)-1] {
		for _, succ := range b.Succs {
			succs = append(succs, succ.Instrs[0])
		}
		/*	dominees := b.Dominees()
			for _, d := range dominees {
				succs = append(succs, d.Instrs[0])
			} */
	} else {
		for j, k := range b.Instrs {
			if k == i {
				succs = append(succs, b.Instrs[j+1])
			}
		}
	}
	return succs
}

// matchParams match the params of the caller to the values used from the callee
// A function (callee) has a field Params which offers the params.
// matchParams has a param with the params of the caller and the lattice.
// For matching we assume that the order of the params between callee and caller is equal.
func matchParams(pcaller []ssa.Value, lcaller lattice.Latticer, callee *ssa.Function, isClosure bool) lattice.Latticer {
	// Get the Params (normal call) or freevars(closure) of the callee
	var pcallee []ssa.Value
	if isClosure {
		fvs := callee.FreeVars
		pcallee = make([]ssa.Value, len(fvs))
		for i, fv := range fvs {
			pcallee[i] = fv
		}
	} else {
		ps := callee.Params
		pcallee = make([]ssa.Value, len(ps))
		for i, p := range ps {
			pcallee[i] = p
		}
	}

	// Generate the new lattice which should be set
	var l lattice.Latticer
	if isPointer {
		l = taint.NewLatticePointer(0, valToPtr)
	} else {
		l = taint.NewLattice(0)
	}

	// Match the params
	for i, val := range pcaller {
		l.SetVal(pcallee[i], lcaller.GetVal(val))
	}

	return l
}

func getContext(n *ContextCallSite, s ssa.Instruction) *ContextCallSite {
	var c *ContextCallSite
	c = nil
	call, isCall := s.(ssa.CallInstruction)
	noUpdate := true
	if s.Parent() == n.Context().Method() {
		for _, ccs := range ccsPool {
			if ccs.Node() == s {
				if isCall {
					ins := ccs.GetIn()
					args := call.Common().Args
					for _, in := range args {
						nin := n.GetOut().GetVal(in)
						nccs := ins.GetVal(in)
						eq, _ := nin.Equal(nccs)
						if !eq {
							noUpdate = false
						}
					}
					if noUpdate {
						c = ccs
					}
				} else {
					c = ccs
				}
			}
		}
	}
	return c
}

func getContextChannel(n *ContextCallSite, s ssa.Instruction) *ContextCallSite {
	if s.Parent() == n.Context().Method() {
		for _, ccs := range ccsPool {
			if ccs.Node() == s {
				return ccs
			}
		}
	}
	return nil
}

func updateNewVC(vc *ValueContext, ps []ssa.Value, c *ContextCallSite, callee *ssa.Function, call ssa.Instruction) {
	vcInContexts, knownContext, err := vcs.vcKnown(vc, ps)
	if err != nil {
		errors.Wrap(err, "failed to check whether value context is known")
	}
	// Add a transition from the caller context to the callee context (= known context to avoid duplicates) - line 27
	callerContext := c.Context()
	if knownContext == nil {
		knownContext = vc
	}
	calleeContext := knownContext
	// line 27 - add a new transition (if not already existing)
	NewTransition(callerContext, calleeContext, call)

	// line 28-32
	if vcInContexts {
		// Line 29 calleContext.ExitValue() as function call in the next line
		// Match return value of callerContext with y Line 30
		b1 := matchRetStatements(calleeContext.ExitValue(), callerContext.ExitValue(), calleeContext.Method(), call)
		// Line 31 : Assumption no side effects occurs
		b2 := c.GetOut()
		// Line 32: Set OUT(X, n) <- LUP(b1, b2)
		lup, err := b1.LeastUpperBound(b2)
		if err != nil {
			errors.Wrap(err, "")
		}
		vc.exitValue = lup
		c.SetOut(lup)
	}
}

// checks whether the node of c is a return.
// If the node is a return statement:
// - update the exit lattice of the context
// - if the exit lattice has changed: get all transitions to the node and add them to the worklist
func checkAndHandleReturn(c *ContextCallSite) {
	ret := checkReturn(c)
	if ret {
		handleReturn(c)
	}
}

// checkReturn returns true if c's node is a *ssa.Return statement.
func checkReturn(c *ContextCallSite) bool {
	_, ok := c.Node().(*ssa.Return)
	return ok
}

// handleReturn updates the exit lattice of the contextes
// and adds the contexts to the worklist which has a transfer to the node.
func handleReturn(c *ContextCallSite) {
	// update the output value of the value context
	c.Context().NewExitValue(c.GetOut())
	// add transitions with target comatchntext equal to the current value context to the worklist.
	ccss := getTransToX(c)
	for _, d := range ccss {
		log.Printf("d %s\n", d.String())
		worklist.Add(d)
	}
}

func logging() {
	log.Printf("### Contexts (unsorted!): %d ###\n", vcs.len())
	log.Printf("%s\n", vcs.String())
	log.Printf("### Transitions: %d ###\n", len(transitions))
	for _, t := range transitions {
		log.Printf("%s\n", t.String())
	}
	log.Printf("### Maximum elements in worklist: %d ### \n", maxElemsIn)

	stat.Printf("#contexts, %d,", vcs.len())
	stat.Printf("#transitions, %d,", len(transitions))
	stat.Printf("#maxWorklist, %d", maxElemsIn)
}
