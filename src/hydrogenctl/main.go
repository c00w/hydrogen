package main

import (
	"flag"
	"fmt"

	"liblithium"
)

func main() {
	flag.Parse()

	c, err := liblithium.NewClient()
	if err != nil {
		fmt.Print(err)
		return
	}

	switch flag.Arg(0) {
	case "getbalance":
		r := c.GetBalance(flag.Arg(1))
		fmt.Print(r)
	case "transfer":
		b := uint64(0)
		fmt.Sscan(flag.Arg(2), b)
		r := c.SendMoney(flag.Arg(1), b)
		fmt.Print(r)
	default:
		fmt.Print(`Usage:
hydrogenctl getbalance <account>
hydrogenctl transfer <destination> <amount>`)
	}
}
