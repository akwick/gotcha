package worklist

import (
	"goretech/tools/ssaHelper"
	"strconv"
	"testing"

	"golang.org/x/tools/go/ssa"

	"github.com/stretchr/testify/assert"
)

type idomRes struct {
	block   int
	idomNil bool
	idom    []int
}

type node struct {
	block int
	instr int
}

type succnode struct {
	n node
	s []node
}

func TestGetSuccessors(t *testing.T) {
	t.Skip()
	n00 := node{block: 0, instr: 0}
	n01 := node{block: 0, instr: 1}
	n02 := node{block: 0, instr: 2}
	n03 := node{block: 0, instr: 3}
	n04 := node{block: 0, instr: 4}
	n10 := node{block: 1, instr: 0}
	n11 := node{block: 1, instr: 1}
	n12 := node{block: 1, instr: 2}
	n13 := node{block: 1, instr: 3}
	n20 := node{block: 2, instr: 0}
	n21 := node{block: 2, instr: 1}
	n30 := node{block: 3, instr: 0}
	n31 := node{block: 3, instr: 1}
	n32 := node{block: 3, instr: 2}
	n33 := node{block: 3, instr: 3}
	tests := []struct {
		path  string
		src   []string
		succs []succnode
	}{
		{
			path: "goretech/analysis",
			src:  []string{"../tests/examplbeCode/hello2.go"},
			succs: []succnode{
				succnode{
					n: n00,
					s: []node{n01},
				},
				succnode{
					n: n01,
					s: []node{n02},
				},
				succnode{
					n: n02,
					s: []node{n03},
				},
				succnode{
					n: n03,
					s: []node{n04},
				},
				succnode{
					n: n04,
					s: []node{n10, n30},
				},
				succnode{
					n: n10,
					s: []node{n11},
				},
				succnode{
					n: n11,
					s: []node{n12},
				},
				succnode{
					n: n12,
					s: []node{n13},
				},
				succnode{
					n: n13,
					s: []node{n20},
				},
				succnode{
					n: n20,
					s: []node{n21},
				},
				succnode{
					n: n21,
					s: []node{},
				},
				succnode{
					n: n30,
					s: []node{n31},
				},
				succnode{
					n: n31,
					s: []node{n32},
				},
				succnode{
					n: n32,
					s: []node{n33},
				},
				succnode{
					n: n33,
					s: []node{n20},
				},
			},
		},
	}
	for _ = range tests {
		/*		helper, err := ssaHelper.NewSsaHelper(test.path, test.src)
				if assert.Nil(t, err) {
					mainPkg := helper.GetMainPackage()
					mainPkg.Build()
					mainFunc := mainPkg.Func("main")
					blocks := mainFunc.Blocks

				}
		*/
	}
}

func TestIdoms(t *testing.T) {
	tests := []struct {
		path  string
		src   []string
		idoms []idomRes
	}{
		{
			path: "goretech/analysis",
			src:  []string{"../tests/exampleCode/hello2.go"},
			idoms: []idomRes{
				idomRes{
					block:   0,
					idomNil: true,
					idom:    []int{},
				}, idomRes{
					block:   1,
					idomNil: false,
					idom:    []int{0},
				}, idomRes{
					block:   2,
					idomNil: false,
					idom:    []int{0},
				}, idomRes{
					block:   3,
					idomNil: false,
					idom:    []int{0},
				},
			},
		}, {
			path: "goretech/analysis",
			src:  []string{"../tests/exampleCode/idom.go"},
			idoms: []idomRes{
				idomRes{
					block:   0,
					idomNil: true,
					idom:    []int{},
				}, idomRes{
					block:   1,
					idomNil: false,
					idom:    []int{0},
				}, idomRes{
					block:   2,
					idomNil: false,
					idom:    []int{0},
				}, idomRes{
					block:   3,
					idomNil: false,
					idom:    []int{0},
				}, idomRes{
					block:   4,
					idomNil: false,
					idom:    []int{2, 0},
				}, idomRes{
					block:   5,
					idomNil: false,
					idom:    []int{2, 0},
				}, idomRes{
					block:   6,
					idomNil: false,
					idom:    []int{2, 0},
				},
			},
		},
	}

	for _, test := range tests {
		helper, err := ssaHelper.NewSsaHelper(test.path, test.src)
		if assert.Nil(t, err) {
			mainPkg := helper.GetMainPackage()
			mainPkg.Build()
			mainFunc := mainPkg.Func("main")
			blocks := mainFunc.Blocks
			var idomb []*ssa.BasicBlock
			expI := test.idoms
			for i, b := range blocks {
				idomb = getIdomsBlocks(b.Idom())
				expect := expI[i]
				if expect.idomNil {
					assert.Nil(t, idomb)
				} else {
					if assert.NotNil(t, idomb, "Block "+strconv.Itoa(i)) {
						assert.Equal(t, len(expect.idom), len(idomb))
						for i := range idomb {
							assert.Equal(t, expect.idom[i], idomb[i].Index)
						}
					}
				}
			}
		}
	}
}
