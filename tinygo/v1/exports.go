package v1

import (
	"bytes"
	"errors"
	"log"
	"os"
	"syscall"

	v1net "github.com/refraction-networking/watm/tinygo/v1/net"
	"github.com/refraction-networking/watm/wasip1"
)

func readInboundConfig() (config []byte, err error) {
	// check if /conf/inbound.cfg exists
	file, err := os.Open("/conf/inbound.cfg")
	if err != nil {
		return nil, syscall.EACCES
	}
	return readConfig(file)
}

func readOutboundConfig() (config []byte, err error) {
	// check if /conf/outbound.cfg exists
	file, err := os.Open("/conf/outbound.cfg")
	if err != nil {
		return nil, syscall.EACCES
	}
	return readConfig(file)
}

func readConfig(file *os.File) (config []byte, err error) {
	// read the config file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Println("readConfig: (*bytes.Buffer).ReadFrom:", err)
		return nil, syscall.EIO
	}

	config = buf.Bytes()

	// close the file
	if err := file.Close(); err != nil {
		return config, syscall.EIO
	}

	return config, nil
}

//export watm_ctrlpipe_v1
func _ctrlpipe(ctrlFd int32) int32 {
	ctrlConn = v1net.RebuildTCPConn(ctrlFd)
	if err := ctrlConn.SetNonBlock(true); err != nil {
		log.Printf("dial: ctrlConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	return 0 // ESUCCESS
}

//export watm_userpipe_v1
func _userpipe(userFd int32) int32 {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	sourceConn = v1net.RebuildTCPConn(userFd)
	if err := sourceConn.SetNonBlock(true); err != nil {
		log.Printf("internal_pipe: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	return 0 // ESUCCESS
}

// _dial
//
//	watm_dial_v1(networkType s32) -> s32
//
//export watm_dial_v1
func _dial(networkType int32) (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	if !globalDialer.locked {
		panic("dialer is not built with any outbound transport")
	}

	// Check if dialer is configurable. If so,
	// pull the config file from the host and configure it.
	configurableDialer := globalDialer.Configurable()
	if configurableDialer != nil {
		config, err := readOutboundConfig()
		if err == nil || config != nil {
			configurableDialer.Configure(config)
		} else if !errors.Is(err, syscall.EACCES) { // EACCES means no config file provided by the host
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	}

	if sourceConn == nil {
		log.Printf("internal_pipe: sourceConn is nil, _internal_pipe must be called first")
		return wasip1.EncodeWATERError(syscall.ENOTCONN) // socket not connected
	}

	var network string = v1net.ToNetworkTypeString(networkType)
	var address string
	var err error
	if address, err = v1net.GetAddrSuggestion(); err != nil {
		log.Printf("dial: v1net.GetAddrSuggestion: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	if globalDialer.wt != nil {
		// call v1net.Dial
		rawNetworkConn, err := v1net.Dial(network, address)
		if err != nil {
			log.Printf("dial: v1net.Dial(%s, %s): %v", network, address, err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = rawNetworkConn.Fd()

		// Note: here we are not setting nonblocking mode on the
		// networkConn -- it depends on the WrappingTransport to
		// determine whether to set nonblocking mode or not.

		// wrap
		remoteConn, err = globalDialer.wt.Wrap(rawNetworkConn)
		if err != nil {
			log.Printf("dial: d.wt.Wrap: %v", err)
			return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
		}
	} else if globalDialer.dt != nil {
		// call dt.Dial
		remoteConn, err = globalDialer.dt.Dial(network, address)
		if err != nil {
			log.Printf("dial: d.dt.Dial(%s, %s): %v", network, address, err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = remoteConn.Fd()
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	remoteConn.SetNonBlock(true) // at this point, it is safe to set non-blocking mode on remoteConn
	workerIdentity = identity_dialer
	return networkFd
}

// _accept
//
//	watm_accept_v1() -> s32
//
//export watm_accept_v1
func _accept() (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	if !globalListener.locked {
		panic("listener is not built with any inbound transport")
	}

	// Check if listener is configurable. If so,
	// pull the config file from the host and configure it.
	configurableListener := globalListener.Configurable()
	if configurableListener != nil {
		config, err := readInboundConfig()
		if err == nil || config != nil {
			configurableListener.Configure(config)
		} else if !errors.Is(err, syscall.EACCES) { // EACCES means no config file provided by the host
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	}

	if sourceConn == nil {
		log.Printf("internal_pipe: sourceConn is nil, _internal_pipe must be called first")
		return wasip1.EncodeWATERError(syscall.ENOTCONN) // socket not connected
	}

	if globalListener.wt != nil {
		var lis v1net.Listener = &v1net.TCPListener{}
		// call v1net.Listener.Accept
		rawNetworkConn, err := lis.Accept()
		if err != nil {
			log.Printf("dial: v1net.Listener.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = rawNetworkConn.Fd()

		// Note: here we are not setting nonblocking mode on the
		// networkConn -- it depends on the WrappingTransport to
		// determine whether to set nonblocking mode or not.

		// wrap
		remoteConn, err = globalListener.wt.Wrap(rawNetworkConn)
		if err != nil {
			log.Printf("dial: d.wt.Wrap: %v", err)
			return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
		}
	} else if globalListener.lt != nil {
		// call v1net.ListeningTransport.Accept
		wrappedNetworkConn, err := globalListener.lt.Accept()
		if err != nil {
			log.Printf("dial: v1net.Listener.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = wrappedNetworkConn.Fd()

		remoteConn = wrappedNetworkConn
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	remoteConn.SetNonBlock(true) // at this point, it is safe to set non-blocking mode on remoteConn
	workerIdentity = identity_listener
	return networkFd
}

// _associate
//
//	watm_associate_v1(networkType s32) -> s32
//
//export watm_associate_v1
func _associate(networkType int32) int32 {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	if !globalRelay.inboundLocked && !globalRelay.outboundLocked {
		panic("relay is not built with either inbound or outbound transport")
	}

	// Check if relay is configurable. If so,
	// pull the config file from the host and configure it.
	configurableInbRelay := globalRelay.InboundConfigurable()
	if configurableInbRelay != nil {
		config, err := readInboundConfig()
		if err == nil || config != nil {
			configurableInbRelay.Configure(config)
		} else if !errors.Is(err, syscall.EACCES) { // EACCES means no config file provided by the host
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	}

	configurableOutRelay := globalRelay.OutboundConfigurable()
	if configurableOutRelay != nil {
		config, err := readOutboundConfig()
		if err == nil || config != nil {
			configurableOutRelay.Configure(config)
		} else if !errors.Is(err, syscall.EACCES) { // EACCES means no config file provided by the host
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	}

	// handle inbound connection
	var err error
	if globalRelay.inboundListeningTransport != nil {
		sourceConn, err = globalRelay.inboundListeningTransport.Accept()
		if err != nil {
			log.Printf("dial: relay.ListeningTransport.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	} else { // we first accept the inbound connection, then wrap it if there's a wrapping transport set for it
		sourceConn, err = (&v1net.TCPListener{}).Accept()
		if err != nil {
			log.Printf("dial: v1net.TCPListener.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}

		if globalRelay.inboundWrappingTransport != nil {
			sourceConn, err = globalRelay.inboundWrappingTransport.Wrap(sourceConn.(*v1net.TCPConn))
			if err != nil {
				log.Printf("dial: r.inboundWrappingTransport.Wrap: %v", err)
				return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
			}
		}
	}

	// handle outbound connection
	var network string = v1net.ToNetworkTypeString(networkType)
	var address string
	if address, err = v1net.GetAddrSuggestion(); err != nil {
		log.Printf("dial: v1net.GetAddrSuggestion: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}
	if globalRelay.outboundDialingTransport != nil {
		remoteConn, err = globalRelay.outboundDialingTransport.Dial(network, address)
		if err != nil {
			log.Printf("dial: r.outboundDialingTransport.Dial: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	} else { // we first dial the outbound connection, then wrap it if there's a wrapping transport set for it
		remoteConn, err = v1net.Dial(network, address)
		if err != nil {
			log.Printf("dial: v1net.Dial: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}

		if globalRelay.outboundWrappingTransport != nil {
			remoteConn, err = globalRelay.outboundWrappingTransport.Wrap(remoteConn.(*v1net.TCPConn))
			if err != nil {
				log.Printf("dial: r.outboundWrappingTransport.Wrap: %v", err)
				return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
			}
		}
	}

	// set non-blocking mode on both connections
	if err := sourceConn.SetNonBlock(true); err != nil {
		log.Printf("dial: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}
	if err := remoteConn.SetNonBlock(true); err != nil {
		log.Printf("dial: remoteConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	workerIdentity = identity_relay
	return 0
}

// _start
//
//	watm_start_v1() -> s32
//
//export watm_start_v1
func _start() int32 {
	if workerIdentity == identity_uninitialized {
		log.Println("worker: uninitialized")
		return wasip1.EncodeWATERError(syscall.ENOTCONN) // socket not connected
	}
	log.Printf("worker: working as %s", identityStrings[workerIdentity])
	return worker()
}
