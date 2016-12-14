# Jenkins

We have a [Jenkins CI:](https://envisage.ifi.uio.no:8080/jenkins/view/Vs-dev/job/GoRETech/) which shows the statics of our testcases. 


# How to run tests on your local machine

We assume that you have installed Go correctly and are in the directory: '$GOPATH/src/goretech/analysis'.
More details about the installation is in the RunAnalysis.md file.

## Run all tests

We have a script which runs all the test packages from the goretech/analysis directory.
Ensure that the file (test.sh) has the executable flag. If not use chmod to change that.

```
$.\test.sh
```

## Run only one package

If you want to run only test from one package change into that directory and run go test.

```
$ cd tests
$ go test
```

## Run only one special test case

If you want to run only one special test case in one package, you have to add the -run parameter to go run. The Argument of the parameter is the name of the test case. The tool runs all tests which match to the argument.

```
$ cd tests
$ go test -run="Channel"
```

## What are convey tests?

Some test cases include the name convey.
These test cases uses the [goconvey library]() to execute the test.
GoConvey extends the existing go tools with an Behavior-driven Development framework for writing tests.

## How I have to read the tests messages on the command line?

A normal test case of Go first shows the name and than the messages.
A convey test handles it the other way around.
It is easy to find the start of a convey tests, because it start with a colorful representation of the assertions.
