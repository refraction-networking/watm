package net

import (
	"github.com/refraction-networking/watm/tinygo/v1/wasiimport"
	"github.com/refraction-networking/watm/wasip1"
)

// Listener is the interface for a generic network listener.
type Listener interface {
	Accept() (Conn, error)
	Close() error
}

// type guard: *TCPListener must implement Listener
var _ Listener = (*TCPListener)(nil)

// TCPListener is a fake TCP listener which calls to the host
// to accept a connection.
//
// By saying "fake", it means that the file descriptor is not
// managed inside the WATM, but by the host application.
type TCPListener struct {
}

func (*TCPListener) Accept() (Conn, error) {
	fd, err := wasip1.DecodeWATERError(wasiimport.WaterAccept())
	if err != nil {
		return nil, err
	}

	return RebuildTCPConn(fd), nil
}

// TCPListener does not really keep a file descriptor for an
// underlying TCP socket, so it does not need to close anything.
func (*TCPListener) Close() error {
	return nil
}
