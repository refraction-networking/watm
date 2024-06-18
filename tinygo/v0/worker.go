package v0

import (
	"io"
	"log"
	"net"
	"runtime"
	"syscall"

	v0net "github.com/refraction-networking/watm/tinygo/v0/net"
	"github.com/refraction-networking/watm/wasip1"
)

type identity uint8

var workerIdentity identity = identity_uninitialized

const (
	identity_uninitialized identity = iota
	identity_dialer
	identity_listener
	identity_relay
)

var identityStrings = map[identity]string{
	identity_dialer:   "dialer",
	identity_listener: "listener",
	identity_relay:    "relay",
}

var sourceConn v0net.Conn // sourceConn is used to communicate between WASM and the host application or a dialing party (for relay only)
var remoteConn v0net.Conn // remoteConn is used to communicate between WASM and a dialed remote destination (for dialer/relay) or a dialing party (for listener only)
var cancelConn v0net.Conn // cancelConn is used to cancel the entire worker.

var workerFn func() int32 = unfairWorker // by default, use unfairWorker for better performance under mostly unidirectional I/O

var readBuf []byte = make([]byte, 16384) // 16k buffer for reading

// WorkerFairness sets the fairness of a worker.
//
// If sourceConn or remoteConn will not work in non-blocking mode,
// it is highly recommended to set fair to true, otherwise it is most
// likely that the worker will block on reading from a blocking
// connection forever and therefore make no progress in the other
// direction.
func WorkerFairness(fair bool) {
	if fair {
		workerFn = fairWorker
	} else {
		workerFn = unfairWorker
	}
}

func worker() int32 {
	defer _import_host_defer()

	if sourceConn == nil || remoteConn == nil || cancelConn == nil {
		log.Println("worker: worker: sourceConn, remoteConn, or cancelConn is nil")
		return wasip1.EncodeWATERError(syscall.EBADF) // bad file descriptor
	}

	return workerFn()
}

// untilError executes the given function until non-nil error is returned
func untilError(f func() error) error {
	var err error
	for err == nil {
		err = f()
	}
	return err
}

// unfairWorker works on all three connections with a priority order
// of cancelConn > sourceConn > remoteConn.
//
// It keeps working on the current connection until it returns an error,
// and if the error is EAGAIN, it switches to the next connection. If the
// connection is not properly set to non-blocking mode, i.e., never returns
// EAGAIN, this function will block forever and never work on a lower priority
// connection. Thus it is called unfairWorker.
func unfairWorker() int32 {
	for {
		pollFd := []pollFd{
			{
				fd:     uintptr(cancelConn.Fd()),
				events: EventFdRead,
			},
			{
				fd:     uintptr(sourceConn.Fd()),
				events: EventFdRead,
			},
			{
				fd:     uintptr(remoteConn.Fd()),
				events: EventFdRead,
			},
		}

		n, err := _poll(pollFd, -1)
		if n == 0 { // TODO: re-evaluate the condition
			if err == nil || err == syscall.EAGAIN {
				runtime.Gosched() // yield the current goroutine
				continue
			}
			log.Println("worker: unfairWorker: _poll:", err)
			return int32(err.(syscall.Errno))
		}

		// 1st priority: cancelConn
		_, err = cancelConn.Read(readBuf)
		if !(err == syscall.EAGAIN) {
			if err == io.EOF || err == nil {
				log.Println("worker: unfairWorker: cancelConn is closed")
				return wasip1.EncodeWATERError(syscall.ECANCELED) // operation canceled
			}
			log.Println("worker: unfairWorker: cancelConn.Read:", err)
			return wasip1.EncodeWATERError(syscall.EIO) // input/output error
		}

		// 2nd priority: sourceConn
		if err := untilError(func() error {
			nRead, readErr := sourceConn.Read(readBuf)
			if readErr != nil {
				if readErr != syscall.EAGAIN {
					log.Println("worker: unfairWorker: sourceConn.Read:", readErr)
				}
				return readErr
			}

			nWritten, writeErr := remoteConn.Write(readBuf[:nRead])
			if writeErr != nil {
				log.Println("worker: unfairWorker: remoteConn.Write:", writeErr)
				return writeErr
			}

			if nRead != nWritten {
				log.Printf("worker: unfairWorker: nRead != nWritten")
				return syscall.EMSGSIZE // message too long to fit in send buffer even after auto partial write
			}

			return nil
		}); err != syscall.EAGAIN { // silently ignore EAGAIN
			if err == io.EOF {
				log.Println("worker: unfairWorker: sourceConn is closed")
				return wasip1.EncodeWATERError(0) // success, no error
			}
			if errno, ok := err.(syscall.Errno); ok {
				return wasip1.EncodeWATERError(errno)
			}
			return wasip1.EncodeWATERError(syscall.EIO) // input/output error
		}

		// 3rd priority: remoteConn
		if err := untilError(func() error {
			nRead, readErr := remoteConn.Read(readBuf)
			if readErr != nil {
				if readErr != syscall.EAGAIN {
					log.Println("worker: unfairWorker: remoteConn.Read:", readErr)
				}
				return readErr
			}

			nWrite, writeErr := sourceConn.Write(readBuf[:nRead])
			if writeErr != nil {
				log.Println("worker: unfairWorker: sourceConn.Write:", writeErr)
				return writeErr
			}

			if nRead != nWrite {
				log.Printf("worker: unfairWorker: nRead != nWrite")
				return syscall.EMSGSIZE // message too long to fit in send buffer even after auto partial write
			}

			return nil
		}); err != syscall.EAGAIN { // silently ignore EAGAIN
			if err == io.EOF {
				log.Println("worker: unfairWorker: remoteConn is closed")
				return wasip1.EncodeWATERError(0) // success, no error
			}
			if errno, ok := err.(syscall.Errno); ok {
				return wasip1.EncodeWATERError(errno)
			}
			return wasip1.EncodeWATERError(syscall.EIO) // input/output error
		}
	}
}

