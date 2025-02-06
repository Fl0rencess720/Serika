package compressor

type CompressType uint8

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)

type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}

var Compressors = map[CompressType]Compressor{
	Raw:    &RawCompressor{},
	Gzip:   &GzipCompressor{},
	Snappy: &SnappyCompressor{},
	Zlib:   &ZlibCompressor{},
}
