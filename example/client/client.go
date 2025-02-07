package main

import (
	"crypto/tls"
	"fmt"

	"github.com/Fl0rencess720/Serika/client"
	"github.com/Fl0rencess720/Serika/compressor"
	"github.com/Fl0rencess720/Serika/serializer"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

func main() {
	client, err := client.NewClient("tcp", "127.0.0.1:9001",
		client.WithCompressor(compressor.Raw),
		client.WithSerializer(serializer.JSON),
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)
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
