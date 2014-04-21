
using Go = import "../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("util");
$Go.import("util");

@0xa6c896598a08b1dd;

struct Time {
    seconds @0: UInt64;
    nanoSeconds @1: UInt32;
}
