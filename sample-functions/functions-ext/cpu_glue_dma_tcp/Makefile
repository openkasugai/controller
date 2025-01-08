# Copyright 2024 NTT Corporation, FUJITSU LIMITED

# binary name
APP = glue

# build mode setting
USE_MODULE_FPGA := 0
USE_OLD_FUNC_API := 0

# all source are stored in SRCS-y
SRCS-y := src/main.c src/bcdbg.c src/bcdbg_thread.c src/common_func.c src/glue_func.c src/glue.c
INCLUDE := -I./include
LIBS := -lpci

# lib
LIBDIR=../../../src/submodules/fpga-software/lib

# libfpga
LIBFPGADIR = $(LIBDIR)/libfpga
INCLUDE += -I$(LIBFPGADIR)/include

###ifeq ($(USE_OLD_FUNC_API), 1)
###	CFLAGS += -DOLD_FUNC_API
###endif
###
###ifeq ($(USE_MODULE_FPGA), 1)
###	CFLAGS += -DMODULE_FPGA
###endif

# JSON
JSONDIR = $(LIBDIR)/JSON
SRCS-y += $(JSONDIR)/parson/parson.c
INCLUDE += -I$(JSONDIR)/parson

# DPDK
PKGCONF ?= env PKG_CONFIG_PATH=$(LIBDIR)/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config

# Build using pkg-config variables if possible
ifneq ($(shell $(PKGCONF) --exists libdpdk && echo 0),0)
$(error "no installation of DPDK found")
endif

all: shared
.PHONY: shared static
shared: build/$(APP)-shared
	ln -sf $(APP)-shared build/$(APP)
static: build/$(APP)-static
	ln -sf $(APP)-static build/$(APP)

PC_FILE := $(shell $(PKGCONF) --path libdpdk 2>/dev/null)
CFLAGS += -O3 $(shell $(PKGCONF) --cflags libdpdk)
# Add flag to allow experimental API as l2fwd uses rte_ethdev_set_ptype API
CFLAGS += -DALLOW_EXPERIMENTAL_API -I../include $(INCLUDE)

# libfpga
LDFLAGS += $(shell PKG_CONFIG_PATH=$(LIBDIR)/build/pkgconfig pkg-config --libs libfpga)
CFLAGS += $(shell PKG_CONFIG_PATH=$(LIBDIR)/build/pkgconfig pkg-config --cflags libfpga)

# DPDK
LDFLAGS_SHARED = $(LDFLAGS) $(shell PKG_CONFIG_PATH=$(LIBDIR)/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --libs libdpdk)
LDFLAGS_STATIC = $(LDFLAGS) $(shell PKG_CONFIG_PATH=$(LIBDIR)/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --static --libs libdpdk)
CFLAGS += $(shell PKG_CONFIG_PATH=$(LIBDIR)/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig pkg-config --cflags libdpdk)

ifeq ($(MAKECMDGOALS),static)
# check for broken pkg-config
ifeq ($(shell echo $(LDFLAGS_STATIC) | grep 'whole-archive.*l:lib.*no-whole-archive'),)
$(warning "pkg-config output list does not contain drivers between 'whole-archive'/'no-whole-archive' flags.")
$(error "Cannot generate statically-linked binaries with this version of pkg-config")
endif
endif

#build/$(APP)-shared: $(SRCS-y) $(OBJS) Makefile $(PC_FILE) | build
#	$(CC) $(CFLAGS) $(SRCS-y) $(OBJS) $(LIBS) -o $@ $(LDFLAGS) $(LDFLAGS_SHARED)
#
#build/$(APP)-static: $(SRCS-y) $(OBJS) Makefile $(PC_FILE) | build
#	$(CC) $(CFLAGS) $(SRCS-y) $(OBJS) $(LIBS) -o $@ $(LDFLAGS) $(LDFLAGS_STATIC)

build/$(APP)-shared: $(SRCS-y) Makefile $(PC_FILE) | build
	$(CC) $(CFLAGS) $(SRCS-y) $(LIBS) -o $@ $(LDFLAGS) $(LDFLAGS_SHARED)

build/$(APP)-static: $(SRCS-y) Makefile $(PC_FILE) | build
	$(CC) $(CFLAGS) $(SRCS-y) $(LIBS) -o $@ $(LDFLAGS) $(LDFLAGS_STATIC)

build:
	@mkdir -p $@

.PHONY: clean
clean:
	rm -f build/$(APP) build/$(APP)-static build/$(APP)-shared
	test -d build && rmdir -p build || true
