#/bin/bash

# Copyright 2022 NTT Corporation, FUJITSU LIMITED 


if [ $# -ne 5 ]; then
  echo "Usage: ./start_gst_sender.sh <movie_file> <host_ip> <start_port> <session_num> <sleep_time>"
  exit 1
fi


movie_file=$1
host_ip=$2
start_port=$3
session_num=$4
sleep_time=$5


for ((port=$start_port; port < $(($start_port+$session_num)); port+=1)); do
    /bin/bash gst_child_process.sh $movie_file $host_ip $port &
    sleep $sleep_time
done

