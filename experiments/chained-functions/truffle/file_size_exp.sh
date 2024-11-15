#!/bin/bash

# Declare an array of strings
my_array=("func-b" "func-a-b-direct")
my_files=("file_1K.txt" "file_1M.txt" "file_2M.txt" "file_4M.txt" "file_10M.txt" "file_20M.txt" "file_40M.txt")
#my_files=("file_4M.txt" "file_10M.txt" "file_20M.txt" "file_40M.txt")

set -e


for file in "${my_files[@]}"; do
  total_time=0
   for i in {1..5}; do
       is_running=true
       while "$is_running"; do
           for fn in "${my_array[@]}"; do
               # shellcheck disable=SC2216
               kubectl_cmd=$(kubectl get pod -l function.knative.dev/name=$fn --field-selector=status.phase==Running -o jsonpath='{.items[0].metadata.name}' > /dev/null 2>&1 && echo "OK" || echo "NOK")
               if [ "$kubectl_cmd" == "OK" ]  ; then
                  #echo "function $fn still running. sleep 10"
                  sleep 10  # Return false (true)
               else
                  #echo "function $fn not running anymore."
                  is_running=false
               fi
           done
       done

       echo "start experiment for size $file index $i"
       start=$(date +%s%N)
       duration=$(curl -s --trace-time "http://10.152.183.128" -H "Host: func-a-b-direct.default.svc.cluster.local" -d @storage/$file -H "x-target: func-b")
       float_duration=$(printf "%.0f" "$duration")
       total_time=$((total_time + float_duration))
       echo "End Duration: $float_duration ms"
       # Sleep while function is running
    done
   echo "Avg time for $file: $((total_time/5))"
done