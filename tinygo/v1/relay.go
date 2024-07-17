package v1

import v1net "github.com/refraction-networking/watm/tinygo/v1/net"

type relay struct {
	// inbound: pick one
	inboundWrappingTransport  WrappingTransport
	inboundListeningTransport ListeningTransport
	inboundLocked             bool

	// outbound: pick one
	outboundWrappingTransport WrappingTransport
	outboundDialingTransport  DialingTransport
	outboundLocked            bool
}

// InboundConfigurable returns non-nil if the relay is built with a
// configurable inbound transport.
func (r *relay) InboundConfigurable() Configurable {
	if r.inboundWrappingTransport != nil {
		if wt, ok := r.inboundWrappingTransport.(Configurable); ok {
			return wt
		}
	} else if r.inboundListeningTransport != nil {
		if lt, ok := r.inboundListeningTransport.(Configurable); ok {
			return lt
		}
	}

	return nil
}

// OutboundConfigurable returns non-nil if the relay is built with a
// configurable outbound transport.
func (r *relay) OutboundConfigurable() Configurable {
	if r.outboundWrappingTransport != nil {
		if wt, ok := r.outboundWrappingTransport.(Configurable); ok {
			return wt
		}
	} else if r.outboundDialingTransport != nil {
		if dt, ok := r.outboundDialingTransport.(Configurable); ok {
			return dt
		}
	}

	return nil
}

var globalRelay relay

// BuildRelayWithInboundTransport arms the relay with an inbound
// transport that is used to accept inbound connections on a local
// address and provide high-level application layer protocol over the
// accepted connection.
//
// Outbound transport must be set as well before or after calling
// this function to complete the relay configuration. Otherwise,
// the outbound transport will be plain.
func BuildRelayWithInboundTransport(anyTransport interface{}) {
	switch t := anyTransport.(type) {
	case WrappingTransport:
		BuildRelayWithInboundWrappingTransport(t)
	case ListeningTransport:
		BuildRelayWithInboundListeningTransport(t)
	default:
		panic("transport type not supported")
	}
}

func BuildRelayWithOutboundTransport(anyTransport interface{}) {
	switch t := anyTransport.(type) {
	case WrappingTransport:
		BuildRelayWithOutboundWrappingTransport(t)
	case DialingTransport:
		BuildRelayWithOutboundDialingTransport(t)
	default:
		panic("transport type not supported")
	}
}

func BuildRelayWithInboundWrappingTransport(wt WrappingTransport) {
	if globalRelay.inboundLocked {
		panic("relay is locked")
	}

	globalRelay.inboundWrappingTransport = wt
	globalRelay.inboundListeningTransport = nil
	globalRelay.inboundLocked = true
}

func BuildRelayWithInboundListeningTransport(lt ListeningTransport) {
	if globalRelay.inboundLocked {
		panic("relay is locked")
	}

	globalRelay.inboundListeningTransport = lt
	lt.SetListener(&v1net.TCPListener{})
	globalRelay.inboundWrappingTransport = nil
	globalRelay.inboundLocked = true
}

func BuildRelayWithOutboundWrappingTransport(wt WrappingTransport) {
	if globalRelay.outboundLocked {
		panic("relay is locked")
	}

	globalRelay.outboundWrappingTransport = wt
	globalRelay.outboundDialingTransport = nil
	globalRelay.outboundLocked = true
}

func BuildRelayWithOutboundDialingTransport(dt DialingTransport) {
	if globalRelay.outboundLocked {
		panic("relay is locked")
	}

	globalRelay.outboundDialingTransport = dt
	dt.SetDialer(v1net.Dial)
	globalRelay.outboundWrappingTransport = nil
	globalRelay.outboundLocked = true
}
