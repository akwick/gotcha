// A package for a easier handling of creating the SSA representation of a program
package ssabuilder

import (
	//"bytes"

	"strings"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	//"log"
)

//private struct for handling the variables
//Initializing must be done with the help of the public NewSsaHelper function
type ssaHelper struct {
	prog    *ssa.Program // Aggregation
	mainPkg *ssa.Package // Aggregation
	//	fileName string       //Aggregation
}

func NewSsaHelper(path string, sourcefiles []string) (*ssaHelper, error) {
	var conf loader.Config
	sourcefilesstring := strings.Join(sourcefiles, ", ")
	conf.CreateFromFilenames(path, sourcefilesstring)
	//	conf.CreateFromFilenames("github.com/rev
	//el/cmd/revel", "../../../go/src/github.com/revel/cmd/revel/new.go", "../../../go/src/github.com/revel/cmd/revel/rev.go", "../../../go/src/github.com/revel/cmd/revel/util.go", "../../../go/src/github.com/revel/cmd/revel/run.go", "../../../go/src/github.com/revel/cmd/revel/build.go", "../../../go/src/github.com/revel/cmd/revel/package.go", "../../../go/src/github.com/revel/cmd/revel/clean.go", "../../../go/src/github.com/revel/cmd/revel/test.go")
	iprog, err := conf.Load()
	if err != nil {
		return nil, err
	}
	// file, err := conf.ParseFile(fileName, nil)
	// if err != nil {
	// 	fmt.Errorf("%v", err)
	// 	return nil, err
	// }
	//
	// conf.CreateFromFiles("main", file)
	// iprog, err := conf.Load()
	// if err != nil {
	// 	fmt.Errorf("%v", err)
	// 	return nil, err
	// }

	prog := ssautil.CreateProgram(iprog, ssa.SanityCheckFunctions)
	mainPkg := prog.Package(iprog.Created[0].Pkg)
	//	prog.BuildAll()
	prog.Build()

	retValue := ssaHelper{prog, mainPkg}
	//retValue := ssaHelper{prog, mainPkg, fileName}
	return &retValue, nil
}

func (ssah *ssaHelper) GetProgram() *ssa.Program {
	return ssah.prog
}

func (ssah *ssaHelper) GetMainPackage() *ssa.Package {
	return ssah.mainPkg
}

//func (ssa *SsaForCallgraph) GetProgram()iprog
//// Returns a *ssa.Programm.
//// The file is parsed against nil.
//func GetProgram(fileName string) (*ssa.Program, error) {
//	var buf bytes.Buffer
//	logger := log.New(&buf, "logger: ", log.Lshortfile)
//
//	var conf loader.Config
//	file, err := conf.ParseFile(fileName, nil)
//	if err != nil {
//		logger.Print(err)
//		fmt.Print(&buf)
//		return nil, err
//	}
//
//	conf.CreateFromFiles("main", file)
//
//	iprog, err := conf.Load()
//	if err != nil {
//		logger.Print(err)
//		fmt.Print(&buf)
//		return nil, err
//	}
//
//	prog := ssautil.CreateProgram(iprog, ssa.SanityCheckFunctions)
//
//	fmt.Print(&buf)
//
//	return prog, nil
//}
//
//func GetConfig(fileName string) (*loader.Program, error) {
//	var conf loader.Config
//	file, err := conf.ParseFile(fileName, nil)
//	if err != nil {
//		fmt.Errorf("%v", err)
//		return nil, err
//	}
//
//	conf.CreateFromFiles("main", file)
//	retProgram, err := conf.Load()
//
//	if err != nil {
//		fmt.Errorf("%v%", err)
//		return nil, err
//	}
//
//	return retProgram, nil
//}
