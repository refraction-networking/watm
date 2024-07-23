package net

import (
	"syscall"
	"unsafe"

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
	closed bool
}

func (tl *TCPListener) Accept() (Conn, error) {
	if tl.closed {
		return nil, syscall.EINVAL
	}

	var fd int32
	err := wasip1.ErrnoToError(wasiimport.WaterAccept(
		unsafe.Pointer(&fd),
	))
	if err != nil {
		return nil, err
	}

	return RebuildTCPConn(fd), nil
}

// TCPListener does not really keep a file descriptor for an
// underlying TCP socket, so it does not need to close anything.
func (tl *TCPListener) Close() error {
	if tl.closed {
		return syscall.EINVAL
	}
	tl.closed = true
	return nil
}
