package chonker

import (
	"fmt"
	"log"
	"syscall"
)

var (
	ErrChonker       = fmt.Errorf("chonker error")
	ErrNotDirectory  = fmt.Errorf("%w: not a directory", ErrChonker)
	ErrAlreadyExists = fmt.Errorf("%w: already exists", ErrChonker)
	ErrNoEntry       = fmt.Errorf("%w: no entry", ErrChonker)
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
	case ErrNotDirectory:
		return syscall.ENOTDIR
	case ErrAlreadyExists:
		return syscall.EEXIST
	case ErrNoEntry:
		return syscall.ENOENT
	default:
		return syscall.EIO
	}
}
