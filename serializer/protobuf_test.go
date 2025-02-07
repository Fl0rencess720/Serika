package serializer

import (
	"testing"

	"github.com/Fl0rencess720/Serika/serializer/serializer_proto"
	"github.com/stretchr/testify/assert"
)

func TestProtobufSerializer(t *testing.T) {
	serializer := &ProtobufSerializer{}

	// 创建一个测试消息
	originalMsg := &serializer_proto.TestMessage{
		Name: "Alice",
		Age:  25,
	}

	// 测试 Encode
	encodedData, err := serializer.Encode(originalMsg)
	assert.NoError(t, err)
	assert.NotNil(t, encodedData)

	// 测试 Decode
	decodedMsg := &serializer_proto.TestMessage{}
	err = serializer.Decode(encodedData, decodedMsg)
	assert.NoError(t, err)
	assert.Equal(t, originalMsg.Name, decodedMsg.Name)
	assert.Equal(t, originalMsg.Age, decodedMsg.Age)
}
