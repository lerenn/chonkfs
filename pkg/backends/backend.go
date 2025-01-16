//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=backend.go -destination=backend.mock.gen.go -package backends

package backends

import (
	"context"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type Directory interface {
	// Self

	ListEntries(ctx context.Context) ([]fuse.DirEntry, error)
	GetAttributes(ctx context.Context) (fuse.Attr, error)
	SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error

	// Child nodes
	RenameNode(ctx context.Context, name string, newParent Directory, newName string) error

	// Child directories

	CreateDirectory(ctx context.Context, name string) (Directory, error)
	GetDirectory(ctx context.Context, name string) (Directory, error)
	RemoveDirectory(ctx context.Context, name string) error

	// Child files

	CreateFile(ctx context.Context, name string) (File, error)
	GetFile(ctx context.Context, name string) (File, error)
	RemoveFile(ctx context.Context, name string) error
}

type File interface {
	GetAttributes(ctx context.Context) (fuse.Attr, error)
	SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error
	Read(ctx context.Context, start, end int) ([]byte, error)
	WriteCache(ctx context.Context, data []byte, off int) (written int, errno error)
	Sync(ctx context.Context) error
}
