package libhydrogen

import (
    "testing"
    "time"

    "libhydrogen/message"
    "libnode"
    "util"
)

func TestHydrogen(t *testing.T) {
    key1 := util.GenKey()
    key2 := util.GenKey()

    n1 := libnode.NewNode("node1", key1, "location1")
    n2 := libnode.NewNode("node2", key2, "location2")

    l := NewLedger()
    l.AddEntry("node1", util.KeyString(key1), "location1")
    l.AddEntry("node2", util.KeyString(key2), "location2")

    n1.Listen("localhost:4005")
    n2.Connect("localhost:4005")

    now := time.Now()

    b1 := NewBlockTimer(time.Second, now)
    b2 := NewBlockTimer(time.Second, now)

    h1 := NewHydrogen(l, b1)
    h2 := NewHydrogen(l, b2)

    NewMessagePasser(n1, key1, h1)
    NewMessagePasser(n2, key2, h2)

    tc1 := make(chan []message.Vote)
    tc2 := make(chan []message.Vote)

    h1.newblock = tc1
    h2.newblock = tc2

    <-tc1
    <-tc2

    v1 := <-tc1
    v2 := <-tc2

    if len(v1) != 2 {
        t.Log(v1[0].Authorization().Account())
        t.Fatalf("Not enough votes %s", len(v1))
    }

    if len(v2) != 2 {
        t.Log(v1[0].Authorization().Account())
        t.Fatalf("Not enough votes %s", len(v2))
    }

}

