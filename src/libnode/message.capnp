
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

struct TransactionChange {
    destination @0: Data;
    amount @1: UInt64;
}

struct LocationChange {
    location @0: Text;
}

struct KeyChange {
    newkeys @0: List(Data);
}

struct Change {

    account @0: Data;
    authorization @1: List(Signature);
    created @2: Time;

    type :union {
         transaction @3 :TransactionChange;
         location @4 :LocationChange;
         key @5 :KeyChange;
    }
}

struct Vote {
    votes @0: List(Change);
    time @1: Time;
    signature @2: List(Signature);
}

struct Message {
    union {
        vote @0: Vote;
        change @1: Change;
    }
}

