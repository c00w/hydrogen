.PHONY: capn all

GOPATH := $(CURDIR)
PATH := ${PATH}:${GOPATH}/bin
PLATFORM := linux_amd64

all: bin/hydrogend

test: all
	go test libnode -race
	go test libhydrogen -race

src/libhydrogen/message/message.capnp.go: src/libhydrogen/message/message.capnp bin/capnpc-go
	capnp compile -ogo src/libhydrogen/message/message.capnp

bin/capnpc-go:
	go install github.com/glycerine/go-capnproto/capnpc-go

dependencies:
	go get -u github.com/glycerine/go-capnproto
	go get -u github.com/glycerine/go-capnproto/capnpc-go

pkg/${PLATFORM}/libnode.a: src/libnode/*.go
	go fmt libnode
	go install libnode

pkg/${PLATFORM}/libhydrogen/message.a: src/libhydrogen/message/message.capnp.go src/libhydrogen/message/*.go
	go fmt libhydrogen/message
	go install libhydrogen/message

pkg/${PLATFORM}/libhydrogen.a: pkg/${PLATFORM}/libnode.a pkg/${PLATFORM}/libhydrogen/message.a
	go fmt libhydrogen
	go install libhydrogen

bin/hydrogend: pkg/${PLATFORM}/libnode.a pkg/${PLATFORM}/libhydrogen.a src/hydrogend/*.go
	go fmt hydrogend
	go install hydrogend

clean:
	rm -r bin
	rm -r pkg



