package net

import (
	"github.com/refraction-networking/watm/tinygo/v1/wasiimport"
	"github.com/refraction-networking/watm/wasip1"
)

// Dial dials a remote host for a network connection.
func Dial(network, address string) (Conn, error) {
	fd, err := wasip1.DecodeWATERError(wasiimport.WaterDial(
		ToNetworkTypeInt(network),
		makeIOVec([]byte(address)), 1,
	))
	if err != nil {
		return nil, err
	}

	return RebuildTCPConn(fd), nil
}

// GetAddrSuggestion should be invoked by a dialer when it
// could not determine which address to dial. This function
// must be called before watm_dial_v1 returns.
func GetAddrSuggestion() (string, error) {
	var addrBuf []byte = make([]byte, 256)
	n, err := wasip1.DecodeWATERError(wasiimport.WaterGetAddrSuggestion(makeIOVec(addrBuf), 1))
	if err != nil {
		return "", err
	}

	return string(addrBuf[:n]), nil
}
