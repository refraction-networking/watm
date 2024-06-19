package main

import (
	v1 "github.com/refraction-networking/watm/tinygo/v1"
)

func init() {
	v1.WorkerFairness(false)
	v1.SetReadBufferSize(1024) // 1024B buffer for copying data
	v1.BuildDialerWithWrappingTransport(&PlainWrappingTransport{})
	v1.BuildListenerWithWrappingTransport(&PlainWrappingTransport{})
	v1.BuildRelayWithWrappingTransport(&PlainWrappingTransport{}, v1.RelayWrapRemote)
	v1.BuildFixedDialerWithFixedDialingTransport(&PlainFixedDialingTransport{})
}

func main() {}
