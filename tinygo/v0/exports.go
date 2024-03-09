package v0

import (
	"bytes"
	"errors"
	"log"
	"os"
	"syscall"

	v0net "github.com/refraction-networking/watm/tinygo/v0/net"
	"github.com/refraction-networking/watm/wasip1"
)

// Export the WATM version indicator.
//
//	gaukas: I noticed that in Rust we can export a const variable
//	but here in Go we have to export a function instead. Luckily
//	in our standard we are not checking against its type but only
//	the name.
//
//export _water_v0
func _water_v0() {}

//export _water_init
func _water_init() int32 {
	// Check if dialer/listener/relay is configurable. If so,
	// pull the config file from the host and configure them.
	dct := d.ConfigurableTransport()
	lct := l.ConfigurableTransport()
	// rct := r.ConfigurableTransport()
	if dct != nil || lct != nil /* || rct != nil */ {
		config, err := readConfig()
		if err == nil || config != nil {
			if dct != nil {
				dct.Configure(config)
			}

			if lct != nil {
				lct.Configure(config)
			}

			// if rct != nil {
			// 	rct.Configure(config)
			// }
		} else if !errors.Is(err, syscall.EACCES) { // EACCES means no config file provided by the host
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
	}

	// TODO: initialize the dialer, listener, and relay
	d.Initialize()
	l.Initialize()
	// r.Initialize()

	return 0 // ESUCCESS
}

func readConfig() (config []byte, err error) {
	fd, err := wasip1.DecodeWATERError(_import_pull_config())
	if err != nil {
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "config")
	if file == nil {
		return nil, syscall.EBADF
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

//export _water_cancel_with
func _water_cancel_with(cancelFd int32) int32 {
	cancelConn = v0net.RebuildTCPConn(cancelFd)
	if err := cancelConn.(v0net.Conn).SetNonBlock(true); err != nil {
		log.Printf("dial: cancelConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	return 0 // ESUCCESS
}

//export _water_dial
func _water_dial(internalFd int32) (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	// wrap the internalFd into a v0net.Conn
	sourceConn = v0net.RebuildTCPConn(internalFd)
	err := sourceConn.(*v0net.TCPConn).SetNonBlock(true)
	if err != nil {
		log.Printf("dial: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	if d.wt != nil {
		// call v0net.Dial
		rawNetworkConn, err := v0net.Dial("", "")
		if err != nil {
			log.Printf("dial: v0net.Dial: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = rawNetworkConn.Fd()

		// Note: here we are not setting nonblocking mode on the
		// networkConn -- it depends on the WrappingTransport to
		// determine whether to set nonblocking mode or not.

		// wrap
		remoteConn, err = d.wt.Wrap(rawNetworkConn)
		if err != nil {
			log.Printf("dial: d.wt.Wrap: %v", err)
			return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
		}
		// TODO: implement _water_dial with DialingTransport
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	workerIdentity = identity_dialer
	return networkFd
}

//export _water_accept
func _water_accept(internalFd int32) (networkFd int32) {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	// wrap the internalFd into a v0net.Conn
	sourceConn = v0net.RebuildTCPConn(internalFd)
	err := sourceConn.(*v0net.TCPConn).SetNonBlock(true)
	if err != nil {
		log.Printf("dial: sourceConn.SetNonblock: %v", err)
		return wasip1.EncodeWATERError(err.(syscall.Errno))
	}

	if d.wt != nil {
		var lis v0net.Listener = &v0net.TCPListener{}
		// call v0net.Listener.Accept
		rawNetworkConn, err := lis.Accept()
		if err != nil {
			log.Printf("dial: v0net.Listener.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}
		networkFd = rawNetworkConn.Fd()

		// Note: here we are not setting nonblocking mode on the
		// networkConn -- it depends on the WrappingTransport to
		// determine whether to set nonblocking mode or not.

		// wrap
		remoteConn, err = d.wt.Wrap(rawNetworkConn)
		if err != nil {
			log.Printf("dial: d.wt.Wrap: %v", err)
			return wasip1.EncodeWATERError(syscall.EPROTO) // protocol error
		}
		// TODO: implement _water_accept with ListeningTransport
	} else {
		return wasip1.EncodeWATERError(syscall.EPERM) // operation not permitted
	}

	workerIdentity = identity_listener
	return networkFd
}

//export _water_associate
func _water_associate() int32 {
	if workerIdentity != identity_uninitialized {
		return wasip1.EncodeWATERError(syscall.EBUSY) // device or resource busy (worker already initialized)
	}

	if r.wt != nil {
		var err error
		var lis v0net.Listener = &v0net.TCPListener{}
		sourceConn, err = lis.Accept()
		if err != nil {
			log.Printf("dial: v0net.Listener.Accept: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}

		remoteConn, err = v0net.Dial("", "")
		if err != nil {
			log.Printf("dial: v0net.Dial: %v", err)
			return wasip1.EncodeWATERError(err.(syscall.Errno))
		}

		if r.wrapSelection == RelayWrapRemote {
			// wrap remoteConn
			remoteConn, err = r.wt.Wrap(remoteConn.(*v0net.TCPConn))
			// set sourceConn, the not-wrapped one, to non-blocking mode
			sourceConn.(*v0net.TCPConn).SetNonBlock(true)
		} else {
			// wrap sourceConn
			sourceConn, err = r.wt.Wrap(sourceConn.(*v0net.TCPConn))
			// set remoteConn, the not-wrapped one, to non-blocking mode
			remoteConn.(*v0net.TCPConn).SetNonBlock(true)
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

//export _water_worker
func _water_worker() int32 {
	if workerIdentity == identity_uninitialized {
		log.Println("worker: uninitialized")
		return wasip1.EncodeWATERError(syscall.ENOTCONN) // socket not connected
	}
	log.Printf("worker: working as %s", identityStrings[workerIdentity])
	return worker()
}
