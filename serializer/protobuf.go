package serializer

type ProtobufSerializer struct {
}

func (s *ProtobufSerializer) Encode(v interface{}) ([]byte, error) {
	return nil, nil
}

func (s *ProtobufSerializer) Decode(data []byte, v interface{}) error {
	return nil
}
func (s *ProtobufSerializer) GetSerializerType() SerializerType {
	return PROTOBUF
}
