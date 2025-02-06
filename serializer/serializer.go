package serializer

type SerializerType uint8

const (
	JSON SerializerType = iota
	PROTOBUF
)

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	GetSerializerType() SerializerType
}

var Serializers = map[SerializerType]Serializer{
	JSON:     &JSONSerializer{},
	PROTOBUF: &ProtobufSerializer{},
}
