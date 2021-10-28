#!/bin/bash - 

target=${1:?}
problem=0
echo "Checking target - $target"
while true; do 
  curl -s -m1 $1 -k --show-error -o /dev/null
  if [[ ! $? == 0 ]]; then
    if [[ $problem == 0 ]]; then
      echo "Problem detected"
      problem=1
      startTime=$(date +%s)
    fi
  elif [[ $problem == 1 ]] && [[ $? == 0 ]]; then
    problem=0
    echo "Recovery Detected. $(($(date +%s)-$startTime)) seconds elapsed"
  fi
done

