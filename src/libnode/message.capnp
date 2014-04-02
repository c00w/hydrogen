
using Go = import "../github.com/jmckaskill/go-capnproto/go.capnp";
$Go.package("libnode");
$Go.import("libnode");

@0xf5151606f45c93ef;

struct Signature {
    key @0 : Data;
    signature @1 : Data;
}

struct Time {
    #UTC
    seconds @0: UInt64;
    nanoSeconds @1: UInt32;
}

struct Change {

    type @0 : UInt8;
    account @1: Data;
    authorization @2: List(Signature);
    created @3: Time;

    newValue @4: Data;
    # Is the new value for everything but transactions,
    # Is the amount transfered for transactions.

    destination @5: Data;
    # Only Used for transactions

}

struct Message {
    votes @0: List(Change);
    time @1: Time;
    signature @2: List(Signature);
}
