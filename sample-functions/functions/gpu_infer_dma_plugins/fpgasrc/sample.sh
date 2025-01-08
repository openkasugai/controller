#!/bin/sh

export FPGA_DEV=/dev/xpcie_21330621T049
export FILE_PREFIX=file_prefix
export WIDTH=448
export HEIGHT=448
export CONNECTOR_ID=connector_id_tx_0

##export SHMEM_SECONDARY=1
##export DEBUG_MODE=1

export GST_PLUGIN_PATH=$PWD

[ -z "$FPGA_DEV" ]     && echo "environment variable FPGA_DEV must be specified" && exit
[ -z "$FILE_PREFIX" ]  && echo "environment variable FILE_PREFIX must be specified" && exit
[ -z "$WIDTH" ]        && echo "environment variable WIDTH must be specified" && exit
[ -z "$HEIGHT" ]       && echo "environment variable HEIGHT must be specified" && exit
[ -z "$CONNECTOR_ID" ] && echo "environment variable CONNECTOR_ID must be specified" && exit

export FRAMERATE=${FRAMERATE:-30}

gst-launch-1.0 \
    -e \
    fpgasrc ! \
    "video/x-raw, format=(string)BGR, width=$WIDTH, height=$HEIGHT, framerate=$FRAMERATE/1" ! \
    timeoverlay ! \
    queue ! \
    nvvideoconvert ! \
    'video/x-raw(memory:NVMM), format=(string)RGBA' ! \
    m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1 width=$WIDTH height=$HEIGHT ! \
    queue ! \
    nvinfer config-file-path= /opt/nvidia/deepstream/deepstream/samples/configs/deepstream-app/config_infer_primary.txt batch-size=1 unique-id=1 ! \
    queue ! \
    nvdsosd process-mode=1 ! \
    queue ! \
    nvvideoconvert ! \
    videoconvert ! \
    x264enc ! \
    video/x-h264, stream-format=byte-stream ! \
    queue ! \
    h264parse ! \
    qtmux ! \
    perf ! \
    filesink location=/tmp/test_fpga_${CONNECTOR_ID}.mp4 sync=1 > /tmp/${CONNECTOR_ID}.log
