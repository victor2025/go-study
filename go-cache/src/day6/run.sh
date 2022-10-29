#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
time_start=$(date +%s)
for i in {1..1000};
do
    curl "http://localhost:9999/api?key=Tom" &
done
time_end=$(date +%s)
time_diff=$(expr $time_end - $time_start)
echo "spend: $time_diff seconds"
wait