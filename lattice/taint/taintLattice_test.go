package taint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTlBottomElem tests whether the function LestElement works as expected.
/*
test object: t1 -> Empty | t2 -> Tainted | t3 -> Untainted | t5 -> Both
test result: new element with: t1 -> Empty | t2 -> Empty | t3 -> Empty | t5 -> Empty
test object should still the same.
test parameters: all element should be empty
number of elements in the new lattice should be 4
all names should occur once
the original lattice is the same as before calling
*/
func TestTlBottomElem(t *testing.T) {
	expectedLength := 4
	var l Lattice
	// L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
	l = getLatticeMock1()
	leastElem, err := l.LeastElement()
	assert.Nil(t, err)
	if assert.IsType(t, (Lattice)(nil), leastElem) {
		leastElemTaint := leastElem.(Lattice)
		// check whether all elements are empty
		for _, el := range leastElemTaint {
			assert.Equal(t, Uninitialized, el)
		}
		// check whether the length is equal to 4.
		assert.Equal(t, expectedLength, len(leastElemTaint))

		// check whether the names are equal
		t1 := 0
		t2 := 0
		t3 := 0
		t5 := 0
		for i := range leastElemTaint {
			switch i.Name() {
			case "t1":
				t1++
			case "t2":
				t2++
			case "t3":
				t3++
			case "t5":
				t5++
			default:
				t.Errorf("The name of the value should be t1, t2, t3 or t5, but is: %s", i.Name())
			}
		}
		assert.Equal(t, 1, t1)
		assert.Equal(t, 1, t2)
		assert.Equal(t, 1, t3)
		assert.Equal(t, 1, t5)
		checkNoChangeL1(t, l)
	}
}

// TestTlLeastUpperBound tests whether the lup of two distinct lattices are build correctly
/*
test object 1: t1 -> Empty | t2 -> Tainted | t3 -> Untainted | t5 -> Both
test object 2: t1 -> Tainted | t2 -> Untainted | t3 -> Empty | t4 -> Untainted
test result: t1 -> Tainted | t2 -> Both | t3 -> Untainted | t4 -> Untainted | t5 -> Both
to1 should be the same after calling the function
to2 should be the same after valling the function
the length should be 5.
*/
func TestTlLeastUpperBound(t *testing.T) {
	// L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
	l1 := getLatticeMock1()
	// L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
	l2 := getLatticeMock2()
	wanted := []struct {
		ssaVal   string
		taintVal Value
	}{
		{"t1", Tainted},
		{"t2", Both},
		{"t3", Untainted},
		{"t4", Untainted},
		{"t5", Both},
	}
	lub, err := l1.LeastUpperBound(l2)
	if assert.Nil(t, err) {
		if assert.IsType(t, (Lattice)(nil), lub) {
			lubTaint := lub.(Lattice)
			// check length
			assert.Equal(t, len(wanted), len(lubTaint))
			// check against the expected lattice
			for ssaVal, val := range lubTaint {
				for _, want := range wanted {
					if ssaVal.String() == want.ssaVal {
						assert.Equal(t, want.taintVal, val)
					}
				}
			}
			checkNoChangeL1(t, l1)
			checkNoChangeL2(t, l2)
		}
	}
}

// TestTlLeastUpperBoundEqualLattices tests whether the lup of two equal lattices is the same.
/*
test object 1: t1 -> Empty | t2 -> Tainted | t3 -> Untainted | t5 -> Both
test object 2: t1 -> Empty | t2 -> Tainted | t3 -> Untainted | t5 -> Both
test result: t1 -> Empty | t2 -> Tainted | t3 -> Untainted | t5 -> Both
The test result should be equal to test object 1.
Test object 1 should not change after the call.
Test object 2 should not change after the call.
*/
func TestTlLeastUpperBoundEqualLattices(t *testing.T) {
	l1 := getLatticeMock1()
	l2 := getLatticeMock1()
	lub, err := l1.LeastUpperBound(l2)
	if assert.Nil(t, err) {
		if assert.IsType(t, (Lattice)(nil), lub) {
			lubL, _ := lub.(Lattice)
			lattices := []Lattice{l1, l2, lubL}
			for _, lat := range lattices {
				eq, err := lat.Equal(getLatticeMock1())
				assert.Nil(t, err)
				assert.True(t, eq)
			}
			checkNoChangeL1(t, l1)
			checkNoChangeL1(t, l2)
		}
	}
}

