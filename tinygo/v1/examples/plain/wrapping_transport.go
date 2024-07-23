package main

import (
	v1 "github.com/refraction-networking/watm/tinygo/v1"
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

// type guard: PlainWrappingTransport must implement [v1.WrappingTransport].
var _ v1.WrappingTransport = (*PlainWrappingTransport)(nil)

type PlainWrappingTransport struct {
}

func (*PlainWrappingTransport) Wrap(conn v1net.Conn) (v1net.Conn, error) {
	return &PlainConn{conn}, conn.SetNonBlock(true) // must set non-block, otherwise will block on read and lose fairness
}

// PlainConn simply passes through the underlying Conn.
type PlainConn struct {
	v1net.Conn // embedded Conn
}