// like unfairWorker, fairWorker also works on all three connections with a priority order
// of cancelConn > sourceConn > remoteConn.
//
// But different from unfairWorker, fairWorker spend equal amount of turns on each connection
// for calling Read. Therefore it has a better fairness than unfairWorker, which may still
// make progress if one of the connection is not properly set to non-blocking mode.
func fairWorker() int32 {
	for {
		pollFd := []pollFd{
			{
				fd:     uintptr(cancelConn.Fd()),
				events: EventFdRead,
			},
			{
				fd:     uintptr(sourceConn.Fd()),
				events: EventFdRead,
			},
			{
				fd:     uintptr(remoteConn.Fd()),
				events: EventFdRead,
			},
		}

		n, err := _poll(pollFd, -1)
		if n == 0 { // TODO: re-evaluate the condition
			if err == nil || err == syscall.EAGAIN {
				runtime.Gosched() // yield the current goroutine
				continue
			}
			log.Println("worker: fairWorker: _poll:", err)
			return int32(err.(syscall.Errno))
		}

		// 1st priority: cancelConn
		_, err = cancelConn.Read(readBuf)
		if !(err == syscall.EAGAIN) {
			if err == io.EOF || err == nil {
				log.Println("worker: fairWorker: cancelConn is closed")
				return wasip1.EncodeWATERError(syscall.ECANCELED) // operation canceled
			}
			log.Println("worker: fairWorker: cancelConn.Read:", err)
			return wasip1.EncodeWATERError(syscall.EIO) // input/output error
		}

		// 2nd priority: sourceConn -> remoteConn
		if err := copyOnce(
			"remoteConn", // dstName
			"sourceConn", // srcName
			remoteConn,   // dst
			sourceConn,   // src
			readBuf); err != nil {
			if err == io.EOF {
				return wasip1.EncodeWATERError(0) // success, no error
			}
			if errno, ok := err.(syscall.Errno); ok {
				return wasip1.EncodeWATERError(errno)
			}
			return wasip1.EncodeWATERError(syscall.EIO) // other input/output error
		}

		// 3rd priority: remoteConn -> sourceConn
		if err := copyOnce(
			"sourceConn", // dstName
			"remoteConn", // srcName
			sourceConn,   // dst
			remoteConn,   // src
			readBuf); err != nil {
			if err == io.EOF {
				return wasip1.EncodeWATERError(0)
			}
			if errno, ok := err.(syscall.Errno); ok {
				return wasip1.EncodeWATERError(errno)
			}
			return wasip1.EncodeWATERError(syscall.EIO) // other input/output error
		}
	}
}

func copyOnce(dstName, srcName string, dst, src net.Conn, buf []byte) error {
	if len(buf) == 0 {
		buf = make([]byte, 16384) // 16k buffer for reading
	}

	nRead, readErr := src.Read(buf)
	if !(readErr == syscall.EAGAIN) { // if EAGAIN, do nothing and return
		if readErr == io.EOF {
			log.Printf("worker: copyOnce: EOF on %s", srcName)
			return io.EOF
		} else if readErr != nil {
			log.Printf("worker: copyOnce: %s.Read: %v", srcName, readErr)
			return syscall.EIO // input/output error
		}

		nWritten, writeErr := dst.Write(buf[:nRead])
		if writeErr != nil {
			log.Printf("worker: copyOnce: %s.Write: %v", dstName, writeErr)
			return syscall.EIO // no matter input/output error or EAGAIN we cannot retry async write yet
		}

		if nRead != nWritten {
			log.Printf("worker: copyOnce: %s.nRead != %s.nWritten", srcName, dstName)
			return syscall.EIO // input/output error
		}
	}

	return nil
}
