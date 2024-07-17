package main

import (
	v1 "github.com/refraction-networking/watm/tinygo/v1"
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

// type guard
var _ v1.ListeningTransport = (*ReverseListeningTransport)(nil)

type ReverseListeningTransport struct {
	listener v1net.Listener
}

func (lt *ReverseListeningTransport) SetListener(listener v1net.Listener) {
	lt.listener = listener
}

func (lt *ReverseListeningTransport) Accept() (v1net.Conn, error) {
	conn, err := lt.listener.Accept() // accept the connection
	if err != nil {
		return nil, err
	}

	return &ReverseConn{conn}, conn.SetNonBlock(true) // must set non-block, otherwise will block on read and lose fairness
}
