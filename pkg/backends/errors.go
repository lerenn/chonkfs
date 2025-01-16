package backends

import (
	"fmt"
	"log"
	"syscall"
)

var (
	ErrBackend                = fmt.Errorf("backend error")
	ErrReadEndBeforeReadStart = fmt.Errorf("%w: end before start", ErrBackend)
	ErrNotDirectory           = fmt.Errorf("%w: not a directory", ErrBackend)
	ErrAlreadyExists          = fmt.Errorf("%w: already exists", ErrBackend)
	ErrNoEntry                = fmt.Errorf("%w: no entry", ErrBackend)
)

type ToSyscallErrnoOptions struct {
	Logger *log.Logger
}

func ToSyscallErrno(err error, opts ToSyscallErrnoOptions) syscall.Errno {
	// If no error, returns before anything happens
	if err == nil {
		return syscall.Errno(0)
	}

	// Logs if requested
	if opts.Logger != nil {
		opts.Logger.Printf("ToSyscallErrno(err=%v)\n", err)
	}

	// Change error to errno
	switch err {
	case ErrReadEndBeforeReadStart:
		return syscall.EINVAL
	case ErrNotDirectory:
		return syscall.ENOTDIR
	case ErrAlreadyExists:
		return syscall.EEXIST
	case ErrNoEntry:
		return syscall.ENOENT
	default: // Default to EBADMSG
		return syscall.EBADMSG
	}
}
