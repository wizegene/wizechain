package core

import "github.com/wizegene/wizechain/chaincore/tools"

func CreateSeedForKey() ([]byte, error) {
	// todo persist seed
	seed, err := tools.NewSeed()
	return seed, err
}

func CreateMasterKey(keySeed []byte) *tools.Key {

	masterKey, _ := tools.NewMasterKeyWithCurve(keySeed)
	return masterKey

}

func CreateChildKey(masterKey *tools.Key, index uint32) (*tools.Key, error) {
	ck, err := masterKey.NewChildKey(index)
	if err != nil {
		ck, err = masterKey.NewChildKey(index + 1)
	}
	return ck, err
}

func GetChildPubKey(childKey *tools.Key) *tools.Key {

	pub := childKey.PublicKey()
	return pub
}

func verifyKey(data string) (*tools.Key, error) {
	return tools.B58Deserialize(data)
}
