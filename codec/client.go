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

func NewClientCodec(compressor compressor.Compressor, serializer serializer.Serializer) *ClientCodec {
	return &ClientCodec{
		Compressor: compressor,
		Serializer: serializer,
	}
}

func (c *ClientCodec) EncodeRequest(args interface{}, h *protocol.Header, b *protocol.Body) ([]byte, error) {
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

func (c *ClientCodec) DecodeResponse(data []byte, h interface{}, b interface{}) error {
	header := h.(*protocol.Header)
	body := b.(*protocol.Body)
	headerLen := data[0]
	err := header.Unmarshall(data[1 : 1+headerLen])
	if err != nil {
		return err
	}
	payload, err := c.Compressor.Unzip(data[1+headerLen:])
	if err != nil {
		return err
	}
	body.Payload = payload
	return nil
}
