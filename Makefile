.PHONY: capn all

GOPATH := $(CURDIR)
PATH := ${PATH}:${GOPATH}/bin
PLATFORM := linux_amd64

all: bin/hydrogend

test: pkg/${PLATFORM}/libnode.a pkg/${PLATFORM}/libhydrogen.a
	go test libnode
	go test libhydrogen

src/libhydrogen/message.capnp.go: src/libhydrogen/message.capnp bin/capnpc-go
	capnp compile -ogo src/libhydrogen/message.capnp

bin/capnpc-go:
	go install github.com/glycerine/go-capnproto/capnpc-go

dependencies:
	go get -u github.com/glycerine/go-capnproto
	go get -u github.com/glycerine/go-capnproto/capnpc-go

pkg/${PLATFORM}/libnode.a: src/libnode/*.go
	go fmt libnode
	go install libnode

pkg/${PLATFORM}/libhydrogen.a: src/libhydrogen/*.go src/libhydrogen/message.capnp.go
	go fmt libhydrogen
	go install libhydrogen

bin/hydrogend: pkg/${PLATFORM}/libnode.a pkg/${PLATFORM}/libhydrogen.a src/hydrogend/*.go
	go fmt hydrogend
	go install hydrogend



