
# CPU filter/resize processing module [cpu_filter_resize]

## Overview

* Filter/resize processing module for CPU
* Can be used as functions that make up dataflows
* Works as follows:
  * Video streaming server -> CPU decode -> TCP -> CPU filter/resize -> TCP
* Only supports TCP payload format of frame data (frame header + data frame)

## Requirement for software

software | verified version |
------------|-------------|
 Python | 3.8.10, 3.10.12 |
 opencv-python | 4.9.0.80 |

## Usage

### Command

```sh
$ python fr.py --help

NAME
    fr.py

SYNOPSIS
    fr.py <flags>

FLAGS
    --in_port=IN_PORT
        Default: 8888
    --out_addr=OUT_ADDR
        Default: '127.0.0.1'
    --out_port=OUT_PORT
        Default: 9999
    --in_width=IN_WIDTH
        Default: 3840
    --in_height=IN_HEIGHT
        Default: 2160
    --out_width=OUT_WIDTH
        Default: 1280
    --out_height=OUT_HEIGHT
        Default: 1280
    -s, --sockbufsize=SOCKBUFSIZE
        Default: 4194304
    -l, --loglevel=LOGLEVEL
        Default: 'INFO'
```

#### Example
```sh
$ python fr.py --in_port=8888 --out_port=9999 --out_addr=10.38.119.20 --in_width=1920 --in_height=1080 --out_width=960 --out_height=540 
```

Changing log level to DEBUG
```sh
$ python fr.py --in_port=8888 --out_port=9999 --out_addr=10.38.119.20 --in_width=1920 --in_height=1080 --out_width=960 --out_height=540  -l DEBUG 
```

## Run CPU filter/resize on K8s or host

### Prerequisites

#### Changing kernel parameters

When receiving 4K streams, net.core.rmem_max should be set to maximum of net.ipv4.tcp_rmem.
```sh
$ sysctl net.core.rmem_default
net.core.rmem_default = 212992
$ sysctl net.core.rmem_max
net.core.rmem_max = 212992
$ sysctl net.ipv4.tcp_rmem
net.ipv4.tcp_rmem = 1024        131072  6291456
$ sudo sysctl -w net.core.rmem_max=6291456
net.core.rmem_max = 6291456
$ sysctl net.core.rmem_max
net.core.rmem_max = 6291456
$ sudo vi /etc/sysctl.conf
  (add following line to persist configuration changes)
   net.core.rmem_max = 6291456
```

####  Checking out the source

```sh
$ git clone https://github.com/openkasugai/controller.git
$ cd controller/
$ git config -f .gitmodules submodule.src/submodules/fpga-software.url https://github.com/openkasugai/hardware-drivers.git
$ git submodule sync
$ git submodule update --init --recursive
```

### 1. Run CPU filter/resize on K8s

#### Build container image

Build container image on server where pod is to be deployed
```sh
$ cd controller/sample-functions/functions/cpu_filter_resize
$ sudo buildah bud -t cpu_filter_resize:1.0.0 -f containers/cpu/Dockerfile
```
You can also download container image from ghcr.
```sh
$ sudo buildah pull ghcr.io/openkasugai/controller/cpu_filter_resize:1.0.0
$ sudo buildah tag ghcr.io/openkasugai/controller/cpu_filter_resize:1.0.0 localhost/cpu_filter_resize:1.0.0
```

#### Edit and apply Pod manifest

Change environment variables and args as needed
```sh
vi deploy/cpu_filter_resize.yaml
kubectl apply -f deploy/cpu_filter_resize.yaml
```

### Run CPU filter/resize on host

#### Install Python

Example below shows how to install Python v3.10.12.

```sh
$ git clone https://github.com/pyenv/pyenv.git ~/.pyenv
$ cd ~/.pyenv && src/configure && make -C src
$ echo 'export PYENV_ROOT="$HOME/.pyenv"' >> ~/.bashrc
$ echo 'command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.bashrc
$ echo 'eval "$(pyenv init -)"' >> ~/.bashrc
$ cat ~/.bashrc 
$ source ~/.bashrc
$ sudo apt update
$ sudo apt install build-essential libssl-dev zlib1g-dev libbz2-dev \
  libreadline-dev libsqlite3-dev curl libncursesw5-dev xz-utils \
  tk-dev libxml2-dev libxmlsec1-dev libffi-dev liblzma-dev
$ pyenv install 3.10.12
$ pyenv global 3.10.12 
$ pyenv versions
  system
* 3.10.12 (set by /home/ubuntu/.pyenv/version)
$ python -V
Python 3.10.12
```

#### Create Python virtual environment and install packages
```sh
$ cd controller/sample-functions/functions/cpu_filter_resize
$ python -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements.txt
```

#### Run command

Refer to "Usage"

###### Copyright 2024 NTT Corporation, FUJITSU LIMITED
