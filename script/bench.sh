#!/bin/bash

looprun() {
  for f in \
    $infile \
    $infile \
    $infile \
    $infile \
    $infile \
    $infile
  do
    arth -m $f
  done
}

infile=testdata/bench.txt

echo median.sh vs arth
time bash script/median.sh $infile > /dev/null
time arth -m $infile > /dev/null

echo ================================

echo arth -m goroutine vs loop arth -m
time arth -m $infile $infile $infile > /dev/null
time looprun > /dev/null
