package wasip1

import (
	"fmt"
	"syscall"
)

// DecodeWATERError converts a error code returned by WATER API
// into a syscall.Errno or a higher-level error in Go.
//
// It automatically detects whether the error code is a WATER error
// or a success code (positive). In case of a success code, it
// returns the code itself and a nil error.
//
// Deprecated: starting from WATM v1 API, returned errno is always
// a ground truth syscall.Errno and not multiplexed with other positive
// return values. Positive return values are always set via a pointer
// as a parameter to the function.
func DecodeWATERError(errorCode int32) (n int32, err error) {
	if errorCode >= 0 {
		n = errorCode // such that when error code is 0, it will return 0, nil
	} else {
		errorCode = -errorCode // flip the sign
		if errno, ok := mapErrno2Syscall[errno(errorCode)]; ok {
			// if the negative of the error code is a valid Errno, then it is a valid WATERErrno.
			err = errno

			// TODO: convert some special error codes to higher-level errors.
		} else {
			// otherwise, it is an unknown error code.
			err = fmt.Errorf("unknown WATERErrno %d", errorCode)
		}
	}
	return
}

// EncodeWATERError converts a syscall.Errno (positive) into a error code
// returned by WATER API (negative).
//
// Deprecated: starting from WATM v1 API, returned errno is always a ground
// truth syscall.Errno and not multiplexed with other positive return values.
func EncodeWATERError(errno syscall.Errno) int32 {
	if errno == 0 {
		return 0
	}

	// first find the corresponding Errno (there might
	// be missing Errno in the map, which means they
	// are not supported)
	if foundErrno, ok := mapSyscall2Errno[errno]; ok {
		// then convert it to the negative value of itself
		return -int32(foundErrno)
	}
	// if the errno is not found, then it is an unknown error
	return -int32(syscall.ENOSYS)
}
