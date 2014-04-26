
using Go = import "../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("util");
$Go.import("util");

@0xa6c896598a08b1dd;

struct Time {
    seconds @0: UInt64;
    nanoSeconds @1: UInt32;
}

struct P521Key {
    d @0: Data;
    x @1: Data;
    y @2: Data;
}
