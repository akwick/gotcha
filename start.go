package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/akwick/gotcha/worklist"
)

func main() {
	// Explains how to use the command line flags
	flag.Usage = func() {
		fmt.Printf("Usage and defaults of %s: \n", os.Args[0])
		fmt.Printf("The flags allpkgs, path and ssf are optional. \n")
		fmt.Printf("The flag sourceFilesFlag is mandatory.\n")
		flag.PrintDefaults()
	}

	var ssf = flag.String("ssf", "./sourcesAndSinks.txt", "Changes the file which holds the sources and sinks")
	var path = flag.String("path", "github.com/akwick/gotcha", "The path to the .go-files starting at $GOPATH/src: e.g. the path for $GOPATH/src/example/example.go will be example")
	var sourceFilesFlag sourcefiles
	var allpkgs = flag.Bool("allpkgs", false, "If it is set all packages of the source file will be analyzed, else only the main package.")
	var pkgs = flag.String("pkgs", "", "Specify some packages in addition to the main package which should be analyzed.")
	var ptr = flag.Bool("ptr", true, "If is is set we perfom a pointer analysis, else not")
	flag.Var(&sourceFilesFlag, "src", "comma-seperated list of .go-files which should be analzed")
	flag.Parse()

	if !srcFlagCalled {
		flag.PrintDefaults()
	} else {
		// analyse with the given values
		err := worklist.DoAnalysis(*path, sourceFilesFlag, *ssf, *allpkgs, *pkgs, *ptr)
		if err != nil {
			switch err := err.(type) {
			case *worklist.ErrInFlows:
				fmt.Printf("err.NumberOfFlows: %d \n", err.NumberOfFlows())
				fmt.Printf("err.Error() %s\n", err.Error())
			default:
				fmt.Printf("Errors: %+v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Gongrats. Gotcha has not found an error.\n")
			fmt.Printf("Your parameters are: \n")
			fmt.Printf("path: %s\n", *path)
			fmt.Printf("source file: %s\n", sourceFilesFlag)
			fmt.Printf("sources and sinks file: %s\n", *ssf)
			os.Exit(0)
		}
	}
}

// sourcefiles is a flag type which handles multiple .go-Files
type sourcefiles []string

var srcFlagCalled = false

// String is part of the flag.Value interface
// The output will be used in diagnostics
func (s *sourcefiles) String() string {
	return fmt.Sprint(*s)
}

// Set is part of the flag.Value interface and set the flag value
// Only files with suffix .go are allowed.
func (s *sourcefiles) Set(value string) error {
	srcFlagCalled = true
	for _, file := range strings.Split(value, ",") {
		isGoFile := strings.HasSuffix(file, ".go")
		if !isGoFile {
			errorMessage := file + "is not a .go file"
			return errors.New(errorMessage)
		}
		*s = append(*s, file)
	}
	return nil
}
