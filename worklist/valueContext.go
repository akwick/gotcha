package worklist

import (
	"github.com/akwick/gotcha/lattice"
	"github.com/akwick/gotcha/lattice/taint"
	"log"
	"strconv"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/ssa"
)

// ValueContext is a struct to identify a value context
// A value context consists of a call, the entry and exit lattice.
// Further a value context has a unique id which helps to distinguish between different value contexts.
type ValueContext struct {
	*valueContext
}

type valueContext struct {
	VcIdentifier
	exitValue lattice.Latticer
	id        int
	params    []ssa.Value
}

// VcIdentifier provides all the necessary methods to compare two value contexts
type VcIdentifier interface {
	// GetIn returns the entry lattice of a value context
	GetIn() lattice.Latticer
	// SetIn updates the entry lattice of a value context with l
	SetIn(l lattice.Latticer)
	// GetFunction returns the function of a value context
	GetFunction() *ssa.Function
	// SetFunction sets the function of a value context to f
	SetFunction(f *ssa.Function)
	// Equal returns true if the function and the entry lattice of two value contextes are equal
	Equal(v VcIdentifier) bool
}

type vcident struct {
	in       lattice.Latticer
	function *ssa.Function
}

func (v vcident) GetIn() lattice.Latticer {
	return v.in
}

// SetIn sets the in lattice to l
// TODO implement update
func (v *vcident) SetIn(l lattice.Latticer) {
	v.in = l
}

func (v vcident) GetFunction() *ssa.Function {
	return v.function
}

func (v *vcident) SetFunction(f *ssa.Function) {
	v.function = f
}

func (v vcident) Equal(i VcIdentifier) bool {
	ineq, _ := v.GetIn().Equal(i.GetIn())
	return ineq && v.GetFunction() == i.GetFunction()
}

var idCounter int
var notTypeTaintLattice = "%s is not of type taint.Lattice"

// String returns a readable version of a value context
func (v *ValueContext) String() string {
	return " [" + strconv.Itoa(v.id) + "] Method: " + v.GetFunction().String() + " \n entryValue (Lattice) : " + v.GetIn().String() + " \n exitValue (Lattice) : " + v.ExitValue().String()
}

// GetValueContext returns a new value context, if the context is not known already.
// If the context is already known, the function returns the known context.
// If a new context is created, the in lattice will be set.
func GetValueContext(callee *ssa.Function, pcaller []ssa.Value, lcaller lattice.Latticer, isClosure bool) (*ValueContext, error) {
	_ = "breakpoint"
	if callee == nil {
		return nil, errors.New("don't expect nil as callee in function GetValueContext")
	}

	lin := matchParams(pcaller, lcaller, callee, isClosure)

	// check whether a context is already created
	// If we found a equal context (method and in lattice equal), we will return the lattice
	{
		c, err := vcs.findInContext(callee, lin)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		if c != nil {
			return c, nil
		}
	}

	// check whether the function matches to a source or sink
	// -> don't create a new value context
	{
		var sinkSources []*taint.Data
		sinkSources = taint.Sources
		sinkSources = append(sinkSources, taint.Sinks...)
		for i, s := range sinkSources {
			// interface signature is equal
			if s.IsInterface() {
				if callee.Signature.String() == s.GetSig() {
					return nil, nil
				}
			} else {
				if i == len(sinkSources)-1 {
					log.Printf("s string %s\n", s.String())
				}
				if callee.Signature.String()+" "+callee.String() == s.String() {
					return nil, nil
				}
			}
		}
	}

	_ = "breakpoint"
	// we don't found context -> create a new context
	// add context to context list and initialize the context
	vc := newValueContext(callee)
	vc.SetIN(lin)
	vcs.addToContext(vc)
	log.Printf("  new vc: %s", vc.String())
	initContextVC(callee, vc)
	return vc, nil
}

type VCIfer interface {
	addToContext(vc *ValueContext)
	findInContext(callee *ssa.Function, lin lattice.Latticer) (*ValueContext, error)
	len() int
	vcKnown(vc *ValueContext, params []ssa.Value) (bool, *ValueContext, error)
	String() string
}

type VCS struct {
	ctx map[VcIdentifier]*ValueContext
}

var vcs VCIfer

func (v VCS) len() int {
	return len(v.ctx)
}

// vcKnown returns either return true or false when vc is already known (the function and entry lattice is eqal).
// Check:
// 1) is the potential value context in the map
// 2) are the in values equal (TODO verify how to become rid of the second check)
// Return values: If check 1 && 2 is successfull: true, the known context and no error.
// Else: If check 1 successful && check 2 returns false: result of check 2, the known context and no error.
// If check 2 returns an error the function returns an error (and false, nil) too.
func (v VCS) vcKnown(vc *ValueContext, params []ssa.Value) (bool, *ValueContext, error) {

	//vcidentier is not a comparble type in Go
	//var knownContext *ValueContext = v.ctx[vc.VcIdentifier]
	var knownContext *ValueContext
	for k, c := range v.ctx {
		if k.Equal(vc.VcIdentifier) {
			knownContext = c
		}
	}
	if knownContext == nil {
		return false, nil, nil
	}
	entryValEq, err := knownContext.INEqual(vc, params)
	if err != nil {
		return false, nil, err
	}
	return entryValEq, knownContext, nil
}

