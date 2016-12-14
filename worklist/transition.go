package worklist

import (
	"golang.org/x/tools/go/ssa"
)

// Transition represents a transition from a value context to another value context.
// The transition is caused by a node which is a call in the context causing the change to the tagetContext.
type Transition struct {
	*transition
}

type transition struct {
	context       *ValueContext
	targetContext *ValueContext
	node          ssa.Instruction
}

// Equal tests whether two transitions are equal
func (t1 *Transition) Equal(t2 *Transition) bool {
	equalContext := t1.context.Equal(t2.context)
	equalTargetContext := t1.targetContext.Equal(t2.targetContext)
	equalnode := t1.node == t1.node
	return equalContext && equalTargetContext && equalnode
}

// NewTransition returns a Transition.
// context is the valuecontext of the "start" node.
// callSite is the ssa.Function which "goes" to the other context.
// targetContext is the valuecontext which is the "goal" of the flow.
// NewTransition appends the new transition to the transitions slice.
func NewTransition(context *ValueContext, targetContext *ValueContext, node ssa.Instruction) {
	if context == nil || targetContext == nil || node == nil {
		return
	}
	// don't create a transition if we don't analyze the context because we know the abstract results.
	if context.Equal(targetContext) {
		return
	}
	t := &Transition{&transition{context: context, targetContext: targetContext, node: node}}
	// only add a new transition to the transition list
	for _, a := range transitions {
		if a.Equal(t) {
			return
		}
	}
	transitions = append(transitions, t)
}

// String returns a readable string of a transition.
func (t1 *Transition) String() string {
	return "Context: " + t1.context.String() + "  \n node: " + t1.node.String() + " | targetContext: " + t1.targetContext.String()
}

// Node returns the node of the transition
func (t1 *Transition) Node() ssa.Instruction {
	return t1.node
}

// Context returns the context of the caller
func (t1 *Transition) Context() *ValueContext {
	return t1.context
}

// TagerContext returns the context of the callee
func (t1 *Transition) TargetContext() *ValueContext {
	return t1.targetContext
}
