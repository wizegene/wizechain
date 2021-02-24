package core

import (
	"chaincore/tools"
)

var (
	AddrPrefixWallet  = [2]byte{0x77, 0x32}
	AddrPrefixTX      = [2]byte{0x74, 0x31}
	AddrPrefixBlock   = [2]byte{0x62, 0x31}
	AddrPrefixContent = [2]byte{0x63, 0x63}
)

var LocalPrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x03}
var LivePrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x01}
var TestPrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x02}
var CurrencySymbol4Byte = [4]byte{0x57, 0x49, 0x5a, 0x45}

const (
	HDPurpose = "48"
	HDNetwork = "Wizechain"
)

var HDRoleOwner = `0x0`
var HDRoleActive = `0x1`
var HDRoleMemo = `0x3`
var HDRoleValidator = `0x4`

func CreateHDHierarchy(purpose, network, index, role, keyIndex string) string {

	s := "m/" + purpose
	s += "'/" + network
	s += "'/" + index
	s += "'/" + role
	s += "'/" + keyIndex + "'"

	return s

}

type Address struct {
	network  [4]byte
	prefix   [2]byte
	checksum [4]byte
	pubKey   *tools.Key
	key      string
	rootKey  *tools.Key
	seed     []byte
	index    uint32
}

type IAddress interface {
	Create(network [4]byte, prefix [2]byte) *Address
	lock() bool
	unlock() bool
	anonymize() bool
	createWalletAddress(key *tools.Key) *Address
	createTxAddress(key *tools.Key) *Address
	createBlockAddress(key *tools.Key) *Address
	createContentAddress(key *tools.Key) *Address
	Validate() bool
	ToString() string
}

type AddressManager struct {
	IAddress
}

var AM AddressManager

func (am AddressManager) Create(network [4]byte, prefix [2]byte) *Address { return &Address{} }

func (am AddressManager) lock() bool { return false }

func (am AddressManager) unlock() bool { return false }

func (am AddressManager) anonymize() bool { return false }

func (am AddressManager) createWalletAddress(masterAddr *Address) *Address {

	wAddress := &Address{}
	wAddress.index = masterAddr.index + 1
	child, _ := masterAddr.rootKey.NewChildKey(wAddress.index)
	wAddress.pubKey = child.PublicKey()
	wAddress.key = Encode(child.PublicKey().Key, WizegeneAlphabet)
	wAddress.rootKey = masterAddr.rootKey
	wAddress.prefix = masterAddr.prefix
	wAddress.seed = masterAddr.seed

	return wAddress
}

func (am AddressManager) createTxAddress(key *tools.Key) *Address { return &Address{} }

func (am AddressManager) createBlockAddress(key *tools.Key) *Address { return &Address{} }

func (am AddressManager) createContentAddress(key *tools.Key) *Address { return &Address{} }

func (am AddressManager) Validate() bool { return false }

func (am AddressManager) ToString() string {
	return ""
}

func (am *AddressManager) NewAddressRing(network [4]byte, prefix [2]byte) *Address {

	// we initialize a new keyring with a seed. seed needs to be stored off chain
	seed, _ := CreateSeedForKey()

	// from the seed we create a master key for the new address keyring
	master, _ := tools.NewMasterKeyWithCurve(seed)
	a := &Address{}

	a.network = network
	a.prefix = prefix
	a.seed = seed

	// its a new keyring so the root pubkey will be the master public key
	pubKey := master.PublicKey()
	a.rootKey = master

	a.pubKey = pubKey
	return AM.createWalletAddress(a)
	/*switch am.prefix {
	case AddrPrefixWallet:
		AM.createWalletAddress(am)
		break
	case AddrPrefixTX:
		AM.createTxAddress(pubKey)
		break
	case AddrPrefixBlock:
		AM.createBlockAddress(pubKey)
		break
	case AddrPrefixContent:
		AM.createContentAddress(pubKey)
		break
	default:
		break
	}*/

}

func (am *Address) ToString() string {
	return am.key
}

func (am *Address) SeedToMnemonic() string {
	return ""
}
