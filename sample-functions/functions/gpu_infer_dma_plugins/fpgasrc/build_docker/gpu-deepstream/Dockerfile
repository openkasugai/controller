## *****************************************************************
## * Copyright 2024 NTT Corporation, FUJITSU LIMITED
## *****************************************************************

FROM nvcr.io/nvidia/deepstream:7.0-gc-triton-devel

ARG SRC_PATH=/opt/nvidia/deepstream/deepstream-7.0/openkasugai-controller
ARG LIB_PATH=${SRC_PATH}/src/submodules/fpga-software/lib
ARG FUNCTION_PATH=${SRC_PATH}/sample-functions/functions/gpu_infer_dma_plugins
ARG CUR_FUNCTION_PATH=controller/sample-functions/functions/gpu_infer_dma_plugins

# install GStreamer
RUN ./user_additional_install.sh

# copy source code (fpgasrc)

COPY controller ${SRC_PATH}
WORKDIR ${SRC_PATH}

# check build environment

RUN sample-functions/functions/gpu_infer_dma_plugins/fpgasrc/build_docker/gpu-deepstream/check_gpus.sh

# build libfpga

RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y python3-pip
RUN apt-get install -y pkg-config
RUN pip3 install meson ninja
RUN pip3 install pyelftools
RUN apt-get install -y libnuma-dev
RUN apt-get install -y udev libpci-dev libpciaccess-dev
RUN cd ${LIB_PATH} && make clean && make dpdk-uninstall && make dpdk && make mcap && make json && make shared

# build fpgasrc

RUN cd ${FUNCTION_PATH}/fpgasrc && make clean && make prepare && make

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:${LIB_PATH}/DPDK/dpdk/lib/x86_64-linux-gnu:${LIB_PATH}/build
ENV GST_PLUGIN_PATH=${FUNCTION_PATH}/fpgasrc

RUN apt update && apt install -y build-essential
RUN apt install -y m4 autoconf automake libtool git vim
RUN apt install -y tcpdump iproute2 iputils-ping net-tools telnet dnsutils

WORKDIR /opt
    
RUN git clone https://github.com/RidgeRun/gst-perf.git \
    && cd gst-perf/ \
    && git checkout 01df710 \
    && libtoolize -c \
    && ./autogen.sh \
    && ./configure --prefix /usr/ --libdir /usr/lib/x86_64-linux-gnu/ \
    && make && make install

RUN git clone https://github.com/marcoslucianops/DeepStream-Yolo.git

ENV CUDA_DEVICE_ORDER=PCI_BUS_ID

WORKDIR /opt/DeepStream-Yolo

RUN CUDA_VER=12.2 make -C nvdsinfer_custom_impl_Yolo
RUN wget https://github.com/AlexeyAB/darknet/releases/download/darknet_yolo_v4_pre/yolov4-p6.weights
RUN wget https://raw.githubusercontent.com/AlexeyAB/darknet/master/cfg/yolov4-p6.cfg
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/labels.txt labels.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV3_tiny.txt config_infer_primary_yoloV3_tiny.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6_th020_025.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6_th020_040.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6_th020_070.txt
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/config_infer_primary_yolo.patch config_infer_primary_yolo.patch
RUN patch -p7 < config_infer_primary_yolo.patch
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/generate_engine_file.sh generate_engine_file.sh
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/find_gpu.sh find_gpu.sh
RUN ./generate_engine_file.sh A100 config_infer_primary_yoloV4_p6_th020_040.txt model_b1_gpu0_fp16.engine

WORKDIR /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo

RUN CUDA_VER=12.2 make -C nvdsinfer_custom_impl_Yolo
RUN wget https://raw.githubusercontent.com/pjreddie/darknet/master/cfg/yolov3-tiny.cfg
RUN wget https://pjreddie.com/media/files/yolov3-tiny.weights
RUN cp /opt/DeepStream-Yolo/config_infer_primary_yoloV3_tiny.txt config_infer_primary_yoloV3_tiny.txt
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/labels.txt labels.txt
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/generate_engine_file.sh generate_engine_file.sh
COPY ${CUR_FUNCTION_PATH}/fpgasrc/build_docker/gpu-deepstream/find_gpu.sh find_gpu.sh
RUN ./generate_engine_file.sh T4 config_infer_primary_yoloV3_tiny.txt model_b1_gpu0_int8.engine

WORKDIR /opt/DeepStream-Yolo

RUN apt install -y less
RUN apt install -y emacs
RUN apt install -y numactl
