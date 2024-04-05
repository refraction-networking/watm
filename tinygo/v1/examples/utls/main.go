package main

import v1 "github.com/refraction-networking/watm/tinygo/v1"

func init() {
	v1.BuildDialerWithWrappingTransport(&UTLSClientWrappingTransport{})
	// v1.BuildListenerWithWrappingTransport(&UTLSClientWrappingTransport{})
	// v1.BuildRelayWithWrappingTransport(&UTLSClientWrappingTransport{}, v0.RelayWrapRemote)
}

func main() {}
