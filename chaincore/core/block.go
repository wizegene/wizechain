package core

type IBlock interface {
	serialize()
	unSerialize()
	getLastProof()
	getProof()
	getValidator()
	getStakeValue()
	getAddress()
	getNonce()
	getTransactions()
	getSize() uint32
	isOrphan() bool
	isLocked() bool
	lock()
	unlock()
	verify()
	validate()
	isGenesis() bool
}

type Block struct {

}

type BlockHeader struct {

}
