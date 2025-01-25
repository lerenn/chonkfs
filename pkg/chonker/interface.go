package chonker

import (
	"context"
)

// Directory is the structure of chonker regrouping the files.
type Directory interface {
	// Attributes (optional)

	GetAttributes(ctx context.Context) (DirectoryAttributes, error)
	SetAttributes(ctx context.Context, attr DirectoryAttributes) error

	// Children directories

	CreateDirectory(ctx context.Context, name string) (Directory, error)
	GetDirectory(ctx context.Context, name string) (Directory, error)
	RemoveDirectory(ctx context.Context, name string) error
	ListDirectories(ctx context.Context) ([]string, error)
	RenameDirectory(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error

	// Children files

	CreateFile(ctx context.Context, name string, chunkSize int) (File, error)
	GetFile(ctx context.Context, name string) (File, error)
	RemoveFile(ctx context.Context, name string) error
	ListFiles(ctx context.Context) ([]string, error)
	RenameFile(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error
}

// File is the structure of chonker making the link between a file and its chunks.
type File interface {
	// Attributes (optional)

	GetAttributes(ctx context.Context) (FileAttributes, error)
	SetAttributes(ctx context.Context, attr FileAttributes) error

	// Data

	Read(ctx context.Context, dest []byte, off int) ([]byte, error)
	Size(ctx context.Context) (int, error)
	Sync(ctx context.Context) error
	Truncate(ctx context.Context, size int) error
	Write(ctx context.Context, data []byte, off int, opts WriteOptions) (written int, errno error)
}

// DirectoryAttributes contains the directory attributes.
type DirectoryAttributes struct {
}

// FileAttributes contains the file attributes.
type FileAttributes struct {
	Size int
}

// WriteOptions represents the options usable for writing.
type WriteOptions struct {
	Truncate bool
	Append   bool
}
