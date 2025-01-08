#!/bin/bash
# Copyright 2024 NTT Corporation, FUJITSU LIMITED

function usage() {
cat <<EOF
  Usage: $0 [version]
    Options:
      [version] : Set the image version.
EOF
  exit 1
}

if [[ $# != 1 ]]; then
  usage
fi

VERSION=$1

############################################
# Parameter Settings
############################################
# Created Image Name
IMAGE_NAME=cpu_decode:${VERSION}

# Repository root directory path
ROOT_DIR=../../../../

# cpu_decode directory path
CPU_DECODE_DIR=sample-functions/functions/cpu_decode/

# FPGA library directory path
FPGA_SOFTWARE_DIR=src/submodules/fpga-software/


############################################
# Create Image
############################################
# Go to repository root directory
cd $ROOT_DIR

# Run docker build
docker build . -t ${IMAGE_NAME} -f ${CPU_DECODE_DIR}/docker/Dockerfile --build-arg CPU_DECODE_DIR=${CPU_DECODE_DIR} --build-arg FPGA_SOFTWARE_DIR=${FPGA_SOFTWARE_DIR}
