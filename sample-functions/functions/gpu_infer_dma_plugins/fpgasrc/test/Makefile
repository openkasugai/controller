# Copyright 2024 NTT Corporation, FUJITSU LIMITED

SRC_DIR=/opt/nvidia/deepstream/deepstream-7.0/fpga-software
CFLAGS += -I../include -I..

# libfpga
LDFLAGS += $(shell PKG_CONFIG_PATH=${SRC_DIR}/lib/build/pkgconfig pkg-config --libs libfpga)
CFLAGS += $(shell PKG_CONFIG_PATH=${SRC_DIR}/lib/build/pkgconfig pkg-config --cflags libfpga)

# DPDK
LDFLAGS += $(shell PKG_CONFIG_PATH=${SRC_DIR}/lib/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --libs libdpdk)
CFLAGS += $(shell PKG_CONFIG_PATH=${SRC_DIR}/lib/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --cflags libdpdk)

all: controller tester recver

controller: controller.c lib_controller.c
	$(CC) $(CFLAGS) controller.c lib_controller.c $(LDFLAGS)  -o $@

tester: tester.c
	$(CC) $(CFLAGS) tester.c ../util.c $(LDFLAGS) -o $@

recver: recver.c
	$(CC) $(CFLAGS) recver.c ../util.c $(LDFLAGS) -o $@

.PHONY: clean
clean:
	rm -f *.o controller tester recver
