#=================================================
# Copyright 2025 NTT Corporation , FUJITSU LIMITED
#=================================================

# APP name
APP_NAME      := fpga-chk-connection
# Version(Major.Minor.Revision-Patch)
APP_MAJOR_VER := 1
APP_MINOR_VER := 0
APP_REVISION  := 0
APP_PATCH     := 0
APP_VERSION   := v.$(APP_MAJOR_VER).$(APP_MINOR_VER).$(APP_REVISION)-$(APP_PATCH)

# Target directory
TARGET_DIR ?= bin
# APP path
APP := $(TARGET_DIR)/$(APP_NAME)

# Set APP source files
SRC_DIR := src
SRCS    := $(wildcard $(SRC_DIR)/*.c)
# Set APP include files
INC_DIR := include
HEADERS := $(wildcard $(INC_DIR)/*.h)

# Set the path for libfpga
LIBFPGA_LIB_DIR ?= /usr/local/lib/fpgalib
LIBFPGA_INC_DIR ?= /usr/local/include/fpgalib
# Set the path for DPDK
DPDK_LIB_DIR ?= /usr/local/lib/fpgalib/dpdk/lib/x86_64-linux-gnu
# Set the path for parson
PARSON_DIR ?= 3rdparty/parson
PARSON_SRCS ?= $(PARSON_DIR)/parson.c

# Set CFLAGS(APP definition)
CFLAGS += -DAPP_NAME=\"$(APP_NAME)\"
CFLAGS += -DAPP_MAJOR_VER=$(APP_MAJOR_VER)
CFLAGS += -DAPP_MINOR_VER=$(APP_MINOR_VER)
CFLAGS += -DAPP_REVISION=$(APP_REVISION)
CFLAGS += -DAPP_PATCH=$(APP_PATCH)
# Set CFLAGS(Options)
CFLAGS += -I$(INC_DIR)
CFLAGS += -O3
CFLAGS += -Wall
# Set CFLAGS(libfpga)
CFLAGS += -I$(LIBFPGA_INC_DIR) -DSUPPORT_OTHER_THAN_LIBFPGA_API
# Set CFLAGS(parson)
CFLAGS += -I$(PARSON_DIR)

# Set LDFLAGS(libfpga)
LDFLAGS += -L$(LIBFPGA_LIB_DIR) -lfpga -lpciaccess -rdynamic -lstdc++
# Set LDFLAGS(DPDK)
DPDK_LIBS := librte_eal librte_mempool librte_log librte_telemetry librte_kvargs
LDFLAGS += -L$(DPDK_LIB_DIR) $(patsubst %,-l:%.a,$(DPDK_LIBS)) -lnuma

# Suppress command line text, if you need, remove '@'
CC := @$(CC)


# COMMAND
.PHONY: all default static shared clean remove-parson clean-all

default: static
all: clean default
static: $(APP)

$(APP): $(SRCS) $(PARSON_SRCS) $(HEADERS) Makefile | $(TARGET_DIR)
	$(CC) $(CFLAGS) -o $@ $(SRCS) $(PARSON_SRCS) $(LDFLAGS)
	@echo " * Build App : $@[$(APP_VERSION)]"
	@ln -sf $@ $(APP_NAME)
	@echo " * Create link : $(APP_NAME)"

clean:
	@rm -f $(APP_NAME)
	@echo " * Clean $(APP_NAME)"
	@rm -f $(APP)
	@echo " * Clean $(APP)"
	@test -d $(TARGET_DIR) && rmdir $(TARGET_DIR) || true
	@echo " * Clean $(TARGET_DIR)"

remove-parson:
	@rm -rf $(PARSON_DIR)
	@echo " * Clean $(PARSON_DIR)"

clean-all: clean remove-parson

$(TARGET_DIR):
	@mkdir -p $@

$(PARSON_SRCS): | $(PARSON_DIR)
$(PARSON_DIR):
	@git clone https://github.com/kgabis/parson.git $@
