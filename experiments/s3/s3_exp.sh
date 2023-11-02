#!/bin/bash

# Declare an array of strings
my_array=("func-b-s3" "func-a-s3")
my_files=("file_1K.txt" "file_1M.txt" "file_2M.txt" "file_4M.txt" "file_10M.txt" "file_20M.txt" "file_20M.txt")

set -e

invocations=$1

for file in "${my_files[@]}"; do
   total_time=0
   for i in {1..10}; do
       # Sleep while function is running
        is_running=true
        while "$is_running"; do
            for fn in "${my_array[@]}"; do
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
       curl -s --trace-time "http://10.152.183.128" -H "Host: func-a-s3.default.svc.cluster.local" -d @storage/$file -o output &
       target_fn_result=$(curl -s --trace-time "http://10.152.183.128" -H "Host: func-b-s3.default.svc.cluster.local")
       #end=$(date +%s% N)
       end=$(date -d "$target_fn_result" +%s%N)
       duration="$(($(($end-$start))/1000000))"
       total_time=$(($total_time+$duration))
       echo "End Duration: $duration ms"

   done
   echo "Avg time for $file: $((total_time/10))"
done
