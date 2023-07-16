package main

import (
	"flag"
	"fmt"
	topSenders "github.com/blockpane/topsenders"
	"log"
	"time"
)

func main() {
	flag.StringVar(&topSenders.ApiURL, "api", "", "REQUIRED: API URL, ex: http://localhost:1317")
	flag.StringVar(&topSenders.RpcUrl, "rpc", "", "REQUIRED: RPC URL, ex: http://localhost:26657")
	flag.DurationVar(&topSenders.Interval, "interval", 2*time.Second, "optional: polling interval, ex: 2s, 250ms")
	flag.IntVar(&topSenders.LookBack, "blocks", 100, "optional: max number of blocks to show")
	flag.IntVar(&topSenders.Top, "top", 20, "optional: how many top senders to show")

	flag.Parse()

	if topSenders.ApiURL == "" || topSenders.RpcUrl == "" {
		fmt.Print("\n---------------------------------------\n")
		fmt.Print("ERROR: -api and -rpc flags are required\n")
		fmt.Print("---------------------------------------\n\n")
		flag.PrintDefaults()
		return
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	topSenders.Txs()
}
