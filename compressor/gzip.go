package compressor

type GzipCompressor struct {
}

func (c *GzipCompressor) Zip(data []byte) ([]byte, error) {
	return nil, nil
}

func (c *GzipCompressor) Unzip(data []byte) ([]byte, error) {
	return nil, nil
}
func (c *GzipCompressor) GetCompressorType() CompressType {
	return Gzip
}
