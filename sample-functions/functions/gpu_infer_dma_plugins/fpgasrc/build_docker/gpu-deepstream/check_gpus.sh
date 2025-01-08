#!/bin/bash

cd `dirname $0`

GPU_ID_A100=`./find_gpu.sh A100`
GPU_ID_T4=`./find_gpu.sh T4`
if [ $GPU_ID_A100 -eq -1 -a $GPU_ID_T4 -eq -1 ]; then
    echo "no target GPUs not found: building container image for GPU inference app requires either A100 or T4"
    exit 1
fi
