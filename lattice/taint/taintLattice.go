package taint

import (
	"goretech/analysis/lattice"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/ssa"
)

// Lattice is data structure (map) which maps a ssa.Value to the abstract value.
// ssa.Value is an interface -> no pointer needed because it's already a pointer.
type Lattice map[ssa.Value]Value

// Several constant values for throwing errors with a formated string
var (
	onlyTaintLatAcc   = "taintLattice: can handle only taint.Lattice, but get %s"
	onlyTaintValueAcc = "taintLattice: can handle only taint.Value, but get %s"
	failedEqual       = "failed on operation %s.Equal(%s)"
	failedGreater     = "failed on operation %s.Greater(%s)"
)

// NewLattice returns a new lattice with length len.
func NewLattice(len int) Lattice {
	l := make(map[ssa.Value]Value, len)
	return l
}

// getTaintLattice returns the Lattice instance of an Latticer instance
// Support Lattice and LatticePointer
// For other instances of Latticer, the function will return an error.
func getTaintLattice(l lattice.Latticer) (Lattice, error) {
	ltaint, ok := l.(Lattice)
	if ok {
		return ltaint, nil
	} else {
		lp, ptr := l.(lattice.Pter)
		if ptr {
			lpl := lp.GetLat()
			ltaint, ok = lpl.(Lattice)
			if ok {
				return ltaint, nil
			}
		}
		return nil, errors.Errorf(onlyTaintLatAcc, reflect.TypeOf(l))
	}
}

// LeastUpperBound creates the leastUpperBound of two Lattices.
// In the current implementation an error (ErrOnlyTaintLatAcc) will be returned if l2 is not of type taint.Lattice
// The function supports all intances of Latticer which returns a lattice when called with getTaintLattice.
func (l1 Lattice) LeastUpperBound(l2 lattice.Latticer) (lattice.Latticer, error) {
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		return nil, errors.Wrapf(err, onlyTaintLatAcc, reflect.TypeOf(l2taint))
	}

	var added bool
	// identify the bigger lattice to optimize the iterations thorugh the lattices
	var smallerL, biggerL Lattice
	if len(l1) > len(l2taint) {
		smallerL = l2taint.DeepCopy().(Lattice)
		biggerL = l1.DeepCopy().(Lattice)
	} else {
		smallerL = l1.DeepCopy().(Lattice)
		biggerL = l2taint.DeepCopy().(Lattice)
	}

RangeLO:
	for ssaValO, valO := range biggerL {
		added = false
		// prevent that the algorithm operates on a nil value
		if ssaValO == nil {
			continue RangeLO
		}
	RangeLI:
		for ssaValI, valI := range smallerL {
			// prevent that the algorithm operates on a nil value
			if ssaValI == nil {
				continue RangeLI
			}
			// found a match between two ssa values -> build the lup of them
			if ssaValO == ssaValI {
				lub, err := valO.LeastUpperBound(valI)
				if err != nil {
					return nil, err
				}
				var lubTaint Value
				var ok bool
				if lubTaint, ok = lub.(Value); !ok {
					return nil, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lub))
				}
				// update the value with the newly computed abstract value
				biggerL[ssaValO] = lubTaint
				added = true
				delete(smallerL, ssaValI)
				// no further match in the inner lattice should occur -> stop loop
				continue RangeLO
			}
		}
		if !added {
			biggerL[ssaValO] = valO
		}
	}

	// elements which are in the inner lattice but not in the outer lattice
	// works as the elements which are in both are deleted in the previous loop
	for ssaValI, valI := range smallerL {
		biggerL[ssaValI] = valI
	}
	return biggerL, nil
}

// GreatestLowerBound computes the greatest lower bound of two lattices.
// In the current implementation an error (ErrOnlyTaintLatAcc) will be returned if l2 is not of type taint.Lattice
// The function supports all intances of Latticer which returns a lattice when called with getTaintLattice.
// TODO optimized version is in lub -> merge into one function with higher order funcitons
func (l1 Lattice) GreatestLowerBound(l2 lattice.Latticer) (lattice.Latticer, error) {
	retMap := NewLattice(len(l1))
	// Deep copy of l1 into the return map
	for ssaValThis, valueThis := range l1 {
		retMap[ssaValThis] = valueThis
	}
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		errors.Wrap(err, "failed GLB get taint lattice")
	}
	var added bool
	for ssaValL2, valL2 := range l2taint {
		// Iterating over l1 to have the chance to compare the ssa.Values.
		added = false
		for ssaValL1, valL1 := range retMap {
			// Comparising with Name () - but its propably not unique
			// TODO better solution
			if ssaValL2.Name() == ssaValL1.Name() {
				var glbTaint Value
				glb, err := valL1.GreatestLowerBound(valL2)
				if err != nil {
					return nil, errors.Wrapf(err, "%s greatest lower bound (%s) failed", valL1.String(), valL2.String())
				}
				var ok bool
				if glbTaint, ok = glb.(Value); !ok {
					return nil, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(glb))
				}
				retMap[ssaValL1] = glbTaint
				added = true
			}
		}
		if !added {
			retMap[ssaValL2] = valL2
		}
	}
	retMapNew := make(Lattice)
	for key, item := range retMap {
		retMapNew[key] = item
	}
	return retMapNew, nil
}

