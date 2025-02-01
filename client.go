package main

import (
	"fmt"

	"github.com/Fl0rencess720/suzuRPC/protocol"
)

func main() {
	header := protocol.Header{
		CompressType: 1,
		Method:       "love_taffy",
		ID:           1883,
		Len:          1023,
		Checksum:     666,
	}
	res := header.Mashall()
	fmt.Printf("res: %v\n", res)
}
