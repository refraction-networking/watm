package net

import "syscall"

// writeFD writes data to the file descriptor fd. When a partial write occurs,
// it will continue with the remaining data until all data is written or an
// error occurs. If no progress is made in a single write call, it will return
// syscall.EIO.
//
// It is ported from (*FD).Write in golang/go/src/internal/poll/fd_unix.go
func writeFD(fd uintptr, p []byte) (int, error) {
	var nn int
	for {
		n, err := ignoringEINTRIO(syscall.Write, syscallFd(fd), p[nn:])
		if n > 0 {
			nn += n
		}
		if nn == len(p) {
			return nn, err
		}
		if err != nil {
			return nn, err
		}
		if n == 0 {
			return nn, syscall.EIO
		}

		// // TODO: retry if EAGAIN or no progress?
		// if n == 0 {
		// 	noprogress++
		// }
		// if noprogress == 10 {
		// 	return nn, syscall.EIO
		// }
		// runtime.Gosched()
	}
}

// ignoringEINTRIO is like ignoringEINTR, but just for IO calls.
func ignoringEINTRIO(fn func(fd syscallFd, p []byte) (int, error), fd syscallFd, p []byte) (int, error) {
	for {
		n, err := fn(fd, p)
		if err != syscall.EINTR {
			return n, err
		}
	}
}
