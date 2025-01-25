package storage

import (
	"context"
)

// Directory represents a directory in the storage.
type Directory interface {
	// Directories

	CreateDirectory(ctx context.Context, name string) (Directory, error)
	GetDirectory(ctx context.Context, name string) (Directory, error)
	ListDirectories(ctx context.Context) (map[string]Directory, error)
	Info(ctx context.Context) (DirectoryInfo, error)
	RemoveDirectory(ctx context.Context, name string) error
	RenameDirectory(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error

	// Files

	GetFile(ctx context.Context, name string) (File, error)
	ListFiles(ctx context.Context) (map[string]File, error)
	CreateFile(ctx context.Context, name string, chunkSize int) (File, error)
	RemoveFile(ctx context.Context, name string) error
	RenameFile(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error
}

// DirectoryInfo represents the information of a directory.
type DirectoryInfo struct {
}

// File represents a file in the storage.
type File interface {
	Size(ctx context.Context) (int, error)
	WriteChunk(ctx context.Context, chunkIndex int, start int, end *int, data []byte) (int, error)
	ReadChunk(ctx context.Context, chunkIndex int, data []byte, start int, end *int) (int, error)
	ChunksCount(ctx context.Context) (int, error)
	ResizeChunksNb(ctx context.Context, size int) error
	LastChunkSize(ctx context.Context) (int, error)
	ResizeLastChunk(ctx context.Context, size int) (changed int, err error)
	Info(ctx context.Context) (FileInfo, error)
}

// FileInfo represents the information of a file.
type FileInfo struct {
	ChunkSize int
}
