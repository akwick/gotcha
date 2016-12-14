package worklist

//
// import (
// 	"go/token"
// 	"testing"
//
// 	"golang.org/x/tools/go/ssa"
// )
//
// func TestAddElem(t *testing.T) {
// 	wlm := NewWlList()
// 	vc := NewValueContext(nil)
// 	instr := new(ssaInstrMock)
// 	instr.N = "i1"
// 	ccs := NewContextCallSite(vc, instr)
// 	wlm.Add(ccs)
//
// 	if wlm.Len() != 1 {
// 		t.Errorf("The length of wlm should be 1 but is %d", wlm.Len())
// 	}
//
// 	wlm.Add(ccs)
// 	if wlm.Len() != 1 {
// 		t.Errorf("The length of wlm should be 1 but is %d", wlm.Len())
// 	}
// }
//
// type ssaInstrMock struct {
// 	N    string
// 	inst ssa.Instruction
// }
//
// func (m ssaInstrMock) String() string {
// 	return m.N
// }
// func (m ssaInstrMock) Parent() *ssa.Function {
// 	return nil
// }
// func (m ssaInstrMock) Block() *ssa.BasicBlock {
// 	return nil
// }
// func (m ssaInstrMock) Operands(rands []*ssa.Value) []*ssa.Value {
// 	return rands
// }
// func (m ssaInstrMock) Pos() token.Pos {
// 	return 0
// }
