# Copyright 2024 NTT Corporation, FUJITSU LIMITED 

FROM ubuntu:22.04

RUN apt update -y && apt install -y build-essential
RUN apt install -y m4 autoconf automake libtool git vim
RUN apt install -y tcpdump iproute2 iputils-ping net-tools telnet dnsutils

# --------
#  gstreamer Install
# --------
RUN apt update
RUN apt -y install libgstreamer1.0-0
RUN apt -y install gstreamer1.0-plugins-base
RUN apt -y install gstreamer1.0-plugins-good
RUN apt -y install gstreamer1.0-plugins-bad
RUN apt -y install gstreamer1.0-tools
RUN apt -y install gstreamer1.0-x
RUN apt -y install gstreamer1.0-alsa
RUN apt -y install gstreamer1.0-gl
RUN apt -y install gstreamer1.0-gtk3
RUN apt -y install gstreamer1.0-qt5
RUN apt -y install gstreamer1.0-pulseaudio
RUN apt -y install libgstreamer1.0-dev
RUN apt -y install libgstreamer-plugins-base1.0-dev

WORKDIR /opt

RUN git clone https://github.com/RidgeRun/gst-perf.git \
    && cd gst-perf/ \
    && git checkout 01df710 \
    && libtoolize -c \
    && ./autogen.sh \
    && ./configure --prefix /usr/ --libdir /usr/lib/x86_64-linux-gnu/ \
    && make && make install
