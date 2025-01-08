#!/bin/bash

export GLUEENV_DPDK_FILE_PREFIX=file_prefix
export GLUEENV_FPGA_DEV_NAME=/dev/xpcie_21330621T049
export GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID=connector_id_tx_0
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:../../lib/DPDK/dpdk/lib/x86_64-linux-gnu/
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:../../lib/build/

./build/glue 127.0.0.1:5000 448 448