/*
test object 1:   L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
test object 2:   L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
test result:     L: {t1 -> Empty, t2 -> Empty, t3 -> Empty, t4 -> Untainted, t5 -> Both}
Test object 1 should not change after the call.
Test object 2 should not change after the call.
*/
func TestTlGreatestLowerBound(t *testing.T) {
	l1 := getLatticeMock1()
	l2 := getLatticeMock2()
	wanted := []struct {
		ssaVal   string
		taintVal Value
	}{
		{"t1", Uninitialized},
		{"t2", Uninitialized},
		{"t3", Uninitialized},
		{"t4", Untainted},
		{"t5", Both},
	}
	glb, err := l1.GreatestLowerBound(l2)
	if assert.Nil(t, err) {
		if assert.IsType(t, (Lattice)(nil), glb) {
			glbTaint, _ := glb.(Lattice)
			assert.Equal(t, len(wanted), len(glbTaint))
			for ssaVal, val := range glbTaint {
				for _, want := range wanted {
					if ssaVal.String() == want.ssaVal {
						assert.Equal(t, want.taintVal, val)
					}
				}
			}
		}
		checkNoChangeL1(t, l1)
		checkNoChangeL2(t, l2)
	}
}

/*
Test object 1:  L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 2:  L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
test object 3:  L3: {t2 -> Empty, t3 -> Empty, t5 -> Tainted}
Test results:
L1 < L2 : false
L2 < L1 : false
L3 < L1 : true
L1 < L3 : fasle
All test objects should not change after the call
*/
func TestTlLess(t *testing.T) {
	l1 := getLatticeMock1()
	l2 := getLatticeMock2()
	l3 := getLatticeMock3()

	wanted := []struct {
		leftL  Lattice
		rightL Lattice
		less   bool
		origL  Lattice
	}{
		{l1, l2, false, getLatticeMock1()},
		{l2, l1, false, getLatticeMock2()},
		{l3, l1, true, getLatticeMock3()},
		{l1, l3, false, getLatticeMock1()},
	}

	for _, want := range wanted {
		less, err := want.leftL.Less(want.rightL)
		assert.Nil(t, err)
		assert.Equal(t, want.less, less)
		checkNoChange(t, want.leftL, want.origL)
	}

}

/*
Test object 1: L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 2: L4: {t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 3: L5: {}
Test object 4: L6: {}
Test object 5: L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
Test results:
l1 = l4 : true
l4 = l1 : true
l5 = l1 : false
l1 = l5 : false
l5 = l6 : true
l6 = l5 : true
l2 = l1 : false
l1 = l2 : false
All test objects should not change after the call
*/
func TestTlEqual(t *testing.T) {
	l1 := getLatticeMock1()
	l4 := getLatticeMock4()
	l5 := make(Lattice)
	l6 := make(Lattice)
	l2 := getLatticeMock2()

	wanted := []struct {
		leftL  Lattice
		rightL Lattice
		equal  bool
		origL  Lattice
	}{
		{l1, l4, true, getLatticeMock1()},
		{l4, l1, true, getLatticeMock4()},
		{l5, l1, false, make(Lattice)},
		{l1, l5, false, getLatticeMock1()},
		{l5, l6, true, make(Lattice)},
		{l6, l5, true, make(Lattice)},
		{l2, l1, false, getLatticeMock2()},
		{l1, l2, false, getLatticeMock1()},
	}

	for _, want := range wanted {
		eq, err := want.leftL.Equal(want.rightL)
		assert.Nil(t, err)
		assert.Equal(t, want.equal, eq)
		checkNoChange(t, want.leftL, want.origL)
	}
}

/*
Test object 1:  L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 2:  L4: {t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 3:  L3: {t2 -> Empty, t3 -> Empty, t5 -> Tainted}
Test object 4:  L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
Test result:
L1 <= l4 : true
l1 <= l2 : false
l2 <= l1 : false
l3 <= l1 : true
All test objects should not change after the call
*/
func TestTlLessEqual(t *testing.T) {
	l1 := getLatticeMock1()
	l4 := getLatticeMock4()
	l2 := getLatticeMock2()
	l3 := getLatticeMock3()

	wanted := []struct {
		leftL  Lattice
		rightL Lattice
		lesseq bool
		origL  Lattice
	}{
		{l1, l4, true, getLatticeMock1()},
		{l1, l2, false, getLatticeMock1()},
		{l2, l1, false, getLatticeMock2()},
		{l3, l1, true, getLatticeMock3()},
	}

	for _, want := range wanted {
		lesseq, err := want.leftL.LessEqual(want.rightL)
		assert.Nil(t, err)
		assert.Equal(t, want.lesseq, lesseq)
		checkNoChange(t, want.leftL, want.origL)
	}
}

