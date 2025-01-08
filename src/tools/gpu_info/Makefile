.PHONY: all
all: build run

.PHONY: build
build:
	go mod tidy
	go fmt ./...
	go build .

.PHONY: debug
debug:
	go mod tidy
	go fmt ./...
	go build -x .

.PHONY: run
run:
	sudo ./dcgm-gpu-info
