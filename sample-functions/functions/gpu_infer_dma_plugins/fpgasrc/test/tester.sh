#!/bin/sh

export FPGA_DEV=/dev/xpcie_21330621T049
export INPUT_WIDTH=448
export INPUT_HEIGHT=448
export CONNECTOR_ID=connector_id_rx_0

[ -z "$FPGA_DEV" ]       && echo "environment variable FPGA_DEV must be specified" && exit
[ -z "$INPUT_WIDTH" ]    && echo "environment variable INPUT_WIDTH must be specified" && exit
[ -z "$INPUT_HEIGHT" ]   && echo "environment variable INPUT_HEIGHT must be specified" && exit
[ -z "$CONNECTOR_ID" ]   && echo "environment variable CONNECTOR_ID must be specified" && exit

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:../../../../lib/DPDK/dpdk/lib/x86_64-linux-gnu
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:../../../../lib/build
./tester
