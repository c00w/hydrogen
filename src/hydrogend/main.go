package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"libhelium"
	"libhydrogen"
	"libnode"
	"util"
)

var bootstrap = flag.Bool("bootstrap", false, "Whether this is the first client on the network")
var connect = flag.String("connect", "hydrogen.daedrum.net:2674", "the default hydrogen node to bootstrap off")
var debug = flag.Bool("debug", false, "Show debug output")
var keypath = flag.String("keypath", "/etc/hydrogend/server.key", "Path to key file")
var listenaddress = flag.String("address", "", "IP address to listen on. Should be internet routable")

func main() {
	flag.Parse()

	if len(*listenaddress) == 0 {
		flag.PrintDefaults()
		return
	}

	log.Println("Loading Key")
	key, err := util.LoadKey(*keypath)

	if os.IsNotExist(err) {
		log.Println("Key does not exist, Generating Key")
		err = util.GenerateKey(*keypath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Loading Key")
		key, err = util.LoadKey(*keypath)
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Key loaded")
	log.Printf("Identity is %x\n", util.KeyString(key))

	n := libnode.NewNode(key, *listenaddress)
	h := libhydrogen.NewHydrogen(n)
	libhelium.NewServer(n, h)

	if !*bootstrap {
		log.Printf("Downloading network state")
		l, err := libhelium.Connect(n, *connect)
		if err != nil {
			log.Fatal(err)
		}
		h.AddLedger(l)
	} else {
		log.Printf("Boostrapping")
		l := libhydrogen.NewLedger()
		l.AddEntry(util.KeyString(key), *listenaddress, 1<<63)
		h.AddLedger(l)
	}

	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)
	for _ = range c {
		log.Printf("Quitting\n")
		os.Exit(0)
	}
}
