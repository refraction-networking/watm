package net

import (
	"github.com/refraction-networking/watm/tinygo/v1/wasiimport"
	"github.com/refraction-networking/watm/wasip1"
)

// DialFixed connects to a pre-determined-by-runtime address and returns a Conn.
func DialFixed() (Conn, error) {
	fd, err := wasip1.DecodeWATERError(wasiimport.WaterDialFixed())
	if err != nil {
		return nil, err
	}

	return RebuildTCPConn(fd), nil
}

// Dial dials a remote host for a network connection.
func Dial(network, address string) (Conn, error) {
	fd, err := wasip1.DecodeWATERError(wasiimport.WaterDial(
		makeIOVec([]byte(network)), 1,
		makeIOVec([]byte(address)), 1,
	))
	if err != nil {
		return nil, err
	}

	return RebuildTCPConn(fd), nil
}
