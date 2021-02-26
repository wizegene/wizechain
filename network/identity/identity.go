package identity

import (
	"github.com/WizegeneFoundation/wizechain/crypto"
	"crypto/ecdsa"
	"github.com/libp2p/go-libp2p-core/peer"
)

type PrivKey ecdsa.PrivateKey
type PubKey ecdsa.PublicKey

type ID peer.ID

func Generate() *PrivKey {
	res, _ := crypto.GenerateECDSAKeys()
}