package core

import (
	slip10 "github.com/lmars/go-slip10"
)

func CreateSeedForKey() ([]byte, error) {
	// todo persist seed
	seed, err := slip10.NewSeed()
	return seed, err
}

func CreateMasterKey(keySeed []byte) *slip10.Key {

	masterKey, _ := slip10.NewMasterKeyWithCurve(keySeed, slip10.CurveP256)
	return masterKey

}

func CreateChildKey(masterKey *slip10.Key, index uint32) (*slip10.Key, error) {
	ck, err := masterKey.NewChildKey(index)
	if err != nil {
		ck, err = masterKey.NewChildKey(index + 1)
	}
	return ck, err
}

func GetChildPubKey(childKey *slip10.Key) *slip10.Key {

	pub := childKey.PublicKey()
	return pub
}
