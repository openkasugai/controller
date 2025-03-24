.PHONY: all
all: build run

.PHONY: build
build:
	go mod tidy
	go fmt ./...
	go vet ./...
	go build .

.PHONY: run
run:
	./FPGACheckPerDF

.PHONY: t1-run
t1-run: build run
