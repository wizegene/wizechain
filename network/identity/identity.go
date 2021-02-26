package identity

import (
	"crypto/ecdsa"
	"encoding/json"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"math/big"

	"github.com/libp2p/go-libp2p-core/peer"
	"io/ioutil"
)

type PrivKey ecdsa.PrivateKey

type PubKey ecdsa.PublicKey

type ID peer.ID

func Generate() *PrivKey {
	res, _ := GenerateECDSAKeys()
	WriteKeyToFile(res)
	return (*PrivKey)(res)
}

func WriteKeyToFile(priv *ecdsa.PrivateKey) {

	marshalled, _ := json.Marshal(priv)

	_ = ioutil.WriteFile("./.keys/priv.key", marshalled, 0600)

}

func ReadKeyFromFile(filePath string) *ecdsa.PrivateKey {

	var priv ecdsa.PrivateKey
	data, _ := ioutil.ReadFile(filePath)
	_ = json.Unmarshal(data, &priv)
	return &priv

}

func (pk *PrivKey) PeerID() peer.ID {
	priv, _, _ := pk.Libp2p()
	pid, _ := peer.IDFromPrivateKey(priv)
	return pid
}

func (pk *PrivKey) ID() ID {
	return ID(pk.PeerID())
}

func (pk *PrivKey) Libp2p() (p2pcrypto.PrivKey, p2pcrypto.PubKey, error) {
	priv, pub, err := p2pcrypto.ECDSAKeyPairFromKey((*ecdsa.PrivateKey)(pk))
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil

}

func UnmarshalPubKey(b []byte) (p2pcrypto.PubKey, error) {

	pub, err := p2pcrypto.UnmarshalPublicKey(b)
	if err != nil {
		return nil, err
	}

	return pub, nil

}

func (pub *PubKey) Libp2p() p2pcrypto.PubKey {
	return (*p2pcrypto.Secp256k1PublicKey)(pub)
}

func (pub *PubKey) Wize() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(pub)
}

func (pub *PubKey) PeerID() peer.ID {
	pid, _ := peer.IDFromPublicKey(pub.Libp2p())
	return pid
}

func (pub *PubKey) ID() ID {
	return pub.ID()
}

func (id *ID) Big() *big.Int {
	x := big.NewInt(0)
	b := []byte(*id)
	x.SetBytes(b[:32])
	return x
}

func (id *ID) Pretty() string {
	return (*peer.ID)(id).Pretty()
}

func StringToID(s string) ID {
	return ID(s)
}
