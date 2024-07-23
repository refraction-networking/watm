package wasip1

import "fmt"

func ErrnoToError(errno uint32) error {
	if errno == 0 {
		return nil
	}

	if syscallErrno, ok := mapErrno2Syscall[errno]; ok {
		return syscallErrno
	}
	return fmt.Errorf("unknown errno %d", errno)
}
