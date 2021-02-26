package protocol

import (
	"WizechainFoundation/wizechain/network/identity"
	"github.com/WizegeneFoundation/Cryptography/libs"
)

func IDFromBatch(batch [][]byte) identity.ID {
	h := libs.GetKeccak256(batch...).String()
}
