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

// TODO: gaukas: I feel this function can be hugely optimized.
// It is not necessary for us to check and configure
// all the transports, but maybe only the one to be used.
// Should we consider moving the configuration part to the
// role-specific functions? (e.g. _dial, _accept, _associate)
//
//export watm_init_v1
func _init() int32 {
	// Check if dialer/listener/relay is configurable. If so,
	// pull the config file from the host and configure them.
	dct := globalDialer.ConfigurableTransport()
	fdct := globalFixedDialer.ConfigurableTransport()
	lct := globalListener.ConfigurableTransport()
	rct := globalRelay.ConfigurableTransport()
	if dct != nil || lct != nil /* || rct != nil */ {
		config, err := readConfig()
		if err == nil || config != nil {
			if dct != nil {
				dct.Configure(config)
			}

			if fdct != nil {
				fdct.Configure(config)
			}

			if lct != nil {
				lct.Configure(config)
			}

			if rct != nil {
				rct.Configure(config)
			}
		} else if !errors.Is(err, syscall.EACCES) { // EACCES means no config file provided by the host
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	}

	// TODO: initialize the dialer, listener, and relay
	globalDialer.Initialize()
	globalFixedDialer.Initialize()
	globalListener.Initialize()
	globalRelay.Initialize()

	return 0 // ESUCCESS
}

func readConfig() (config []byte, err error) {
	// check if /conf/watm.cfg exists
	file, err := os.Open("/conf/watm.cfg")
	if err != nil {
		return nil, syscall.EACCES
	}

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

//export watm_dial_v1
func _dial(internalFd int32) (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	// wrap the internalFd into a v1net.Conn
	sourceConn = v1net.RebuildTCPConn(internalFd)
	err := sourceConn.(*v1net.TCPConn).SetNonBlock(true)
	if err != nil {
		log.Printf("dial: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	if globalDialer.wt != nil {
		// call v1net.Dial
		rawNetworkConn, err := v1net.DialFixed()
		if err != nil {
			log.Printf("dial: v1net.Dial: %v", err)
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
		// TODO: implementation using DialingTransport
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	workerIdentity = identity_dialer
	return networkFd
}

//export watm_dial_fixed_v1
func _dial_fixed(internalFd int32) (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	// wrap the internalFd into a v1net.Conn
	sourceConn = v1net.RebuildTCPConn(internalFd)
	err := sourceConn.(*v1net.TCPConn).SetNonBlock(true)
	if err != nil {
		log.Printf("dial: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	if globalFixedDialer.fdt != nil {
		globalFixedDialer.fdt.SetDialer(v1net.Dial)
		remoteConn, err = globalFixedDialer.fdt.DialFixed()
		if err != nil {
			log.Printf("dial: c.dt.Dial: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = remoteConn.Fd()
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	workerIdentity = identity_connector
	return networkFd
}

//export watm_accept_v1
func _accept(internalFd int32) (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	// wrap the internalFd into a v1net.Conn
	sourceConn = v1net.RebuildTCPConn(internalFd)
	err := sourceConn.(*v1net.TCPConn).SetNonBlock(true)
	if err != nil {
		log.Printf("dial: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
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
		globalListener.lt.SetListener(&v1net.TCPListener{})
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

	workerIdentity = identity_listener
	return networkFd
}

//export watm_associate_v1
func _associate() int32 {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	if globalRelay.wt != nil {
		var err error
		var lis v1net.Listener = &v1net.TCPListener{}
		sourceConn, err = lis.Accept()
		if err != nil {
			log.Printf("dial: v1net.Listener.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}

		remoteConn, err = v1net.DialFixed()
		if err != nil {
			log.Printf("dial: v1net.Dial: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}

		if globalRelay.wrapSelection == RelayWrapRemote {
			// wrap remoteConn
			remoteConn, err = globalRelay.wt.Wrap(remoteConn.(*v1net.TCPConn))
			// set sourceConn, the not-wrapped one, to non-blocking mode
			sourceConn.(*v1net.TCPConn).SetNonBlock(true)
		} else {
			// wrap sourceConn
			sourceConn, err = globalRelay.wt.Wrap(sourceConn.(*v1net.TCPConn))
			// set remoteConn, the not-wrapped one, to non-blocking mode
			remoteConn.(*v1net.TCPConn).SetNonBlock(true)
		}
		if err != nil {
			log.Printf("dial: r.wt.Wrap: %v", err)
			return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
		}
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	workerIdentity = identity_relay
	return 0
}

//export watm_start_v1
func _start() int32 {
	if workerIdentity == identity_uninitialized {
		log.Println("worker: uninitialized")
		return wasip1.EncodeWATERError(syscall.ENOTCONN) // socket not connected
	}
	log.Printf("worker: working as %s", identityStrings[workerIdentity])
	return worker()
}
