#!/bin/bash

infile=testdata/bench.txt
time bash script/median.sh $infile
time arth -m $infile
