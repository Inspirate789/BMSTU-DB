#!/bin/bash

rm -f log/cpu.csv

for (( size=2500; size<=25000; size+=2500 ))
do
  echo "size $size in processing"
  go test -bench=. -cpuprofile=cpu.out -benchtime=100x -size=$size > log/tmp
  nums=( $(cat log/tmp) )
  echo "${nums[20]},${nums[24]},${nums[35]},${nums[46]},${nums[57]},${nums[61]},${nums[72]},${nums[83]},${nums[94]}" >> log/cpu.csv
  echo "size $size done"
  # sleep 5m
done

rm -f log/tmp