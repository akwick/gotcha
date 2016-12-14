package taint

import (
	"github.com/akwick/gotcha/lattice"
	"reflect"
	"testing"
)

func TestTlvImplLv(t *testing.T) {
	var tl lattice.Valuer
	tl = new(Value)
	// Checking with reflection
	tlType := reflect.TypeOf(tl)
	interfaceType := reflect.TypeOf((*lattice.Valuer)(nil)).Elem()
	impl := tlType.Implements(interfaceType)
	if !impl {
		t.Errorf("%v implements not Valuer", tl)
	}
}

func TestTlImplL(t *testing.T) {
	var l lattice.Latticer
	l = new(Lattice)
	lType := reflect.TypeOf(l)
	interfaceType := reflect.TypeOf((*lattice.Latticer)(nil)).Elem()
	impl := lType.Implements(interfaceType)
	if !impl {
		t.Errorf("%v implements not Latticer", t)
	}
}

func TestTlvBottomElem(t *testing.T) {
	v := new(Value)
	bottomElem := v.BottomElement()
	tlv := bottomElem.(Value)
	if tlv != 0 {
		t.Errorf("The value of Bottom Element should be 0")
	}
}

func TestTlvTopElem(t *testing.T) {
	v := new(Value)
	topElem := v.TopElement()
	tlv := topElem.(Value)
	if tlv != 3 {
		t.Errorf("The value of Top Element should be 3.")
	}
}

func TestTlvLeastUpperBound(t *testing.T) {
	// Value 0
	v := Value(0)
	cases := []struct {
		in, want Value
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, c := range cases {
		got, _ := v.LeastUpperBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}

	// Value 1
	v = Value(1)
	cases = []struct {
		in, want Value
	}{
		{0, 1},
		{1, 1},
		{2, 3},
		{3, 3},
	}
	for _, c := range cases {
		got, _ := v.LeastUpperBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}

	// Value 2
	v = Value(2)
	cases = []struct {
		in, want Value
	}{
		{0, 2},
		{1, 3},
		{2, 2},
		{3, 3},
	}
	for _, c := range cases {
		got, _ := v.LeastUpperBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}

	// Value 3
	v = Value(3)
	cases = []struct {
		in, want Value
	}{
		{0, 3},
		{1, 3},
		{2, 3},
		{3, 3},
	}
	for _, c := range cases {
		got, _ := v.LeastUpperBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}
}

func TestTlvGratestLowerBound(t *testing.T) {
	// Value 0
	v := Value(0)
	cases := []struct {
		in, want Value
	}{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}
	for _, c := range cases {
		got, _ := v.GreatestLowerBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}
	// Value 1
	v = Value(1)
	cases = []struct {
		in, want Value
	}{
		{0, 0},
		{1, 1},
		{2, 0},
		{3, 1},
	}
	for _, c := range cases {
		got, _ := v.GreatestLowerBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}
	// Value 2
	v = Value(2)
	cases = []struct {
		in, want Value
	}{
		{0, 0},
		{1, 0},
		{2, 2},
		{3, 2},
	}
	for _, c := range cases {
		got, _ := v.GreatestLowerBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}
	// Value 3
	v = Value(3)
	cases = []struct {
		in, want Value
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, c := range cases {
		got, _ := v.GreatestLowerBound(c.in)
		if got != c.want {
			t.Errorf("LeastUpperbound(%s) with a %s-value is %s, want %s", c.in, v, got, c.want)
		}
	}
}

func TestTlvLess(t *testing.T) {
	var tv Value
	cases := []struct {
		tlv1, tlv2 Value
		wanted     bool
	}{
		{0, 0, false},
		{0, 1, true},
		{0, 2, true},
		{0, 3, true},
		{1, 0, false},
		{1, 1, false},
		{1, 2, false},
		{1, 3, true},
		{2, 0, false},
		{2, 1, false},
		{2, 2, false},
		{2, 3, true},
		{3, 0, false},
		{3, 1, false},
		{3, 2, false},
		{3, 3, false},
	}
	for _, c := range cases {
		tv = Value(c.tlv1)
		wanted := c.wanted
		got, err := tv.Less(c.tlv2)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got != wanted {
			t.Errorf("%s.Less(%s) returns %t, want %t", tv, c.tlv2, got, wanted)
		}
	}
}
func TestTlvEqual(t *testing.T) {
	var tv Value
	cases := []struct {
		tlv1, tlv2 Value
		wanted     bool
	}{
		{0, 0, true},
		{0, 1, false},
		{0, 2, false},
		{0, 3, false},
		{1, 0, false},
		{1, 1, true},
		{1, 2, false},
		{1, 3, false},
		{1, 4, false},
		{2, 0, false},
		{2, 1, false},
		{2, 2, true},
		{2, 3, false},
		{3, 0, false},
		{3, 1, false},
		{3, 2, false},
		{3, 3, true},
	}
	for _, c := range cases {
		tv = Value(c.tlv1)
		got, err := tv.Equal(c.tlv2)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got != c.wanted {
			t.Errorf("%s.Equal(%s) returns %t, want %t", tv, c.tlv2, got, c.wanted)
		}
	}
}

func TestTlvLessEqual(t *testing.T) {
	var tv Value
	cases := []struct {
		tlv1, tlv2 Value
		wanted     bool
	}{
		{0, 0, true},
		{0, 1, true},
		{0, 2, true},
		{0, 3, true},
		{1, 0, false},
		{1, 1, true},
		{1, 2, false},
		{1, 3, true},
		{2, 0, false},
		{2, 1, false},
		{2, 2, true},
		{2, 3, true},
		{3, 0, false},
		{3, 1, false},
		{3, 2, false},
		{3, 3, true},
	}
	for _, c := range cases {
		tv = Value(c.tlv1)
		got, err := tv.LessEqual(c.tlv2)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got != c.wanted {
			t.Errorf("%s.LessEqual(%s) returns %t, want %t", tv, c.tlv2, got, c.wanted)
		}
	}
}
func TestTlvGreater(t *testing.T) {
	var tv Value
	cases := []struct {
		tlv1, tlv2 Value
		wanted     bool
	}{
		{0, 0, false},
		{0, 1, false},
		{0, 2, false},
		{0, 3, false},
		{1, 0, true},
		{1, 1, false},
		{1, 2, false},
		{1, 3, false},
		{2, 0, true},
		{2, 1, false},
		{2, 2, false},
		{2, 3, false},
		{3, 0, true},
		{3, 1, true},
		{3, 2, true},
		{3, 3, false},
	}
	for _, c := range cases {
		tv = Value(c.tlv1)
		got, err := tv.Greater(c.tlv2)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got != c.wanted {
			t.Errorf("%s.Greater(%s) returns %t, want %t", tv, c.tlv2, got, c.wanted)
		}
	}
}
func TestTlvGreaterEqual(t *testing.T) {
	var tv Value
	cases := []struct {
		tlv1, tlv2 Value
		wanted     bool
	}{
		{0, 0, true},
		{0, 1, false},
		{0, 2, false},
		{0, 3, false},
		{1, 0, true},
		{1, 1, true},
		{1, 2, false},
		{1, 3, false},
		{2, 0, true},
		{2, 1, false},
		{2, 2, true},
		{2, 3, false},
		{3, 0, true},
		{3, 1, true},
		{3, 2, true},
		{3, 3, true},
	}
	for _, c := range cases {
		tv = Value(c.tlv1)
		got, err := tv.GreaterEqual(c.tlv2)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got != c.wanted {
			t.Errorf("%s.GreaterEqual(%s) returns %t, want %t", tv, c.tlv2, got, c.wanted)
		}
	}
}
