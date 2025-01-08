FROM golang:1.23 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY api api
COPY main.go main.go

RUN go fmt .
RUN go vet .
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o calcapp main.go

# FROM alpine:3.20.0
FROM ubuntu:24.04
WORKDIR /
COPY --from=builder /workspace/calcapp .
COPY config.yaml.tmpl .
