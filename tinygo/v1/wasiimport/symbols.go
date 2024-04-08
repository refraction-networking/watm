package wasiimport

import "unsafe"

func WaterDial(
	networkIOVs unsafe.Pointer, networkIOVsLen size,
	addressIOVs unsafe.Pointer, addressIOVsLen size,
) (fd int32) {
	return water_dial(networkIOVs, networkIOVsLen, addressIOVs, addressIOVsLen)
}

func WaterDialFixed() (fd int32) {
	return water_dial_fixed()
}

func WaterAccept() (fd int32) {
	return water_accept()
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
