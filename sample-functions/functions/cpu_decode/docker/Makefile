# Copyright 2024 NTT Corporation, FUJITSU LIMITED

# ====================================================
# Build Settings
# ====================================================
# --- Common Build Settings -------------- #
# Application log standard output
APPLOG_PRINT := 1

# Modular FPGA support
MODULE_FPGA := 1

# For debugging: Standard output of the execution progress log
DEBUG_PRINT := 0


# --- Build Settings for FPGA ------------ #
# DPDK Secondary Process Mode
DPDK_SECONDARY_PROC_MODE := 1


# For stand-alone execution: Allocate FPGA SRC shared memory
ALLOCATE_SRC_SHMEM_MODE := 1

# For single execution: issue ENQUEUE DMA request to FPGA
EXEC_FPGA_ENQUEUE_MODE := 1

# For single execution: issue DEQUEUE DMA request to FPGA
EXEC_FPGA_DEQUEUE_MODE := 0


# For standalone execution: FPGA device settings [fpga_dev_init]
CONTROL_FPGA_DEV_INIT := 1

# For single execution: FPGA function settings [fpga_function_config, fpga_function_init, fpga_function_set]
CONTROL_FPGA_FUNC_INIT := 0

# For single execution: FPGA ENQUEUE DMA (DMA_HOST_TO_DEV) channel setting [fpga_lldma_init, fpga_lldma_finish]
CONTROL_FPGA_ENQUEUE_LLDMA_INIT := 0

# For single execution: FPGA function chain connection setting [fpga_chain_connect]
CONTROL_FPGA_CHAIN_CONNECT := 0

# For stand-alone execution: Exit FPGA device [fpga_finish] # Valid only when CONTROL_FPGA_DEV_INIT=1
CONTROL_FPGA_FINISH := 0

# For stand-alone execution: End of FPGA function [fpga_function_finish] * Valid only when CONTROL_FPGA_FUNC_INIT=1
CONTROL_FPGA_FUNC_FINISH := 0

# for single execution: FPGA shared memory end [fpga_shmem_finish]
CONTROL_FPGA_SHMEM_FINISH := 1


# For debugging: FPGA output data header to standard output # Valid only when EXEC_FPGA_DEQUEUE_MODE=1
DEBUG_PRINT_FPGA_OUT_HEADER := 0

# For debugging: Save frames of FPGA output data as mp4 movie files * Valid only when EXEC_FPGA_DEQUEUE_MODE=1
DEBUG_FPGA_OUT_IMAGE_TO_MP4 := 0


# ====================================================
# SETUP
# ====================================================
# binary name
APP = cpu_decode

# all source are stored in SRCS
SRCS := ./src/cpu_decode.cpp ./src/tp_func.cpp ./src/logging.cpp ./src/recv_header_log.cpp
SRCD := ../../../src/submodules/fpga-software/lib

DRIVERDIR := $(SRCD)/../driver

# libfpga
LIBFPGA := $(SRCD)/build/libfpga.a

# lib
LDFLAGS := -lpthread -lpciaccess

# OpenCV
export PKG_CONFIG_PATH=/usr/local/opencv-3.4.3/lib/pkgconfig
PKGCONFOCV := `pkg-config --cflags opencv pkg-config --libs opencv`

# for dpdk config
PKGCONF ?= env PKG_CONFIG_PATH=`cd $(SRCD)/DPDK; pwd`/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config

# optimization and include directory setup
CFLAGS += -O3 -I./include -I$(SRCD)/libfpga/include -I$(SRCD)/libptu/lldma/include -I$(DRIVERDIR) -Wno-unused-result

ifeq ($(APPLOG_PRINT),1)
CFLAGS += -DAPPLOG_PRINT
endif

ifeq ($(MODULE_FPGA),1)
CFLAGS += -DMODULE_FPGA
endif

ifeq ($(DPDK_SECONDARY_PROC_MODE),1)
CFLAGS += -DDPDK_SECONDARY_PROC_MODE
endif

