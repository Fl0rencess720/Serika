package compressor

type ZlibCompressor struct {
}

func (c *ZlibCompressor) Zip(data []byte) ([]byte, error) {
	return nil, nil
}

func (c *ZlibCompressor) Unzip(data []byte) ([]byte, error) {
	return nil, nil
}
func (c *ZlibCompressor) GetCompressorType() CompressType {
	return Zlib
}
