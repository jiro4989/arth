#!/bin/bash

looprun() {
  for f in \
    $infile \
    $infile \
    $infile
  do
    arth -m $f
  done
}

infile=testdata/bench.txt

echo median.sh vs arth
echo median.sh
time bash script/median.sh $infile > /dev/null
echo
echo arth
time arth -m $infile > /dev/null
echo ================================

echo loop arth vs arth goroutine
echo loop arth
time looprun > /dev/null
echo
echo arth goroutine
time arth -m $infile $infile $infile > /dev/null
