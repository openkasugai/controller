#!/bin/bash -x
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

output=$(lsmod | grep xpcie)
if [ -n "$output" ]; then
  echo $output
  sudo rmmod xpcie
fi

BIT_DIR=$HOME/hardware-design/example-design/bitstream
cd $BIT_DIR
sudo mcap -x 903f -s 1f:00.0 -p OpenKasugai-fpga-example-design-1.0.0-2.bit

FPGA_DRIVER_DIR=$HOME/controller/src/submodules/fpga-software/driver
cd $FPGA_DRIVER_DIR
sudo insmod xpcie.ko
lsmod |grep xpcie
ls -l /dev/xpcie*
