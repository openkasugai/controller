
# GPU inference module (DMA version) [gpu_infer_dma]

## Checking out the source

```sh
$ git clone https://github.com/openkasugai/controller.git
$ cd controller/
$ git config -f .gitmodules submodule.src/submodules/fpga-software.url https://github.com/openkasugai/hardware-drivers.git
$ git submodule sync
$ git submodule update --init --recursive
```
## How to build container image

```sh
cp controller/sample-functions/functions/gpu_infer_dma_plugins/fpgasrc/build_docker/gpu-deepstream/Dockerfile .
sudo buildah bud --runtime=/usr/bin/nvidia-container-runtime -t gpu_infer_dma:1.0.0 -f Dockerfile
```

Container image is not available on ghcr according to GStreamer licensing terms.

----
Copyright 2024 NTT Corporation , FUJITSU LIMITED
