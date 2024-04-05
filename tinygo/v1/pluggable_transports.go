package v1

import (
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

// WrappingTransport is the most basic transport type. It wraps
// a [v1net.Conn] into another [v1net.Conn] by providing some
// high-level application layer protocol.
type WrappingTransport interface {
	// Wrap wraps a v1net.Conn into another v1net.Conn with a protocol
	// wrapper layer.
	//
	// The returned v1net.Conn is NOT by default set to non-blocking.
	// It is the responsibility of the transport to make it
	// non-blocking by calling v1net.Conn.SetNonblock. This is to
	// allow some transport to perform blocking operations such as
	// TLS handshake.
	//
	// The transport SHOULD provide non-blocking v1net.Conn.Read
	// operation on the returned v1net.Conn if possible, otherwise
	// the worker may block on reading from a blocking connection.
	// And it is highly recommended to pass all funtions other than
	// Read and Write to the underlying v1net.Conn from the underlying
	// dialer function.
	Wrap(v1net.Conn) (v1net.Conn, error)
}

// Any transport requiring a dialer to be set before use should
// implement this interface.
type RequireDialer interface {
	// SetDialer sets the dialer function that is used to dial
	// a remote address.
	//
	// In v0, the input parameter of the dialer function is
	// unused inside the WATM, given the connection is always
	// managed by the host application.
	//
	// The returned v1net.Conn is NOT by default set to non-blocking.
	// It is the responsibility of the transport to make it
	// non-blocking by calling v1net.Conn.SetNonblock. This is to
	// allow some transport to perform blocking operations such as
	// TLS handshake.
	SetDialer(dialer func(network, address string) (v1net.Conn, error))
}

// DialingTransport is a transport type that can be used to dial
// a remote address and provide high-level application layer
// protocol over the dialed connection.
type DialingTransport interface {
	// A dialer must be set for a DialingTransport to be used.
	RequireDialer

	// Dial connects to the given remote address and returns a
	// v1net.Conn that provides high-level application layer
	// protocol over the dialed connection.
	//
	// The transport SHOULD provide non-blocking v1net.Conn.Read
	// operation on the returned v1net.Conn if possible, otherwise
	// the worker may block on reading from a blocking connection.
	// And it is highly recommended to pass all funtions other than
	// Read and Write to the underlying v1net.Conn from the underlying
	// dialer function.
	Dial(network, address string) (v1net.Conn, error)
}

// FixedDialingTransport is a transport type that can be used to dial
// a otherwise-determined (unknown to its immediate caller) address and
// provide high-level application layer protocol over the dialed connection.
type FixedDialingTransport interface {
	// A dialer must be set for a FixedDialingTransport to be used.
	RequireDialer

	// DialFixed connects to a remote address determined by the
	// FixedDialingTransport and returns a v1net.Conn that provides
	// high-level application layer protocol over the dialed connection.
	//
	// The transport SHOULD provide non-blocking v1net.Conn.Read
	// operation on the returned v1net.Conn if possible, otherwise
	// the worker may block on reading from a blocking connection.
	// And it is highly recommended to pass all funtions other than
	// Read and Write to the underlying v1net.Conn from the underlying
	// dialer function.
	DialFixed() (v1net.Conn, error)
}

// ListeningTransport is a transport type that can be used to
// accept incoming connections on a local address and provide
// high-level application layer protocol over the accepted
// connection.
type ListeningTransport interface {
	// SetListener sets the listener that is used to accept
	// incoming connections.
	//
	// The returned v1net.Conn is not by default non-blocking.
	// It is the responsibility of the transport to make it
	// non-blocking if required by calling v1net.Conn.SetNonblock.
	SetListener(listener v1net.Listener)

	// Accept accepts an incoming connection and returns a
	// net.Conn that provides high-level application layer
	// protocol over the accepted connection.
	//
	// The transport SHOULD provide non-blocking v1net.Conn.Read
	// operation on the returned v1net.Conn if possible, otherwise
	// the worker may block on reading from a blocking connection.
	// And it is highly recommended to pass all funtions other than
	// Read and Write to the underlying v1net.Conn from the underlying
	// dialer function.
	Accept() (v1net.Conn, error)
}

// ConfigurableTransport is a transport type that can be configured
// with a config file in the form of a byte slice.
type ConfigurableTransport interface {
	Configure(config []byte) error
}
