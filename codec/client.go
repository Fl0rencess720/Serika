package codec

import (
	"github.com/Fl0rencess720/suzuRPC/compressor"
	"github.com/Fl0rencess720/suzuRPC/protocol"
	"github.com/Fl0rencess720/suzuRPC/serializer"
)

type CodecClient struct {
	compressor compressor.Compressor
	serializer serializer.Serializer
}

func NewCodecClient(compressor compressor.Compressor, serializer serializer.Serializer) *CodecClient {
	return &CodecClient{
		compressor: compressor,
		serializer: serializer,
	}
}

func (c *CodecClient) EncodeRequest(h *protocol.Header, b *protocol.Body) ([]byte, error) {

	return nil, nil
}

func (c *CodecClient) DecodeResponse(data []byte, v interface{}) error {
	return nil
}
