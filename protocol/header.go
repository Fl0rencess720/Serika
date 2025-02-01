package protocol

import (
	"encoding/binary"
	"sync"

	"github.com/Fl0rencess720/suzuRPC/compressor"
)

const (
	magicNumber   byte = 0x75
	MaxHeaderSize      = 36
	Uint32Size         = 4
	Uint16Size         = 2
)

type Header struct {
	sync.RWMutex
	MagicNumber  byte
	CompressType compressor.CompressType
	Method       string
	ID           uint64
	Len          uint32
	Checksum     uint32
}

func (h *Header) Mashall() []byte {
	// 加读锁，header只读状态
	h.RLock()
	defer h.RUnlock()
	h.MagicNumber = magicNumber
	byteHeader := make([]byte, MaxHeaderSize+len(h.Method))
	idx := 0
	// [idx:] 其实就是 [0]
	binary.LittleEndian.PutUint16(byteHeader[idx:], uint16(h.MagicNumber))
	idx += Uint16Size
	binary.LittleEndian.PutUint16(byteHeader[idx:], uint16(h.CompressType))
	idx += Uint16Size
	idx += putString(byteHeader[idx:], h.Method)
	idx += binary.PutUvarint(byteHeader[idx:], h.ID)
	idx += binary.PutUvarint(byteHeader[idx:], uint64(h.Len))
	binary.LittleEndian.PutUint32(byteHeader[:idx], h.Checksum)
	idx += Uint32Size
	return byteHeader[:idx]
}

func putString(header []byte, s string) int {
	idx := 0
	idx += binary.PutUvarint(header, uint64(len(s)))
	copy(header[idx:], s)
	idx += len(s)
	return idx
}
