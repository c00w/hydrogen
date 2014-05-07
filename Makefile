.PHONY: all test capnp clean dependencies

GOPATH := $(CURDIR)
PATH := ${PATH}:${GOPATH}/bin
PLATFORM := linux_amd64

pkgs = libnode libhydrogen/message libhydrogen libhelium hydrogend util \
	   liblithium hydrogenctl libberyllium

capnp = src/libhelium/account.capnp \
	src/libhydrogen/message/message.capnp \
	src/util/util.capnp \
	src/liblithium/command.capnp \
	src/libberyllium/command.capnp \

all: capnp
	go fmt $(pkgs)
	go install $(pkgs)

test: all
	go test $(pkgs)

capnp: bin/capnpc-go
	capnp compile -ogo ${capnp}

bin/capnpc-go:
	go install github.com/glycerine/go-capnproto/capnpc-go

dependencies:
	go get -u github.com/glycerine/go-capnproto
	go get -u github.com/glycerine/go-capnproto/capnpc-go

clean:
	rm -r bin
	rm -r pkg

