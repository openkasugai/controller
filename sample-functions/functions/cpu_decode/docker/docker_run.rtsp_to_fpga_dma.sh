#!/bin/bash
# Copyright 2024 NTT Corporation, FUJITSU LIMITED

docker run \
  -it \
  -d \
  --net=host \
  --device=/dev/xpcie_2133072BM03P \
  --privileged \
  -v=/var/run/dpdk:/var/run/dpdk \
  -v=/dev/hugepages:/dev/hugepages \
  -e DECENV_APPLOG_LEVEL=4 \
  -e DECENV_VIDEOSRC_PROTOCOL="RTSP" \
  -e DECENV_VIDEOSRC_PORT=8554 \
  -e DECENV_VIDEOSRC_IPA="192.168.0.51" \
  -e DECENV_FRAME_FPS=5.0 \
  -e DECENV_FRAME_WIDTH=3840 \
  -e DECENV_FRAME_HEIGHT=2160 \
  -e DECENV_OUTDST_PROTOCOL="DMA" \
  -e DECENV_VIDEO_CONNECT_LIMIT=0 \
  -e DECENV_DPDK_FILE_PREFIX="0" \
  -e DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID="connector_id_0" \
  -e DECENV_FPGA_SRC_SHMEM_ADDR=0x17ff29c00 \
  -e DECENV_FPGA_DEV_NAME="/dev/xpcie_2133072BM03P" \
  -e DECENV_FPGA_CH_ID=0 \
  -e DECENV_FPGA_OUT_FRAME_WIDTH=1280 \
  -e DECENV_FPGA_OUT_FRAME_HEIGHT=1280 \
  --name cpu_decode_00 \
  cpu_decode
