package serializer

import "encoding/json"

type JSONSerializer struct {
}

func (s *JSONSerializer) Encode(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// v当前需要是指针类型
func (s *JSONSerializer) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
