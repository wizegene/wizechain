package core

import (
	"encoding/hex"
	"encoding/json"
	"github.com/beevik/ntp"
	"github.com/minio/blake2b-simd"
	"log"
	"sync"
)

type IBlock interface {
	serialize() []byte
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
	CreateBlankBlock(chainId string, lastHash []byte, lastBlockId uint32, version uint32) *Block
}

type Block struct {
	mu           *sync.RWMutex
	Header       *BlockHeader
	BlockVersion uint32
	Tx           []Transaction
	Payload      []byte
	// are sidechain blocks
	RightBlocks []Block
	// are parent chain blocks may be null if not a sidechain
	LeftBlocks           []Block
	Created              int64
	ElapsedSinceCreation uint64
}

type BlockHeader struct {
	Id         uint32
	ChainID    string
	Version    uint32
	Height     uint32
	Hash       string
	MerkleRoot []byte
	Dna        string
	PrevHash   []byte
	Timestamp  int64
}

func (b *Block) CreateBlankBlock(
	chainId string, lastHash []byte, lastBlockId uint32, version uint32) *Block {

	// we get the UTC time from the NTP server to make sure the chain is using a neutral / trustable time

	timestamp, err := ntp.Time("0.ca.pool.ntp.org")
	if err != nil {
		log.Fatalf("time server (NTP) error! : %v", err)
	}
	created := timestamp.UnixNano()
	block := &Block{}
	blockID := lastBlockId + 1
	block.Created = created
	bh := &BlockHeader{}
	bh.Id = blockID
	bh.ChainID = chainId
	bh.Version = version
	bh.Height = blockID
	dna := GetDNA(5000).D
	bh.Dna = hex.EncodeToString(dna)
	bh.PrevHash = lastHash
	block.Header = bh
	block.BlockVersion = 1
	block.Tx = make([]Transaction, 0)
	block.RightBlocks = make([]Block, 0)
	block.LeftBlocks = make([]Block, 0)
	now, _ := ntp.Time("pool.ntp.org")
	block.ElapsedSinceCreation = uint64(now.UnixNano() - created)
	block.Payload = make([]byte, 0)

	return block

}

func (b *Block) serialize() []byte {

	sBlock, _ := json.Marshal(b)
	return sBlock

}

func (b *Block) setBlockHash() {
	blockBytes := b.serialize()
	hash := blake2b.New512()
	hash.Write(blockBytes)
	b.Header.Hash = hex.EncodeToString(hash.Sum(nil))
}

func (b *Block) isOrphan() bool {
	if len(b.Tx) <= 0 {
		return true
	}
	return false
}

func (b *Block) setPayload() {
	b.Payload = []byte(Encode(b.serialize(), WizegeneAlphabet))
}

func CreateNewBlock(chainId string, lastHash []byte, lastBlockId uint32, version uint32) []byte {
	B := &Block{}
	newBlock := B.CreateBlankBlock(chainId, lastHash, lastBlockId, version)
	newBlock.setBlockHash()
	newBlock.setPayload()
	return newBlock.serialize()
}

func CreateRawBlock(chainId string, lastHash []byte, lastBlockId uint32, version uint32) *Block {
	B := &Block{}
	newBlock := B.CreateBlankBlock(chainId, lastHash, lastBlockId, version)
	newBlock.setBlockHash()
	newBlock.setPayload()
	return newBlock
}
