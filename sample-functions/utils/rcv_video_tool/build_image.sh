#!/bin/bash

# Copyright 2024 NTT Corporation, FUJITSU LIMITED


#PROXY_SERVER="http://<user>:<pass>@<proxy_address>:8080"

sudo buildah bud \
  # --build-arg http_proxy=$PROXY_SERVER \
  # --build-arg https_proxy=$PROXY_SERVER \
  # --build-arg HTTP_PROXY=$PROXY_SERVER \
  # --build-arg HTTPS_PROXY=$PROXY_SERVER \
  -t rcv_video_tool:latest ./
