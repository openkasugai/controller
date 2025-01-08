## *****************************************************************
## * Copyright 2024 NTT Corporation, FUJITSU LIMITED
## *****************************************************************

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
	build-essential \
	software-properties-common \
	cmake \
	zip \
	git \
	wget \
	python3-pip \
	pkg-config \
	libnuma-dev \
	gcc \
	cppcheck \
	libboost-dev \
	vim \
	gdb \
	libpci-dev \
	libpciaccess-dev

RUN pip install \
	gcovr \
	cpplint \
	meson \
	ninja \
	pyelftools

# Copy openkasugai-controller
COPY controller /opt/openkasugai-controller

## Advance preparation
# Hugepage

# FPGA Library
WORKDIR /opt/openkasugai-controller/src/submodules/fpga-software/lib/
RUN make clean && make dpdk-uninstall && make dpdk && make mcap && make json && make shared

# OpneCV
RUN apt install g++
RUN apt install cmake
RUN apt install -y ffmpeg

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Tokyo

# Driver

## Application
# Build Glue Application
WORKDIR /opt/openkasugai-controller/sample-functions/functions-ext/cpu_glue_dma_tcp
RUN make
RUN ln -s build/glue .
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/openkasugai-controller/src/submodules/fpga-software/lib/DPDK/dpdk/lib/x86_64-linux-gnu:/opt/openkasugai-controller/src/submodules/fpga-software/lib/build/
