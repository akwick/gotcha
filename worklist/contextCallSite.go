package worklist

import (
	"goretech/analysis/lattice"
	"goretech/analysis/lattice/taint"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/ssa"
)

// ContextCallSite is a data structure which holds a value context,
// a ssa.Instruction and the in and out lattice for this node.
type ContextCallSite struct {
	*contextCallSite
}

type contextCallSite struct {
	*ValueContext
	node ssa.Instruction
	in   lattice.Latticer
	out  lattice.Latticer
}

// NewContextCallSite creates a new contextcallsite.
// context is the wanted value context
// node is the ssa.Instruction for the contextcallsite
// The in and out lattice is set to empty
func NewContextCallSite(context *ValueContext, node ssa.Instruction) *ContextCallSite {
	if context == nil || node == nil {
		return nil
	}
	var l1, l2 lattice.Latticer
	if isPointer {
		l1 = taint.NewLatticePointer(0, valToPtr)
		l2 = taint.NewLatticePointer(0, valToPtr)
	} else {
		l1 = taint.NewLattice(0)
		l2 = taint.NewLattice(0)
	}
	ccs := &ContextCallSite{&contextCallSite{ValueContext: context, node: node, in: l1, out: l2}}

	// add ccs to ccsPool
	ccsPool = append(ccsPool, ccs)
	return ccs
}

// Equal tests whether two context callsites are equal
func (c *ContextCallSite) Equal(ccs2 *ContextCallSite) bool {
	equalContext := c.Context().Equal(ccs2.Context())
	equalNode := c.node == ccs2.node
	equalIn, _ := c.in.Equal(ccs2.in)
	equalOut, _ := c.in.Equal(ccs2.in)
	return equalContext && equalNode && equalIn && equalOut
}

// Node returns the ssa.Instruction of c
func (c *ContextCallSite) Node() ssa.Instruction {
	return c.node
}

// Context returns the valueContext of c
func (c *ContextCallSite) Context() *ValueContext {
	return c.ValueContext
}

// String returns a string representation of c.
func (c *ContextCallSite) String() string {
	s := "context: " + c.Context().String() + " \n"
	s += "node " + c.node.String() + "\n"
	s += "in" + c.GetIn().String() + "\n"
	s += "out" + c.GetOut().String() + "\n"
	return s
	//	return "Context: " + c.context.String() + " \n node: " + c.node.String() + " \n in Lattice " + c.GetIn().String() + " \n out Lattice " + c.GetOut().String()
}

// SetIn sets l to the input lattice for c.
// new_in = LUP(old_in, l)
// TODO: handle the error.
func (c *ContextCallSite) SetIn(l lattice.Latticer) {
	c.in, _ = c.in.LeastUpperBound(l)
}

// SetOut sets l to the output lattice for c.
// new_out = LUP(old_out, l)
// TODO: handle the error.
func (c *ContextCallSite) SetOut(l lattice.Latticer) {
	c.out, _ = c.GetOut().LeastUpperBound(l)
}

// GetOut returns the out lattice of context c
func (c *ContextCallSite) GetOut() lattice.Latticer {
	return c.out
}

// GetIn returns the in lattice of context c
func (c *ContextCallSite) GetIn() lattice.Latticer {
	return c.in
}

// CheckAndHandleChange checks whether a value context has changed.
// In the case the value context has changed,
// it will add the successory of n's node to the worklist.
func (c *ContextCallSite) CheckAndHandleChange() error {
	change, err := checkChange(c)
	if change {
		handleChange(c)
	}
	return err
}

// checkChange checks whether the input lattice of n is equal to the output lattice of n
// returns true if the context has changed, false if it has not changed.
func checkChange(n *ContextCallSite) (bool, error) {
	change, err := n.GetIn().Equal(n.GetOut())
	if err != nil {
		return false, errors.Wrapf(err, "equal failed with \n   %s\n    %s", n.GetIn().String(), n.GetOut().String())
	}
	return !change, nil
}

// handleChange adds the successors of n's node to the worklist.
func handleChange(n *ContextCallSite) {
	worklist.AddSucc(n)
}