// LeastElement returns a Lattice in which each element of l1 is set to the
// lowest abstract value of Valuer being Empty
func (l1 Lattice) LeastElement() (lattice.Latticer, error) {
	retMap := make(Lattice)
	for ssaVal := range l1 {
		retMap[ssaVal] = Uninitialized
	}
	return retMap, nil
}

// Less computes whether l1 is less than l2.
// Throws an error if l2 is not of type taint.Lattice.
func (l1 Lattice) Less(l2 lattice.Latticer) (bool, error) {
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		return false, errors.Wrap(err, "failed less get taint lattice")
	}
	for ssaValL1, valL1 := range l1 {
		// Iterating over l2 to have the chance to compare the ssa.Values.
		var visited bool
		visited = false
		for ssaValL2, valL2 := range l2taint {
			// Comparising with Name () - but its propably not unique
			// TODO better solution
			if ssaValL2.Name() == ssaValL1.Name() {
				less, err := valL1.Less(valL2)
				if err != nil {
					return false, errors.Wrap(err, "failed less of taint values")
				}
				if !less {
					return false, nil
				}
				visited = true
			}
		}
		if !visited && valL1 != Uninitialized {
			return false, nil
		}
	}
	return true, nil
}

// Equal compares l1 against l2 on equality.
// Throws an error if l2 is not of type taint.Lattice.
func (l1 Lattice) Equal(l2 lattice.Latticer) (equal bool, err error) {
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		return false, errors.Wrap(err, "failed equal get taint lattice")
	}
	equal = true
	for ssaValL1, valL1 := range l1 {
		// Iterating over l2 to have the chance to compare the ssa.Values.
		var visited bool
		visited = false
		for ssaValL2, valL2 := range l2taint {
			if ssaValL2 == ssaValL1 {
				equal, err = valL1.Equal(valL2)
				if err != nil {
					return false, errors.Wrapf(err, failedEqual, valL1.String(), valL2.String())
				}
				if !equal {
					return false, nil
				}
				visited = true
			}
		}
		if !visited && valL1 != Uninitialized {
			return false, nil
		}
	}
	for ssaValL2, valL2 := range l2taint {
		var visited bool
		visited = false
		for ssaValL1, valL1 := range l1 {
			if ssaValL2 == ssaValL1 {
				equal, err = valL2.Equal(valL1)
				if err != nil {
					return false, errors.Wrapf(err, failedEqual, valL2.String(), valL1.String())
				}
				if !equal {
					return false, nil
				}
				visited = true
			}
		}
		if !visited && valL2 != Uninitialized {
			return false, nil
		}
	}
	return
}

// LessEqual computes: l1 <= l2.
// Throws an error if t2 is not of type taint.Lattice.
func (l1 Lattice) LessEqual(l2 lattice.Latticer) (lesseq bool, err error) {
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	lesseq, err = l1.Equal(l2taint)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	if lesseq {
		return true, nil
	}
	lesseq, err = l1.Less(l2)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	return lesseq, nil
}

// Greater computes: l1 > l2.
// Throws an error if t2 is not of type taint.Lattice.
func (l1 Lattice) Greater(l2 lattice.Latticer) (bool, error) {
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	for ssaValL1, valL1 := range l1 {
		// Iterating over l2 to have the chance to compare the ssa.Values.
		var visited bool
		visited = false
		for ssaValL2, valL2 := range l2taint {
			// Comparising with Name () - but its propably not unique
			// TODO better solution
			if ssaValL2.Name() == ssaValL1.Name() {
				greater, err := valL1.Greater(valL2)
				if err != nil {
					return false, errors.Wrapf(err, failedGreater, valL1.String(), valL2.String())
				}
				if !greater {
					return false, nil
				}
				visited = true
			}
		}
		if !visited && valL1 != Uninitialized {
			return false, nil
		}
	}
	return true, nil
}

// GreaterEqual computes: l1 >= l2
func (l1 Lattice) GreaterEqual(l2 lattice.Latticer) (greatereq bool, err error) {
	l2taint, err := getTaintLattice(l2)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	greatereq, err = l1.Equal(l2taint)
	if err != nil || greatereq {
		return
	}
	greatereq, err = l1.Greater(l2taint)
	return
}

// BottomLattice sets for all elements of the lattice the value to empty.
func (l1 Lattice) BottomLattice() lattice.Latticer {
	retLattice := make(Lattice)
	for ssaValL1 := range l1 {
		retLattice[ssaValL1] = Uninitialized
	}
	return retLattice
}

