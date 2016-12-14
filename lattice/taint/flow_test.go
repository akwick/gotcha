package taint

import (
	"go/types"

	"github.com/stretchr/testify/assert"

	"github.com/akwick/gotcha/lattice"
	"testing"

	"golang.org/x/tools/go/ssa"
)

//func TestTLImplFF(t *testing.T) {
//	var tff transferFunction.Semanticer
//	tff = new(Lattice)
//	tffType := reflect.TypeOf(tff)
//	interfaceType := reflect.TypeOf((*transferFunction.Semanticer)(nil)).Elem()
//	impl := tffType.Implements(interfaceType)
//	if !impl {
//		t.Errorf("%v doesn't implement transferFunctioner", tff)
//	}
//}
//

var retVal lattice.Valuer

// A flow function to a allocation should return always untainted
func TestAlloc(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	ssaAlloc := new(ssa.Alloc)
	l := make(Lattice)
	l[mock1] = Tainted
	expected := Untainted

	ff := l.TransferFunction(ssaAlloc, nil)
	retVal, err = ff(l[ssaAlloc])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != expected {
		t.Errorf("Returned value should: %s, but is %s", expected, retVal)
	}

	// Set the value to tainted to check that the LUP is not build
	l[ssaAlloc] = Tainted
	expected = Untainted
	ff = l.TransferFunction(ssaAlloc, nil)
	retVal, err = ff(l[ssaAlloc])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != expected {
		t.Errorf("Returned value should: %s, but is %s", expected, retVal)
	}
}

