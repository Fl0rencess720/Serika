package main

import (
	"fmt"

	"github.com/Fl0rencess720/Serika/client"
	"github.com/Fl0rencess720/Serika/compressor"
	"github.com/Fl0rencess720/Serika/registry"
	"github.com/Fl0rencess720/Serika/serializer"
	consulAPI "github.com/hashicorp/consul/api"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

func NewDiscovery(consulAddress string) (registry.ServiceDiscovery, error) {
	config := consulAPI.DefaultConfig()
	config.Address = consulAddress
	consulClient, err := consulAPI.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &registry.ConsulServiceDiscovery{ConsulClient: consulClient}, nil
}

func main() {
	discovery, err := NewDiscovery("127.0.0.1:8500")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	metadata, err := discovery.DiscoveryWithHeathCheck("example", nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("metadata: %v\n", *metadata)
	client, err := client.NewClient(metadata.Network, metadata.Address,
		client.WithCompressor(compressor.Raw),
		client.WithSerializer(serializer.JSON),
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	reply := new(Reply)
	args := &Args{A: 10, B: 4}
	err = client.Call("Arith", "Mul", args, reply)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("reply.C: %v\n", reply.C)
	select {}
}
