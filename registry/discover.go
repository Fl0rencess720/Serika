package registry

import (
	"github.com/Fl0rencess720/Serika/server"
	consulAPI "github.com/hashicorp/consul/api"
)

type ServiceDiscovery interface {
	Discovery(serviceID string, o *consulAPI.QueryOptions) (*server.Metadata, error)
	DiscoveryWithHeathCheck(serviceID string, o *consulAPI.QueryOptions) ([]*server.Metadata, error)
}
