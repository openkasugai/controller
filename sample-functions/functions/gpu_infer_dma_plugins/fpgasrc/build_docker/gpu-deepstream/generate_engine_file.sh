#!/bin/bash

export GST_PLUGIN_PATH=$PWD
cd `dirname $0`

GPU=$1
CONFIG_FILE=$2
ENGINE_FILE=$3

GPU_ID=`./find_gpu.sh $GPU`
if [ $GPU_ID -ne -1 ]; then
    echo "$GPU found"
else
    echo "$GPU not found"
    exit
fi

if [ $GPU == "A100" ]; then
    WIDTH=1280
    HEIGHT=1280
elif [ $GPU == "T4" ]; then
    WIDTH=416
    HEIGHT=416
fi
FRAMERATE=${FRAMERATE:-30}

gst-launch-1.0 \
    -e \
    videotestsrc ! \
    "video/x-raw, format=(string)BGR, width=$WIDTH, height=$HEIGHT, framerate=$FRAMERATE/1" ! \
    timeoverlay ! \
    queue ! \
    nvvideoconvert gpu-id=$GPU_ID ! \
    'video/x-raw(memory:NVMM), format=(string)RGBA' ! \
    m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1 width=$WIDTH height=$HEIGHT gpu-id=$GPU_ID ! \
    queue ! \
    nvinfer config-file-path=$CONFIG_FILE batch-size=1 model-engine-file=$ENGINE_FILE gpu-id=$GPU_ID ! \
    fakesink &

PID=$!
echo "generating engine file..."
NEW_ENGINE_FILE=`echo $ENGINE_FILE | sed s/gpu0/gpu$GPU_ID/`
while [ ! -e $NEW_ENGINE_FILE ]
do
      sleep 1
done
# wait for engine file to be written out
while true
do
    CHANGED=`stat $NEW_ENGINE_FILE --format=%Y`
    NOW=`date +%s`
    DIFF=`expr $NOW - $CHANGED`
    if [ $DIFF -gt 10 ]; then
	break
    fi
    sleep 1
done
if [ $NEW_ENGINE_FILE != $ENGINE_FILE  ]; then
    mv $NEW_ENGINE_FILE $ENGINE_FILE
fi
kill $PID
echo "engine file generation done!"
