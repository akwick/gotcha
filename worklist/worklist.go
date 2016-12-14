package worklist

import (
	"fmt"
	"github.com/akwick/gotcha/lattice/taint"
	"github.com/akwick/gotcha/ssabuilder"
	"github.com/akwick/gotcha/transferFunction"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var logFile *os.File
var now string

// Variables for the creation of SSA and the pointer information
var ssaProg *ssa.Program
var conf *pointer.Config
var pta *pointer.Result
var valToPtr map[ssa.Value]pointer.Pointer
var mainFunc *ssa.Function

// Variables to for running the analysis
var transitions []*Transition
var worklist *WlList

// contains all context call sites which are produced to build the union, but not the nicest and efficients solution.
// Other idea: Map[Context] to Array of ValueContexts or Map[Node] to ContextCallSites
var ccsPool []*ContextCallSite
var errFlows *ErrInFlows

// Error messages
var failLUP = "buildLUP(%s, %s) failed"

// Log stats
var stat *log.Logger

func init() {
	now = time.Now().String()
	setLogger()
}

// Logging
func setLogger() {
	var err error
	logFile, err = os.OpenFile(now+"_Log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("%v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile)
	// TODO defer handling
	//	defer logFile.Close()

	logStat, err := os.OpenFile(now+"_Log.stat", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//	defer logStat.Close()
	if err != nil {
		fmt.Printf("failed creating stat file: %v", err)
	}
	stat = log.New(logStat, "", 0)
}

// DoAnalysis handles the worklist algorithm.
// path is the relative path starting from $GOPATH
// sourcefiles are the source file which should be analyzed.
// sourceAndSinkFile is the file which contains the sources and sinks.
// error will be returned if an error occurs during execution.
// If a flow from a source to a sink occurs, the information will be packed
// into an error of type ErrInFlows.
// Comments with line xx refers to the algorithm of worklist (Figure 1)
func DoAnalysis(path string, sourcefiles []string, sourceAndSinkFile string, allpkgs bool, pkgs string, ptr bool) error {
	// Initialization (includes line 13)
	wlInit(path, sourcefiles, sourceAndSinkFile, allpkgs, pkgs, ptr)

	// Worklist algorithm
	for !worklist.Empty() {
		// Take the first context callsite and remove it from the worklist (line 15)
		nextccs := worklist.RemoveFirst()
		// lupPred = LUP(of all nextccs.pred.Out)
		// IN(nextccs) = LUP(IN(nextccs), lupPred)
		if err := updateEntryContext(nextccs); err != nil {
			return errors.Wrapf(err, "failed updateEntryContext(%s)", nextccs.String())
		}
		// flow the node Line 23-39
		if err := flow(nextccs); err != nil {
			return errors.Wrap(err, "error at function flow")
		}
		// check whether the lattice changed through the flow Line 39
		// If the lattice changes -> add all nodes to the worklist which are
		// influenced by the node (Line 39-42).
		if err := nextccs.CheckAndHandleChange(); err != nil {
			return errors.Wrap(err, "error at function CheckAndHandleChange")
		}
		// Handle return nodes: A return node transfers it's out lattice to the value context.
		// Every transition from X' -> X has to be added to the worklist
		// Line 44-49
		checkAndHandleReturn(nextccs)
	}
	// Return flows if flows occurs
	if errFlows.NumberOfFlows() > 0 {
		return errFlows
	}
	// Log contexts and transitions
	logging()
	return nil
}

// flow handles the flow of a node
// Distinguishes between a normal call and a method call which enforces a change of the context.
func flow(c *ContextCallSite) error {
	// Check, whether c.in implements the Semanticer interface.
	// (Else a flow is not possible because methods are missing.)
	lIn, ok := c.GetIn().(transferFunction.Semanticer)
	if !ok {
		log.Fatalf("%v throws an error because it doesn't implement the interface transferFunction.Semanticer", c.in)
		return errors.New("The in lattice of contextcallsite doesn't implement the interfache transferFunction.Semanticer")
	}

	var ff transferFunction.PlainFF
	// Flow the node [a normal as well as a call node]
	// The flow function can also handle on not ssa.Values, but the lattice needs a ssa.Value
	switch n := c.node.(type) {
	case ssa.Value, *ssabuilder.Send, *ssa.Store, *ssa.Go, *ssa.Defer:
		// Get the flow function for n upon the lattice and handle the case that ff = nil
		ff = lIn.TransferFunction(n, pta)
		if ff == nil {
			em := "the returned flow function is nil. Called upon Lattice: % \n  with node: %s and pta %v\n"
			log.Fatalf(em, lIn, c.Node().String(), pta)
			return errors.Errorf(em, lIn, c.Node().String(), pta)
		}

		// Get the correct value for the flow function (for Send and ssa.Store the ssa.Value is within the type)
		// and the value which should be changed through the flow.
		// In the case of a *Send the out value is the channel.
		var valin, valout ssa.Value
		switch n := n.(type) {
		case ssa.Value:
			valin = n
			valout = n
		case *ssabuilder.Send:
			valin = n.X
			valout = n.Chan
		case *ssa.Store:
			valin = n.Val
			valout = n.Val
		case *ssa.Go:
			valin = n.Common().Value
			valout = n.Common().Value
		case *ssa.Defer:
			valin = n.Common().Value
			valout = n.Common().Value
		}
		valn := c.GetIn().GetVal(valin)

		// Flow the flow function and handle the errors
		ffval, err := ff(valn)
		if err != nil {
			switch err := err.(type) {
			case taint.ErrLeak:
				errFlows.add(&err)
			default:
				return errors.Wrapf(err, "failed call ff with %s", valn.String())
			}
			err = nil
		}

		// Build the lup of the current value which should be set and the value from the flow function
		lupval, err := ffval.LeastUpperBound(c.GetIn().GetVal(valout))
		if err != nil {
			log.Fatalf("%s\n", err.Error())
			return errors.Wrapf(err, "failed call LeastUpperBound with %s and %s", ffval.String(), c.GetIn().GetVal(valout))
		}

		// Set the out lattice for the value
		c.SetOut(c.GetIn())
		c.GetOut().SetVal(valout, lupval)
	case *ssa.Jump, *ssa.Return:
		// ToDo improve - current problem: can't pass n as value for the lattice (is *ssa.Instruction)
		c.SetOut(c.GetIn())
	}

	// Get information whether the node is a call or a closure
	// Line 23
	switch call := c.node.(type) {
	// handle in next case
	case ssa.CallInstruction, *ssa.MakeClosure, *ssa.Call:
		var staticCallee *ssa.Function
		var paramsCaller []ssa.Value
		var vc *ValueContext
		var err error
		switch call := call.(type) {
		case ssa.CallInstruction:
			// Get the static callee of the node
			// Line 24 - 26
			var callCom *ssa.CallCommon
			callCom = call.Common()
			staticCallee = callCom.StaticCallee()
			// *Builtin or any other value indicating a dynamically dispatched function call
			if staticCallee != nil {
				// staticCalle is the targetMethod of the call
				vc, err = GetValueContext(staticCallee, callCom.Args, c.GetIn(), false)
				if err != nil {
					return errors.Wrap(err, "failed get value context")
				}
			}
		case *ssa.MakeClosure:
			fn, ok := call.Fn.(*ssa.Function)
			if !ok {
				errors.Errorf("unexpected: call(%s).Fn should be of type *ssa.Function", call)
			}
			vc, err = GetValueContext(fn, call.Bindings, c.GetIn(), true)
			if err != nil {
				return errors.Wrap(err, "")
			}

			staticCallee = fn
			paramsCaller = call.Bindings
		}
		// the parameters of the call are added to the lattice -> updating the out lattice too
		// SetOut builds the LUP
		c.SetOut(c.GetIn()) // local flow
		if vc != nil {
			// TODO rename
			updateNewVC(vc, paramsCaller, c, staticCallee, call)
		}
	}

	return nil
}
