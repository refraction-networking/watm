package main

import (
	v1 "github.com/refraction-networking/watm/tinygo/v1"
)

func init() {
	v1.WorkerFairness(false)
	v1.SetReadBufferSize(1024) // 1024B buffer for copying data
	v1.BuildDialerWithDialingTransport(&PlainDialingTransport{})
	v1.BuildListenerWithListeningTransport(&PlainListeningTransport{})
	v1.BuildRelayWithOutboundWrappingTransport(&PlainWrappingTransport{})
}

func main() {}
