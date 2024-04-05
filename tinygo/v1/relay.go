package v1

type RelayWrapSelection bool

const (
	RelayWrapRemote RelayWrapSelection = false
	RelayWrapSource RelayWrapSelection = true
)

type relay struct {
	wt            WrappingTransport
	wrapSelection RelayWrapSelection
	// lt ListeningTransport
	// dt DialingTransport
}

func (r *relay) ConfigurableTransport() ConfigurableTransport {
	if r.wt != nil {
		if wt, ok := r.wt.(ConfigurableTransport); ok {
			return wt
		}
	}

	// if r.lt != nil {
	// 	if lt, ok := r.lt.(ConfigurableTransport); ok {
	// 		return lt
	// 	}
	// }

	// if r.dt != nil {
	// 	if dt, ok := r.dt.(ConfigurableTransport); ok {
	// 		return dt
	// 	}
	// }

	return nil
}

func (r *relay) Initialize() {
	// TODO: allow initialization on relay
}

var globalRelay relay

// BuildRelayWithWrappingTransport arms the relay with a
// [WrappingTransport] that is used to wrap a [v1net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// The caller MUST keep in mind that the [WrappingTransport] is
// used to wrap the connection to the remote address, not the
// connection from the source address (the dialing peer).
// To reverse this behavior, i.e., wrap the inbounding connection,
// set wrapSelection to [RelayWrapSource].
//
// Mutually exclusive with [BuildRelayWithListeningDialingTransport].
func BuildRelayWithWrappingTransport(wt WrappingTransport, wrapSelection RelayWrapSelection) {
	globalRelay.wt = wt
	globalRelay.wrapSelection = wrapSelection
	// r.lt = nil
	// r.dt = nil
}

// BuildRelayWithListeningDialingTransport arms the relay with a
// [ListeningTransport] that is used to accept incoming connections
// on a local address and provide high-level application layer
// protocol over the accepted connection, and a [DialingTransport]
// that is used to dial a remote address and provide high-level
// application layer protocol over the dialed connection.
//
// Mutually exclusive with [BuildRelayWithWrappingTransport].
func BuildRelayWithListeningDialingTransport(lt ListeningTransport, dt DialingTransport) {
	// TODO: implement BuildRelayWithListeningDialingTransport
	// r.lt = lt
	// r.dt = dt
	// r.wt = nil
	panic("BuildRelayWithListeningDialingTransport: not implemented")
}
