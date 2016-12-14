package taint

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// ssaValMock is a simple struct which implements all methods required for the interface ssa.Value.
// The functions return nil for complex types.
type ssaValMock struct {
	N string
}

func (m ssaValMock) Name() string {
	return m.N
}
func (m ssaValMock) String() string {
	return m.N
}
func (m ssaValMock) Parent() *ssa.Function {
	return nil
}
func (m ssaValMock) Referrers() *[]ssa.Instruction {
	return nil
}
func (m ssaValMock) Type() types.Type {
	return nil
}
func (m ssaValMock) Pos() token.Pos {
	return 0
}

// Setting up different combinations of mocks for testing
var (
	m1 = new(ssaValMock)
	m2 = new(ssaValMock)
	m3 = new(ssaValMock)
	m4 = new(ssaValMock)
	m5 = new(ssaValMock)
)

func initM() {
	m1.N = "t1"
	m2.N = "t2"
	m3.N = "t3"
	m4.N = "t4"
	m5.N = "t5"
}

// Building Lattices with the Mocks
// L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
func getLatticeMock1() Lattice {
	initM()
	l1 := make(Lattice)
	l1[m1] = Uninitialized
	l1[m2] = Tainted
	l1[m3] = Untainted
	l1[m5] = Both
	return l1
}

// L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
func getLatticeMock2() Lattice {
	initM()
	l1 := make(Lattice)
	l1[m1] = Tainted
	l1[m2] = Untainted
	l1[m3] = Uninitialized
	l1[m4] = Untainted
	return l1
}

// L3: {t2 -> Empty, t3 -> Empty, t5 -> Tainted}
func getLatticeMock3() Lattice {
	initM()
	l1 := make(Lattice)
	l1[m3] = Uninitialized
	l1[m2] = Uninitialized
	l1[m5] = Tainted
	return l1
}

// L4: {t2 -> Tainted, t3 -> Untainted, t5 -> Both}
func getLatticeMock4() Lattice {
	initM()
	l1 := make(Lattice)
	l1[m3] = Untainted
	l1[m2] = Tainted
	l1[m5] = Both
	return l1
}
