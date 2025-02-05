package codec

import (
	"github.com/Fl0rencess720/suzuRPC/compressor"
	"github.com/Fl0rencess720/suzuRPC/protocol"
	"github.com/Fl0rencess720/suzuRPC/serializer"
)

type ClientCodec struct {
	Compressor compressor.Compressor
	Serializer serializer.Serializer
}

func NewCodecClient(compressor compressor.Compressor, serializer serializer.Serializer) *ClientCodec {
	return &ClientCodec{
		Compressor: compressor,
		Serializer: serializer,
	}
}

func (c *ClientCodec) EncodeRequest(h *protocol.Header, b *protocol.Body) ([]byte, error) {

	return nil, nil
}

func (c *ClientCodec) DecodeResponse(data []byte, v interface{}) error {
	return nil
}
