//go:build wasi || wasip1

package wasiimport

import "unsafe"

//go:wasmimport env water_dial
//go:noescape
func water_dial(
	networkType uint32,
	addressCiovs unsafe.Pointer, addressCiovsLen size,
	fd unsafe.Pointer,
) errno

//go:wasmimport env water_accept
//go:noescape
func water_accept(fd unsafe.Pointer) errno

//go:wasmimport env water_get_addr_suggestion
//go:noescape
func water_get_addr_suggestion(
	addressIovs unsafe.Pointer, addressIovsLen size,
	nread unsafe.Pointer,
) errno

// TODO: remove this once tinygo provides wrapper for fd_fdstat_set
//
//go:wasmimport wasi_snapshot_preview1 fd_fdstat_set_flags
//go:noescape
func fd_fdstat_set_flags(fd int32, flags uint32) uint32

// TODO: remove this once tinygo provides wrapper for fd_fdstat_get
//
//go:wasmimport wasi_snapshot_preview1 fd_fdstat_get
//go:noescape
func fd_fdstat_get(fd int32, buf unsafe.Pointer) uint32

// TODO: remove this once tinygo provides wrapper for poll_oneoff
//
//go:wasmimport wasi_snapshot_preview1 poll_oneoff
//go:noescape
func poll_oneoff(in, out unsafe.Pointer, nsubscriptions size, nevents unsafe.Pointer) errno
