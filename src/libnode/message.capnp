
using Go = import "../github.com/jmckaskill/go-capnproto/go.capnp";
$Go.package("libnode");
$Go.import("libnode");

@0xf5151606f45c93ef;

struct Signature {
    key @0: Data;
    signature @1: Data;
}

struct Time {
    #UTC
    seconds @0: UInt64;
    nanoSeconds @1: UInt32;
}

struct TransactionChange {
    source @0: Data;
    destination @1: Data;
    amount @2: UInt64;
}

struct LocationChange {
    account @0: Data;
    location @1: Text;
}

struct KeyChange {
    account @0: Data;
    newkeys @1: List(Data);
}

struct DropChange {
    account @0: Data;
}

enum TimeVote {
    constant @0;
    increase @1;
    decrease @2;
}

struct TimeChange {
    vote @0: TimeVote;
}

struct Change {

    authorization @0: List(Signature);
    created @1: Time;

    type :union {
         transaction @2 :TransactionChange;
         location @3 :LocationChange;
         key @4: KeyChange;
         drop @5: DropChange;
         time @6: TimeChange;

    }
}

struct Vote {
    votes @0: List(Change);
    time @1: Time;
    authorization @2: List(Signature);
}

struct Message {
    union {
        vote @0: Vote;
        change @1: Change;
    }
}

