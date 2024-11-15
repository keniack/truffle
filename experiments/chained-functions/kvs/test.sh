#!/bin/bash

start=1699048961.8026214
end=1699048964.0246968

duration_seconds=$(echo "$end - $start" | bc)
milliseconds=$(printf "%.0f" "$(echo "($duration_seconds) * 1000" | bc)")

echo "Duration: $seconds seconds and $milliseconds milliseconds"