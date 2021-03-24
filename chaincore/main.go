package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/wizegene/wizechain/chaincore/core/p2p"

	"github.com/wizegene/wizechain/chaincore/core"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	/*lastHashDummy := make([]byte, 64)
	hash := blake2b.New256()
	hash.Write(lastHashDummy)

	info := core.GetDogeBalance()

	addr := core.AM.NewAddressRing(LocalPrefix4Byte, [2]byte{0xff, 0xff})

	fmt.Println(addr.ToString())
	fmt.Printf("doge value: %v\n", info.Usd_value)
	fmt.Printf("liquidities: %s\n", info.Liquidity_balance)
	fbal, _ := strconv.ParseFloat(info.Liquidity_balance, 10)

	fmt.Printf("liq value: USD$ %v\n", float32(fbal)*info.Usd_value)*/
	w := core.NewWizeChain("d8ba206f-32e6-4a23-98ee-25f944f5d2fa", "1")

	m := p2p.NewMasterNode()
	spew.Dump(m.MasterName)
	spew.Dump(w)

}

// 4 byte network magics

var LocalPrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x03}
var LivePrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x01}
var TestPrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x02}
var CurrencySymbol4Byte = [4]byte{0x57, 0x49, 0x5a, 0x45}

// Script Tags as Bytes

var Script_4Byte_Amount_Prefix = [4]byte{0x77, 0x41, 0x6d, 0x74}
var Script_Cmd_Prefix = [2]byte{0x77, 0x43}
var Script_Elements_Sep = 0x5f
var Script_Tx_FlagByte = 0x7e
var Script_PrePadding = []byte("______")
var Script_End = []byte("_E_")
