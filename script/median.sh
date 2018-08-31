#!/bin/bash

n=$(cat $1 | wc -l)
n=$((n / 2))

sort -n $1 | awk '
BEGIN{ min = 0 }

{
  if ($1 < min) {
    min = $1
  }
  if (max < $1) {
    max = $1
  }
  sum += $1
}

NR=='"$n"'{median=$1}

END{
  avg = sum / NR
  print "count", "min", "max", "sum", "avg", "median"
  print NR, min, max, sum, avg, median
}'

