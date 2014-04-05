
using Go = import "../../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("message");
$Go.import("message");

@0xf5151606f45c93ef;

struct Authorization {
    account @0: Data;
    signatures @1: List(KeySignature);
}

struct KeySignature{
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

enum RateVote {
    constant @0;
    increase @1;
    decrease @2;
}

struct RateChange {
    vote @0: RateVote;
}

struct Change {

    authorization @0: KeySignature;
    created @1: Time;

    type :union {
         transaction @2 :TransactionChange;
         location @3 :LocationChange;
         key @4: KeyChange;
         drop @5: DropChange;
         time @6: RateChange;

    }
}

struct Vote {
    votes @0: List(Change);
    time @1: Time;
    authorization @2: KeySignature;
}

struct Message {
    payload :union {
        vote @0: Vote;
        change @1: Change;
    }
    authchain @2: List(KeySignature);
}

