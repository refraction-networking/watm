// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/refraction-networking/watm/tinygo/v1/wasiimport"
)

// WASI network poller.
//
// WASI preview 1 includes a poll_oneoff host function that behaves similarly
// to poll(2) on Linux. Like poll(2), poll_oneoff is level triggered. It
// accepts one or more subscriptions to FD read or write events.
//
// Major differences to poll(2):
// - the events are not written to the input entries (like pollfd.revents), and
//   instead are appended to a separate events buffer. poll_oneoff writes zero
//   or more events to the buffer (at most one per input subscription) and
//   returns the number of events written. Although the index of the
//   subscriptions might not match the index of the associated event in the
//   events buffer, both the subscription and event structs contain a userdata
//   field and when a subscription yields an event the userdata fields will
//   match.
// - there's no explicit timeout parameter, although a time limit can be added
//   by using "clock" subscriptions.
// - each FD subscription can either be for a read or a write, but not both.
//   This is in contrast to poll(2) which accepts a mask with POLLIN and
//   POLLOUT bits, allowing for a subscription to either, neither, or both
//   reads and writes.
//
// Since poll_oneoff is similar to poll(2), the implementation here was derived
// from netpoll_aix.go.

const (
	EventFdRead  uint16 = iota + 1 // readable event
	EventFdWrite                   // writeable event
)

var (
	evts []event
	subs []subscription
)

type pollFd struct {
	fd      uintptr
	events  uint16
	revents uint16 // todo
}

func _poll(fds []pollFd, maxTimeout int64) (nevents int32, err error) {
	// Unlike poll(2), WASI's poll_oneoff doesn't accept a timeout directly. To
	// prevent it from blocking indefinitely, a clock subscription with a
	// timeout field needs to be submitted. Reserve a slot here for the clock
	// subscription, and set fields that won't change between poll_oneoff calls.

	subs = make([]subscription, 1, 128)
	evts = make([]event, 0, 128)

	timeout := &subs[0]
	eventtype := timeout.u.eventtype()
	*eventtype = eventtypeClock
	clock := timeout.u.subscriptionClock()
	clock.id = clockMonotonic
	clock.precision = 1e3

	// construct the subscriptions
	for _, fd := range fds {
		sub := subscription{}
		sub.userdata = userdata(fd.fd)
		subscription := sub.u.subscriptionFdReadwrite()
		subscription.fd = int32(fd.fd)
		if fd.events == EventFdRead {
			eventtype := sub.u.eventtype()
			*eventtype = eventtypeFdRead
		} else if fd.events == EventFdWrite {
			eventtype := sub.u.eventtype()
			*eventtype = eventtypeFdWrite
		} else {
			panic(fmt.Sprintf("invalid event type: %d", fd.events))
		}
		subs = append(subs, sub)
	}

	// If maxTimeout >= 0, we include a subscription of type Clock that we use as
	// a timeout. If maxTimeout < 0, we omit the subscription and allow poll_oneoff
	// to block indefinitely.
	pollsubs := subs
	if maxTimeout >= 0 {
		timeout := &subs[0]
		clock := timeout.u.subscriptionClock()
		clock.timeout = uint64(maxTimeout)
	} else {
		pollsubs = subs[1:]
	}

	if len(pollsubs) == 0 {
		return 0, nil
	}

	evts = evts[:len(pollsubs)]
	for i := range evts {
		evts[i] = event{}
	}

retry:
	errno := wasiimport.PollOneoff(unsafe.Pointer(&pollsubs[0]), unsafe.Pointer(&evts[0]), uint32(len(pollsubs)), unsafe.Pointer(&nevents))
	if errno != 0 {
		if errno != uint32(syscall.EINTR) {
			return 0, syscall.Errno(errno)
		}

		// If a timed sleep was interrupted, just return to
		// let the caller retry.
		if maxTimeout > 0 {
			return 0, syscall.EAGAIN
		}
		goto retry
	}

	// go through all events and see if any event.error is not ESUCCESS
	lastEvtError := uint16(0)
	for i, evt := range evts {
		if evt.error != 0 {
			fds[i].revents = evt.error
			lastEvtError = evt.error
		}
	}
	if lastEvtError != 0 {
		return int32(nevents), syscall.Errno(lastEvtError)
	}

	return int32(nevents), nil
}

func Poll(conns []Conn, events []uint16) (nevents int32, revents []uint16, err error) {
	if len(conns) != len(events) {
		return 0, nil, syscall.EINVAL
	}

	fds := make([]pollFd, len(conns))
	for i, conn := range conns {
		fds[i] = pollFd{
			fd:     uintptr(conn.Fd()),
			events: events[i],
		}
	}

	nevents, err = _poll(fds, -1)
	if err != nil {
		return nevents, nil, err
	}

	revents = make([]uint16, len(conns))
	for i, fd := range fds {
		revents[i] = fd.revents
	}

	return nevents, revents, nil
}
