package wire

import "bytes"

const maxSize = uint64(0x02000000)
const maxMem = uint(5000000)

// Serializer ...
type Serializer interface {

}

type Writer struct {
	Buf bytes.Buffer
	mempool map[int][]byte
}

func (w Writer) writeCompact() {

}