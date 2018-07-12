# Go parameters

GOPATH ?= $(HOME)/go

#This is how we want to name the binary output
BINARY=go-loginAuth

all: test build

build:
	cd $(GOPATH)/src; go install github.com/puneethreddy20/go-loginAuth

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

deps:
	go get -t ./...