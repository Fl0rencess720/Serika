package registry

import (
	"context"

	"github.com/Fl0rencess720/Serika/server"
)

type ServiceRegister interface {
	Register(ctx context.Context, server *server.Server) error
	Deregister(ctx context.Context, server *server.Server) error
	UpdateTTL(checkID string) error
}
