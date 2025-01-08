
# CPU copy branch processing module [cpu_copy_branch]

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
cd controller/sample-functions/functions-ext/cpu_copy_branch/app
g++ -o copy_branch copy_branch.cpp
```

### 2.Run application

`./copy_branch addres num_branches dst_addresses buf_size`

- Args:
    - `address`: IP:Port this application uses
    - `num_branches` : Number of branches
    - `dst_addresses` : List of IP:Port for all destinations
    - `buf_size` : Buffer size (Byte)

- Example:
    - `./copy_branch 0.0.0.0:8000 3 "127.0.0.1:8001,127.0.0.2:8002,127.0.0.3:8003" 1024`

## How to build container image

```sh
cp controller/sample-functions/functions-ext/cpu_copy_branch/build_docker/Dockerfile .
sudo buildah bud -t cpu_copy_branch:1.0.0 -f Dockerfile
```

Container image is not available on ghcr.

----
Copyright 2024 NTT Corporation , FUJITSU LIMITED
