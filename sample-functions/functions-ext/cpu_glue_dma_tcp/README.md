
# CPU glue processing module [cpu_glue_dma_tcp]

## Checking out the source

```sh
$ git clone https://github.com/openkasugai/controller.git
$ cd controller/
$ git config -f .gitmodules submodule.src/submodules/fpga-software.url https://github.com/openkasugai/hardware-drivers.git
$ git submodule sync
$ git submodule update --init --recursive
```

## How to run application

### 1.Build application

```sh
cd controller/sample-functions/functions-ext/cpu_glue_dma_tcp
make
```

### 2.Run application

`./build/glue dst_address width height`
- Args:
    - `dst_address` : IP:Port to send data to
    - `width` : width of frame to be transferred
    - `height` : height of frame to be transferred

- Environment variables
    - `file prefix` : used by DPDK
      - name : GLUEENV_DPDK_FILE_PREFIX
    - `device name` : name of FPGA device
      - name : GLUEENV_FPGA_DEV_NAME
    - `connector ID` : used by FPGA to setup DMA
      - name : GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID

- Example:
    - `./build/glue 127.0.0.1:8001 1280 1280`

## How to build container image

```sh
cp controller/sample-functions/functions-ext/cpu_glue_dma_tcp/build_docker/Dockerfile .
sudo buildah bud -t cpu_glue_dma_tcp:1.0.0 -f Dockerfile
```

Container image is not available on ghcr.

----
Copyright 2024 NTT Corporation , FUJITSU LIMITED
