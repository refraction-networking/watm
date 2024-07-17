//go:build wasi || wasip1

package wasiimport

import "unsafe"

// Import the Runtime-provided dial function,
// which takes iovs for network and address and
// returns a file descriptor for the dialed connection.
//
// Host is not expected to write to the ciovs.
//
// Used by dialers and relays.
//
//go:wasmimport env water_dial
//go:noescape
func water_dial(
	networkType int32,
	addressCiovs unsafe.Pointer, addressCiovsLen size,
) (fd int32)

// Import the Runtime-provided accept function,
// which returns a file descriptor for the next incoming
// connection of a listener managed by the Runtime.
//
// Used by listeners and relays.
//
//go:wasmimport env water_accept
//go:noescape
func water_accept() (fd int32)

// Import the Runtime-provided get address suggestion function,
// which does a scatter read of the suggested network address into
// the provided iovec.
//
// Before returning, the host is expected to write the address
// as a byte array to the iovec.
//
// Returns the length of the address written.
//
//go:wasmimport env water_get_addr_suggestion
//go:noescape
func water_get_addr_suggestion(
	addressIovs unsafe.Pointer, addressIovsLen size,
) (n int32)

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
