package client

import (
	"crypto/tls"
	"net"
)

func (c *Client) connect(network, address string) error {
	conn, err := newTCPConn(c, network, address)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func newTCPConn(c *Client, network, address string) (net.Conn, error) {
	if c.options.TLSConfig != nil {
		dialer := &net.Dialer{Timeout: c.options.dialTimeout}
		tlsConn, err := tls.DialWithDialer(dialer, network, address, c.options.TLSConfig)
		if err != nil {
			return nil, err
		}
		return tlsConn, nil
	}
	return net.DialTimeout(network, address, c.options.dialTimeout)
}
