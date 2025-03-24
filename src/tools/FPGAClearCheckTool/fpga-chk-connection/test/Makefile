#=================================================
# Copyright 2025 NTT Corporation , FUJITSU LIMITED
#=================================================

# APP name
APP_NAME      := setup_connection

# APP path
APP := $(APP_NAME)

# Set APP source files
SRC_DIR := ./src
SRCS    := $(wildcard $(SRC_DIR)/*.cpp)

# Set the path for libfpga
LIBFPGA_LIB_DIR ?= /usr/local/lib/fpgalib
LIBFPGA_INC_DIR ?= /usr/local/include/fpgalib
# Set the path for DPDK
DPDK_LIB_DIR ?= /usr/local/lib/fpgalib/dpdk/lib/x86_64-linux-gnu

# Set CFLAGS(Options)
CFLAGS += -O3
CFLAGS += -Wall
# Set CFLAGS(libfpga)
CFLAGS += -I$(LIBFPGA_INC_DIR) -DAPP_NAME=\"$(APP_NAME)\"

# Set LDFLAGS(libfpga)
LDFLAGS += -L$(LIBFPGA_LIB_DIR) -lfpga -lpciaccess -ldl -rdynamic -lstdc++
# Set LDFLAGS(DPDK)
DPDK_LIBS := librte_eal librte_mempool librte_log librte_telemetry librte_kvargs
LDFLAGS += -L$(DPDK_LIB_DIR) $(patsubst %,-l:%.a,$(DPDK_LIBS)) -lnuma

# Suppress command line text, if you need, remove '@'
CC := @$(CC)


# COMMAND
.PHONY: all default static clean

default: static
all: clean default
static: $(APP)

$(APP): $(SRCS) Makefile
	$(CC) $(CFLAGS) -o $@ $(SRCS) $(LDFLAGS)
	@echo " * Build App : $@"

clean:
	@rm -f $(APP)
	@echo " * Clean $(APP)"
