package server

import (
	"crypto/tls"
	"net"
)

func (s *Server) makeListener(network, address string) (net.Listener, error) {
	var err error
	var ln net.Listener
	if s.Options.TLSConfig != nil {
		ln, err = tls.Listen(network, address, s.Options.TLSConfig)
		if err != nil {
			return nil, err
		}
	} else {
		ln, err = net.Listen(network, address)
		if err != nil {
			return nil, err
		}
	}
	return ln, nil
}
