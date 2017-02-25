package taint

import (
	"github.com/akwick/gotcha/lattice"
	"github.com/akwick/gotcha/transferFunction"
	"log"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"

	"go/types"

	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

func returnID(currentValue lattice.Valuer) (lattice.Valuer, error) {
	return currentValue, nil
}

func returnUntainted(currentValue lattice.Valuer) (lattice.Valuer, error) {
	return Untainted, nil
}

func returnTainted(currentValue lattice.Valuer) (lattice.Valuer, error) {
	return Tainted, nil
}

func returnLUP(lupVal lattice.Valuer) func (lattice.Valuer) (lattice.Valuer, error) {
	return func (currentValue lattice.Valuer)(lattice.Valuer, error) { return currentValue.LeastUpperBound(lupVal) }
}

// returnLUPTaint is similar to returnLUP.
// In contrast to returnLUP, returnLUPTaint can handle error messages.
// An information flow is packed into the ErrInFlow struct which implements the error interface.
func returnLUPTaint (lupVal lattice.Valuer) func (lattice.Valuer) (lattice.Valuer, error) {
     return func (currentValue lattice.Valuer) (lattice.Valuer, error) {
	tempVal, errCall := currentValue.LeastUpperBound(lupVal)
	if errCall != nil {
		return lupVal, errCall
	}
	return tempVal, err
     }
}

func returnError(currentValue lattice.Valuer) (lattice.Valuer, error) {
	return currentValue, err
}

// ErrInFlow is an error type which holds an information flow.
// It can be used to report flows from sources to sinks to the user.
type ErrLeak struct {
	Call ssa.CallCommon
	Args []ssa.Value
	Err  error
}

func (e ErrLeak) Error() (s string) {
	s = "The function with signature: "
	callCom := e.Call
	if callCom.Signature() != nil {
		s += callCom.Signature().String()
	}
	if callCom.StaticCallee() != nil {
		s += callCom.StaticCallee().Name()
	}
	s += " is reached by at minimum one tainted argument: \n"
	for i, arg := range e.Args {
		if i > 1 {
			s += " | "
		}
		s += arg.String()
		s += " - " + arg.Name() + " of type: " + arg.Type().String() + " "
		pos := arg.Pos()
		if arg.Parent() != nil && arg.Parent().Prog != nil && arg.Parent().Prog.Fset != nil {
			fileset := arg.Parent().Prog.Fset
			location := fileset.File(pos).Line(pos)
			filepath := fileset.File(pos).Name()
			s += " at: " + filepath + " near to :" + strconv.Itoa(location)
		}
	}
	s += "\n"
	return
}

// NewErrInFlow returns a error of type ErrInFlow.
func NewErrInFlow(c *ssa.CallCommon, a []ssa.Value, e error) error {
	return ErrLeak{Call: *c, Args: a, Err: e}
}

//var taint = false
var err error

// TransferFunction handels a normal transfer of an instruction.
// Returns nil if an error occurs.
func (l Lattice) TransferFunction(node ssa.Instruction, ptr *pointer.Result) transferFunction.PlainFF {
	// Handling the different possibilities for a ssa.Instruction
	//log.Printf("ptr is nil %t", ptr == nil)
	switch nType := node.(type) {
	// Handle all cases which returns only the id
	// *ssa.MakeClosure returns only the id, becuase it's a ~function~ call which creates a new context
	case *ssa.DebugRef, *ssa.Jump, *ssa.MakeClosure, *ssa.Panic, *ssa.Return:
		return returnID
		// Handle the cases which operates on one ssa.Value e.g. Type.X and requires a LUP
		// A allocation should set a value to untainted
		// *ssa.MakeInterface is not listed because there construct a new type based upon another type
	case *ssa.Alloc, *ssa.MakeChan, *ssa.MakeMap, *ssa.MakeSlice:
		return returnUntainted
	case *ssa.ChangeInterface, *ssa.ChangeType, *ssa.Convert, *ssa.Extract, *ssa.Field, *ssa.FieldAddr, *ssa.Index, *ssa.MakeInterface, *ssa.Next, *ssa.Range, *ssa.Send, *ssa.Slice, *ssa.TypeAssert, *ssa.UnOp:
		var valX ssa.Value
		switch xType := node.(type) {
		case *ssa.ChangeInterface:
			valX = xType.X
		case *ssa.ChangeType:
			valX = xType.X
		case *ssa.Convert:
			valX = xType.X
		case *ssa.Extract:
			valX = xType.Tuple
		case *ssa.Field:
			valX = xType.X
		case *ssa.FieldAddr:
			valX = xType.X
		case *ssa.Index:
			valX = xType.X
		case *ssa.MakeInterface:
			valX = xType.X
		case *ssa.MakeMap:
			valX = xType.Reserve
		case *ssa.Next:
			valX = xType.Iter
		case *ssa.Range:
			valX = xType.X
		case *ssa.Send:
			valX = xType.X
		case *ssa.Slice:
			valX = xType.X
		case *ssa.TypeAssert:
			valX = xType.X
		case *ssa.UnOp:
			valX = xType.X
			/*	if xType.Op != token.MUL {
					valX = xType.X
				} else {
					return ptrUnOp(xType, l, ptr)
				} */
		}
		return returnLUP(l.GetVal(valX))
		// Handle the cases which operates on two ssa.Values
	case *ssa.BinOp, *ssa.IndexAddr, *ssa.Lookup:
		var val1, val2 ssa.Value
		switch xType := node.(type) {
		case *ssa.BinOp:
			val1 = xType.X
			val2 = xType.Y
		case *ssa.IndexAddr:
			val1 = xType.X
			val2 = xType.Index
		case *ssa.Lookup:
			val1 = xType.X
			val2 = xType.Index
		}
		lVal1 := l.GetVal(val1)
		lVal2 := l.GetVal(val2)
		var retLatValuer lattice.Valuer
		if retLatValuer, err = lVal1.LeastUpperBound(lVal2); err != nil {
			return returnError
		}
		return returnLUP(retLatValuer.(Value))
		// Handle the cases which operates upon an slice
	case *ssa.Phi:
		valArr := nType.Edges
		var v lattice.Valuer
		v = Uninitialized
		for _, ssaVal := range valArr {
			val := l.GetVal(ssaVal)
			if v, err = v.LeastUpperBound(val); err != nil {
				return returnError
			}
		}
		return returnLUP(v.(Value))
		// Hande the case with three ssa.Values
	case *ssa.MapUpdate:
		// updates a values in Map[Key] to value
		valMap := nType.Map
		valKey := nType.Key
		valValue := nType.Value
		lValMap := l.GetVal(valMap)
		lValKey := l.GetVal(valKey)
		var lupMapKey lattice.Valuer
		lupMapKey, err = lValMap.LeastUpperBound(lValKey)
		if err != nil {
			return returnError
		}
		lValValue := l.GetVal(valValue)
		var lup3Val lattice.Valuer
		if lup3Val, err = lValValue.LeastUpperBound(lupMapKey); err != nil {
			return returnError
		}
		return returnLUP(lup3Val.(Value))
		// Handle calls
	case *ssa.Call, *ssa.Defer, *ssa.Go:
		ff := checkAndHandleSourcesAndsinks(node, l, false)
		if ff == nil {
			return returnID
		} else {
			return ff
		}
	case *ssa.If:
		// handling the ssa representation for an if statement
		return returnID
	case *ssa.RunDefers:
		// pops and invokes the defered calls
	case *ssa.Select:
		// testing whether one of the specified sent or received states is entered
		// TODO describe behaivour for concurrency
	case *ssa.Store:
		// *t1 = t0
		// Addr is of type FieldAddr
		//log.Printf("ssa.Store \n")
		//log.Printf("  Addr: %s | Value: %s\n", nType.Addr, nType.Val)
		value := nType.Val
		return returnLUP(l.GetVal(value))
	default:
		return returnID
	}
	return returnID
}

func lupRetValues(r *ssa.Return, exitLat lattice.Latticer) (lattice.Valuer, error) {
	vals := r.Results
	var lupValer lattice.Valuer
	for _, val := range vals {
		v := exitLat.GetVal(val)
		lupValer, err = lupValer.LeastUpperBound(v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed %s.LUP(%s)", lupValer.String(), v.String())
		}
	}
	return lupValer, nil
}

func getSignature(c ssa.CallCommon) (signature, staticCallee, iSignature string) {
	// [vs: can this signature be an interface here? If yes, do we miss "up-casted" sources?]
	// Get the signatuer and iterate through the sources
	signature = c.Signature().String()
	if c.StaticCallee() != nil {
		staticCallee = c.StaticCallee().String()
	}

	var sigSlice []string
	if types.IsInterface(c.Signature().Underlying()) {
		// Signature for an interface does not contain names for the parameters
		sigI := c.Signature().String()
		// splits a string into a string slice.
		// Each element in the slice consists only of an letter, a number, [, ], (,) or a *
		f := func(c rune) bool {
			r1, _ := utf8.DecodeRuneInString("[")
			r2, _ := utf8.DecodeRuneInString("]")
			r3, _ := utf8.DecodeRuneInString("*")
			r4, _ := utf8.DecodeRuneInString("(")
			r5, _ := utf8.DecodeRuneInString(")")
			return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != r1 && c != r2 && c != r3 && c != r4 && c != r5
		}
		sigSlice = strings.FieldsFunc(sigI, f)
		for i, s := range sigSlice {
			paramsReady := false
			retReady := true
			retStart := 1
			if i == 0 {
				if !strings.Contains(s, "func()") {
					// func(name -> func(
					withoutName := strings.SplitAfter(s, "(")
					sigSlice[i] = withoutName[0]
					paramsReady = true
				}
			}

			if !paramsReady {
				if i%2 != 0 {
					paramsReady = strings.Contains(s, ")")
					if paramsReady {
						retReady = false
						retStart = i + 1
					}
				} else {
					sigSlice[i] = ""
				}
			}

			if !retReady {
				if !strings.Contains(s, ")") {
					if (i-retStart)%2 != 0 {
						sigSlice[i] = ""
					}
				}
			}
		}
	}
	for _, s := range sigSlice {
		iSignature += s
	}
	return
}

func handleSinkDetection(c ssa.CallCommon, l lattice.Latticer, ptr bool) {
	//log.Printf("HandleSinkDetection: for: %s with lattice: %s \n", c.String(), l.String())
	val := c.Value
	argsErr := []ssa.Value{}
	if l.GetVal(val) == Tainted || l.GetVal(val) == Both {
		argsErr = append(argsErr, val)
	}
	// [vs] Sometimes Args contains Value! Be careful, at least comment?
	// See: https://godoc.org/golang.org/x/tools/go/ssa#Call.Value
	args := c.Args

	for _, arg := range args {
		if l.GetVal(arg) == Tainted || l.GetVal(arg) == Both {
			argsErr = append(argsErr, arg)
		}
	}
	if ptr {
		lptr, ok := l.(*LatticePointer)
		if ok {
			for _, arg := range args {
				valsPtsTo := lptr.GetSSAValMayAlias(arg)
				for _, v := range valsPtsTo {
					if v.Name() == "t2" {
						log.Printf("t2 aliases %s\n", arg.Name())
					}
					if lptr.GetVal(v) == Tainted || lptr.GetVal(v) == Both {
						argsErr = append(argsErr, arg)
					}
				}
			}
		}
	}
	// Handle the case that one parameter is a variable

	if len(argsErr) != 0 {
		err = NewErrInFlow(&c, argsErr, nil)
	}
}
