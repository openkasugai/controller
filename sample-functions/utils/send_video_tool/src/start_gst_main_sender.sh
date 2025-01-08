#!/bin/bash

# Copyright 2022 NTT Corporation, FUJITSU LIMITED 


sleep_time=$1

./start_gst_sender.sh /opt/video/sample_720p.mp4 127.0.0.1 988 2 ${sleep_time:-3} &

#sleep ${sleep_time:-3}
#./start_gst_sender.sh /opt/video/Pexels_Videos_2053100_4K_1_conv_15fps_5Mbps.mp4 127.0.0.1 1000 9 ${sleep_time:-3} &

#sleep ${sleep_time:-3}
#./start_gst_sender.sh /opt/video/Pexels_Videos_2053100_4K_1_conv_15fps_5Mbps.mp4 127.0.0.1 1100 1 ${sleep_time:-3} &

#sleep ${sleep_time:-3}
#./start_gst_sender.sh /opt/video/Pexels_Videos_2053100_4K_1_conv_5fps_2Mbps.mp4 127.0.0.1 1200 8 ${sleep_time:-3} &


