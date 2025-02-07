package main

import (
	"crypto/tls"
	"fmt"

	"github.com/Fl0rencess720/Serika/server"
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
	cert, err := tls.LoadX509KeyPair("../../server.pem", "../../server.key")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	server := server.NewServer(server.WithTLSConfig(
		&tls.Config{Certificates: []tls.Certificate{cert}}),
	)
	err = server.Register("Arith", new(Arith))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("start server on 127.0.0.1:9001......")
	err = server.Serve("tcp", ":9001")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
