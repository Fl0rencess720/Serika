package protocol

import (
	"encoding/binary"
	"sync"

	"github.com/Fl0rencess720/Serika/compressor"
	"github.com/Fl0rencess720/Serika/serializer"
)

const (
	magicNumber   byte = 0x08
	MaxHeaderSize int  = 38
	Uint32Size    int  = 4
	Uint16Size    int  = 2
)

type Header struct {
	sync.RWMutex
	MagicNumber    byte
	Status         byte // 0: success, 1: fail
	CompressType   compressor.CompressType
	SerializerType serializer.SerializerType
	ServicePath    string
	ServiceMethod  string
	ID             uint64
	PayloadLen     uint32
	Checksum       uint32
}

func GetMagicNumber() byte {
	return magicNumber
}

func (h *Header) Marshall() []byte {
	// 加读锁，header只读状态
	h.RLock()
	defer h.RUnlock()
	h.MagicNumber = magicNumber
	byteHeader := make([]byte, MaxHeaderSize+len(h.ServiceMethod)+len(h.ServicePath))
	idx := 0
	byteHeader[idx] = h.MagicNumber
	idx++
	byteHeader[idx] = h.Status
	idx++
	binary.LittleEndian.PutUint16(byteHeader[idx:], uint16(h.CompressType))
	idx += Uint16Size
	binary.LittleEndian.PutUint16(byteHeader[idx:], uint16(h.SerializerType))
	idx += Uint16Size
	idx += putString(byteHeader[idx:], h.ServiceMethod)
	idx += putString(byteHeader[idx:], h.ServicePath)
	idx += binary.PutUvarint(byteHeader[idx:], h.ID)
	idx += binary.PutUvarint(byteHeader[idx:], uint64(h.PayloadLen))
	binary.LittleEndian.PutUint32(byteHeader[idx:], h.Checksum)
	idx += Uint32Size
	return byteHeader[:idx]
}

func (h *Header) Unmarshall(data []byte) error {
	h.Lock()
	defer h.Unlock()
	idx, size := 0, 0
	h.MagicNumber = data[idx]
	idx++
	h.Status = data[idx]
	idx++
	h.CompressType = compressor.CompressType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size
	h.SerializerType = serializer.SerializerType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size
	h.ServiceMethod, size = readString(data[idx:])
	idx += size
	h.ServicePath, size = readString(data[idx:])
	idx += size
	h.ID, size = binary.Uvarint(data[idx:])
	idx += size
	length, size := binary.Uvarint(data[idx:])
	h.PayloadLen = uint32(length)
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
