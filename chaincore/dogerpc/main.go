package main

import (
	rpcclient "github.com/stevenroose/go-bitcoin-core-rpc"

	"log"
)

const DOGE_RPCAPI_IP = "134.122.26.240"
const DOGE_RPCAPI_PORT = "22556"

func main() {

	connCfg := &rpcclient.ConnConfig{
		Host: DOGE_RPCAPI_IP + ":" + DOGE_RPCAPI_PORT,
		User: "wizecoin",
		Pass: "thisismywizecoindogecoinnode",
	}

	client, err := rpcclient.New(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	blockCount, err := client.GetInfo()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %s", blockCount)

}
