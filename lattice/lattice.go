// Package lattice holds the interfaces and concrete implementations for
// lattices which are used during the analysis.
package lattice

import (
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

// Latticer describes the lattice interface.
// A lattice contains ssa.Values which points to a Valuer
type Latticer interface {
	// LeastUpperBound returns the lub of two lattices
	LeastUpperBound(l2 Latticer) (Latticer, error)
	// GreatestLowerBound returns the glb of two lattices
	GreatestLowerBound(l2 Latticer) (Latticer, error)
	// LeastElement returns the least element (_|_)
	LeastElement() (Latticer, error)
	// Less returns true if l1 < l2
	Less(l2 Latticer) (bool, error)
	// Equal returns true if l1 == l2
	Equal(l2 Latticer) (bool, error)
	// LessEqual returns true if l1 <= l2
	LessEqual(l2 Latticer) (bool, error)
	// Greater returns true if l1 > l2
	Greater(l2 Latticer) (bool, error)
	// GreaterEqual returns true if l1 >= l2
	GreaterEqual(l2 Latticer) (bool, error)
	// BottomLattice returns the lowest lattice
	BottomLattice() Latticer
	// DeepCopy returns a deep copy of the lattice
	DeepCopy() Latticer
	// String returns a readable string of the lattice
	String() string
	// GetVal returns the abstract value of key of the lattice
	GetVal(key ssa.Value) Valuer
	// SetVal sets key to val in the lattice
	SetVal(key ssa.Value, val Valuer) error
}

// Pter describes the methods a pointer lattice needs
type Pter interface {
	// Embedd Latticer interface
	Latticer
	// GetPtr returns the pointer.Pointer for the key
	GetPtr(key ssa.Value) pointer.Pointer
	// SetPtr sets the ptr as the pointer for the key
	SetPtr(key ssa.Value, ptr pointer.Pointer)
	// GetLat returns an instance of Latticer
	GetLat() Latticer
	// GetPTrs returns a map of all variables and the corresponding pointer
	GetPtrs() map[ssa.Value]pointer.Pointer
}

// Valuer defines the methods of a value inside a Latticer
type Valuer interface {
	// BottomElement returns the bottom element of a lattice value
	BottomElement() Valuer
	// TopElement returns the top element of a lattice value
	TopElement() Valuer
	// LeastUpperBound returns the the lub of lv1 and lv2
	LeastUpperBound(lv2 Valuer) (Valuer, error)
	// GreatestLowerBound returns the the glb of lv1 and lv2
	GreatestLowerBound(lv2 Valuer) (Valuer, error)
	// Less returns true if lv1 < lv2
	Less(lv2 Valuer) (bool, error)
	// Equal returns true if lv1 == lv2
	Equal(lv2 Valuer) (bool, error)
	// LessEqual returns true if lv1 <= lv2
	LessEqual(lv2 Valuer) (bool, error)
	// Greater returns true if lv1 > lv2
	Greater(lv2 Valuer) (bool, error)
	// GreaterEqual returns true if lv1 >= lv2
	GreaterEqual(lv2 Valuer) (bool, error)
	// String returns a readable string of the value
	String() string
}