ifeq ($(ALLOCATE_SRC_SHMEM_MODE),1)
CFLAGS += -DALLOCATE_SRC_SHMEM_MODE
endif
ifeq ($(EXEC_FPGA_ENQUEUE_MODE),1)
CFLAGS += -DEXEC_FPGA_ENQUEUE_MODE
endif
ifeq ($(EXEC_FPGA_DEQUEUE_MODE),1)
CFLAGS += -DEXEC_FPGA_DEQUEUE_MODE
endif

ifeq ($(CONTROL_FPGA_DEV_INIT),1)
CFLAGS += -DCONTROL_FPGA_DEV_INIT
endif
ifeq ($(CONTROL_FPGA_FUNC_INIT),1)
CFLAGS += -DCONTROL_FPGA_FUNC_INIT
endif
ifeq ($(CONTROL_FPGA_ENQUEUE_LLDMA_INIT),1)
CFLAGS += -DCONTROL_FPGA_ENQUEUE_LLDMA_INIT
endif
ifeq ($(CONTROL_FPGA_CHAIN_CONNECT),1)
CFLAGS += -DCONTROL_FPGA_CHAIN_CONNECT
endif
ifeq ($(CONTROL_FPGA_FINISH),1)
CFLAGS += -DCONTROL_FPGA_FINISH
endif
ifeq ($(CONTROL_FPGA_FUNC_FINISH),1)
CFLAGS += -DCONTROL_FPGA_FUNC_FINISH
endif
ifeq ($(CONTROL_FPGA_SHMEM_FINISH),1)
CFLAGS += -DCONTROL_FPGA_SHMEM_FINISH
endif

ifeq ($(DEBUG_PRINT),1)
CFLAGS += -DDEBUG_PRINT
endif
ifeq ($(DEBUG_PRINT_FPGA_OUT_HEADER),1)
CFLAGS += -DDEBUG_PRINT_FPGA_OUT_HEADER
endif
ifeq ($(DEBUG_FPGA_OUT_IMAGE_TO_MP4),1)
CFLAGS += -DDEBUG_FPGA_OUT_IMAGE_TO_MP4
endif

# ====================================================
# DPDK 
# ====================================================
# Build using pkg-config variables if possible
ifneq ($(shell $(PKGCONF) --exists libdpdk && echo 0),0)
$(error "no installation of DPDK found")
endif

.PHONY: all
.PHONY: shared static
all: clean shared
shared: build/$(APP)-shared
	ln -sf $(APP)-shared build/$(APP)
static: build/$(APP)-static
	ln -sf $(APP)-static build/$(APP)

PC_FILE := $(shell $(PKGCONF) --path libdpdk 2>/dev/null)
CFLAGS += $(shell $(PKGCONF) --cflags libdpdk)
# Add flag to allow experimental API as l2fwd uses rte_ethdev_set_ptype API
CFLAGS += -DALLOW_EXPERIMENTAL_API
LDFLAGS_SHARED = $(shell $(PKGCONF) --libs libdpdk)
LDFLAGS_STATIC = $(shell $(PKGCONF) --static --libs libdpdk)

ifeq ($(MAKECMDGOALS),static)
# check for broken pkg-config
ifeq ($(shell echo $(LDFLAGS_STATIC) | grep 'whole-archive.*l:lib.*no-whole-archive'),)
$(warning "pkg-config output list does not contain drivers between 'whole-archive'/'no-whole-archive' flags.")
$(error "Cannot generate statically-linked binaries with this version of pkg-config")
endif
endif


# ====================================================
# COMMAND
# ====================================================

# delete command
.PHONY: clean
clean:
	@rm -f build/$(APP) build/$(APP)-static build/$(APP)-shared
	@test -d build && rmdir -p build && echo "clean app" || true


build/$(APP)-shared: $(SRCS) Makefile $(PC_FILE) $(LIBFPGA) | build
	@$(CXX) $(CFLAGS) $(SRCS) $(LIBFPGA) $(PKGCONFOCV) -o $@ $(LDFLAGS) $(LDFLAGS_SHARED)

build/$(APP)-static: $(SRCS) Makefile $(PC_FILE) $(LIBFPGA) | build
	@$(CXX) $(CFLAGS) $(SRCS) $(LIBFPGA) $(PKGCONFOCV) -o $@ $(LDFLAGS) $(LDFLAGS_STATIC)

build:
	@mkdir -p $@
		
	
	
