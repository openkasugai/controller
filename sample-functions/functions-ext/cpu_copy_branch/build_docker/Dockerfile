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
	libpci-dev

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
WORKDIR /opt/openkasugai-controller/src/submodules/fpga-software/lib
RUN make mcap
RUN make json

# OpneCV
RUN apt install g++
RUN apt install cmake
RUN apt install -y ffmpeg

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Tokyo

## Application
# Build Copy_branch Application
WORKDIR /opt/openkasugai-controller/sample-functions/functions-ext/cpu_copy_branch/app
RUN g++ -o copy_branch copy_branch.cpp
WORKDIR /opt/openkasugai-controller/sample-functions/functions-ext/cpu_copy_branch
RUN ln -s app/copy_branch .
