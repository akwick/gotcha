// Contains helper method for initializing the algorithm
package worklist

import (
	"fmt"
	"github.com/akwick/gotcha/lattice/taint"
	"github.com/akwick/gotcha/ssabuilder"
	"go/token"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

var isPointer bool
var allpkgs bool
var contextpkgs []*ssa.Package

//Initializations
func initSSAandPTA(path string, sourcefiles []string, sourceAndSinkFile string, pkgs string) (*ssa.Function, error) {
	// First generating a ssa with source code to get the main function
	//log.Printf("sourcefile: %s\n", sourcefiles[0])
	mainpkg, err := ssabuilder.Build(path, sourcefiles)
	if err != nil {
		return nil, err
	}
	mainpkg.Build()
	mainFunc = mainpkg.Func("main")
	if mainFunc == nil {
		return nil, errors.New("no main() function found!")
	}
	vcs = &VCS{ctx: make(map[VcIdentifier]*ValueContext, 0)}
	worklist = NewWlList()
	ccsPool = make([]*ContextCallSite, 0)
	transitions = make([]*Transition, 0)
	// Initialize the Sources and Sinks Slices with the help of the sources and sinks file
	taint.Read(sourceAndSinkFile)
	log.Printf("sources: %v", taint.Sources)

	// Add only the packages which are defined by the arguemnts
	// Flag allpkgs analyze all possible packages. If the flag is not set, it could be that a certain amound of packages are defined in pkgs.
	if allpkgs {
		contextpkgs = mainFunc.Prog.AllPackages()
	} else {
		log.Printf("mainpkg %v", mainpkg)
		contextpkgs = []*ssa.Package{mainpkg}
		//log.Printf("Name of mainpkg: %s\n", contextpkgs[0].String())
		// If manuelly pkgs are passed with the pkgs argument: Add these packages.
		if pkgs != "" {
			for _, pkg := range strings.Split(pkgs, ",") {
				p := mainFunc.Prog.ImportedPackage(pkg)
				// p is nil if no ssa package for the string pkg is created
				// TODO "improve" handling (=/= ignoring)
				if p != nil {
					contextpkgs = append(contextpkgs, p)
				} else {
					fmt.Printf("Pkg [%s] is unknown in %s", pkg, mainFunc.String())
				}
			}
		}
	}
	log.Printf("Analyze: %d : packages(%v)\n", len(contextpkgs), contextpkgs)
	stat.Printf("#packages, %d,", len(contextpkgs))

	// Setup and analyze pointers
	// pointer analysis needs a package with a main function
	setupPTA([]*ssa.Package{mainpkg})
	pta, err = pointer.Analyze(conf)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	setPtrMap([]*ssa.Package{mainpkg})

	// Replace ssa.Send with Send
	ssabuilder.ReplaceSend(contextpkgs)

	return mainFunc, nil
}

func initContext(ssaFun *ssa.Function) {
	_, _ = GetValueContext(ssaFun, []ssa.Value{}, nil, false)
}

// An analysis state for the main function will be created.
// The context gets a value context for the main function. The entry and exit value is a empty lattice.
// The worklist consists of all instructions of the main function.
func initContextVC(ssaFun *ssa.Function, vc *ValueContext) {
	pkg := ssaFun.Package()
	analyze := false
	// check whether the pkg is defined within the packages which should analyzed
ctxtfor:
	for _, p := range contextpkgs {
		if p == pkg {
			analyze = true
			break ctxtfor
		}
	}
	// only add the blocks and instructions if the package should be analyzed.
	if analyze {
		ssaFun.WriteTo(logFile)
		for _, block := range ssaFun.Blocks {
			for _, instr := range block.Instrs {
				// build a new context call site for every instruction within the main value context
				c := NewContextCallSite(vc, instr)
				worklist.Add(c)
				//	ccsPool = append(ccsPool, c)
			}
		}
	}
}

// setupPTA creates the config type for a pointer analysis.
func setupPTA(mains []*ssa.Package) {
	conf = &pointer.Config{
		Mains:          mains,
		BuildCallGraph: false,
	}

	addQueries(mains)
}

// addQueries addbuild queries for all functions of the packages
func addQueries(mains []*ssa.Package) {
	for _, pkg := range mains {
		for name, memb := range pkg.Members {
			if memb.Token() == token.FUNC {
				f := pkg.Func(name)
				for _, b := range f.Blocks {
					for _, i := range b.Instrs {
						val, ok := i.(ssa.Value)
						if ok {
							ok, ptrv := taint.IsPointerVal(val)
							if ok {
								conf.AddQuery(ptrv)
							}
							ok, indptrv := taint.IsIndirectPtr(val)
							if ok {
								conf.AddIndirectQuery(indptrv)
							}
						}
					}
				}
				//log.Print(")\n")
			}
		}
	}
}

// setPtrMap sets once the map for all functions which are available for the analysis
func setPtrMap(pkgs []*ssa.Package) {
	valToPtr = make(map[ssa.Value]pointer.Pointer)
	// 1) For each package which is imported: Get all functions
	// 2) For each function: Get all basic blocks
	// 3) For each basic block: Get all instructions
	// 4) For each instruction: Add it to the map
	for _, pkg := range pkgs {
		//log.Printf("pkg.String: %s", pkg.String())
		pkgs := pkg.String()
		for _, p := range contextpkgs {
			if ok := strings.Contains(pkgs, p.String()); ok {
				membrs := pkg.Members
				for name, memb := range membrs {
					if memb.Token() == token.FUNC {
						if memb.Name() != "Clearenv" {
							f := pkg.Func(name)
							for _, b := range f.Blocks {
								for _, i := range b.Instrs {
									val, ok := i.(ssa.Value)
									if ok {
										ok = pointer.CanPoint(val.Type())
										if ok {
											ptr := pta.Queries[val]
											valToPtr[val] = ptr
										} else {
											// can not point -> nothing to do
										}
										if val.Type() != nil {
											//log.Printf("(%s) vall: %s\n", val.Parent().Name(), val.Name()+" = "+val.String())

											if _, ok := val.(*ssa.Range); !ok {
												if val.Type().Underlying() != nil {
													tp, ok := val.Type().Underlying().(*types.Pointer)
													if ok {
														ok = pointer.CanPoint(tp.Elem())
														if ok {
															ptr := pta.IndirectQueries[val]
															valToPtr[val] = ptr
														} else {
															// tp.Elem() is not a pointer -> nothing to do
														}
													}
												}
											}
										}
									} else {
										// not all instructions are values.
										// the pointer analysis needs a value as parameter.
									}
								}
							}
						} else {
							// ignore -> will crash
						}
					} else {
						// not interested in VAR, CONST and TYPE
					}
				}
			} else {
				// pkg should not be analyzed (is not part of contextpkgs)
			}
		}
	}
}

func wlInit(path string, sourcefiles []string, sourceAndSinkFile string, allpackages bool, pkgs string, ptranalysis bool) {
	//setLogger()
	var err error
	// Set pkg and pointer variable first (used by initSSAandPTA)
	allpkgs = allpackages
	isPointer = ptranalysis
	mainFunc, err = initSSAandPTA(path, sourcefiles, sourceAndSinkFile, pkgs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		log.Fatalf("Fatal error: %s\n", err.Error())
	}
	initContext(mainFunc)
	errFlows = newErrInFlows()
}
