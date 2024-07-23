package wasiimport

import "unsafe"

func WaterDial(
	networkType uint32,
	addressCiovs unsafe.Pointer, addressCiovsLen size,
	fd unsafe.Pointer,
) errno {
	return water_dial(
		networkType,
		addressCiovs, addressCiovsLen,
		fd)
}

func WaterAccept(fd unsafe.Pointer) errno {
	return water_accept(fd)
}

func WaterGetAddrSuggestion(
	addressIovs unsafe.Pointer, addressIovsLen size,
	nread unsafe.Pointer,
) errno {
	return water_get_addr_suggestion(addressIovs, addressIovsLen, nread)
}

func FdFdstatSetFlags(fd int32, flags uint32) uint32 {
	return fd_fdstat_set_flags(fd, flags)
}

func FdFdstatGet(fd int32, buf unsafe.Pointer) uint32 {
	return fd_fdstat_get(fd, buf)
}

func PollOneoff(in, out unsafe.Pointer, nsubscriptions size, nevents unsafe.Pointer) errno {
	return poll_oneoff(in, out, nsubscriptions, nevents)
}
