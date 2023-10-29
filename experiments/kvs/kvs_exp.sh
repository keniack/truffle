#!/bin/bash

# Declare an array of strings
my_array=("func-b-kvs" "func-a-kvs")
my_files=("file_1K.txt" "file_1M.txt" "file_2M.txt" "file_4M.txt" "file_10M.txt" "file_20M.txt")

set -e


for file in "${my_files[@]}"; do
   echo "start experiment for size $file"
   start=$(date +%s%N)
   curl -s --trace-time "http://10.152.183.128" -H "Host: func-a-kvs.default.svc.cluster.local" -d @storage/$file -H "x-target: func-b"
   target_fn_result=$(curl -s --trace-time "http://10.152.183.128" -H "Host: func-b-kvs.default.svc.cluster.local" -d @storage/$file -H "x-target: func-b")
   #end=$(date +%s% N)
   end=$(date -d "$target_fn_result" +%s%N)
   echo "End Duration: $(($(($end-$start))/1000000)) ms"
   # Sleep while function is running
   is_running=true
   while "$is_running"; do
       for fn in "${my_array[@]}"; do
               # shellcheck disable=SC2216
               kubectl_cmd=$(kubectl get pod -l function.knative.dev/name=$fn --field-selector=status.phase==Running -o jsonpath='{.items[0].metadata.name}' > /dev/null 2>&1 && echo "OK" || echo "NOK")
               if [ "$kubectl_cmd" == "OK" ]  ; then
                  echo "function $fn still running. sleep 10"
                  sleep 10  # Return false (true)
               else
                  echo "function $fn not running anymore."
                  is_running=false
               fi
       done
   done
done