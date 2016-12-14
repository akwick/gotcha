package transferFunction

import (
	"goretech/analysis/lattice"

	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

//Semanticer is an interface for a transfer function
//A transfer function describes the change in a Lattice caused by an expression.
//In our implementation, an expression is represented by a ssa.Instruction.
type Semanticer interface {
	// TransferFunction returns a PlainFF which describes the change of a lattice.Valuer caused by node
	TransferFunction(node ssa.Instruction, pointers *pointer.Result) PlainFF
}

// PlainFF describes a plain flow function without any connection to an instruction.
type PlainFF func(lattice.Valuer) (lattice.Valuer, error)
