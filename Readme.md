# Gotcha - **Go T**aint **Ch**eck **A**nalysis

## Requirements
0. If the `$GOPATH` is not set, then: export GOPATH=$HOME/go (this directory is your workspace, create it somewhere)
1. '' go get github.com/akwick/gotcha ``



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

We have several tests which ensure some functionality of our analysis.
The results are available via [Jenkins](https://envisage.ifi.uio.no:8080/jenkins/view/Vs-dev/job/GoRETech/)
Are more detailed descriptions about running tests on your machine are in the file [*tests.md*](https://github.com/akwick/gotcha/blob/master/tests.md)

# Debug the program

The repository has a small shell script which can build a debug file.
A reference for the commands is in the [repository of godebug](https://github.com/mailgun/godebug).

```
$ ./debug.sh
$ ./analysis.debug -src="fileyouwanttodebug"
```
