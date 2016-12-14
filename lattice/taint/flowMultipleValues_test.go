package taint

import (
	"testing"

	"golang.org/x/tools/go/ssa"
)

func TestMultipleSSAValues(t *testing.T) {
	l := make(Lattice)
	// ssa.BinOp: Tainted * Empty
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	ssaBinOp := new(ssa.BinOp)
	ssaBinOp.X = mock1
	ssaBinOp.Y = mock2
	l[mock1] = Tainted

	// ssa.ChangeType
	mock3 := new(ssaValMock)
	mock3.N = "t3"
	ssaChangeType := new(ssa.ChangeType)
	ssaChangeType.X = mock3
	l[mock3] = Untainted

	// ssa.Store
	mock4 := new(ssaValMock)
	mock4.N = "t4"
	mock5 := new(ssaValMock)
	mock5.N = "t5"
	ssaStore := new(ssa.Store)
	ssaStore.Val = mock4
	ssaStore.Addr = mock5
	l[mock4] = Untainted
	l[mock5] = Tainted

	var retLat Lattice
	ff := l.TransferFunction(ssaBinOp, nil)
	retVal, err = ff(l[ssaBinOp])
	if err != nil {
		t.Error(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	ff = l.TransferFunction(ssaChangeType, nil)
	retVal, err = ff(l[ssaChangeType])
	if err != nil {
		t.Error(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
	ff = l.TransferFunction(ssaStore, nil)
	retVal, err = ff(l[ssaStore.Val])
	if err != nil {
		t.Error(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Untainted, retVal)
	}

	t.Log(retLat)
}
