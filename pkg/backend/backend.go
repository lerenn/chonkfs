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
	CreateChildDirectory(ctx context.Context, name string) (Directory, syscall.Errno)
	GetChildDirectory(ctx context.Context, name string) (Directory, syscall.Errno)
	ListDirectoryEntries(ctx context.Context) ([]fuse.DirEntry, syscall.Errno)

	CreateChildFile(ctx context.Context, name string) (File, syscall.Errno)
	GetChildFile(ctx context.Context, name string) (File, syscall.Errno)
}

type File interface {
	Read(ctx context.Context, off int64) ([]byte, syscall.Errno)
	Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno)
}
