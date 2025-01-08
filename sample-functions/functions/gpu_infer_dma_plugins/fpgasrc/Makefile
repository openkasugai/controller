# Copyright 2024 NTT Corporation, FUJITSU LIMITED

SRC_DIR=/opt/nvidia/deepstream/deepstream-7.0/openkasugai-controller
LIB_DIR=${SRC_DIR}/src/submodules/fpga-software/lib

#ARG FUNCTION_PATH=${SRC_DIR}/sample-functions/functions/gpu_infer_dma_plugins

CFLAGS += -I../include

# driver
CFLAGS += -I${SRC_DIR}/src/submodules/fpga-software/driver

# libfpga
LDFLAGS += $(shell PKG_CONFIG_PATH=${LIB_DIR}/build/pkgconfig pkg-config --libs libfpga)
CFLAGS += $(shell PKG_CONFIG_PATH=${LIB_DIR}//build/pkgconfig pkg-config --cflags libfpga)

# DPDK
LDFLAGS += $(shell PKG_CONFIG_PATH=${LIB_DIR}/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --libs libdpdk)
CFLAGS += $(shell PKG_CONFIG_PATH=${LIB_DIR}/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --cflags libdpdk)

# GStreamer
gstfpgasrc.so: CFLAGS += -fPIC -shared
gstfpgasrc.so: CFLAGS += $(shell pkg-config --cflags gstreamer-1.0 gstreamer-base-1.0)
gstfpgasrc.so: LDFLAGS += $(shell pkg-config --libs gstreamer-1.0 gstreamer-base-1.0)

all: gstfpgasrc.so

prepare: gstfpgasrc.patch
	rm -rf gst-plugins-bad
	git clone https://github.com/GStreamer/gst-plugins-bad.git
	cd gst-plugins-bad && git checkout refs/tags/1.19.2 && cd ..
	gst-plugins-bad/tools/gst-element-maker fpgasrc basesrc
	rm -rf gst-plugins-bad
	patch -p4 < gstfpgasrc.patch

gstfpgasrc.so: gstfpgasrc.c lib_worker.c ${LIB_DIR}/build/libfpga.so
	$(CC) $(CFLAGS) gstfpgasrc.c lib_worker.c util.c $(LDFLAGS) -o $@

.PHONY: clean
clean:
	rm -f *.o *.so
