#=================================================
# Copyright 2024 NTT Corporation, FUJITSU LIMITED
#=================================================

# ====================================================
# User Definition
# ====================================================
# target directory
BUILD_DIR ?= build

# Set the direcotry of header files for libfpga
LIBFPGA_INC_DIR ?= /usr/local/include/fpgalib

THIRD_PARTY_DIR := 3rdparty
PARSON_DIR      := $(THIRD_PARTY_DIR)/parson
ENABLE_INCLUDE_PARSON ?= 0

# suppress command line text, if you need, remove '@'
CC    := @$(CC)
AR    := @$(AR)

# ====================================================
# Library Definition
# ====================================================
# library name
LIB := libfpgadb

# v(Major.Minor.Revision-Patch)
LIBRARY_NAME      := $(LIB)
LIBRARY_MAJOR_VER :=1
LIBRARY_MINOR_VER :=0
LIBRARY_REVISION  :=0
LIBRARY_PATCH     :=0
LIBRARY_VERSION   := v.$(LIBRARY_MAJOR_VER).$(LIBRARY_MINOR_VER).$(LIBRARY_REVISION)-$(LIBRARY_PATCH)

# ====================================================
# Source files
# ====================================================
# Set directory for source files
SRC_DIR := src
INC_DIR := include

SRCS := $(wildcard $(SRC_DIR)/*.c)
ifeq ($(ENABLE_INCLUDE_PARSON),1)
SRCS-3rdpaty := $(PARSON_DIR)/parson.c
endif
SRCS += $(SRCS-3rdpaty)

OBJS := $(addsuffix .o,$(SRCS))

HEADERS := $(wildcard $(INC_DIR)/*.h)
HEADERS-3rdpaty := $(PARSON_DIR)/parson.h

LIB_STATIC := $(BUILD_DIR)/$(LIBRARY_NAME).a
LIB_SHARED := $(BUILD_DIR)/$(LIBRARY_NAME).so

# ====================================================
# setting CFLAGS
# ====================================================
CFLAGS += -DLIBRARY_NAME=\"$(LIBRARY_NAME)\"
CFLAGS += -DLIBRARY_MAJOR_VER=$(LIBRARY_MAJOR_VER)
CFLAGS += -DLIBRARY_MINOR_VER=$(LIBRARY_MINOR_VER)
CFLAGS += -DLIBRARY_REVISION=$(LIBRARY_REVISION)
CFLAGS += -DLIBRARY_PATCH=$(LIBRARY_PATCH)
CFLAGS += -O3
CFLAGS += -Wall
CFLAGS += -I$(INC_DIR)
CFLAGS += -I$(LIBFPGA_INC_DIR)
CFLAGS += -I$(PARSON_DIR)

# ====================================================
# COMMAND
# ====================================================
.PHONY: all default static shared clean

default: static
all: clean default
static: $(LIB_STATIC)
shared: $(LIB_SHARED)
clean:
	@rm -f $(OBJS) $(LIB_STATIC) $(LIB_SHARED)
	@rm -rf $(BUILD_DIR)/include
	@test -d $(BUILD_DIR) && rmdir $(BUILD_DIR) || true
	@echo " * Clean library"
	@rm -rf $(PARSON_DIR)
	@echo " * Clean 3rdparty library"

# ====================================================
# Definitions
# ====================================================
$(OBJS): $(SRCS)
$(SRCS): $(HEADERS) | $(HEADERS-3rdpaty)

%.c.o: %.c
	$(CC) -c $(CFLAGS) -o $@ $<

$(LIB_STATIC): $(OBJS) | $(BUILD_DIR)
	$(AR) rcs $@ $^
	@cp -RT $(INC_DIR) $(BUILD_DIR)/include
	@echo " * Build static library : $@[$(LIBRARY_VERSION)]"

$(LIB_SHARED): $(SRCS) | $(BUILD_DIR)
	$(CC) -shared $(CFLAGS) -fPIC -o $@ $^
	@cp -RT $(INC_DIR) $(BUILD_DIR)/include
	@echo " * Build shared library : $@[$(LIBRARY_VERSION)]"

$(BUILD_DIR):
	@mkdir -p $@

$(PARSON_DIR)/parson.h: $(PARSON_DIR)
$(PARSON_DIR)/parson.c: $(PARSON_DIR)
$(PARSON_DIR):
	@git clone https://github.com/kgabis/parson.git $@
