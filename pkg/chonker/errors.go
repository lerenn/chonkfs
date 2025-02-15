package chonker

import (
	"errors"
	"fmt"
	"log"
	"syscall"
)

var (
	// ErrChonker regroups errors from chonker.
	ErrChonker = fmt.Errorf("chonker error")
	// ErrNotDirectory happens when the requested entry is not a directory.
	ErrNotDirectory = fmt.Errorf("%w: not a directory", ErrChonker)
	// ErrAlreadyExists happens when an already existing entry is making the operation fails.
	ErrAlreadyExists = fmt.Errorf("%w: already exists", ErrChonker)
	// ErrNoEntry happens when the requested entry doesn't exist.
	ErrNoEntry = fmt.Errorf("%w: no entry", ErrChonker)
)

// ToSyscallErrnoOptions is the options for ToSyscallErrno.
type ToSyscallErrnoOptions struct {
	Logger *log.Logger
}

// ToSyscallErrno turns a chonker error into a syscall.Errno.
func ToSyscallErrno(err error, opts ToSyscallErrnoOptions) syscall.Errno {
	// If no error, returns before anything happens
	if err == nil {
		return syscall.Errno(0)
	}

	// Check that the error is wrapped by ErrChonker
	if !errors.Is(err, ErrChonker) {
		// Default to EIO
		return syscall.EIO
	}

	// Logs if requested
	if opts.Logger != nil {
		opts.Logger.Printf("ToSyscallErrno(err=%v)\n", err)
	}

	// Change error to errno
	switch {
	case errors.Is(err, ErrNotDirectory):
		return syscall.ENOTDIR
	case errors.Is(err, ErrAlreadyExists):
		return syscall.EEXIST
	case errors.Is(err, ErrNoEntry):
		return syscall.ENOENT
	default:
		fmt.Println("2")
		return syscall.EIO
	}
}
