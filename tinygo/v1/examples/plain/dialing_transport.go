package main

import v1net "github.com/refraction-networking/watm/tinygo/v1/net"

type PlainFixedDialingTransport struct {
	dialer func(network, address string) (v1net.Conn, error)
}

func (fdt *PlainFixedDialingTransport) SetDialer(dialer func(network, address string) (v1net.Conn, error)) {
	fdt.dialer = dialer
}

func (fdt *PlainFixedDialingTransport) DialFixed() (v1net.Conn, error) {
	conn, err := fdt.dialer("tcp", "localhost:7700") // TODO: hardcoded address, any better idea?
	if err != nil {
		return nil, err
	}

	return &PlainConn{conn}, conn.SetNonBlock(true) // must set non-block, otherwise will block on read and lose fairness
}
