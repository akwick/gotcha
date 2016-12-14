package ssabuilder

import (
	"go/token"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Build returns the main package for given sourcefiles.
func Build(path string, sourcefiles []string) (*ssa.Package, error) {
	var conf loader.Config
	srcfs := strings.Join(sourcefiles, ", ")
	conf.CreateFromFilenames(path, srcfs)

	lprog, err := conf.Load()
	if err != nil {
		return nil, errors.Errorf("fail to load config of path: %s and sourcefiles: %s", path, srcfs)
	}

	prog := ssautil.CreateProgram(lprog, ssa.SanityCheckFunctions)
	mainPkg := prog.Package(lprog.Created[0].Pkg)
	prog.Build()

	return mainPkg, nil
}

// ReplaceSend replaces the ssa.Send with the Send implemented in this package
func ReplaceSend(pkgs []*ssa.Package) {
	chToFuncs := findChannels(pkgs)
	for _, pkg := range pkgs {
		for name, memb := range pkg.Members {
			if memb.Token() == token.FUNC {
				f := pkg.Func(name)
				for _, b := range f.Blocks {
					for n, i := range b.Instrs {
						val, ok := i.(*ssa.Send)
						if ok {
							replace := &Send{&send{val, chToFuncs[val.Chan]}}
							b.Instrs[n] = replace
						}
					}
				}
			}
		}
	}
}

// fincChannels finds for all channels the corresponding call instructions
func findChannels(mains []*ssa.Package) map[ssa.Value][]ssa.CallInstruction {
	var callCom *ssa.CallCommon
	chfuncs := make(map[ssa.Value][]ssa.CallInstruction, 0)
	for _, pkg := range mains {
		for name, memb := range pkg.Members {
			if memb.Token() == token.FUNC {
				f := pkg.Func(name)
				for _, b := range f.Blocks {
					for _, i := range b.Instrs {
						callCom = nil
						switch it := i.(type) {
						case *ssa.Go:
							callCom = it.Common()
						case *ssa.Defer:
							callCom = it.Common()
						case *ssa.Call:
							callCom = it.Common()
						}
						if callCom != nil {
						args:
							for _, v := range callCom.Args {
								mc, ok := v.(*ssa.MakeChan)
								i, _ := i.(ssa.CallInstruction)
								if ok {
									calls := chfuncs[mc]
									if calls == nil {
										calls = make([]ssa.CallInstruction, 0)
									}
									chfuncs[mc] = append(calls, i)
									continue args
								}
								// TODO: find better solution
								underly := v.Type().Underlying()
								isChan := strings.Contains(underly.String(), "chan")
								if isChan {
									calls := chfuncs[v]
									if calls == nil {
										calls = make([]ssa.CallInstruction, 0)
									}
									chfuncs[v] = append(calls, i)
								}
							}
						}
					}
				}
			}
		}
	}
	return chfuncs
}
