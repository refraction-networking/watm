package main

import v1 "github.com/refraction-networking/watm/tinygo/v1"

func init() {
	v1.BuildDialerWithDialingTransport(&ReverseDialingTransport{})
	v1.BuildListenerWithListeningTransport(&ReverseListeningTransport{})
	v1.BuildRelayWithOutboundWrappingTransport(&ReverseWrappingTransport{})
}

func main() {}
