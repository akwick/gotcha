#!/bin/bash
#Builds the analysis.debug file with some instrumentations

godebug build -instrument="github.com/akwick/gotcha/lattice/taint,github.com/akwick/gotcha/lattice,github.com/akwick/gotcha/worklist"
