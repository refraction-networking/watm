package v0

type dialer struct {
	wt WrappingTransport
	// dt DialingTransport
}

func (d *dialer) ConfigurableTransport() ConfigurableTransport {
	if d.wt != nil {
		if wt, ok := d.wt.(ConfigurableTransport); ok {
			return wt
		}
	}

	// if d.dt != nil {
	// 	if dt, ok := d.dt.(ConfigurableTransport); ok {
	// 		return dt
	// 	}
	// }

	return nil
}

func (d *dialer) Initialize() {
	// TODO: allow initialization on dialer
}

var d dialer

// BuildDialerWithWrappingTransport arms the dialer with a
// [WrappingTransport] that is used to wrap a [v0net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// Mutually exclusive with [BuildDialerWithDialingTransport].
func BuildDialerWithWrappingTransport(wt WrappingTransport) {
	d.wt = wt
	// d.dt = nil
}

// BuildDialerWithDialingTransport arms the dialer with a
// [DialingTransport] that is used to dial a remote address and
// provide high-level application layer protocol over the dialed
// connection.
//
// Mutually exclusive with [BuildDialerWithWrappingTransport].
func BuildDialerWithDialingTransport(dt DialingTransport) {
	// TODO: implement BuildDialerWithDialingTransport
	// d.dt = dt
	// d.wt = nil
	panic("BuildDialerWithDialingTransport: not implemented")
}
