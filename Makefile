BIN_PATH := bin
BIN_NAME := weave

build:
	go build -o ${BIN_PATH}/${BIN_NAME} main.go

lint:
	golangci-lint run

reltest:
	goreleaser build --snapshot --clean

deps:
	go mod download
	go mod tidy

clean:
	rm -rf ${BIN_PATH}
