package wasiimport

import "unsafe"

func WaterDial(
	networkType int32,
	addressCiovs unsafe.Pointer, addressCiovsLen size,
) (fd int32) {
	return water_dial(networkType, addressCiovs, addressCiovsLen)
}

func WaterAccept() (fd int32) {
	return water_accept()
}

func WaterGetAddrSuggestion(
	addressIovs unsafe.Pointer, addressIovsLen size,
) (n int32) {
	return water_get_addr_suggestion(addressIovs, addressIovsLen)
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
