## *****************************************************************
## * Copyright 2022 NTT Corporation, FUJITSU LIMITED
## *****************************************************************

FROM nvcr.io/nvidia/deepstream:7.0-gc-triton-devel

ARG SRC_PATH=/opt/nvidia/deepstream/deepstream-7.0/sample-functions

# install GStreamer
RUN ./user_additional_install.sh

# copy source code (fpga_depayloader)
RUN mkdir -p ${SRC_PATH}/functions
COPY functions/gpu_infer_tcp_plugins ${SRC_PATH}/functions/gpu_infer_tcp_plugins
RUN rm -rf ${SRC_PATH}/functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker
WORKDIR ${SRC_PATH}

# build fpga_depayloader
RUN cd functions/gpu_infer_tcp_plugins/fpga_depayloader/ && ./build_app.sh

ENV GST_PLUGIN_PATH=${SRC_PATH}/functions/gpu_infer_tcp_plugins/fpga_depayloader
ENV WIDTH=1280
ENV HEIGHT=1280

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

WORKDIR /opt

RUN git clone https://github.com/marcoslucianops/DeepStream-Yolo.git

WORKDIR /opt/DeepStream-Yolo

ENV CUDA_DEVICE_ORDER=PCI_BUS_ID

RUN CUDA_VER=12.2 make -C nvdsinfer_custom_impl_Yolo
RUN wget https://github.com/AlexeyAB/darknet/releases/download/darknet_yolo_v4_pre/yolov4-p6.weights
RUN wget https://raw.githubusercontent.com/AlexeyAB/darknet/master/cfg/yolov4-p6.cfg
# COPY yolov4-p6.cfg /opt/DeepStream-Yolo/yolov4-p6.cfg

COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/labels.txt /opt/DeepStream-Yolo/labels.txt

RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV3_tiny.txt config_infer_primary_yoloV3_tiny.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6_th020_025.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6_th020_040.txt
RUN cp /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV2.txt config_infer_primary_yoloV4_p6_th020_070.txt
COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/config_infer_primary_yolo.patch config_infer_primary_yolo.patch
RUN patch -p7 < config_infer_primary_yolo.patch
COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/generate_engine_file.sh generate_engine_file.sh
COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/find_gpu.sh find_gpu.sh
RUN ./generate_engine_file.sh A100 config_infer_primary_yoloV4_p6_th020_040.txt model_b1_gpu0_fp16.engine

WORKDIR /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo

RUN CUDA_VER=12.2 make -C nvdsinfer_custom_impl_Yolo
RUN wget https://raw.githubusercontent.com/pjreddie/darknet/master/cfg/yolov3-tiny.cfg
RUN wget https://pjreddie.com/media/files/yolov3-tiny.weights

RUN cp /opt/DeepStream-Yolo/config_infer_primary_yoloV3_tiny.txt /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/config_infer_primary_yoloV3_tiny.txt
COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/labels.txt /opt/nvidia/deepstream/deepstream-7.0/sources/objectDetector_Yolo/labels.txt
COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/generate_engine_file.sh generate_engine_file.sh
COPY functions/gpu_infer_tcp_plugins/fpga_depayloader/build_docker/gpu-deepstream/find_gpu.sh find_gpu.sh
RUN ./generate_engine_file.sh T4 config_infer_primary_yoloV3_tiny.txt model_b1_gpu0_int8.engine

WORKDIR /opt/DeepStream-Yolo

RUN apt install -y less
RUN apt install -y emacs
RUN apt install -y numactl
