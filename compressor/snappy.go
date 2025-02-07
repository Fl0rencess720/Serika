package compressor

import (
	"github.com/golang/snappy"
)

type SnappyCompressor struct {
}

func (c *SnappyCompressor) Zip(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func (c *SnappyCompressor) Unzip(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}

func (c *SnappyCompressor) GetCompressorType() CompressType {
	return Snappy
}
