package protocol

type Body struct {
	Metadata map[string]string
	Payload  []byte
}
