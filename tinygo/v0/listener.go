package v0

type listener struct {
	wt WrappingTransport
	// lt ListeningTransport
}

func (l *listener) ConfigurableTransport() ConfigurableTransport {
	if l.wt != nil {
		if wt, ok := l.wt.(ConfigurableTransport); ok {
			return wt
		}
	}

	// if l.lt != nil {
	// 	if lt, ok := l.lt.(ConfigurableTransport); ok {
	// 		return lt
	// 	}
	// }

	return nil
}

func (l *listener) Initialize() {
	// TODO: allow initialization on listener
}

var l listener

// BuildListenerWithWrappingTransport arms the listener with a
// [WrappingTransport] that is used to wrap a [v0net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// Mutually exclusive with [BuildListenerWithListeningTransport].
func BuildListenerWithWrappingTransport(wt WrappingTransport) {
	l.wt = wt
	// l.lt = nil
}

// BuildListenerWithListeningTransport arms the listener with a
// [ListeningTransport] that is used to accept incoming connections
// on a local address and provide high-level application layer
// protocol over the accepted connection.
//
// Mutually exclusive with [BuildListenerWithWrappingTransport].
func BuildListenerWithListeningTransport(lt ListeningTransport) {
	// TODO: implement BuildListenerWithListeningTransport
	// l.lt = lt
	// l.wt = nil
	panic("BuildListenerWithListeningTransport: not implemented")
}