func TestBinOp(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"

	ssaBinOp := new(ssa.BinOp)
	ssaBinOp.X = mock1
	ssaBinOp.Y = mock2

	l := make(Lattice)
	l[mock1] = Tainted
	l[mock2] = Uninitialized

	ff := l.TransferFunction(ssaBinOp, nil)
	retVal, err = ff(l[ssaBinOp])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	l[ssaBinOp] = Untainted
	ff = l.TransferFunction(ssaBinOp, nil)
	retVal, err = ff(l[ssaBinOp])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestCall(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	mock3 := new(ssaValMock)
	mock3.N = "t3"

	packageForMethod := types.NewPackage("packagePathName", "packageName")
	varForSignature := types.NewVar(2, packageForMethod, "Name of Var", nil)
	tupleForSignature := types.NewTuple(varForSignature)
	tupleForSignature2 := types.NewTuple(varForSignature, varForSignature)
	signatureForMethod := types.NewSignature(varForSignature, tupleForSignature, tupleForSignature2, false)
	method := types.NewFunc(2, packageForMethod, "Name of method", signatureForMethod)

	ssaCallCommon := new(ssa.CallCommon)
	ssaCallCommon.Value = mock1
	ssaCallCommon.Method = method
	ssaCallCommon.Args = []ssa.Value{mock2, mock3}

	ssaCall := new(ssa.Call)
	ssaCall.Call = *ssaCallCommon

	t.Log(ssaCall.Call.Signature().String())
	t.Log(ssaCall.Call.StaticCallee())

	td := &Data{&taintData{sig: ssaCall.Call.Signature().String()}}
	Sources = append(Sources, td)

	l := make(Lattice)
	l[ssaCall] = Uninitialized

	ff := l.TransferFunction(ssaCall, nil)
	retVal, err = ff(l[ssaCall])
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(retVal)
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	l2 := make(Lattice)
	l2[ssaCall] = Untainted

	t.Log(l2[ssaCall])
	ff = l.TransferFunction(ssaCall, nil)
	retVal, err = ff(l2[ssaCall])
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%s | %s", retVal, l2[ssaCall])
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestCallTainted(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	mock3 := new(ssaValMock)
	mock3.N = "t3"

	packageForMethod := types.NewPackage("packagePathName", "packageName")
	varForSignature := types.NewVar(2, packageForMethod, "Name of Var", nil)
	tupleForSignature := types.NewTuple(varForSignature)
	tupleForSignature2 := types.NewTuple(varForSignature, varForSignature)
	signatureForMethod := types.NewSignature(varForSignature, tupleForSignature, tupleForSignature2, false)
	method := types.NewFunc(2, packageForMethod, "Name of method", signatureForMethod)

	ssaCallCommon := new(ssa.CallCommon)
	ssaCallCommon.Value = mock1
	ssaCallCommon.Method = method
	ssaCallCommon.Args = []ssa.Value{mock2, mock3}

	ssaCall := new(ssa.Call)
	ssaCall.Call = *ssaCallCommon

	td := &Data{&taintData{sig: ssaCall.Call.Signature().String()}}
	Sinks = append(Sinks, td)

	l := make(Lattice)
	l[ssaCall] = Uninitialized
	l[mock2] = Tainted

	ff := l.TransferFunction(ssaCall, nil)
	retVal, err = ff(l[ssaCall])
	switch err := err.(type) {
	case ErrLeak:
		break
	case nil:
		break
	default:
		t.Error(err.Error())
	}
	//	if err != nil {
	//		t.Error(err.Error())
	//	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	l2 := make(Lattice)
	l2[ssaCall] = Untainted
	l2[mock2] = Tainted

	ff = l2.TransferFunction(ssaCall, nil)
	retVal, err = ff(l2[ssaCall])
	switch err := err.(type) {
	case ErrLeak, nil:
		break
	default:
		t.Error(err.Error())
	}

	t.Logf("lVal: %v \n", l2[ssaCall.Call.Value])
	t.Logf("ssaCall: %v \n", l2[ssaCall])
	for _, arg := range ssaCall.Call.Args {
		t.Logf("lArg: %v \n", l[arg])
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
	// According to delv the value of Taint is true
	// TODO find a solution for testing
	//if !taint {
	//	t.Errorf("taint should be true, but is %t", ttaint)
	//}
}

func TestChangeInterface(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"

	ssaChangeInterface := new(ssa.ChangeInterface)
	ssaChangeInterface.X = mock1

	l := make(Lattice)
	l[mock1] = Uninitialized

	ff := l.TransferFunction(ssaChangeInterface, nil)
	retVal, err = ff(l[ssaChangeInterface])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Uninitialized {
		t.Errorf("Returned value should: %s, but is %s", Uninitialized, retVal)
	}

	l[ssaChangeInterface] = Tainted
	ff = l.TransferFunction(ssaChangeInterface, nil)
	retVal, err = ff(l[ssaChangeInterface])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	l[mock1] = Untainted
	ff = l.TransferFunction(ssaChangeInterface, nil)
	retVal, err = ff(l[ssaChangeInterface])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestChangeType(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t4"
	l := make(Lattice)
	l[mock1] = Both
	ssaChangeType := new(ssa.ChangeType)
	ssaChangeType.X = mock1
	ff := l.TransferFunction(ssaChangeType, nil)
	retVal, err = ff(l[ssaChangeType])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
	t.Logf("retVal: %s\n", retVal)
	t.Logf("Lattice: %v\n", retVal)

	l[mock1] = Untainted
	ff = l.TransferFunction(ssaChangeType, nil)
	retVal, err = ff(l[ssaChangeType])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}

	l[ssaChangeType] = Tainted
	ff = l.TransferFunction(ssaChangeType, nil)
	retVal, err = ff(l[ssaChangeType])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestConvert(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t3"
	mock2 := new(ssaValMock)
	mock2.N = "t6"
	l := make(Lattice)
	l[mock1] = Untainted
	l[mock2] = Tainted
	ssaConvert := new(ssa.Convert)
	ssaConvert.X = mock1
	ff := l.TransferFunction(ssaConvert, nil)
	retVal, err = ff(l[ssaConvert])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Untainted, retVal)
	}
	ssaConvert.X = mock2
	ff = l.TransferFunction(ssaConvert, nil)
	retVal, err = ff(l[ssaConvert])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestDefer(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	mock3 := new(ssaValMock)
	mock3.N = "t3"

	packageForMethod := types.NewPackage("packagePathName", "packageName")
	varForSignature := types.NewVar(2, packageForMethod, "Name of Var", nil)
	tupleForSignature := types.NewTuple(varForSignature)
	tupleForSignature2 := types.NewTuple(varForSignature, varForSignature)
	signatureForMethod := types.NewSignature(varForSignature, tupleForSignature, tupleForSignature2, false)
	method := types.NewFunc(2, packageForMethod, "Name of method", signatureForMethod)

	ssaCallCommon := new(ssa.CallCommon)
	ssaCallCommon.Value = mock1
	ssaCallCommon.Method = method
	ssaCallCommon.Args = []ssa.Value{mock2, mock3}

	ssaDefer := new(ssa.Defer)
	ssaDefer.Call = *ssaCallCommon

	td := &Data{&taintData{sig: ssaDefer.Call.Signature().String()}}
	Sources = append(Sources, td)

	l := make(Lattice)
	l[ssaDefer.Call.Value] = Uninitialized

	ff := l.TransferFunction(ssaDefer, nil)
	retVal, err = ff(l[ssaDefer.Call.Value])
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(retVal)
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	l2 := make(Lattice)
	l2[ssaDefer.Call.Value] = Untainted

	t.Log(l2[ssaDefer.Call.Value])
	ff = l.TransferFunction(ssaDefer, nil)
	retVal, err = ff(l2[ssaDefer.Call.Value])
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%s | %s", retVal, l2[ssaDefer.Call.Value])
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestExtract(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t9"
	l := make(Lattice)
	l[mock1] = Tainted
	ssaExtract := new(ssa.Extract)
	ssaExtract.Tuple = mock1
	ssaExtract.Index = 1
	ff := l.TransferFunction(ssaExtract, nil)
	retVal, err = ff(l[ssaExtract])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
	l[ssaExtract] = Untainted
	ff = l.TransferFunction(ssaExtract, nil)
	retVal, err = ff(l[ssaExtract])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestField(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t9"
	l := make(Lattice)
	l[mock1] = Both
	ssaField := new(ssa.Field)
	ssaField.X = mock1
	ff := l.TransferFunction(ssaField, nil)
	retVal, err = ff(l[ssaField])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestFieldAddr(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t3"
	mock2 := new(ssaValMock)
	mock2.N = "t6"
	l := make(Lattice)
	l[mock1] = Uninitialized
	l[mock2] = Tainted
	ssaFieldAddr := new(ssa.FieldAddr)
	ssaFieldAddr.X = mock1
	ff := l.TransferFunction(ssaFieldAddr, nil)
	retVal, err = ff(l[ssaFieldAddr])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Uninitialized {
		t.Errorf("Returned value should: %s, but is %s", Uninitialized, retVal)
	}
	ssaFieldAddr.X = mock2
	ff = l.TransferFunction(ssaFieldAddr, nil)
	retVal, err = ff(l[ssaFieldAddr])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestGo(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t3"
	mock2 := new(ssaValMock)
	mock2.N = "t1"
	mock3 := new(ssaValMock)
	mock3.N = "t5"

	packageForMethod := types.NewPackage("package Path Name", "package Name")
	varForSignature := types.NewVar(2, packageForMethod, "Name of Var", nil)
	tupleForSignature := types.NewTuple(varForSignature)
	tupleForSignature2 := types.NewTuple(varForSignature, varForSignature)
	signatureForMethod := types.NewSignature(varForSignature, tupleForSignature2, tupleForSignature, false)
	method := types.NewFunc(2, packageForMethod, "Name of method", signatureForMethod)

	ssaCallCommon := new(ssa.CallCommon)
	ssaCallCommon.Value = mock1
	ssaCallCommon.Method = method
	ssaCallCommon.Args = []ssa.Value{mock2, mock3}

	ssaGo := new(ssa.Go)
	ssaGo.Call = *ssaCallCommon

	td := &Data{&taintData{sig: ssaGo.Call.Signature().String()}}
	Sources = append(Sources, td)

	l := make(Lattice)
	l[ssaGo.Call.Value] = Uninitialized

	ff := l.TransferFunction(ssaGo, nil)
	retVal, err = ff(l[ssaGo.Call.Value])
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(retVal)
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

	l2 := make(Lattice)
	l2[ssaGo.Call.Value] = Both

	t.Log(l2[ssaGo.Call.Value])
	ff = l.TransferFunction(ssaGo, nil)
	retVal, err = ff(l2[ssaGo.Call.Value])
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%s | %s", retVal, l2[ssaGo.Call.Value])
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestIndex(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t9"
	l := make(Lattice)
	l[mock1] = Uninitialized
	ssaExtract := new(ssa.Extract)
	ssaExtract.Tuple = mock1
	ssaExtract.Index = 1
	ff := l.TransferFunction(ssaExtract, nil)
	retVal, err = ff(l[ssaExtract])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Uninitialized {
		t.Errorf("Returned value should: %s, but is %s", Uninitialized, retVal)
	}
	l[mock1] = Tainted
	ff = l.TransferFunction(ssaExtract, nil)
	retVal, err = ff(l[ssaExtract])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}

}

func TestIndexAddr(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t3"
	mock2 := new(ssaValMock)
	mock2.N = "t6"
	l := make(Lattice)
	l[mock1] = Tainted
	l[mock2] = Uninitialized
	ssaIndexAdrr := new(ssa.IndexAddr)
	ssaIndexAdrr.X = mock1
	ssaIndexAdrr.Index = mock2
	ff := l.TransferFunction(ssaIndexAdrr, nil)
	retVal, err = ff(l[ssaIndexAdrr])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestLookup(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t3"
	l := make(Lattice)
	l[mock1] = Uninitialized
	l[mock2] = Untainted
	ssaLookup := new(ssa.Lookup)
	ssaLookup.X = mock1
	ssaLookup.Index = mock2
	ff := l.TransferFunction(ssaLookup, nil)
	retVal, err = ff(l[ssaLookup])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Untainted, retVal)
	}
	l[ssaLookup] = Tainted
	ff = l.TransferFunction(ssaLookup, nil)
	retVal, err = ff(l[ssaLookup])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestMakeClosure(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	mock3 := new(ssaValMock)
	mock3.N = "t4"
	mock4 := new(ssaValMock)
	mock4.N = "t6"
	l := make(Lattice)
	l[mock1] = Uninitialized
	l[mock2] = Tainted
	l[mock3] = Uninitialized
	l[mock4] = Untainted
	ssaMakeClosure := new(ssa.MakeClosure)
	valArr := make([]ssa.Value, 4)
	valArr = append(valArr, mock1)
	valArr = append(valArr, mock2)
	valArr = append(valArr, mock3)
	valArr = append(valArr, mock4)
	ssaMakeClosure.Bindings = valArr
	ff := l.TransferFunction(ssaMakeClosure, nil)
	retVal, err = ff(l[ssaMakeClosure])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Uninitialized {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestMakeInterface(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t7"
	l := make(Lattice)
	l[mock1] = Both
	ssaMakeInterface := new(ssa.MakeInterface)
	ssaMakeInterface.X = mock1
	l[ssaMakeInterface] = Untainted
	ff := l.TransferFunction(ssaMakeInterface, nil)
	retVal, err = ff(l[ssaMakeInterface])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestMakeMap(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t7"
	l := make(Lattice)
	l[mock1] = Tainted
	ssaMakeMap := new(ssa.MakeMap)
	ssaMakeMap.Reserve = mock1
	ff := l.TransferFunction(ssaMakeMap, nil)
	retVal, err = ff(l[ssaMakeMap])
	assert.Nil(t, err)
	assert.Equal(t, Untainted, retVal)
}

func TestMapUpdate(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	mock3 := new(ssaValMock)
	mock3.N = "t3"
	l := make(Lattice)
	l[mock1] = Uninitialized
	l[mock2] = Tainted
	l[mock3] = Untainted
	ssaMapUpdate := new(ssa.MapUpdate)
	ssaMapUpdate.Key = mock1
	ssaMapUpdate.Value = mock2
	ssaMapUpdate.Map = mock3
	ff := l.TransferFunction(ssaMapUpdate, nil)
	retVal, err = ff(l[ssaMapUpdate.Value])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
	ff = l.TransferFunction(ssaMapUpdate, nil)
	retVal, err = ff(l[ssaMapUpdate.Map])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
	ff = l.TransferFunction(ssaMapUpdate, nil)
	retVal, err = ff(l[ssaMapUpdate.Key])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestNext(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	l := make(Lattice)
	l[mock1] = Both
	ssaNext := new(ssa.Next)
	ssaNext.Iter = mock1
	ssaNext.IsString = false
	ff := l.TransferFunction(ssaNext, nil)
	retVal, err = ff(l[ssaNext])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestPanic(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t3"
	l := make(Lattice)
	l[mock1] = Untainted
	ssaPanic := new(ssa.Panic)
	ssaPanic.X = mock1
	t.Logf("ssaPanic: %v", ssaPanic)
	ff := l.TransferFunction(ssaPanic, nil)
	retVal, err = ff(l[ssaPanic.X])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Untainted, retVal)
	}
}

func TestPhi(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	mock3 := new(ssaValMock)
	mock3.N = "t4"
	mock4 := new(ssaValMock)
	mock4.N = "t6"
	l := make(Lattice)
	l[mock1] = Tainted
	l[mock2] = Uninitialized
	l[mock3] = Tainted
	l[mock4] = Untainted
	ssaPhi := new(ssa.Phi)
	var valArr []ssa.Value
	valArr = append(valArr, mock1)
	valArr = append(valArr, mock2)
	valArr = append(valArr, mock3)
	valArr = append(valArr, mock4)
	ssaPhi.Edges = valArr
	var v lattice.Valuer
	v = Uninitialized
	t.Logf("Length of ssaPhi.Edges: %d", len(ssaPhi.Edges))
	for _, ssaVal := range valArr {
		val := l[ssaVal]
		v, _ = v.LeastUpperBound(val)
		t.Logf("val: %s | v: %s", val, v)
	}
	ff := l.TransferFunction(ssaPhi, nil)
	retVal, err = ff(l[ssaPhi])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestRange(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t7"
	l := make(Lattice)
	l[mock1] = Untainted
	ssaRange := new(ssa.Range)
	ssaRange.X = mock1
	ff := l.TransferFunction(ssaRange, nil)
	retVal, err = ff(l[ssaRange])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Untainted {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}
}

func TestSlice(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	mock2 := new(ssaValMock)
	mock2.N = "t3"
	l := make(Lattice)
	ssaSlice := new(ssa.Slice)
	ssaSlice.X = mock1
	l[mock1] = Tainted
	l[mock2] = Untainted
	ff := l.TransferFunction(ssaSlice, nil)
	retVal, err = ff(l[ssaSlice])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}

func TestSTore(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t!"
	mock2 := new(ssaValMock)
	mock2.N = "t2"
	l := make(Lattice)
	ssaStore := new(ssa.Store)
	ssaStore.Addr = mock1
	ssaStore.Val = mock2
	l[mock1] = Tainted
	l[mock2] = Untainted
	ff := l.TransferFunction(ssaStore, nil)
	retVal, err = ff(l[ssaStore.Val])
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := Untainted
	if retVal != expected {
		t.Errorf("Returned value should: %s, but is %s", expected, retVal)
	}
}

func TestTypeAssert(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t1"
	l := make(Lattice)
	ssaTypeAssert := new(ssa.TypeAssert)
	ssaTypeAssert.X = mock1
	l[mock1] = Both
	ff := l.TransferFunction(ssaTypeAssert, nil)
	retVal, err = ff(l[ssaTypeAssert])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Both {
		t.Errorf("Returned value should: %s, but is %s", Both, retVal)
	}

}

func TestUnOp(t *testing.T) {
	mock1 := new(ssaValMock)
	mock1.N = "t7"
	l := make(Lattice)
	ssaUnop := new(ssa.UnOp)
	ssaUnop.X = mock1
	l[mock1] = Tainted
	ff := l.TransferFunction(ssaUnop, nil)
	retVal, err = ff(l[ssaUnop])
	if err != nil {
		t.Errorf(err.Error())
	}
	if retVal != Tainted {
		t.Errorf("Returned value should: %s, but is %s", Tainted, retVal)
	}
}
