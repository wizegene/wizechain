package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

)

func GenerateECDSAKeys() (priv *ecdsa.PrivateKey, err error) {
	priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	return
}
