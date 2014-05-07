
using Go = import "../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("libberyllium");
$Go.import("libberyllium");

@0x93c1f39e065bcc96;

struct Request {
    account @0: Data;
}

struct Response {
    ok @0: Bool;
}
