package storage

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
)

// Directory represents a directory in the storage.
type Directory interface {
	// Directories

	CreateDirectory(ctx context.Context, name string) (Directory, error)
	GetDirectory(ctx context.Context, name string) (Directory, error)
	ListDirectories(ctx context.Context) (map[string]Directory, error)
	Info(ctx context.Context) (info.Directory, error)
	RemoveDirectory(ctx context.Context, name string) error
	RenameDirectory(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error

	// Files

	GetFile(ctx context.Context, name string) (File, error)
	ListFiles(ctx context.Context) (map[string]File, error)
	CreateFile(ctx context.Context, name string, chunkSize int) (File, error)
	RemoveFile(ctx context.Context, name string) error
	RenameFile(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error

	// Underlayer

	Underlayer() Directory
}

// File represents a file in the storage.
type File interface {
	WriteChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error)
	ReadChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error)
	ResizeChunksNb(ctx context.Context, size int) error
	ResizeLastChunk(ctx context.Context, size int) (changed int, err error)
	GetInfo(ctx context.Context) (info.File, error)
	Underlayer() File
	ChunkSize() int
}
