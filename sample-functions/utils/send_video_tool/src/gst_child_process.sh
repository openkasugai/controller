#!/bin/bash

# Copyright 2022 NTT Corporation, FUJITSU LIMITED 



movie_file=$1
host_ip=$2
port=$3


#while true; do
gst-launch-1.0 filesrc location=$movie_file ! qtdemux ! video/x-h264 ! h264parse ! rtph264pay config-interval=-1 seqnum-offset=1 ! udpsink host=$host_ip port=${port} buffer-size=2048 > /proc/1/fd/1 

#  gst-launch-1.0 filesrc location=$movie_file ! qtdemux ! video/x-h264 ! h264parse ! rtph264pay config-interval=-1 seqnum-offset=1 ! udpsink host=$host_ip port=${port} buffer-size=2048 >> gst_sender_$port.log 
#done
