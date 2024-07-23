package main

import (
	v1 "github.com/refraction-networking/watm/tinygo/v1"
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

// type guard: ReverseWrappingTransport must implement [v1.WrappingTransport].
var _ v1.WrappingTransport = (*ReverseWrappingTransport)(nil)

type ReverseWrappingTransport struct {
}

func (rwt *ReverseWrappingTransport) Wrap(conn v1net.Conn) (v1net.Conn, error) {
	return &ReverseConn{conn}, conn.SetNonBlock(true) // must set non-block, otherwise will block on read and lose fairness
}
