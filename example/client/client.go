package main

import (
	"fmt"

	"github.com/Fl0rencess720/suzuRPC/client"
	"github.com/Fl0rencess720/suzuRPC/compressor"
	"github.com/Fl0rencess720/suzuRPC/serializer"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

func main() {
	client, err := client.NewClient("tcp", "127.0.0.1:9001", client.WithCompressor(compressor.Raw), client.WithSerializer(serializer.JSON))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	reply := new(Reply)
	args := &Args{A: 3, B: 4}
	err = client.Call("Arith", "Mul", args, reply)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("reply.C: %v\n", reply.C)
}
