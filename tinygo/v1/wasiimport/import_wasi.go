//go:build wasi || wasip1

package wasiimport

import "unsafe"

// Import the Runtime-imported dial function,
// which takes iovs for network and address and
// returns a file descriptor for the dialed connection.
//
//go:wasmimport env water_dial
//go:noescape
func water_dial(
	networkIovs unsafe.Pointer, networkIovsLen size,
	addressIovs unsafe.Pointer, addressIovsLen size,
) (fd int32)

// Import the Runtime-imported fixed dialer function,
// which returns a file descriptor for the dialled connection.
//
//go:wasmimport env water_dial_fixed
//go:noescape
func water_dial_fixed() (fd int32)

// Import the Runtime-imported acceptor function.
//
//go:wasmimport env water_accept
//go:noescape
func water_accept() (fd int32)

// Import wasi_snapshot_preview1's fd_fdstat_set_flags function
// until tinygo supports it.
//
//go:wasmimport wasi_snapshot_preview1 fd_fdstat_set_flags
//go:noescape
func fd_fdstat_set_flags(fd int32, flags uint32) uint32

// Import wasi_snapshot_preview1's fd_fdstat_set_flags function
// until tinygo supports it.
//
//go:wasmimport wasi_snapshot_preview1 fd_fdstat_get
//go:noescape
func fd_fdstat_get(fd int32, buf unsafe.Pointer) uint32

//go:wasmimport wasi_snapshot_preview1 poll_oneoff
//go:noescape
func poll_oneoff(in, out unsafe.Pointer, nsubscriptions size, nevents unsafe.Pointer) errno
