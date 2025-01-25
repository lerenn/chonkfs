package implementers

import (
	"context"
	"fmt"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// FileImplementer is a struct that implements every callback for file
// from github.com/hanwen/go-fuse/v2/fs. It should returns an error if the function
// is called and should be implemented. It is used to check if some non implemented
// calls are called.
type FileImplementer struct{}

// Capabilities that the FSImplementer struct should implements.
var (
	_ fs.FileAllocater       = (*FileImplementer)(nil)
	_ fs.FileFlusher         = (*FileImplementer)(nil)
	_ fs.FileFsyncdirer      = (*FileImplementer)(nil)
	_ fs.FileFsyncer         = (*FileImplementer)(nil)
	_ fs.FileGetattrer       = (*FileImplementer)(nil)
	_ fs.FileGetlker         = (*FileImplementer)(nil)
	_ fs.FileLseeker         = (*FileImplementer)(nil)
	_ fs.FilePassthroughFder = (*FileImplementer)(nil)
	_ fs.FileReaddirenter    = (*FileImplementer)(nil)
	_ fs.FileReader          = (*FileImplementer)(nil)
	_ fs.FileReleasedirer    = (*FileImplementer)(nil)
	_ fs.FileReleaser        = (*FileImplementer)(nil)
	_ fs.FileSeekdirer       = (*FileImplementer)(nil)
	_ fs.FileSetattrer       = (*FileImplementer)(nil)
	_ fs.FileSetlker         = (*FileImplementer)(nil)
	_ fs.FileSetlkwer        = (*FileImplementer)(nil)
	_ fs.FileStatxer         = (*FileImplementer)(nil)
	_ fs.FileWriter          = (*FileImplementer)(nil)
)

//nolint:unparam
func (fi FileImplementer) detectorf(skippable bool, format string, args ...interface{}) {
	if skippable {
		fmt.Printf("SKIPPABLE: FileImplementer."+format+"\n", args...)
	} else {
		fmt.Printf("NOT IMPLEMENTED: FileImplementer."+format+"\n", args...)
	}
}

// Allocate is a file callback.
func (fi FileImplementer) Allocate(_ context.Context, _ uint64, _ uint64, _ uint32) syscall.Errno {
	fi.detectorf(false, "Allocate")
	return syscall.EOPNOTSUPP
}

// Flush is a file callback.
func (fi FileImplementer) Flush(_ context.Context) syscall.Errno {
	fi.detectorf(true, "Flush")
	return syscall.EOPNOTSUPP
}

// Fsyncdir is a file callback.
func (fi FileImplementer) Fsyncdir(_ context.Context, _ uint32) syscall.Errno {
	fi.detectorf(false, "Fsyncdir")
	return syscall.EOPNOTSUPP
}

// Fsync is a file callback.
func (fi FileImplementer) Fsync(_ context.Context, _ uint32) syscall.Errno {
	fi.detectorf(false, "Fsync")
	return syscall.EOPNOTSUPP
}

// Getattr is a file callback.
func (fi FileImplementer) Getattr(_ context.Context, _ *fuse.AttrOut) syscall.Errno {
	fi.detectorf(false, "Getattr")
	return syscall.EOPNOTSUPP
}

// Getlk is a file callback.
func (fi FileImplementer) Getlk(
	_ context.Context,
	_ uint64,
	_ *fuse.FileLock,
	_ uint32,
	_ *fuse.FileLock,
) syscall.Errno {
	fi.detectorf(false, "Getlk")
	return syscall.EOPNOTSUPP
}

// Lseek is a file callback.
func (fi FileImplementer) Lseek(_ context.Context, _ uint64, _ uint32) (uint64, syscall.Errno) {
	fi.detectorf(false, "Lseek")
	return 0, syscall.EOPNOTSUPP
}

// PassthroughFd is a file callback.
func (fi FileImplementer) PassthroughFd() (int, bool) {
	fi.detectorf(true, "PassthroughFd")
	return 0, false
}

// Readdirent is a file callback.
func (fi FileImplementer) Readdirent(_ context.Context) (*fuse.DirEntry, syscall.Errno) {
	fi.detectorf(false, "Readdirent")
	return nil, syscall.EOPNOTSUPP
}

// Read is a file callback.
func (fi FileImplementer) Read(_ context.Context, _ []byte, _ int64) (fuse.ReadResult, syscall.Errno) {
	fi.detectorf(false, "Read")
	return nil, syscall.EOPNOTSUPP
}

// Releasedir is a file callback.
func (fi FileImplementer) Releasedir(_ context.Context, _ uint32) {
	fi.detectorf(true, "Releasedir")
}

// Release is a file callback.
func (fi FileImplementer) Release(_ context.Context) syscall.Errno {
	fi.detectorf(true, "Release")
	return syscall.EOPNOTSUPP
}

// Seekdir is a file callback.
func (fi FileImplementer) Seekdir(_ context.Context, _ uint64) syscall.Errno {
	fi.detectorf(false, "Seekdir")
	return syscall.EOPNOTSUPP
}

// Setattr is a file callback.
func (fi FileImplementer) Setattr(_ context.Context, _ *fuse.SetAttrIn, _ *fuse.AttrOut) syscall.Errno {
	fi.detectorf(false, "Setattr")
	return syscall.EOPNOTSUPP
}

// Setlk is a file callback.
func (fi FileImplementer) Setlk(_ context.Context, _ uint64, _ *fuse.FileLock, _ uint32) syscall.Errno {
	fi.detectorf(false, "Setlk")
	return syscall.EOPNOTSUPP
}

// Setlkw is a file callback.
func (fi FileImplementer) Setlkw(_ context.Context, _ uint64, _ *fuse.FileLock, _ uint32) syscall.Errno {
	fi.detectorf(false, "Setlkw")
	return syscall.EOPNOTSUPP
}

// Statx is a file callback.
func (fi FileImplementer) Statx(_ context.Context, _ uint32, _ uint32, _ *fuse.StatxOut) syscall.Errno {
	fi.detectorf(false, "Statx")
	return syscall.EOPNOTSUPP
}

// Write is a file callback.
func (fi FileImplementer) Write(_ context.Context, _ []byte, _ int64) (written uint32, errno syscall.Errno) {
	fi.detectorf(false, "Write")
	return 0, syscall.EOPNOTSUPP
}
