package client

import (
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
	conn, err := net.DialTimeout(network, address, c.DialTimeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
