.PHONY: all
all: build run

.PHONY: build
build:
	go mod tidy
	go fmt ./...
	go vet ./...
	go build .

.PHONY: debug
debug:
	go mod tidy
	deepcopy-gen --v=9 ./pkg/reference/...
	go fmt ./...
	go vet ./...
	go build -x .

.PHONY: run
run:
	./InfoCollector
 
.PHONY: t1
t1: build t1-run

.PHONY: t1-run
t1-run:
	./InfoCollector
