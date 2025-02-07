package registry

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Fl0rencess720/Serika/server"
	consulAPI "github.com/hashicorp/consul/api"
)

type ConsulServiceRegister struct {
	ConsulClient *consulAPI.Client
}

func (r *ConsulServiceRegister) Register(ctx context.Context, server *server.Server) error {
	address := strings.Split(server.Metadata.Address, ":")[0]
	portStr := strings.Split(server.Metadata.Address, ":")[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}
	registration := &consulAPI.AgentServiceRegistration{
		ID:      server.Metadata.Name,
		Name:    server.Metadata.Name,
		Address: address,
		Port:    port,
		Check: &consulAPI.AgentServiceCheck{
			CheckID:                        "service:" + server.Metadata.Name,
			TTL:                            "15s",
			DeregisterCriticalServiceAfter: "60s",
		},
	}
	if err := r.ConsulClient.Agent().ServiceRegister(registration); err != nil {
		return err
	}
	// 启动TTL心跳
	go func() {
		err := r.UpdateTTL("service:" + server.Metadata.Name)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		ticker := time.NewTicker(5 * time.Second) // 每 5s 发送一次心跳
		defer ticker.Stop()
		for range ticker.C {
			err := r.UpdateTTL("service:" + server.Metadata.Name)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}
	}()
	return nil
}

func (r *ConsulServiceRegister) Deregister(ctx context.Context, server *server.Server) error {
	return nil
}

func (r *ConsulServiceRegister) UpdateTTL(checkID string) error {
	if err := r.ConsulClient.Agent().UpdateTTL(checkID, "pass", "pass"); err != nil {
		return err
	}
	return nil
}

type ConsulServiceDiscovery struct {
	ConsulClient *consulAPI.Client
}

func (c *ConsulServiceDiscovery) Discovery(serviceID string, o *consulAPI.QueryOptions) (*server.Metadata, error) {
	service, _, err := c.ConsulClient.Agent().Service(serviceID, o)
	if err != nil {
		return nil, err
	}
	return &server.Metadata{Name: service.ID, Address: fmt.Sprintf("%s:%d", service.Address, service.Port), Network: "tcp"}, nil
}

func (c *ConsulServiceDiscovery) DiscoveryWithHeathCheck(serviceID string, o *consulAPI.QueryOptions) (*server.Metadata, error) {
	service, _, err := c.ConsulClient.Health().Service(serviceID, "", true, nil)
	if err != nil {
		return nil, err
	}
	return &server.Metadata{Name: service[0].Service.ID, Address: fmt.Sprintf("%s:%d", service[0].Service.Address, service[0].Service.Port), Network: "tcp"}, nil
}
