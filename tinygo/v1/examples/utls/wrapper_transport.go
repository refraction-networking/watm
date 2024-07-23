package main

import (
	"log"

	"github.com/CosmWasm/tinyjson"
	tls "github.com/refraction-networking/utls"
	v1 "github.com/refraction-networking/watm/tinygo/v1"
	"github.com/refraction-networking/watm/tinygo/v1/examples/utls/lib"
	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
)

// type guard: ReverseWrappingTransport must implement [v1.WrappingTransport].
var _ v1.WrappingTransport = (*UTLSClientWrappingTransport)(nil)

type UTLSClientWrappingTransport struct {
	tlsConfig     *tls.Config
	clientHelloID tls.ClientHelloID
}

func (uwt *UTLSClientWrappingTransport) Wrap(conn v1net.Conn) (v1net.Conn, error) {
	if uwt.tlsConfig == nil {
		log.Println("UTLSClientWrappingTransport: tlsConfig is nil, using default config")
		uwt.tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	var emptyClientHelloID tls.ClientHelloID
	if uwt.clientHelloID == emptyClientHelloID {
		uwt.clientHelloID = tls.HelloChrome_Auto
	}

	tlsConn := tls.UClient(conn, uwt.tlsConfig, uwt.clientHelloID)
	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}

	if err := conn.SetNonBlock(true); err != nil {
		return nil, err
	}

	return &UTLSConn{
		Conn:    conn,
		tlsConn: tlsConn,
	}, nil
}

var _ v1.Configurable = (*UTLSClientWrappingTransport)(nil)

func (uwt *UTLSClientWrappingTransport) Configure(config []byte) error {
	configurables := &lib.Configurables{}
	if err := tinyjson.Unmarshal(config, configurables); err != nil {
		return err
	}

	uwt.tlsConfig = configurables.GetTLSConfig()
	uwt.clientHelloID = configurables.GetClientHelloID()

	v1.WorkerFairness(configurables.BackgroundWorkerFairness)
	log.Printf("UTLSClientWrappingTransport: set worker fairness to %v\n", configurables.BackgroundWorkerFairness)
	if configurables.InternalBufferSize > 0 {
		v1.SetReadBufferSize(configurables.InternalBufferSize)
		log.Printf("UTLSClientWrappingTransport: set resizing internal buffer to %d Bytes\n", configurables.InternalBufferSize)
	}

	return nil
}

type UTLSConn struct {
	v1net.Conn // embedded Conn
	tlsConn    *tls.UConn
}

func (uc *UTLSConn) Read(b []byte) (n int, err error) {
	return uc.tlsConn.Read(b)
}

func (uc *UTLSConn) Write(b []byte) (n int, err error) {
	return uc.tlsConn.Write(b)
}
