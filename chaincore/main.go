package main

import (
	"chaincore/core"
	"fmt"
)

func main() {

	alphabet := core.WizegeneAlphabet
	wizeBase58 := core.Encode([]byte("wizegenehello"), alphabet)
	fmt.Println(wizeBase58)
}
