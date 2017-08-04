# Jenkins

We have a [Jenkins CI:](https://envisage.ifi.uio.no:8080/jenkins/view/go test github.com/akwick/gotcha/tests github.com/akwick/gotcha/worklist github.com/akwick/gotcha/lattice/taint
Vs-dev/job/GoRETech/) which shows the statics of our testcases.


# How to run tests on your local machine

## Installations Requirements

First, we assume that you have installed Go correctly on your machine. 
For testing, we use two additional packages (testify/assert and goconvey/convey). One can install bothpackages with the corresponding go get commands:

```
go get github.com/stretchr/testify/assert
go get github.com/smartystreets/goconvey/convey
```

## Run all tests

We have a script which runs all the test packages from the github.com/akwick/gotcha directory.
Ensure that the file (test.sh) has the executable flag. If not use chmod to change the flags.

```
$.\tests.sh
```

An alternative is to run the corresponding go command from the script:

```
go test github.com/akwick/gotcha/tests github.com/akwick/gotcha/worklist github.com/akwick/gotcha/lattice/taint
```


## Run only one package

If you want to run only tests from one package, you should change into that directory and run go test within this directory.

```
$ cd tests
$ go test
```
The alternative is to use the following command:
## Run only one special test case

If you want to run only one special test case in one package, you have to add the -run parameter to the go test command. The Argument of the parameter is the name of the test case. The tool runs all tests which match to the argument.

```
$ cd tests
$ go test -run="Channel"
```

## What are convey tests?

Some test cases include the name convey.
These test cases use the [goconvey library]() to execute the test.
GoConvey extends the existing go tools with a Behavior-driven Development framework for writing tests.

## How I have to read the tests messages on the command line?

A normal test case of Go first shows the name and than the messages.
A convey test handles it the other way around.
It is easy to find the start of a convey tests, because it starts with a colorful representation of the assertions.

## Annotations in the example files

In some of our examples we will use annotations to illustrate the expected report of gotcha.
The annoations should be above a function which is in the set of sinks being the function sink(s string) in most of the cases.

* `@expectedflow: true` : We expect gotcha to report a flow
* `@expectedflow: false` : We do not expect gotcha to report a flow
