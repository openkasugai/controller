## *****************************************************************
## * Copyright 2024 NTT Corporation, FUJITSU LIMITED
## *****************************************************************

FROM ubuntu:22.04

RUN apt update && apt install -y build-essential
RUN apt install -y wget git vim tcpdump iproute2 iputils-ping net-tools telnet dnsutils
RUN apt install -y python3.10-dev python3-pip
RUN update-alternatives --install /usr/bin/python python /usr/bin/python3.10 1

WORKDIR /app

COPY *.py /app/
COPY requirements.txt /app/ 
RUN pip install -r /app/requirements.txt
