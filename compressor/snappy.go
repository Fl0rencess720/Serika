package compressor

type SnappyCompressor struct {
}

func (c *SnappyCompressor) Zip(data []byte) ([]byte, error) {
	return nil, nil
}

func (c *SnappyCompressor) Unzip(data []byte) ([]byte, error) {
	return nil, nil
}
func (c *SnappyCompressor) GetCompressorType() CompressType {
	return Snappy
}
