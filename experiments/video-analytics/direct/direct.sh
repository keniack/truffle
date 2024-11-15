#!/bin/bash

# Declare an array of strings
my_array=("streaming-s3" "decode-s3" "recognition-s3")
my_files=("video_1K.mp4"  "video_1M.mp4" "video_2M.mp4" "video_4M.mp4" "video_10M.mp4" "video_20M.mp4"  "video_40M.mp4")
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
       duration=$(curl -s --trace-time "http://10.152.183.128" -H "Host: streaming.default.svc.cluster.local" --data-binary @storage/$file -H "x-target: func-b" -H "x-filename: video.mp4")
       float_duration=$(printf "%.0f" "$duration")
       total_time=$((total_time + float_duration))
       echo "End Duration: $float_duration ms"
       # Sleep while function is running
    done
   echo "Avg time for $file: $((total_time/5))"
done