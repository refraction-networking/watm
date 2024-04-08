package v1

type listener struct {
	wt WrappingTransport
	lt ListeningTransport
}

func (l *listener) ConfigurableTransport() ConfigurableTransport {
	if l.wt != nil {
		if wt, ok := l.wt.(ConfigurableTransport); ok {
			return wt
		}
	}

	if l.lt != nil {
		if lt, ok := l.lt.(ConfigurableTransport); ok {
			return lt
		}
	}

	return nil
}

func (l *listener) Initialize() {
	// TODO: allow initialization on listener
}

var globalListener listener

// BuildListenerWithWrappingTransport arms the listener with a
// [WrappingTransport] that is used to wrap a [v1net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// Mutually exclusive with [BuildListenerWithListeningTransport].
func BuildListenerWithWrappingTransport(wt WrappingTransport) {
	globalListener.wt = wt
	globalListener.lt = nil
}

// BuildListenerWithListeningTransport arms the listener with a
// [ListeningTransport] that is used to accept incoming connections
// on a local address and provide high-level application layer
// protocol over the accepted connection.
//
// Mutually exclusive with [BuildListenerWithWrappingTransport].
func BuildListenerWithListeningTransport(lt ListeningTransport) {
	globalListener.lt = lt
	globalListener.wt = nil
}
