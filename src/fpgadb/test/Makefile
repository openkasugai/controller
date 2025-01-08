#=================================================
# Copyright 2024 NTT Corporation, FUJITSU LIMITED
#=================================================

src := main.c
app := test

LIBFPGA_LIB_DIR ?= /usr/local/lib/fpgalib
LIBFPGA_INC_DIR ?= /usr/local/include/fpgalib
LIBFPGADB_LIB_DIR ?= ../build
LIBFPGADB_INC_DIR ?= ../build/include

# Setting for libfpga
CFLAGS  += -I$(LIBFPGA_INC_DIR)
LDFLAGS += -L$(LIBFPGA_LIB_DIR) -lfpga -lpciaccess -ldl -rdynamic -lstdc++

# Setting for libfpgadb
CFLAGS  += -I$(LIBFPGADB_INC_DIR)
LDFLAGS += -L$(LIBFPGADB_LIB_DIR) -lfpgadb

.PHONY: build clean all

all: clean build

build: $(app)
$(app): $(src)
	gcc $^ -o $@ $(CFLAGS) $(LDFLAGS)

clean:
	rm -f $(app)
