package protocol

import (
	"WizechainFoundation/wizechain/network/identity"
	"WizechainFoundation/wizechain/network/identity/libs"
	"encoding/hex"
)

func IDFromBatch(batch [][]byte) identity.ID {
	h := libs.GetKeccak256(batch)
	hs := hex.EncodeToString(h)
	return identity.ID(hs)
}

