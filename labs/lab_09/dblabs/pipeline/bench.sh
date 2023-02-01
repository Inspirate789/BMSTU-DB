#!/bin/bash

rm -f log/cpu.csv

for (( rcnt=50; rcnt<=500; rcnt+=50 ))
do
  echo "$rcnt requests in processing"
  go test -bench=. -cpuprofile=cpu.out -benchtime=50x -rcnt=$rcnt > log/tmp
  nums=( $(cat log/tmp) )
  echo "${nums[15]},${nums[19]}" >> log/cpu.csv
  echo "$rcnt requests done"
done

rm -f log/tmp
