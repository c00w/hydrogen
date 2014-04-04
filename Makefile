.PHONY: capn all

GOPATH := $(CURDIR)
PATH := ${PATH}:${GOPATH}/bin
PLATFORM := linux_amd64

all: bin/hydrogend

test: pkg/${PLATFORM}/libnode.a
	go test libnode

src/libnode/message.capnp.go: src/libnode/message.capnp bin/capnpc-go
	capnp compile -ogo src/libnode/message.capnp

bin/capnpc-go:
	go install github.com/glycerine/go-capnproto/capnpc-go

capn:
	echo "Do not execute this command unless you are purposefully updating the version"
	go get -u github.com/glycerine/go-capnproto
	go get -u github.com/glycerine/go-capnproto/capnpc-go

pkg/${PLATFORM}/libnode.a: src/libnode/*.go src/libnode/message.capnp.go
	go install libnode

bin/hydrogend: pkg/${PLATFORM}/libnode.a src/hydrogend/*.go
	go install hydrogend



