package v1

type dialer struct {
	wt WrappingTransport
	dt DialingTransport
}

func (d *dialer) ConfigurableTransport() ConfigurableTransport {
	if d.wt != nil {
		if wt, ok := d.wt.(ConfigurableTransport); ok {
			return wt
		}
	}

	if d.dt != nil {
		if dt, ok := d.dt.(ConfigurableTransport); ok {
			return dt
		}
	}

	return nil
}

func (d *dialer) Initialize() {
	// TODO: allow initialization on dialer
}

var globalDialer dialer

// BuildDialerWithWrappingTransport arms the dialer with a
// [WrappingTransport] that is used to wrap a [v1net.Conn] into
// another [net.Conn] by providing some high-level application
// layer protocol.
//
// Mutually exclusive with [BuildDialerWithDialingTransport].
func BuildDialerWithWrappingTransport(wt WrappingTransport) {
	globalDialer.wt = wt
	globalDialer.dt = nil
}

// BuildDialerWithDialingTransport arms the dialer with a
// [DialingTransport] that is used to dial a remote address and
// provide high-level application layer protocol over the dialed
// connection.
//
// Mutually exclusive with [BuildDialerWithWrappingTransport].
func BuildDialerWithDialingTransport(dt DialingTransport) {
	panic("not implemented until Runtime start passing remote access into WATM")
	// globalDialer.dt = dt
	// globalDialer.wt = nil
}

type fixedDialer struct {
	fdt FixedDialingTransport
}

// ConfigurableTransport returns the ConfigurableTransport of the
// underlying DialingTransport if it implements the interface.
func (f *fixedDialer) ConfigurableTransport() ConfigurableTransport {
	if f.fdt != nil {
		if dt, ok := f.fdt.(ConfigurableTransport); ok {
			return dt
		}
	}

	return nil
}

func (f *fixedDialer) Initialize() {
	// TODO: allow initialization on dialer
}

var globalFixedDialer fixedDialer

// BuildFixedDialerWithFixedDialingTransport arms the fixedDialer with a
// [FixedDialingTransport] that provide high-level application layer protocol
// over the dialed connection.
func BuildFixedDialerWithFixedDialingTransport(fdt FixedDialingTransport) {
	globalFixedDialer.fdt = fdt
}
