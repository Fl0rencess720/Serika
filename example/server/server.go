package main

import (
	"fmt"

	"github.com/Fl0rencess720/suzuRPC/server"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith struct {
}

func (*Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func main() {
	server := server.NewServer()
	err := server.Register("Arith", new(Arith))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("start server on 127.0.0.1:9001......")
	err = server.Serve("tcp", "127.0.0.1:9001")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
