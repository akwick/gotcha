#!/bin/bash
#Executes the tests of this repository.

go test github.com/akwick/gotcha/tests github.com/akwick/gotcha/worklist github.com/akwick/gotcha/lattice/taint
