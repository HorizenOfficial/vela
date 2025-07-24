package communication

import (
	"net"
)

// TCPConnectionFactory implements ConnectionFactory for TCP connections
type TCPConnectionFactory struct {
	addr string
}

// NewTCPConnectionFactory creates a new TCP connection factory
func NewTCPConnectionFactory(addr string) *TCPConnectionFactory {
	return &TCPConnectionFactory{
		addr: addr,
	}
}

// CreateClientConnection creates a TCP client connection
func (f *TCPConnectionFactory) CreateClientConnection() (net.Conn, error) {
	return net.Dial("tcp", f.addr)
}

// CreateServerListener creates a TCP server listener
func (f *TCPConnectionFactory) CreateServerListener() (net.Listener, error) {
	return net.Listen("tcp", f.addr)
}
