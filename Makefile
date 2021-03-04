.PHONY: build clean

CGO_ENABLED=0

all: build

build: test clean
	@go build -o ssr-subscriber ./cmd

test: clean
	@go test ./...

clean:
	@go clean ./...
	@rm -f ./ssr-subscriber
