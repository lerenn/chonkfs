package chonker

import (
	"context"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type Directory interface {
	// Attributes (optional)

	GetAttributes(ctx context.Context) (fuse.Attr, error)
	SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error

	// Children nodes

	ListEntries(ctx context.Context) ([]fuse.DirEntry, error)
	RenameEntry(ctx context.Context, name string, newParent Directory, newName string) error

	// Children directories

	CreateDirectory(ctx context.Context, name string) (Directory, error)
	GetDirectory(ctx context.Context, name string) (Directory, error)
	RemoveDirectory(ctx context.Context, name string) error

	// Children files

	CreateFile(ctx context.Context, name string, chunkSize int) (File, error)
	GetFile(ctx context.Context, name string) (File, error)
	RemoveFile(ctx context.Context, name string) error
}

type File interface {
	// Attributes (optional)

	GetAttributes(ctx context.Context) (fuse.Attr, error)
	SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error

	// Data

	Read(ctx context.Context, data []byte, off int) error
	Size(ctx context.Context) (int, error)
	Sync(ctx context.Context) error
	Truncate(ctx context.Context, size int) error
	Write(ctx context.Context, data []byte, off int, opts WriteOptions) (written int, errno error)
}

type WriteOptions struct {
	Truncate bool
	Append   bool
}
