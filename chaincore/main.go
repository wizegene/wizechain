package main

import (
	"chaincore/core"
	"fmt"
	"github.com/minio/blake2b-simd"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	lastHashDummy := make([]byte, 64)
	hash := blake2b.New256()
	hash.Write(lastHashDummy)

	addr := core.AM.NewAddressRing(LocalPrefix4Byte, [2]byte{0xff, 0xff})

	fmt.Println(addr.ToString())

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
