package main

import (
	"context"
	"fmt"

	"github.com/Fl0rencess720/Serika/registry"
	"github.com/Fl0rencess720/Serika/server"
	consulAPI "github.com/hashicorp/consul/api"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith struct{}

func (*Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func NewRegistrar(consulAddress string) (registry.ServiceRegister, error) {
	config := consulAPI.DefaultConfig()
	config.Address = consulAddress
	consulClient, err := consulAPI.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &registry.ConsulServiceRegister{ConsulClient: consulClient}, nil
}

func main() {
	addr := "127.0.0.1:9004"
	consulAddr := "127.0.0.1:8500"

	registrar, err := NewRegistrar(consulAddr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// cert, err := tls.LoadX509KeyPair("../../server.pem", "../../server.key")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }

	// server := server.NewServer("example", server.WithTLSConfig(
	// 	&tls.Config{Certificates: []tls.Certificate{cert}}),
	// )
	server := server.NewServer("test", "test02")
	server.Metadata.Address = addr
	// 注册服务
	if err = registrar.Register(context.Background(), server); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	// 绑定方法
	err = server.Register("Arith", new(Arith))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("start server on %s......", addr)

	err = server.Serve("tcp", addr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
