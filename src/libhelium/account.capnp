using Go = import "../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("libhelium");
$Go.import("libhelium");

@0xaec058e4daabb567;

struct Ledger {
    accounts @0: List(Account);
    tau @1: Int64;
}

struct Account {
    key @0: Data;
    location @1: Text;
    balance @2: UInt64;
}
