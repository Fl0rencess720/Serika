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

func (h *Header) Unmashall(data []byte) error {
	h.Lock()
	defer h.Unlock()
	idx, size := 0, 0
	h.CompressType = compressor.CompressType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size
	h.Method, size = readString(data[idx:])
	idx += size
	h.ID, size = binary.Uvarint(data[idx:])
	idx += size
	length, size := binary.Uvarint(data[idx:])
	h.Len = uint32(length)
	idx += size
	h.Checksum = binary.LittleEndian.Uint32(data[idx:])
	return nil
}

func putString(data []byte, s string) int {
	idx := 0
	idx += binary.PutUvarint(data, uint64(len(s)))
	copy(data[idx:], s)
	idx += len(s)
	return idx
}

func readString(data []byte) (string, int) {
	idx := 0
	length, size := binary.Uvarint(data)
	idx += size
	str := string(data[idx : idx+int(length)])
	idx += len(str)
	return str, idx
}