func (v *VCS) addToContext(vc *ValueContext) {
	v.ctx[vc.VcIdentifier] = vc
}

// findInContext tries to find a context for callee and lin.
// If a context is found: this will be returned.
// If an error occurs while checking: the error will be returned.
// If nothing is found and no error occurs: (nil, nil) will be returned.
func (v VCS) findInContext(callee *ssa.Function, lin lattice.Latticer) (*ValueContext, error) {
	// create a vci to check against
	vci := &vcident{in: lin, function: callee}
	//knownContext := v.ctx[vci]
	var knownContext *ValueContext
	for k, e := range v.ctx {
		if k.Equal(vci) {
			knownContext = e
		}
	}
	if knownContext == nil {
		return nil, nil
	}
	return knownContext, nil
}

// String returns a string reperesentation for v
func (v VCS) String() string {
	s := ""
	for _, a := range v.ctx {
		s += " " + a.String()
	}
	s += "\n"
	return s
}

// newValueContext initialize a new value context
func newValueContext(method *ssa.Function) *ValueContext {
	if method == nil {
		return nil
	}
	var lEntry, lExit lattice.Latticer
	if isPointer {
		lEntry = taint.NewLatticePointer(0, valToPtr)
		lExit = taint.NewLatticePointer(0, valToPtr)
	} else {
		lEntry = taint.NewLattice(0)
		lExit = taint.NewLattice(0)
	}
	vc := &ValueContext{&valueContext{
		VcIdentifier: &vcident{in: lEntry, function: method},
		exitValue:    lExit, id: idCounter}}
	idCounter++
	return vc
}

// SetIN sets the in value for a new value context.
// params are the parameters of the functions and only those will be added to the entry lattice.
// l is the entry lattice of the node which generates the new value context.
func (v *ValueContext) SetIN(l lattice.Latticer) {
	v.VcIdentifier.SetIn(l)
}

// UpdateOut update the exit lattice of v with the lup of the current value and l.
func (v *ValueContext) UpdateOut(l lattice.Latticer) error {
	old := v.exitValue
	new, err := old.GreatestLowerBound(l)
	if err != nil {
		return err
	}
	v.exitValue = new
	return nil
}

func (v *ValueContext) GetIn() lattice.Latticer {
	return v.VcIdentifier.GetIn()
}

func (v *ValueContext) GetOut() lattice.Latticer {
	return v.exitValue
}

//INEqual tests whether the entry Lattice of two value contexts are equal
// params are the parameter for the function which creates v.
// v2 is an arbitrary value context which should checked against .
func (v *ValueContext) INEqual(v2 *ValueContext, params []ssa.Value) (bool, error) {
	// TODO work with VCIdentifier interface
	// as v2.params are nil if the len is 0 as nothing is inside
	//if v2.params == nil {
	//	return false, nil
	//}
	if len(v2.params) != len(v.params) {
		return false, nil
	}
	if v2.Method() != v.Method() {
		return false, nil
	}
	if v2.params == nil || len(params) != len(v2.params) {
		return false, nil
	}
	if len(params) > 0 {
		for i, val := range params {
			valTaint1 := v.EntryValue().GetVal(val)
			valTaint2 := v.EntryValue().GetVal(v2.params[i])
			if valTaint1 != valTaint2 {
				return false, nil
			}
		}
	}
	return true, nil
}

// GetID returns the id of a valuecontext
func (v *ValueContext) GetID() int {
	return v.id
}

// SameID compares the id of v against the id of v2 and return true if v1.id == v2.id
func (v *ValueContext) SameID(v2 *ValueContext) bool {
	return v.id == v2.id
}

// Equal returns true if the both vale contextes are equal.
// Comparision is based upon the method, the Lattices and the id.
func (v *ValueContext) Equal(v2 *ValueContext) bool {
	equalMethod := v.GetFunction() == v2.GetFunction()
	equalLattice, err := v.GetIn().Equal(v2.GetIn())
	if err != nil {
		return false
	}
	equalID := v.SameID(v2)
	return equalMethod && equalLattice && equalID
}

// Method returns the ssa.Function of a valueContext
func (v *ValueContext) Method() *ssa.Function {
	return v.VcIdentifier.GetFunction()
}

// EntryValue returns the entryValue lattice of a valueContext
func (v *ValueContext) EntryValue() lattice.Latticer {
	return v.VcIdentifier.GetIn()
}

// ExitValue returns the exitValue lattice of a valueContext
func (v *ValueContext) ExitValue() lattice.Latticer {
	return v.exitValue
}

// NewEntryValue updates the entryValue of v with a least upper bound with ev.
func (v *ValueContext) NewEntryValue(ev lattice.Latticer) {
	l, _ := v.VcIdentifier.GetIn().LeastUpperBound(ev)
	v.VcIdentifier.SetIn(l)
}

// NewExitValue updated the exitValue of v with a least upper bound with ev.
func (v *ValueContext) NewExitValue(ev lattice.Latticer) {
	v.exitValue, _ = v.exitValue.LeastUpperBound(ev)
}
