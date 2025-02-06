package codec

import (
	"github.com/Fl0rencess720/Serika/compressor"
	"github.com/Fl0rencess720/Serika/protocol"
	"github.com/Fl0rencess720/Serika/serializer"
)

type ServerCodec struct {
	Compressor compressor.Compressor
	Serializer serializer.Serializer
}

func NewServerCodec(compressor compressor.Compressor, serializer serializer.Serializer) *ServerCodec {
	return &ServerCodec{
		Compressor: compressor,
		Serializer: serializer,
	}
}

func (c *ServerCodec) DecodeRequest(data []byte, h *protocol.Header, b *protocol.Body) error {
	headerLen := data[0]
	err := h.Unmarshall(data[1 : 1+headerLen])
	if err != nil {
		return err
	}
	payload, err := c.Compressor.Unzip(data[1+headerLen:])
	if err != nil {
		return err
	}
	b.Payload = payload
	return nil
}

func (c *ServerCodec) EncodeResponse(args interface{}, h *protocol.Header, b *protocol.Body) ([]byte, error) {
	payload, err := c.Serializer.Encode(args)
	if err != nil {
		return nil, err
	}
	b.Payload = payload
	byteHeader := h.Marshall()
	zippedPayload, err := c.Compressor.Zip(b.Payload)
	if err != nil {
		return nil, err
	}
	data := append(byteHeader, zippedPayload...)
	data = append([]byte{byte(len(byteHeader))}, data...)
	return data, nil
}
