.PHONY: build clean

CGO_ENABLED=0

all: build

build: clean
		@go build -o ssr-subscriber ./cmd/main.go

clean:
	@rm -f ./ssr-subscriber
