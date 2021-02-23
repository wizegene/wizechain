package wire

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/minio/blake2b-simd"
	"io"
)

/**
Message Header
(4) message start
(12) command
(4) size
(4) checksum
*/

const commandSize = 12
const MaxMessagePayload = 1024 * 1024 * 1024 //max is 1gb for a message
const headerSize = 24

type MessageEncoding uint32

type IMessageHeader interface {
	Create(pchMessageStartIn uint, cmd string, messageSizeIn uint)
	GetCommand() string
	isCommandValid() bool
}

type MessageHeader struct {
	magic    [4]byte // 4 bytes
	command  string
	length   uint32
	checksum [4]byte
}

type EncryptionKey []byte

type IMessage interface {
	WizeDecode(io.Reader, uint32, MessageEncoding) error
	WizeEncode(io.Writer, uint32, MessageEncoding) error
	WizeEncrypt(io.Writer, uint32, MessageEncoding, EncryptionKey) error
	WizeDecrypt(io.Reader, uint32, MessageEncoding, EncryptionKey) error
	Command() string
	MaxPayloadLength(uint32) uint32
}

type Message struct {
	IMessage
}

func CreateMessage(w io.Writer, m Message, ver uint32, network [4]byte, encoding MessageEncoding) (int, error) {

	totalBytes := 0
	var command [commandSize]byte
	cmd := "PING"

	if len(cmd) > commandSize {
		str := fmt.Sprintf("command [%s] is too long [max %v]",
			cmd, commandSize)
		return totalBytes, errors.New(fmt.Sprintf("WriteMessage:%s", str))
	}

	copy(command[:], cmd)

	payload := []byte("hello")
	lenp := len(payload)
	if lenp > MaxMessagePayload {
		str := fmt.Sprintf("message payload is too large - encoded "+
			"%d bytes, but maximum message payload is %d bytes",
			lenp, MaxMessagePayload)
		return totalBytes, errors.New(fmt.Sprintf("WriteMessage:%s", str))

	}

	header := MessageHeader{}
	header.magic = network
	header.command = cmd
	header.length = uint32(lenp)

	h := blake2b.New256()
	h.Write(payload)
	chk2 := h.Sum(nil)
	h = blake2b.New256()
	h.Write(chk2)
	chksum := h.Sum(nil)
	copy(header.checksum[:], chksum[0:4])

	hw := bytes.NewBuffer(make([]byte, 0, headerSize))

	binary.Write(hw, binary.BigEndian, header.magic)
	binary.Write(hw, binary.BigEndian, command[:])
	binary.Write(hw, binary.BigEndian, header.length)
	binary.Write(hw, binary.BigEndian, header.checksum[:])
	n, err := w.Write(hw.Bytes())
	if err != nil {
		panic(err)
	}
	totalBytes += n
	n, err = w.Write(payload)
	if err != nil {
		panic(err)
	}
	totalBytes += n

	return totalBytes, nil

}
