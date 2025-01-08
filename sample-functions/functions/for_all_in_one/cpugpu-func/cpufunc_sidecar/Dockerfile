FROM golang:1.22 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go

RUN go fmt .
RUN go vet .
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o sidecar main.go

# FROM alpine:3.20.0
FROM ubuntu:24.04
WORKDIR /
COPY --from=builder /workspace/sidecar .
