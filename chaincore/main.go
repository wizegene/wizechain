package main

import (
	"bytes"
	"chaincore/config"
	"chaincore/core"
	"chaincore/core/wire"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	alphabet := core.WizegeneAlphabet
	wizeBase58 := core.Encode([]byte("wizegenehello"), alphabet)
	fmt.Println(wizeBase58)
	r := bytes.NewBuffer(make([]byte, 0))
	var m wire.Message
	total, _ := wire.CreateMessage(r, m, 1, []byte(config.LocalPrefix), wire.MessageEncoding(1))

	if total > 0 {
		fmt.Println(total)

		b := r.Bytes()
		spew.Dump(b)

	}
}
