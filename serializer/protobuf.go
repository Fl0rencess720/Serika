package serializer

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var ErrInvalidProtobufMessage = errors.New("invalid protobuf message type")

type ProtobufSerializer struct {
}

func (s *ProtobufSerializer) Encode(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtobufMessage
	}
	return proto.Marshal(msg)
}

func (s *ProtobufSerializer) Decode(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return ErrInvalidProtobufMessage
	}
	return proto.Unmarshal(data, msg)
}

func (s *ProtobufSerializer) GetSerializerType() SerializerType {
	return PROTOBUF
}
