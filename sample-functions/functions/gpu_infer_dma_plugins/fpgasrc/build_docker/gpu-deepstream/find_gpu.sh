#!/bin/bash

GPU=$1
nvidia-smi | grep $GPU > /dev/null
if [ $? -eq 0 ]; then
    GPU_ID=`nvidia-smi | grep $GPU | head -1 | sed -e 's/  */ /g' | cut -d ' ' -f 2`
else
    GPU_ID=-1
fi
echo $GPU_ID
