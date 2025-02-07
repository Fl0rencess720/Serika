package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
)

type ZlibCompressor struct {
}

func (c *ZlibCompressor) Zip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func (c *ZlibCompressor) Unzip(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *ZlibCompressor) GetCompressorType() CompressType {
	return Zlib
}
