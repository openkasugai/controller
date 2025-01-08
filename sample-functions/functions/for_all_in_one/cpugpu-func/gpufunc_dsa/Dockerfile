FROM localhost/gpu_infer_tcp:1.0.0 as builder

RUN wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
RUN tar xf go1.22.4.linux-amd64.tar.gz -C /usr/local/
ENV PATH $PATH:/usr/local/go/bin

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go

RUN go fmt .
RUN go vet .
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o wrapper main.go


FROM localhost/gpu_infer_tcp:1.0.0

WORKDIR /
COPY --from=builder /workspace/wrapper .
COPY config-dsa.yaml.tmpl .
