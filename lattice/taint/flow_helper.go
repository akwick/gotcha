package taint

import (
	"go/types"

	"github.com/akwick/gotcha/lattice"
	"github.com/akwick/gotcha/transferFunction"

	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

// IsPointerVal checks whether the value is a pointer value and
// in the positive case it returns the ssa.Value
// For some instructions like a *ssa.Call it is requuired to go deeper.
// TODO check whether it is also required for other cases
func IsPointerVal(i ssa.Value) (canPoint bool, val ssa.Value) {
	if call, ok := i.(*ssa.Call); ok {
		val = call.Common().Value
	} else {
		val = i
	}
	return pointer.CanPoint(val.Type()), val
}

// IsIndirectPtr checks whether the value is an indirect pointer value
// In the positive case the function returns the ssa.Value.
func IsIndirectPtr(i ssa.Value) (canPoint bool, val ssa.Value) {
	if call, ok := i.(*ssa.Call); ok {
		val = call.Common().Value
	} else {
		val = i
	}
	// call function like described in the api
	//	fmt.Printf("val: %s | val.Type() %v | val.Type().Underlying() %v\n", val, val.Type(), val.Type().Underlying())
	_, isRange := val.(*ssa.Range)
	if isRange {
		return false, val
	}
	_, isPointer := val.Type().Underlying().(*types.Pointer)
	if !isPointer {
		return false, val
	}

	return pointer.CanPoint(val.Type().Underlying().(*types.Pointer).Elem()), val
}

// getPointsToSet returns the pointsToset for a given ssa.Value
func getPointsToSet(val ssa.Value, ptrRes *pointer.Result) (pointer.PointsToSet, error) {
	ptsto := ptrRes.Queries[val]
	//	//log.Printf("Value (%s) is pointer: %v", val.String(), ptsto)

	ptstoset := ptsto.PointsTo()

	return ptstoset, nil
}
func getPointsToSetIndirect(val ssa.Value, ptrRes *pointer.Result) (pointer.PointsToSet, error) {
	ptsto := ptrRes.IndirectQueries[val]
	//log.Printf("Value (%s) is pointer: %v", val.String(), ptsto)

	ptstoset := ptsto.PointsTo()

	return ptstoset, nil
}

// LUPOverLabels iterates through all labels of a pointsToSet
// and returns the result of the LUP.
func LUPOverLabels(ptstoset pointer.PointsToSet, l Lattice) (val Value, e error) {
	var ssaVal ssa.Value
	labels := ptstoset.Labels()
	for _, label := range labels {
		ssaVal = label.Value()
		var valI lattice.Valuer
		valI, err = l.GetVal(ssaVal).LeastUpperBound(val)
		val, _ = valI.(Value)
		if e != nil {
			return Uninitialized, e
		}
		if val == Both {
			return
		}
	}
	return
}

func ptrUnOp(e *ssa.UnOp, l lattice.Pter, ptr *pointer.Result) transferFunction.PlainFF {
	value := e.X
	lupVal = l.GetVal(e.X)

	if ptr != nil {
		if ok, valr := IsPointerVal(value); ok {
			q := ptr.Queries[valr]
			labels := q.PointsTo().Labels()
			for _, la := range labels {

				l.GetLat().SetVal(la.Value(), lupVal)
				for ssav, p := range l.GetPtrs() {
					if p.MayAlias(l.GetPtr(la.Value())) {
						l.GetLat().SetVal(ssav, lupVal)
					}
				}
			}
		}

		if ok, valr := IsIndirectPtr(value); ok {
			q := ptr.Queries[valr]
			labels := q.PointsTo().Labels()
			for _, la := range labels {

				l.GetLat().SetVal(la.Value(), lupVal)
				for ssav, p := range l.GetPtrs() {
					if p.MayAlias(l.GetPtr(la.Value())) {
						l.GetLat().SetVal(ssav, lupVal)
					}
				}
			}
		}
	}
	return returnLUP
}

func setAllPointsTo(val lattice.Valuer, l lattice.Latticer, ptss pointer.PointsToSet) {
	labels := ptss.Labels()
	for _, label := range labels {
		ssaVal := label.Value()
		l.SetVal(ssaVal, val)
	}
}

func isSource(c ssa.CallCommon) bool {
	sig, call, iSig := getSignature(c)
	for _, source := range Sources {
		if source.IsInterface() {
			if iSig == source.sig {
				return true
			}
		}
		if sig == source.sig {
			if call == source.callee {
				return true
			}
		}
	}
	return false
}

func isSink(c ssa.CallCommon) bool {
	sig, call, iSig := getSignature(c)
	for _, sink := range Sinks {
		if sink.IsInterface() {
			if iSig == sink.sig {
				return true
			}
		}
		if sig == sink.sig {
			if call == sink.callee {
				return true
			}
		}
	}
	return false
}

func checkAndHandleSourcesAndsinks(c ssa.Instruction, l lattice.Latticer, ptr bool) transferFunction.PlainFF {
	// ensure that err is nil because it is later used to distinguish whether a information flow occurs or not.
	err = nil
	var callCom ssa.CallCommon
	switch xType := c.(type) {
	case *ssa.Call:
		callCom = xType.Call
	case *ssa.Defer:
		callCom = xType.Call
	case *ssa.Go:
		callCom = xType.Call
	default:
		return nil
	}
	// Get the lup of the value
	lupVal = l.GetVal(callCom.Value)
	source := isSource(callCom)
	if source {
		return returnTainted
	}
	sink := isSink(callCom)
	if sink {
		handleSinkDetection(callCom, l, ptr)
	}

	if err != nil {
		return returnLUPTaint
	}

	return nil
}
