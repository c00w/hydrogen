.PHONY: all test capnp clean dependencies

GOPATH := $(CURDIR)
PATH := ${PATH}:${GOPATH}/bin
PLATFORM := linux_amd64

pkgs = libnode libhydrogen/message libhydrogen libhelium hydrogend

all: capnp
	go fmt $(pkgs)
	go install $(pkgs)

test: all
	go test $(pkgs)

capnp: bin/capnpc-go
	capnp compile -ogo src/libhelium/account.capnp
	capnp compile -ogo src/libhydrogen/message/message.capnp
	capnp compile -ogo src/util/time.capnp

bin/capnpc-go:
	go install github.com/glycerine/go-capnproto/capnpc-go

dependencies:
	go get -u github.com/glycerine/go-capnproto
	go get -u github.com/glycerine/go-capnproto/capnpc-go

clean:
	rm -r bin
	rm -r pkg

