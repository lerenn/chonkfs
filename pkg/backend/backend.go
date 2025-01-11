package backend

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type Root interface {
	Directory
}

type Directory interface {
	// Self

	ListEntries(ctx context.Context) ([]fuse.DirEntry, syscall.Errno)

	// Child directories

	CreateDirectory(ctx context.Context, name string) (Directory, syscall.Errno)
	GetDirectory(ctx context.Context, name string) (Directory, syscall.Errno)
	RemoveDirectory(ctx context.Context, name string) syscall.Errno

	// Child files

	CreateFile(ctx context.Context, name string) (File, syscall.Errno)
	GetFile(ctx context.Context, name string) (File, syscall.Errno)
	RemoveFile(ctx context.Context, name string) syscall.Errno
}

type File interface {
	Getattr(ctx context.Context, out *fuse.AttrOut) (errno syscall.Errno)
	Read(ctx context.Context, off int64) ([]byte, syscall.Errno)
	Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno)
}