// DeepCopy copies l1 and returns a new lattice.
func (l1 Lattice) DeepCopy() lattice.Latticer {
	retLattice := make(Lattice)
	for ssaVal, latVal := range l1 {
		retLattice[ssaVal] = latVal
	}
	return retLattice
}

// String returns a string representation of l1.
func (l1 Lattice) String() string {
	s := ""
	for ssaVal, val := range l1 {
		if ssaVal == nil {
			s += "nil : " + val.String() + "|"
		} else {
			if strings.Contains(ssaVal.Name(), "nil") {
				withOutNil := strings.Replace(ssaVal.Name(), "nil", "", -1)
				s += "nil + " + withOutNil + " : " + val.String() + "|"
			} else {
				s += ssaVal.Name() + " : " + val.String() + "|"
			}
		}
	}
	return s
}

// GetVal returns the value of key from the lattice.
// Empty will be returned if the key is not in the lattice.
// In such a case we set manually the value of these element to Empty.
func (l1 Lattice) GetVal(key ssa.Value) lattice.Valuer {
	lVal := l1[key]
	if lVal == Uninitialized {
		// We assume that constant values are untainted -> set them to untainted.
		if _, cnst := key.(*ssa.Const); cnst {
			l1[key] = Untainted
			lVal = Untainted
		} else {
			l1[key] = Uninitialized
		}
	}
	return lVal
}

// SetVal sets the lattice value of key to value.
func (l1 Lattice) SetVal(key ssa.Value, value lattice.Valuer) error {
	val, ok := value.(Value)
	if !ok {
		return errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(value))
	}
	if key != nil {
		l1[key] = val
	}
	return nil
}

// Value represents a taint Value
type Value int

const (
	// Unitialized represents the lowest abstract value.
	Uninitialized Value = iota
	// Tainted represents a tainted abstract value.
	Tainted
	// Untainted represents an untainted abstract value.
	Untainted
	// Both represents the highest abstract value.
	Both
)

// BottomElement returns the lowest element of the lattice.
func (tlv Value) BottomElement() lattice.Valuer {
	return Uninitialized
}

// TopElement returns the highest element of the lattice.
func (tlv Value) TopElement() lattice.Valuer {
	return Both
}

// LeastUpperBound build the least upper bound of tlv and lv.
// Returns an error (ErrOnlyTaintLatAcc) if lv is not a taint.Value.
func (tlv Value) LeastUpperBound(lv lattice.Valuer) (lattice.Valuer, error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return nil, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	return tlv | tlv2, nil
}

// GreatestLowerBound build the greatest lower bound of tlv and lv.
// Returns an error (ErrOnlyTaintValueAcc) if lv is not a taint.Value.
func (tlv Value) GreatestLowerBound(lv lattice.Valuer) (lattice.Valuer, error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return nil, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	return tlv & tlv2, nil
}

// Less returns tlv < lv.
// Returns an error (ErrOnlyTaintValueAcc) if lv is not a taint.Value.
func (tlv Value) Less(lv lattice.Valuer) (bool, error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return false, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	if (tlv == Tainted && tlv2 == Untainted) || (tlv == Untainted && tlv2 == Tainted) {
		return false, nil
	}
	return (tlv < tlv2), nil
}

// Equal returns tlv == lv.
// Returns an error (ErrOnlyTaintValueAcc) if lv is not a taint.Value.
func (tlv Value) Equal(lv lattice.Valuer) (bool, error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return false, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	return tlv == tlv2, nil
}

// LessEqual returns tlv <= lv.
// Returns an error (ErrOnlyTaintValueAcc) if lv is not a taint.Value.
func (tlv Value) LessEqual(lv lattice.Valuer) (lesseq bool, err error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return false, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	lesseq, err = tlv.Less(tlv2)
	lesseq = lesseq || tlv == tlv2
	return
}

// Greater returns tlv > lv.
// Returns an error if lv is not type taint.Value.
func (tlv Value) Greater(lv lattice.Valuer) (bool, error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return false, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	if (tlv == Tainted && tlv2 == Untainted) || (tlv == Untainted && tlv2 == Tainted) {
		return false, nil
	}
	return (tlv > tlv2), nil
}

// GreaterEqual returns tlv >= lv.
// Returns an error if lv is not of type taint.Value.
func (tlv Value) GreaterEqual(lv lattice.Valuer) (greatereq bool, err error) {
	tlv2, ok := lv.(Value)
	if !ok {
		return false, errors.Errorf(onlyTaintValueAcc, reflect.TypeOf(lv))
	}
	greatereq, err = tlv.Greater(tlv2)
	greatereq = greatereq || tlv == tlv2
	return
}

// String returns a string representation of tlv.
func (tlv Value) String() string {
	if tlv == Uninitialized {
		return "Uninitialized"
	}
	if tlv == Tainted {
		return "Tainted"
	}
	if tlv == Untainted {
		return "Untainted"
	}
	if tlv == Both {
		return "Both"
	}
	return "Unknown value"
}
