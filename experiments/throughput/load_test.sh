#!/bin/bash

#define variables

#set -x # run in debug mode

for j in {1..5}; do
  avg_throughtput=0
  DURATION=1 # how long should load be applied ? - in seconds
  TPS=$1 # number of requests per second
  end=$((SECONDS+$DURATION))
  #start load
  file_name=response-times.log+"$TPS"
  for ((i=1;i<=$TPS;i++)); do
      curl -s -X POST "http://10.152.183.70" -H 'x-source: func-b' -H "Host: func-a-b.default.svc.cluster.local" -d 'test' -w '\n%{stderr}' >> $file_name &
  done

  wait
  #end load

  # Receive file name as first argument
  total_time=0
  while read -r line; do
      # Reading line by line
      float_duration=$(printf "%.0f" "$line")
      total_time=$((total_time + float_duration))
  done < $file_name
  total_time=`awk "BEGIN {print $total_time/1000}"`
  echo "Total duration $total_time"
  troughput=`awk "BEGIN {print $TPS/$total_time}"`
  echo "Avg req/sec: $troughput"
  echo "Load test has been completed"

  avg_throughtput=`awk "BEGIN {print $avg_throughtput+$troughput}"`
  rm $file_name
done

avg_throughtput=`awk "BEGIN {print $avg_throughtput/5}"`
echo "Final Throughput req/sec: $avg_throughtput"