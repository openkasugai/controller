FROM ubuntu:22.04

ARG CPU_DECODE_DIR=sample-functions/functions/cpu_decode/
ARG FPGA_SOFTWARE_DIR=src/submodules/fpga-software/

# copy source code (cpu_decode)
COPY ${CPU_DECODE_DIR} /root/${CPU_DECODE_DIR}

# copy source code (FPGA library)
COPY ${FPGA_SOFTWARE_DIR} /root/${FPGA_SOFTWARE_DIR}

# set workdir
WORKDIR /root/

# tool & library
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y \
      wget \
      build-essential \
      python3-pip \
      pkg-config \
      libnuma-dev \
      cmake \
      unzip \
      git \
      libpciaccess-dev

# build & install DPDK
RUN cd ${FPGA_SOFTWARE_DIR}/lib/DPDK && \
    rm -rf ./dpdk* && \
    wget http://fast.dpdk.org/rel/dpdk-23.11.1.tar.xz && \
    tar Jxvf dpdk-23.11.1.tar.xz && \
    pip3 install meson ninja pyelftools && \
    ./build_dpdk.sh && \
    echo /root/${FPGA_SOFTWARE_DIR}/lib/DPDK/dpdk/lib/x86_64-linux-gnu > /etc/ld.so.conf.d/dpdk.conf && \
    ldconfig

# build libfpga
RUN cd ${FPGA_SOFTWARE_DIR}/lib && \
    make mcap-uninstall && \
    make clean && \
    make mcap && \
    make json && \
    make

# build & install OpenCV
RUN cd ${CPU_DECODE_DIR}/docker && \
    bash install_opencv_for_container.sh && \
    apt-get -y clean && \
    rm -rf /var/lib/apt/lists/*

# build cpu_decode
COPY ${CPU_DECODE_DIR}/docker/Makefile ${CPU_DECODE_DIR}
RUN cd ${CPU_DECODE_DIR} && make

