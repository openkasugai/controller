#!/bin/bash
# Copyright 2024 NTT Corporation, FUJITSU LIMITED

############################################
# Common Parameter Settings
############################################
### app log level settings: any environment variable
export DECENV_APPLOG_LEVEL=4

### Input Video Distribution Protocol (RTP or RTSP): required environment variables
export DECENV_VIDEOSRC_PROTOCOL="RTP"

### Input Video Delivery IPv4 Port: Required Environment Variables
export DECENV_VIDEOSRC_PORT=5004

### Frame FPS: required environment variables
export DECENV_FRAME_FPS=5.0

### frame size: required environment variable
export DECENV_FRAME_WIDTH=3840
export DECENV_FRAME_HEIGHT=2160

### Output protocol (DMA or TCP): required environment variables
export DECENV_OUTDST_PROTOCOL="DMA"

### Limit the number of video distribution connections: Optional environment variable
export DECENV_VIDEO_CONNECT_LIMIT=0


############################################
# Setting parameters for video distribution RTSP server
# Used when DECENV_VIDEOSRC_PROTOCOL="RTSP"
############################################
### Video distribution RTSP server IPv4 address: required environment variable
export DECENV_VIDEOSRC_IPA="192.168.0.51"


############################################
# Parameter settings for FPGA DMA output
# Used when DECENV_OUTDST_PROTOCOL="DMA"
############################################
### DPDK file prefix: any environment variable
export DECENV_DPDK_FILE_PREFIX="0"

### FPGA DMA_HOST_TO_DEV connector ID: required environment variable
export DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID="connector_id_0"

### Destination shared memory address (SRC shared memory on the FPGA): any environment variable
export DECENV_FPGA_SRC_SHMEM_ADDR=0x17ff29c00

### FPGA device files: arbitrary environment variables
export DECENV_FPGA_DEV_NAME="/dev/xpcie_2133072BM03P"

### FPGA channel ID: any environment variable
export DECENV_FPGA_CH_ID=0

### FPGA output frame size: arbitrary environment variable
export DECENV_FPGA_OUT_FRAME_WIDTH=1280
export DECENV_FPGA_OUT_FRAME_HEIGHT=1280


############################################
# Setting parameters for TCP output
# Used when DECENV_OUTDST_PROTOCOL="TCP"
############################################
### TCP output destination IP address: required environment variable
export DECENV_OUTDST_IPA="192.168.0.222"

### TCP output destination IP port: required environment variables
export DECENV_OUTDST_PORT=12000


############################################
# Run Application
############################################
sudo -E ./build/cpu_decode

