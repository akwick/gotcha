// taintatticePointer is a taint lattice which also holds the pointer for every ssa.Valuer
package taint

import (
	"github.com/akwick/gotcha/lattice"
	"github.com/akwick/gotcha/transferFunction"
	"go/token"

	"github.com/pkg/errors"

	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

// LatticePointer is an instance of Pter and extends Lattice with pointer information
type LatticePointer struct {
	*latticePointer
}

// latticePointer holds the Lattice and the additional pointer information in the map p
type latticePointer struct {
	l Lattice
	p map[ssa.Value]pointer.Pointer
}

func NewLatticePointer(len int, m map[ssa.Value]pointer.Pointer) *LatticePointer {
	lat := make(map[ssa.Value]Value, len)
	l := &latticePointer{l: lat, p: m}
	return &LatticePointer{l}
}

var (
	onlyTaintLatPotAcc = "taintLatticePointer: can handle only taint.LatticePointer, but not %s"
)

func (l1 *LatticePointer) LeastUpperBound(l2 lattice.Latticer) (lattice.Latticer, error) {
	lat, err := l1.l.LeastUpperBound(l2)
	latTainted, ok := lat.(Lattice)
	if !ok || err != nil {
		return nil, errors.Errorf("failed %s LUP %s", l1.String(), l2.String())
	}
	ptrs := l1.GetPtrs()
	//	ptrs := copyPtrs(l1.GetPtrs())
	l := &latticePointer{l: latTainted, p: ptrs}
	return &LatticePointer{l}, nil
}

func (l1 *LatticePointer) GreatestLowerBound(l2 lattice.Latticer) (lattice.Latticer, error) {
	lat, err := l1.l.GreatestLowerBound(l2)
	latTainted, ok := lat.(Lattice)
	if !ok || err != nil {
		return nil, errors.Errorf("failed %s GLB %s", l1.String(), l2.String())
	}
	ptrs := l1.GetPtrs()
	//ptrs := copyPtrs(l1.GetPtrs())
	l := &latticePointer{l: latTainted, p: ptrs}
	return &LatticePointer{l}, nil
}

func (l1 *LatticePointer) LeastElement() (lattice.Latticer, error) {
	lat, err := l1.l.LeastElement()
	latTainted, ok := lat.(Lattice)
	if !ok || err != nil {
		return nil, errors.Errorf("failed LeastEement(%s) ", l1.String())
	}
	ptrs := l1.GetPtrs()
	//	ptrs := copyPtrs(l1.GetPtrs())
	l := &latticePointer{l: latTainted, p: ptrs}
	return &LatticePointer{l}, nil
}

func (l1 *LatticePointer) Less(l2 lattice.Latticer) (bool, error) {
	return l1.l.Less(l2)
}

func (l1 *LatticePointer) Equal(l2 lattice.Latticer) (bool, error) {
	return l1.l.Equal(l2)
}

func (l1 *LatticePointer) LessEqual(l2 lattice.Latticer) (bool, error) {
	return l1.l.LessEqual(l2)
}

func (l1 *LatticePointer) Greater(l2 lattice.Latticer) (bool, error) {
	return l1.l.Greater(l2)
}

func (l1 *LatticePointer) GreaterEqual(l2 lattice.Latticer) (bool, error) {
	return l1.l.GreaterEqual(l2)
}

func (l1 *LatticePointer) BottomLattice() lattice.Latticer {
	lat := l1.l.BottomLattice()
	// TODO currently no check for ok
	latTainted := lat.(Lattice)
	ptrs := l1.GetPtrs()
	//	ptrs := copyPtrs(l1.GetPtrs())
	l := &latticePointer{l: latTainted, p: ptrs}
	return &LatticePointer{l}
}

func (l1 *LatticePointer) DeepCopy() lattice.Latticer {
	ptrs := copyPtrs(l1.GetPtrs())
	lat := l1.GetLat().(Lattice)
	l := &latticePointer{l: lat, p: ptrs}
	return &LatticePointer{l}
}

func (l1 *LatticePointer) String() string {
	return "LatticePointer: " + l1.l.String()
}

func (l1 *LatticePointer) GetVal(key ssa.Value) lattice.Valuer {
	return l1.l.GetVal(key)
}

func (l1 *LatticePointer) SetVal(key ssa.Value, value lattice.Valuer) error {
	return l1.l.SetVal(key, value)
}

func (l1 *LatticePointer) GetPtr(key ssa.Value) pointer.Pointer {
	return l1.p[key]
}

func (l1 *LatticePointer) SetPtr(key ssa.Value, ptr pointer.Pointer) {
	l1.p[key] = ptr
}

func (l1 *LatticePointer) GetLat() lattice.Latticer {

	if l1 == nil || l1.l == nil {
		m := make(map[ssa.Value]pointer.Pointer)
		NewLatticePointer(0, m)
	}
	return l1.l
}

func (l1 *LatticePointer) GetPtrs() map[ssa.Value]pointer.Pointer {
	return l1.p
}

func (l1 *LatticePointer) SetPtrs(m map[ssa.Value]pointer.Pointer) {
	l1.p = m
}

func (l1 *LatticePointer) GetSSAValMayAlias(v ssa.Value) []ssa.Value {
	rets := make([]ssa.Value, 0)
	vptr := l1.GetPtr(v)
	ptrs := l1.GetPtrs()
	for v, ptr := range ptrs {
		if ptr.MayAlias(vptr) {
			rets = append(rets, v)
		}
	}
	return rets
}

//Semanticer interface

func (l1 *LatticePointer) TransferFunction(node ssa.Instruction, ptr *pointer.Result) transferFunction.PlainFF {
	switch nType := node.(type) {
	case *ssa.UnOp:
		if nType.Op != token.MUL && nType.Op != token.ARROW {
			l := l1.GetLat().(Lattice)
			return l.TransferFunction(node, ptr)
		}
		//handling unop ptrs
		return ptrUnOp(nType, l1, ptr)
	case *ssa.Store:
		// *t1 = t0
		// everything which mayalias the addres should set to the lattice value of val.
		addr := nType.Addr
		lupVal := l1.GetVal(nType.Val)
		if ptr != nil {
			if ok, addrp := IsPointerVal(addr); ok {
				q := ptr.Queries[addrp]
				qset := q.PointsTo()
				lbaels := qset.Labels()
				for _, l := range lbaels {
					l1.GetLat().SetVal(l.Value(), lupVal)
					for ssav, p := range l1.GetPtrs() {
						if p.MayAlias(l1.GetPtr(l.Value())) {
							l1.GetLat().SetVal(ssav, lupVal)
						}
					}
				}
			}
		}
		l1.GetLat().SetVal(addr, lupVal)
	case *ssa.Call, *ssa.Defer, *ssa.Go:
		ff := checkAndHandleSourcesAndsinks(node, l1, true)
		if ff == nil {
			return returnID
		} else {
			return ff
		}
	}
	l := l1.GetLat().(Lattice)
	return l.TransferFunction(node, ptr)
}

// helper function
func copyPtrs(ptrs map[ssa.Value]pointer.Pointer) map[ssa.Value]pointer.Pointer {
	p := make(map[ssa.Value]pointer.Pointer)
	for i, ptr := range ptrs {
		p[i] = ptr
	}
	return p
}
