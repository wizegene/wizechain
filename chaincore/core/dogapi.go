package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const DOGE_TESTNET_API_KEY = "571b-3bd9-8808-a3bd"
const TEST_API_BALANCE_URL = "https://block.io/api/v2/get_address_balance/?api_key=" + DOGE_TESTNET_API_KEY + "&labels=testwize"
const TEST_API_EST_NETFEE = "https://block.io/api/v2/get_network_fee_estimate/?api_key=" + DOGE_TESTNET_API_KEY
const TICKER_URL = "https://api.coingecko.com/api/v3/simple/price?ids=dogecoin&vs_currencies=usd"

type DOGEInfo struct {
	Liquidity_balance string
	Usd_value         float32
}

func (d *DOGEInfo) Balance() {

	resp, err := http.Get(TEST_API_BALANCE_URL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	sb := body
	var data map[string]map[string]string
	json.Unmarshal(sb, &data)

	toint, _ := data["data"]["available_balance"]
	d.Liquidity_balance = toint

}

func (d *DOGEInfo) Value() {

	resp, err := http.Get(TICKER_URL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	sb := body
	var data map[string]map[string]float32
	json.Unmarshal(sb, &data)

	toint, _ := data["dogecoin"]["usd"]
	log.Printf("%v", data)
	d.Usd_value = toint

}

func GetDogeBalance() *DOGEInfo {
	doge := new(DOGEInfo)
	doge.Balance()
	doge.Value()
	return doge
}
