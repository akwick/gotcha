# Run the default analysis

## Requierements
3. Clone the repository to a folder you prefer: `$ git clone https://git.cased.de/git/ec-spride-sse-internal-goretech`
0. Installed go with version >= 1.5
1. If the `$GOPATH` is not set, then: export GOPATH=$HOME/go (this directory is your workspace, create it somewhere)
2. `ln -s PATH_TO_GIT/goretech $GOPATH/src/goretech`  
2. `$ cd $GOPATH/src/github.com/akwick/gotcha`
8. add missing packages via go get
  * `go get golang.org/x/tools/go/ssa`
  * `go get github.com/stretchr/testify/assert`
  * `go get github.com/pkg/errors`
  * `go get github.com/smartystreets/goconvey/convey`
  * if you want to debug: `$ go get github.com/mailgun/godebug`



## Build the analysis

0. cd $GOPATH/src/github.com/akwick/gotcha
1. go build

## Run the analysis

0. ./analysis -path="path to go-files as relative part from $GOPATH/src" -src="path to source code file which should analyzed" -ssf="path to the sources and sinks file"
`./analysis -src="tests/exampleCode/hello.go"`
1. The -src flag is mandatory, the path, ssf, allpkgs, pkgs and ptrflag are optional.
2. The default parameter are:
  - path = github.com/akwick/gotcha
    - It is important to change the path if you are not running our examples.   

  - ssf = ./sourcesAndSinks.txt
  - allpkgs = false
  - pkgs = ""  
  - ptr = true    
3. `./analysis -h` prints a short help for the flags.  

# Test Results

We have several tests which ensures some functionality of our analysis.
The results are available via [Jenkins](https://envisage.ifi.uio.no:8080/jenkins/view/Vs-dev/job/GoRETech/)
Are more detailed descriptions about running tests on your machine are in tests.md

# Debug the program

The repository has a small shell script which can build a debug file.
A reference for the commands is in the [repository of godebug](https://github.com/mailgun/godebug).

```
$ ./debug.sh
$ ./analysis.debug -src="fileyouwanttodebug"
```
