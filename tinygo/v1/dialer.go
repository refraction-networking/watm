package v1

import v1net "github.com/refraction-networking/watm/tinygo/v1/net"

type dialer struct {
	// outbound: pick one
	wt WrappingTransport
	dt DialingTransport

	locked bool
}

func (d *dialer) Configurable() Configurable {
	if d.wt != nil {
		if wt, ok := d.wt.(Configurable); ok {
			return wt
		}
	} else if d.dt != nil {
		if dt, ok := d.dt.(Configurable); ok {
			return dt
		}
	}

	return nil
}

var globalDialer dialer

// BuildDialerWithOutboundTransport arms the dialer with a
// [WrappingTransport] or [DialingTransport] that is used to wrap
// a [v1net.Conn] into another [net.Conn] by providing some
// high-level application layer protocol or dial a remote address
// and provide high-level application layer protocol over the dialed
// connection.
func BuildDialerWithOutboundTransport(anyTransport any) {
	switch t := anyTransport.(type) {
	case WrappingTransport:
		BuildDialerWithWrappingTransport(t)
	case DialingTransport:
		BuildDialerWithDialingTransport(t)
	default:
		panic("transport type not supported")
	}
}

// BuildDialerWithWrappingTransport arms the dialer with a
// [WrappingTransport] that is used to wrap a [v1net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// The caller MUST keep in mind that the [WrappingTransport] is
// used to wrap the outbound connection (to the remote address), not the
// connection from the caller (the client).
//
// Mutually exclusive with [BuildDialerWithDialingTransport].
func BuildDialerWithWrappingTransport(wt WrappingTransport) {
	if globalDialer.locked {
		panic("dialer is locked")
	}

	globalDialer.wt = wt
	globalDialer.dt = nil
	globalDialer.locked = true
}

// BuildDialerWithDialingTransport arms the dialer with a
// [DialingTransport] that is used to dial a remote address and
// provide high-level application layer protocol over the dialed
// connection.
//
// Mutually exclusive with [BuildDialerWithWrappingTransport].
func BuildDialerWithDialingTransport(dt DialingTransport) {
	if globalDialer.locked {
		panic("dialer is locked")
	}

	globalDialer.dt = dt
	dt.SetDialer(v1net.Dial)
	globalDialer.wt = nil
	globalDialer.locked = true
}
