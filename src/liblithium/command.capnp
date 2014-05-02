
using Go = import "../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("liblithium");
$Go.import("liblithium");

@0xf7ca1290d56d8eac;

struct GetBalance {
    account @0: Data;
}

struct SendMoney {
    to @0: Data;
    amount @1: UInt64;
}

struct Command {
    union {
        getbalance @0: GetBalance;
        sendmoney @1: SendMoney;
    }
}

struct Result {
    message @0: Text;
}
