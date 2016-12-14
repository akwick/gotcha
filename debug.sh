#!/bin/bash
#Builds the analysis.debug file with some instrumentations

godebug build -instrument="goretech/analysis/lattice/taint,goretech/analysis/lattice,goretech/analysis/worklist"
