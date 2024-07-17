package main

import (
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

type ReverseConn struct {
	v1net.Conn // embedded Conn
}

func (rc *ReverseConn) Read(b []byte) (n int, err error) {
	tmpBuf := make([]byte, len(b))
	n, err = rc.Conn.Read(tmpBuf)
	if err != nil {
		return 0, err
	}

	// reverse all bytes read successfully so far
	for i := 0; i < n; i++ {
		b[i] = tmpBuf[n-i-1]
	}

	return n, err
}

func (rc *ReverseConn) Write(b []byte) (n int, err error) {
	tmpBuf := make([]byte, len(b))

	// reverse the bytes to be written
	for i := 0; i < len(b); i++ {
		tmpBuf[i] = b[len(b)-i-1]
	}

	return rc.Conn.Write(tmpBuf[:len(b)])
}
