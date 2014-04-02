.PHONY: capn all

GOPATH := $(CURDIR)
PATH := ${PATH}:${GOPATH}/bin
PLATFORM := linux_amd64

all: bin/hydrogend

src/libnode/message.capnp.go: src/libnode/message.capnp
	capnp compile -ogo src/libnode/message.capnp

bin/capnpc-go:
	go install github.com/jmckaskill/go-capnproto/capnpc-go

capn:
	echo "Do not execute this command unless you are purposefully updating the version"
	go get -u github.com/jmckaskill/go-capnproto
	go get -u github.com/jmckaskill/go-capnproto/capnpc-go

pkg/${PLATFORM}/libnode.a: src/libnode/*.go
	go install libnode

bin/hydrogend: pkg/${PLATFORM}/libnode.a src/hydrogend/*.go
	go install hydrogend



