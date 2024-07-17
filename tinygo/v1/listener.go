package v1

import v1net "github.com/refraction-networking/watm/tinygo/v1/net"

type listener struct {
	// incoming: pick one
	wt WrappingTransport
	lt ListeningTransport

	locked bool
}

func (l *listener) Configurable() Configurable {
	if l.wt != nil {
		if wt, ok := l.wt.(Configurable); ok {
			return wt
		}
	} else if l.lt != nil {
		if lt, ok := l.lt.(Configurable); ok {
			return lt
		}
	}

	return nil
}

var globalListener listener

// BuildListenerWithInboundTransport arms the listener with an inbound
// transport that is used to accept inbound connections on a local
// address and provide high-level application layer protocol over the
// accepted connection.
func BuildListenerWithInboundTransport(anyTransport interface{}) {
	switch t := anyTransport.(type) {
	case WrappingTransport:
		BuildListenerWithWrappingTransport(t)
	case ListeningTransport:
		BuildListenerWithListeningTransport(t)
	default:
		panic("transport type not supported")
	}
}

// BuildListenerWithWrappingTransport arms the listener with a
// [WrappingTransport] that is used to wrap a [v1net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// Mutually exclusive with [BuildListenerWithListeningTransport].
func BuildListenerWithWrappingTransport(wt WrappingTransport) {
	if globalListener.locked {
		panic("listener is locked")
	}

	globalListener.wt = wt
	globalListener.lt = nil

	globalListener.locked = true
}

// BuildListenerWithListeningTransport arms the listener with a
// [ListeningTransport] that is used to accept incoming connections
// on a local address and provide high-level application layer
// protocol over the accepted connection.
//
// Mutually exclusive with [BuildListenerWithWrappingTransport].
func BuildListenerWithListeningTransport(lt ListeningTransport) {
	if globalListener.locked {
		panic("listener is locked")
	}

	globalListener.lt = lt
	lt.SetListener(&v1net.TCPListener{})
	globalListener.wt = nil

	globalListener.locked = true
}
