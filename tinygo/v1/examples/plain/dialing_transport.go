package main

import (
	v1 "github.com/refraction-networking/watm/tinygo/v1"
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

// type guard
var _ v1.DialingTransport = (*PlainDialingTransport)(nil)

type PlainDialingTransport struct {
	dialer func(network, address string) (v1net.Conn, error)
}

func (dt *PlainDialingTransport) SetDialer(dialer func(network, address string) (v1net.Conn, error)) {
	dt.dialer = dialer
}

func (dt *PlainDialingTransport) Dial(network, address string) (v1net.Conn, error) {
	conn, err := dt.dialer(network, address) // dial the passed address
	if err != nil {
		return nil, err
	}

	return &PlainConn{conn}, conn.SetNonBlock(true) // must set non-block, otherwise will block on read and lose fairness
}
