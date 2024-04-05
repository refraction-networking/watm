// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"unsafe"
)

// GOARCH=wasm currently has 64 bits pointers, but the WebAssembly host expects
// pointers to be 32 bits so we use this type alias to represent pointers in
// structs and arrays passed as arguments to WASI functions.
//
// Note that the use of an integer type prevents the compiler from tracking
// pointers passed to WASI functions, so we must use KeepAlive to explicitly
// retain the objects that could otherwise be reclaimed by the GC.
type uintptr32 = uint32

// https://github.com/WebAssembly/WASI/blob/a2b96e81c0586125cc4dc79a5be0b78d9a059925/legacy/preview1/docs.md#-size-u32
type size = uint32

// https://github.com/WebAssembly/WASI/blob/a2b96e81c0586125cc4dc79a5be0b78d9a059925/legacy/preview1/docs.md#-errno-variant
type errno = uint32

// https://github.com/WebAssembly/WASI/blob/a2b96e81c0586125cc4dc79a5be0b78d9a059925/legacy/preview1/docs.md#-filesize-u64
type filesize = uint64

// https://github.com/WebAssembly/WASI/blob/a2b96e81c0586125cc4dc79a5be0b78d9a059925/legacy/preview1/docs.md#-timestamp-u64
type timestamp = uint64

// https://github.com/WebAssembly/WASI/blob/a2b96e81c0586125cc4dc79a5be0b78d9a059925/legacy/preview1/docs.md#-clockid-variant
type clockid = uint32

const (
	clockRealtime  clockid = 0
	clockMonotonic clockid = 1
)

type eventtype = uint8

const (
	eventtypeClock eventtype = iota
	eventtypeFdRead
	eventtypeFdWrite
)

type eventrwflags = uint16

const (
	fdReadwriteHangup eventrwflags = 1 << iota
)

type userdata = uint64

// The go:wasmimport directive currently does not accept values of type uint16
// in arguments or returns of the function signature. Most WASI imports return
// an errno value, which we have to define as uint32 because of that limitation.
// However, the WASI errno type is intended to be a 16 bits integer, and in the
// event struct the error field should be of type errno. If we used the errno
// type for the error field it would result in a mismatching field alignment and
// struct size because errno is declared as a 32 bits type, so we declare the
// error field as a plain uint16.
type event struct {
	userdata    userdata
	error       uint16
	typ         eventtype
	p           [5]uint8 // padding to 16 bytes
	fdReadwrite eventFdReadwrite
}

type eventFdReadwrite struct {
	nbytes filesize
	flags  eventrwflags
}

type subclockflags = uint16

const (
	subscriptionClockAbstime subclockflags = 1 << iota
)

type subscriptionClock struct {
	id        clockid
	timeout   timestamp
	precision timestamp
	flags     subclockflags
}

type subscriptionFdReadwrite struct {
	fd int32
}

type subscription struct {
	userdata userdata
	u        subscriptionUnion
}

type subscriptionUnion [5]uint64

func (u *subscriptionUnion) eventtype() *eventtype {
	return (*eventtype)(unsafe.Pointer(&u[0]))
}

func (u *subscriptionUnion) subscriptionClock() *subscriptionClock {
	return (*subscriptionClock)(unsafe.Pointer(&u[1]))
}

func (u *subscriptionUnion) subscriptionFdReadwrite() *subscriptionFdReadwrite {
	return (*subscriptionFdReadwrite)(unsafe.Pointer(&u[1]))
}

// https://github.com/WebAssembly/WASI/blob/a2b96e81c0586125cc4dc79a5be0b78d9a059925/legacy/preview1/docs.md#-iovec-record
type iovec struct {
	buf    uintptr32
	bufLen size
}

func bytesPointer(b []byte) unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(b))
}

func makeIOVec(b []byte) unsafe.Pointer {
	return unsafe.Pointer(&iovec{
		buf:    uintptr32(uintptr(bytesPointer(b))),
		bufLen: size(len(b)),
	})
}
