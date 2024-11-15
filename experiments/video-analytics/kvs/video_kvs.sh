#!/bin/bash

# Declare an array of strings
my_array=("streaming" "decode" "recognition")
my_files=("video_1K.mp4"  "video_1M.mp4" "video_2M.mp4" "video_4M.mp4" "video_10M.mp4" "video_20M.mp4"  "video_40M.mp4")

set -e

invocations=$1

for file in "${my_files[@]}"; do
   total_time=0
   for i in {1..5}; do
       is_running=true
       while "$is_running"; do
            for fn in "${my_array[@]}"; do
                  # shellcheck disable=SC2216
                  kubectl_cmd=$(kubectl get pod -l serving.knative.dev/configuration=$fn --field-selector=status.phase==Running -o jsonpath='{.items[0].metadata.name}' > /dev/null 2>&1 && echo "OK" || echo "NOK")
                  if [ "$kubectl_cmd" == "OK" ]  ; then
                     #echo "function $fn still running. sleep 10"
                     sleep 10  # Return false (true)
                  else
                     #echo "function $fn not running anymore."
                     is_running=false
                  fi
            done
       done
       echo "start experiment for size $file"
       start=$(curl -s --trace-time "http://10.152.183.128" -H "Host: streaming.default.svc.cluster.local" --data-binary @storage/$file -H "x-target: func-b" -H "x-comm: KVS")
       curl -s --trace-time "http://10.152.183.128" -H "Host: decode.default.svc.cluster.local" -d "decode" -H "x-filename: video.mp4" -H "x-comm: KVS"
       end=$(curl -s --trace-time "http://10.152.183.128" -H "Host: recognition.default.svc.cluster.local" -d @storage/$file -H "x-target: func-b" -H "x-filename: image.jpg" -H "x-comm: KVS")
       duration=$(echo "$end - $start" | bc)
       milliseconds=$(printf "%.0f" "$(echo "($duration) * 1000" | bc)")
       #end=$(date -d "$target_fn_result" +%s%N)
       #duration="$(($(($end-$start))/1000000))"
       total_time=$(($total_time+$milliseconds))
       echo "End Duration: $milliseconds ms"
       # Sleep while function is running
   done
   echo "Avg time for $file: $((total_time/5))"
done