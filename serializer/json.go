package serializer

type JSONSerializer struct {
}

func (s *JSONSerializer) Encode(v interface{}) ([]byte, error) {
	return nil, nil
}

func (s *JSONSerializer) Decode(data []byte, v interface{}) error {
	return nil
}
