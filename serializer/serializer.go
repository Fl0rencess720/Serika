package serializer

type Serializer uint8

const (
	JSON Serializer = iota
	PROTOBUF
)

type Serialize interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

var Serializers = map[Serializer]Serialize{
	JSON:     &JSONSerializer{},
	PROTOBUF: &ProtobufSerializer{},
}
