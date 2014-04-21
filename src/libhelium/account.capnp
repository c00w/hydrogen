using Go = import "../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("libhelium");
$Go.import("libhelium");

using import "../util/time.capnp".Time;

@0xaec058e4daabb567;

struct Ledger {
    accounts @0: List(Account);
    created @1: Time;
    tau @2: Int64;
}

struct Account {
    key @0: Data;
    location @1: Text;
    balance @2: UInt64;
}
