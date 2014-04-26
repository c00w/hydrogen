
using Go = import "../../github.com/glycerine/go-capnproto/go.capnp";
$Go.package("message");
$Go.import("message");

using import "../../util/util.capnp".Time;

@0xf5151606f45c93ef;

struct Authorization {
    signatures @0: List(KeySignature);
}

struct Key {
    x @0: Data;
    y @1: Data;
}

struct Signature {
    r @0: Data;
    s @1: Data;
}

struct KeySignature{
    key @0: Key;
    signature @1: Signature;
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

struct DropChange {
    account @0: Data;
}

enum RateVote {
    constant @0;
    increase @1;
    decrease @2;
}

struct Change {

    authorization @0: Authorization;
    created @1: Time;

    type :union {
         transaction @2 :TransactionChange;
         location @3 :LocationChange;
         drop @4: DropChange;
    }
}

struct Vote {
    votes @0: List(Change);
    rate @1: RateVote;
    time @2: Time;
    authorization @3: Authorization;
}

struct Message {
    payload :union {
        vote @0: Vote;
        change @1: Change;
    }
    authChain @2: List(Authorization);
}
