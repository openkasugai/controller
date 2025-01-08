# Build the manager binary
FROM golang:1.23 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY PCIeConnection/go.mod go.mod
COPY PCIeConnection/go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY PCIeConnection/cmd/main.go cmd/main.go
COPY PCIeConnection/api/ api/
COPY PCIeConnection/internal/controller/ internal/controller/
COPY whitebox-k8s-flowctrl/api/v1/common_types.go api/v1/
COPY WBFunction/api/v1/function_common_types.go api/v1/
COPY FPGAFunction/api/v1/childbs_types.go api/v1/
COPY FPGAFunction/api/v1/fpga_types.go api/v1/

# Copy openkasugai-hardware-drivers and build pkg
COPY PCIeConnection/openkasugai-hardware-drivers/ openkasugai-hardware-drivers/
WORKDIR /workspace/openkasugai-hardware-drivers
RUN apt-get update \
    && apt-get install -y build-essential \
        python3-pip pkg-config libnuma-dev vim
RUN pip3 install --break-system-packages meson ninja pyelftools
WORKDIR /workspace/openkasugai-hardware-drivers/lib/DPDK/
RUN wget http://fast.dpdk.org/rel/dpdk-23.11.1.tar.xz \
    && tar Jxvf dpdk-23.11.1.tar.xz
RUN bash build_dpdk.sh
WORKDIR /workspace/openkasugai-hardware-drivers/lib
RUN apt-get update \
    && apt-get install -y libpciaccess-dev
Run make mcap
Run make json
RUN make
WORKDIR /workspace
RUN mkdir /usr/local/include/fpgalib; mkdir /usr/local/lib/fpgalib \
    && cp -a openkasugai-hardware-drivers/lib/libfpga/include/*.h /usr/local/include/fpgalib/. \
    && cp -a openkasugai-hardware-drivers/lib/libptu/include/* /usr/local/include/fpgalib/. \
    && cp -a openkasugai-hardware-drivers/driver/xpcie_device.h /usr/local/include/fpgalib/. \
    && cp -a openkasugai-hardware-drivers/lib/build/libfpga.a /usr/local/lib/fpgalib/. \
    && cp -ra openkasugai-hardware-drivers/lib/libfpga/include/libfpga_internal/ /usr/local/lib/fpgalib/. \
    && ln -s /workspace/openkasugai-hardware-drivers/lib/DPDK/dpdk /usr/local/include/fpgalib/dpdk
ARG PKG_CONFIG_PATH

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
ENV CGO_CFLAGS_ALLOW=-mrtm
RUN go fmt ./...
RUN go vet ./...
RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM gcr.io/distroless/static:nonroot
#FROM gcr.io/distroless/base:nonroot
#WORKDIR /
#COPY --from=builder /workspace/manager .
USER root:root

ENV LD_LIBRARY_PATH "/workspace/openkasugai-hardware-drivers/lib/DPDK/dpdk/lib/x86_64-linux-gnu/"
ENTRYPOINT ["/workspace/manager"]
