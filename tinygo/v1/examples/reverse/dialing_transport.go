package main

import v1net "github.com/refraction-networking/watm/tinygo/v1/net"

type ReverseDialingTransport struct {
	dialer func(network, address string) (v1net.Conn, error)
}

func (dt *ReverseDialingTransport) SetDialer(dialer func(network, address string) (v1net.Conn, error)) {
	dt.dialer = dialer
}

func (dt *ReverseDialingTransport) Dial(network, address string) (v1net.Conn, error) {
	conn, err := dt.dialer(network, address) // dial the passed address
	if err != nil {
		return nil, err
	}

	return &ReverseConn{conn}, conn.SetNonBlock(true) // must set non-block, otherwise will block on read and lose fairness
}