/*
Test object 1:  L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 2:  L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
Test object 3:  L3: {t2 -> Empty, t3 -> Empty, t5 -> Tainted}
Test result:
l1 > l2 : false
l2 > l1 : false
l3 > l1 : false
l1 > l3 : true
All test objects should not change after the call
*/
func TestTlGreater(t *testing.T) {
	l1 := getLatticeMock1()
	l2 := getLatticeMock2()
	l3 := getLatticeMock3()

	wanted := []struct {
		leftL   Lattice
		rightL  Lattice
		greater bool
		origL   Lattice
	}{
		{l1, l2, false, getLatticeMock1()},
		{l2, l1, false, getLatticeMock2()},
		{l3, l1, false, getLatticeMock3()},
		{l1, l3, true, getLatticeMock1()},
	}

	for _, want := range wanted {
		greater, err := want.leftL.Greater(want.rightL)
		assert.Nil(t, err)
		assert.Equal(t, want.greater, greater)
		checkNoChange(t, want.leftL, want.origL)
	}
}

/*
Test object 1:  L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test object 2:  L2: {t1 -> Tainted, t2 -> Untainted, t3 -> Empty, t4 -> Untainted}
Test object 3:  L3: {t2 -> Empty, t3 -> Empty, t5 -> Tainted}
Test object 4:  L4: {t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Test result:
l1 >= l2 : false
l2 >= l1 : false
l3 >= l1 : false
l1 >= l3 : true
l1 >= l4 : true
All test objects should not change after the call
*/
func TestTlGreaterEqual(t *testing.T) {
	l1 := getLatticeMock1()
	l2 := getLatticeMock2()
	l3 := getLatticeMock3()
	l4 := getLatticeMock4()

	wanted := []struct {
		leftL     Lattice
		rightL    Lattice
		greatereq bool
		origL     Lattice
	}{
		{l1, l2, false, getLatticeMock1()},
		{l2, l1, false, getLatticeMock2()},
		{l3, l1, false, getLatticeMock3()},
		{l1, l3, true, getLatticeMock1()},
		{l1, l4, true, getLatticeMock1()},
	}

	for _, want := range wanted {
		greaterEq, err := want.leftL.GreaterEqual(want.rightL)
		assert.Nil(t, err)
		assert.Equal(t, want.greatereq, greaterEq)
		checkNoChange(t, want.leftL, want.origL)
	}
}

/*
Test object 1:  L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
Result: deepcopy=l1.DeepCopy() == l1
Then add an element to object 1
Afterwards: deepcopy != l1
*/
func TestDeepCopy(t *testing.T) {
	l1 := getLatticeMock1()
	dc := l1.DeepCopy()

	if assert.IsType(t, (Lattice)(nil), dc) {
		dcTaint := dc.(Lattice)
		equal, err := dcTaint.Equal(l1)
		assert.Nil(t, err)
		assert.True(t, equal)

		// L1{t1 -> Empty, t2 -> Tainted, t3 -> Untainted,t4 -> Untainted, t5 -> Both}
		m4 := new(ssaValMock)
		m4.N = "t4"
		l1[m4] = Untainted

		equal, err = l1.Equal(dcTaint)
		assert.Nil(t, err)
		assert.False(t, equal)
	}
}

/*
Test object 1:  L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
For all elements in L1 check whether the lattice value is as expected
*/
func TestGetVal(t *testing.T) {
	wanted := []struct {
		ssaVal   string
		taintVal Value
	}{
		{"t1", Uninitialized},
		{"t2", Tainted},
		{"t3", Untainted},
		{"t4", Uninitialized},
		{"t5", Both},
	}

	l1 := getLatticeMock1()
RangeLattice:
	for val := range l1 {
		for _, want := range wanted {
			if val.Name() == want.ssaVal {
				getVal := l1.GetVal(val)
				assert.Equal(t, want.taintVal, getVal)
				continue RangeLattice
			}
		}
	}
}

// checkNoChangeL1 checks whether l is equal to getLatticeMock1()
// L1: {t1 -> Empty, t2 -> Tainted, t3 -> Untainted, t5 -> Both}
func checkNoChangeL1(t *testing.T, l Lattice) {
	lTest := getLatticeMock1()
	checkNoChange(t, l, lTest)
}

func checkNoChangeL2(t *testing.T, l Lattice) {
	lTest := getLatticeMock2()
	checkNoChange(t, l, lTest)
}

func checkNoChange(t *testing.T, l Lattice, lTest Lattice) {
	eq, err := lTest.Equal(l)
	assert.Nil(t, err)
	assert.True(t, eq)
}
