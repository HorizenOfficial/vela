package communication

import (
	"github.com/mdlayher/vsock"
	"net"
)

// VSockConnectionFactory implements ConnectionFactory for v-sock connections
type VSockConnectionFactory struct {
	cid  uint32 // Context ID
	port uint32 // Port number
}

// NewVSockConnectionFactory creates a new v-sock connection factory
func NewVSockConnectionFactory(cid, port uint32) *VSockConnectionFactory {
	return &VSockConnectionFactory{
		cid:  cid,
		port: port,
	}
}

// CreateClientConnection creates a v-sock client connection
func (f *VSockConnectionFactory) CreateClientConnection() (net.Conn, error) {
	return vsock.Dial(f.cid, f.port, nil)
}

// CreateServerListener creates a v-sock server listener
func (f *VSockConnectionFactory) CreateServerListener() (net.Listener, error) {
	return vsock.Listen(f.port, nil)
}
