package main

import v1 "github.com/refraction-networking/watm/tinygo/v1"

func init() {
	v1.BuildDialerWithWrappingTransport(&ReverseWrappingTransport{})
	v1.BuildListenerWithWrappingTransport(&ReverseWrappingTransport{})
	v1.BuildRelayWithWrappingTransport(&ReverseWrappingTransport{}, v1.RelayWrapRemote)
	v1.BuildFixedDialerWithFixedDialingTransport(&ReverseFixedDialingTransport{})
}

func main() {}
