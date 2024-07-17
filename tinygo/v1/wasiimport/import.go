//go:build !wasip1 && !wasi

package wasiimport

import (
	"time"
	"unsafe"
)

var waterDialedFD int32 = -1

func SetWaterDialedFD(fd int32) {
	waterDialedFD = fd
}

func water_dial(
	int32,
	unsafe.Pointer, size,
) (fd int32) {
	return waterDialedFD
}

var waterAcceptedFD int32 = -1

func SetWaterAcceptedFD(fd int32) {
	waterAcceptedFD = fd
}

// This function should be imported from the host in WASI.
// On non-WASI platforms, it mimicks the behavior of the host
// by returning a file descriptor of preset value.
func water_accept() (fd int32) {
	return waterAcceptedFD
}

func water_get_addr_suggestion(
	unsafe.Pointer, size,
) (n int32) {
	return 0
}

// emulate the behavior when no file descriptors are
// ready and the timeout expires immediately.
func poll_oneoff(_, _ unsafe.Pointer, nsubscriptions uint32, nevents unsafe.Pointer) uint32 {
	// wait for a very short period to simulate the polling
	time.Sleep(50 * time.Millisecond)
	*(*uint32)(nevents) = nsubscriptions
	return 0
}

func fd_fdstat_set_flags(fd int32, flags uint32) uint32 {
	panic("not implemented")
}

func fd_fdstat_get(fd int32, buf unsafe.Pointer) uint32 {
	panic("not implemented")
}
