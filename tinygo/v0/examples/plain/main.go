package main

import v0 "github.com/refraction-networking/watm/tinygo/v0"

func init() {
	v0.WorkerFairness(false) // by default, use unfairWorker for better performance
	v0.BuildDialerWithWrappingTransport(&PlainWrappingTransport{})
	v0.BuildListenerWithWrappingTransport(&PlainWrappingTransport{})
	v0.BuildRelayWithWrappingTransport(&PlainWrappingTransport{}, v0.RelayWrapRemote)
}

func main() {}
