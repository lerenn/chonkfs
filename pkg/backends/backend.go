//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=backend.go -destination=backend.mock.gen.go -package backends

package backends

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
	GetAttributes(ctx context.Context, attr *fuse.Attr) syscall.Errno
	SetAttributes(ctx context.Context, in *fuse.SetAttrIn) syscall.Errno

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
	GetAttributes(ctx context.Context, attr *fuse.Attr) syscall.Errno
	SetAttributes(ctx context.Context, in *fuse.SetAttrIn) syscall.Errno
	Read(ctx context.Context, off int64) ([]byte, syscall.Errno)
	WriteCache(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno)
	Sync(ctx context.Context) syscall.Errno
}
