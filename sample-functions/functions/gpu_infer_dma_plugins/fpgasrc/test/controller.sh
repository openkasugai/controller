#!/bin/sh

export FPGA_DEV=/dev/xpcie_21330621T049
export FILE_PREFIX=file_prefix
export CH_ID=0
export CONNECTOR_ID_RX=connector_id_rx_0
export CONNECTOR_ID_TX=connector_id_tx_0
export INPUT_WIDTH=448
export INPUT_HEIGHT=448
export OUTPUT_WIDTH=448
export OUTPUT_HEIGHT=448

[ -z "$FPGA_DEV" ]        && echo "environment variable FPGA_DEV must be specified" && exit
[ -z "$FILE_PREFIX" ]     && echo "environment variable FILE_PREFIX must be specified" && exit
[ -z "$CH_ID" ]           && echo "environment variable CH_ID must be specified" && exit
[ -z "$CONNECTOR_ID_RX" ] && echo "environment variable CONNECTOR_ID_RX must be specified" && exit
[ -z "$CONNECTOR_ID_TX" ] && echo "environment variable CONNECTOR_ID_TX must be specified" && exit
[ -z "$INPUT_WIDTH" ]     && echo "environment variable INPUT_WIDTH must be specified" && exit
[ -z "$INPUT_HEIGHT" ]    && echo "environment variable INPUT_HEIGHT must be specified" && exit
[ -z "$OUTPUT_WIDTH" ]    && echo "environment variable OUTPUT_WIDTH must be specified" && exit
[ -z "$OUTPUT_HEIGHT" ]   && echo "environment variable OUTPUT_HEIGHT must be specified" && exit

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:../../../../lib/DPDK/dpdk/lib/x86_64-linux-gnu
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:../../../../lib/build
./controller
