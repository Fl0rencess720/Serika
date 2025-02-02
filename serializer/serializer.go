package serializer

type Serializer uint8

const (
	JSON Serializer = iota
	PROTOBUF
)
