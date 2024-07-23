package net

import (
	"unsafe"

	"github.com/refraction-networking/watm/tinygo/v1/wasiimport"
	"github.com/refraction-networking/watm/wasip1"
)

// Dial dials a remote host for a network connection.
func Dial(network, address string) (Conn, error) {
	var fd int32
	err := wasip1.ErrnoToError(wasiimport.WaterDial(
		ToNetworkTypeInt(network),
		makeIOVec([]byte(address)), 1,
		unsafe.Pointer(&fd),
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
	var nread size
	err := wasip1.ErrnoToError(wasiimport.WaterGetAddrSuggestion(
		makeIOVec(addrBuf), 1,
		unsafe.Pointer(&nread),
	))
	if err != nil {
		return "", err
	}

	return string(addrBuf[:nread]), nil
}
