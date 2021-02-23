package core

import "sync"

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
	mu *sync.Mutex
	Header *BlockHeader
	BlockVersion uint32
	tx []Transaction
	payload []byte
	// are sidechain blocks
	rightBlocks []Block
	// are parent chain blocks may be null if not a sidechain
	leftBlocks []Block
	created int64
	elapsedSinceCreation uint64
	IBlock
}

type BlockHeader struct {
	Id uint32
	ChainID string
	Version uint32
	Height uint32
	Hash []byte
	MerkleRoot []byte
	Dna []byte
	PrevHash []byte
	Timestamp int64
}

