package protocol

import (
	"WizechainFoundation/wizechain/network/identity"
	"github.com/WizegeneFoundation/wizechain/Cryptography/libs"
)

func IDFromBatch(batch [][]byte) identity.ID {
	h := libs.GetKeccak256(batch...).String()
}
