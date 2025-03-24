.PHONY: all
all: build run

.PHONY: build
build:
	mkdir -p log
	go mod tidy
	go fmt ./...
	go vet ./...
	go build .

.PHONY: run
run:
	./FPGAReconfigurationTool

.PHONY: t1-run
t1-run: build run
