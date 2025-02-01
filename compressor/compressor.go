package compressor

type CompressType uint8

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)
